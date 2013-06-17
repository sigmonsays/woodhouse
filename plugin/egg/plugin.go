package egg
import (
    "strings"
    "fmt"
    "gihub.com/sigmonsays/woodhouse.git/plugin"
    "gihub.com/sigmonsays/woodhouse.git/config"
)

func New(cfg *config.BotConfig) *EggPlugin {
    p := &EggPlugin{
        cfg: cfg,
        TablePlugin: plugin.TablePlugin{
            Table: make(map[string]plugin.CommandFunc),
        },
    }
    p.Table["add"] = p.AddSubCmd
    p.Table["info"] = p.InfoSubCmd
    return p
}

type EggPlugin struct {
    cfg *config.BotConfig
    plugin.TablePlugin
}


func (P *EggPlugin) AddSubCmd(p plugin.PrivMsg, speak chan plugin.PrivMsg) {
    words := strings.Split(p.Message[1:], " ")
    if len(words) > 2 && words[1] == "add" {
        egg := strings.Join(words[2:], " ")
        P.cfg.Eggs = append(P.cfg.Eggs, egg)
    }
}

func (P *EggPlugin) InfoSubCmd(p plugin.PrivMsg, speak chan plugin.PrivMsg) {
    speak <- plugin.NewPrivMsg(p.Channel, fmt.Sprintf("%d eggs", len(P.cfg.Eggs)))
}
