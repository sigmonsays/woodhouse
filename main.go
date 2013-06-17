package main
import (
    "fmt"
    "time"
    "flag"
    "os"
    "strings"
    "crypto/tls"
    "math/rand"
    "bytes"
    "text/template"
    "regexp"
    irc "github.com/fluffle/goirc/client"
    
    "launchpad.net/goyaml"
    "github.com/sigmonsays/woodhouse/config"
    "github.com/sigmonsays/woodhouse/plugin"

    egg_plugin "github.com/sigmonsays/woodhous/plugin/egg"
)


// used to template responses from eggs
type EggMsg struct {
    Name string
}


func main() {
    configFile := ""
    opt_configFile := flag.String("config", "/etc/woodhouse/ircbot.yaml", "bot config file")
    flag.Parse()
    fmt.Println("Starting..")

    
    cfg := config.GetDefaultConfig()
    if opt_configFile != nil {
        err := cfg.LoadYaml(*opt_configFile)
        if err != nil {
            fmt.Println("LoadYaml error", err)
        }
        configFile = *opt_configFile
    }

    config.PrintConfig(cfg)

    c := irc.SimpleClient(cfg.NickName, cfg.UserName, cfg.FullName)

    if strings.Contains(cfg.ServerAddress, ":6697") {
        fmt.Println("Enabling ssl")
        c.SSL = true
        c.SSLConfig = &tls.Config{
            InsecureSkipVerify: true,
        }
    }
    on_connect := func(conn *irc.Conn, line *irc.Line) {
            if len(cfg.OnConnect) > 0 {
                fmt.Printf("OnConnect %#v\n", cfg.OnConnect)
                conn.Raw(cfg.OnConnect)
            }
            conn.Join(cfg.ChannelName) 
        }
    c.AddHandler(irc.CONNECTED, on_connect)

    // when people join
    join := make(chan irc.Line, 10)
    c.AddHandler("JOIN",
        func(conn *irc.Conn, line *irc.Line) { join <- *line })

    quit := make(chan bool)
    c.AddHandler(irc.DISCONNECTED,
        func(conn *irc.Conn, line *irc.Line) { quit <- true })

    privmsg := make(chan irc.Line, 10)
    c.AddHandler("PRIVMSG",
        func(conn *irc.Conn, line *irc.Line) { privmsg <- *line })

    // send to this to make the bot speak
    speak := make(chan plugin.PrivMsg, 10)


    // send to this to process commands (!command ...)
    command := make(chan irc.Line, 10)


    // how often to write running configuration to disk
    sync_interval := time.Tick(cfg.SyncInterval * time.Second)

    handler := plugin.NewCommandHandler(speak)
    handler.PluginDir = cfg.PluginDir
    handler.PluginTimeout = cfg.PluginTimeout

    egg_handler := egg_plugin.New(cfg)
    handler.Register("egg", egg_handler.Dispatch)


    handler.Register("ping", func(p plugin.PrivMsg, speak chan plugin.PrivMsg) {
        speak <- plugin.NewPrivMsg(p.Channel, "pong")
    })

    web := NewWebHandler(cfg, speak)
    if err := web.StartServer(); err != nil {
        fmt.Println("web start error", err)
    }
        
    for {
        fmt.Println("Connecting to", cfg.ServerAddress)

        err := c.Connect(cfg.ServerAddress, cfg.ChannelPassword)
        if err != nil {
            fmt.Println("connect error", cfg.ServerAddress, err)
            time.Sleep(1 * time.Second)
            continue
        }

        var SavedConfig []byte

        command_re := regexp.MustCompile("^([0-9a-z]+)")
        
        for {
            select {
            case <- quit:
                break

            case <- sync_interval:

                d, err := goyaml.Marshal(cfg)
                if err != nil {
                    fmt.Println("error", err)
                    continue
                }

                if bytes.Compare(SavedConfig, d) == 0 {
                    continue
                }

                tmpfile := configFile + ".tmp"
                f, err := os.Create(tmpfile)
                if err != nil {
                    fmt.Println("error", f, err)
                    continue
                }
                f.Write(d)
                f.Close()
                os.Rename(tmpfile, configFile)
                SavedConfig = d
                fmt.Println("Saved configuration updated")

            case p := <- speak:
                fmt.Printf(">>> [%s %s] %s\n", p.Channel, cfg.NickName, p.Message)
                c.Privmsg(p.Channel, p.Message)

            case line := <- command:
                Channel := line.Args[0]
                Msg := line.Args[1]
                words := strings.Split(Msg[1:], " ")
                result := command_re.FindStringSubmatch(words[0])
                if len(result) > 0 {
                    cmdname := result[0]
                    p := plugin.NewPrivMsg(Channel, Msg)
                    handler.Dispatch(cmdname, p)
                }

            case line := <- join:
                if line.Nick == cfg.NickName {
                    continue
                }
                fmt.Printf("%#v joined the chat\n", line.Nick)
                Channel := line.Args[0]
                speak <- plugin.NewPrivMsg(Channel, Greeting(line.Nick))

            case line := <- privmsg:
                fmt.Printf("<<< [%s %s] %s\n", line.Args[0], line.Nick, line.Args[1])
                Channel := line.Args[0]
                Msg := line.Args[1]

                if Channel == cfg.UserName { // direct message
                    speak <- plugin.NewPrivMsg(cfg.ChannelName, Msg)
                    continue
                }

                if strings.Contains(Msg, cfg.UserName) {
                    b := bytes.NewBuffer(nil)
                    t := template.Must(template.New("egg").Parse(cfg.Eggs[rand.Int() % len(cfg.Eggs)]))
                    p := EggMsg{
                        Name: line.Nick,
                    }
                    err := t.Execute(b, p)
                    if err != nil {
                        fmt.Println("template egg error", err)
                    }
                    speak <- plugin.NewPrivMsg(Channel, b.String())
                } else if strings.HasPrefix(Msg, "!") {
                    command <- line
                }
            }
        }
    }
}

