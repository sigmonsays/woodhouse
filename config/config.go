package config
import (
    "fmt"
    "os"
    "bytes"
    "time"
    "launchpad.net/goyaml"
)
type BotConfig struct {

    NickName, UserName, FullName string

    PluginTimeout time.Duration
    SyncInterval time.Duration
    OnConnect string

    ServerAddress string
    ChannelName string
    ChannelPassword string

    PluginDir string
    Eggs []string
}


func (c *BotConfig) LoadDefault() {
    *c = *GetDefaultConfig()
}

func (c *BotConfig) LoadYaml(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }

    b := bytes.NewBuffer(nil)
    _, err = b.ReadFrom(f)
    if err != nil {
        return err
    }

    if err := c.LoadYamlBuffer(b.Bytes()); err != nil  {
        return err
    }

    if err := c.FixupConfig(); err != nil {
        return err
    }

    return nil
}

func (c *BotConfig) FixupConfig() error {

    return nil
}


func (c *BotConfig) LoadYamlBuffer(buf []byte) error {
    err := goyaml.Unmarshal(buf, c)
    if err != nil {
        return err
    }
    return nil
}

func (c *BotConfig) PrintYaml() {
    PrintConfig(c)
}

func GetDefaultConfig() *BotConfig {
    d := &BotConfig{}
    d.PluginDir = "./plugin"

    d.NickName = "steve"
    d.UserName = d.NickName
    d.FullName = d.NickName
    d.PluginTimeout = 2
    d.OnConnect = ""

    d.ServerAddress = "localhost:1025"
    d.ChannelName = "#woodhouse-testing"
    d.SyncInterval = 5

    return d
}

func PrintConfig(conf *BotConfig) {
    d, err := goyaml.Marshal(conf)
    if err != nil {
        fmt.Println("Marshal error", err)
        return
    }
    fmt.Println("-- Configuration --")
    fmt.Println(string(d))
}
