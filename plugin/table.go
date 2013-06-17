package plugin
import (
    "strings"
)

type TablePlugin struct {
    Table map[string] CommandFunc
}

func (P *TablePlugin) LookupTable(p PrivMsg) CommandFunc {
    words := strings.Split(p.Message[1:], " ")
    if len(words) < 2 {
        return nil
    }
    cmd := words[1]
    fun, ok := P.Table[cmd]
    if !ok {
        return nil
    }
    return fun
}

func (P *TablePlugin) Dispatch(p PrivMsg, speak chan PrivMsg) {
    fun := P.LookupTable(p)
    if fun == nil {
        speak <- NewPrivMsg(p.Channel, "Usage: egg [add|info]")
        return
    }
    fun(p, speak)
}
