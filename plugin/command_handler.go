package plugin
import (
    "fmt"
    "time"
    "strings"
    "bytes"
    "os"
    "os/exec"
)
type CommandFunc func (PrivMsg, chan PrivMsg)

type Command struct {
    Run CommandFunc
}


type CommandReply struct {
    Message string
}
func NewCommandHandler(reply chan PrivMsg) *CommandHandler {
    return &CommandHandler{
        Registry: make(map[string]Command),
        Response: reply,
    }
}
type CommandHandler struct {
    PluginTimeout time.Duration
    PluginDir string
    Response chan PrivMsg
    Registry map[string]Command 
}

func (c CommandHandler) Register(name string, cmdf CommandFunc) {
    c.Registry[name] = Command{
        Run: cmdf,
    }
}

func (c CommandHandler) Dispatch(name string, p PrivMsg) {
    go func() {
        cmd, ok := c.Registry[name]
        if !ok {
            c.ScriptPlugin(name, p)
            return
        }
        cmd.Run(p, c.Response)
    }()
}

// plugins receive their input on stdin and the scripts output
// is sent as text to the channel.
// exit code of 0 indicates success, a non-zero exit code is an error state
// and no output is sent
func (c CommandHandler) ScriptPlugin(name string, p PrivMsg) {
    plugin := fmt.Sprintf("%s/%s", c.PluginDir, name)
    cmd := exec.Command(plugin)
    cmd.Dir = c.PluginDir
    cmd.Stdin = bytes.NewBuffer([]byte(p.Message))
    cmd.Stderr = os.Stderr
    output := bytes.NewBuffer(nil)
    cmd.Stdout = output

    finished := make(chan error)
    timeout := time.After(c.PluginTimeout * time.Second)
    go func() {
        err := cmd.Run()
        finished <- err
    }()
    select {
    case <- timeout:
        cmd.Process.Kill()
        fmt.Println("Process timeout error:", plugin)
    case err := <- finished:
        if err == nil {
            for _, response := range strings.Split(output.String(), "\n") {
                if len(response) > 0 {
                    fmt.Println("response:", response)
                    c.Response <- NewPrivMsg(p.Channel, response)
                }
            }
        } else {
            fmt.Printf("command [%#v] error %s", plugin, err)
        }
    }

}
