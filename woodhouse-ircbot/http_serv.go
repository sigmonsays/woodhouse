package main
import (
    "fmt"
    "net/http"
    "github.com/sigmonsays/woodhouse.git/plugin"
    "github.com/sigmonsays/woodhouse.git/config"
)
const HelpText = `help
--------------------------------------------
/eggs
/speak?channel=CHANNEL&message=MESSAGE
`

func (h *WebHandler) handler(w http.ResponseWriter, req *http.Request) {
    if err := req.ParseForm(); err != nil {
        fmt.Println("error parsing form", err)
        return
    }

    if req.Method == "GET" && req.URL.Path == "/speak" {
        channel := "#" + req.FormValue("channel")
        message := req.FormValue("message")
        if channel != "" && message != "" {
            h.Speak <- plugin.NewPrivMsg(channel, message)
            fmt.Fprintf(w, "message sent to channel=%s [%s]\n", channel, message)
        } else {
            fmt.Fprintf(w, "required parameters missing\n")
        }

    } else if req.Method == "GET" && req.URL.Path == "/eggs" {
        for n, egg := range h.Cfg.Eggs {
            fmt.Fprintf(w, "%03d %s\n", n, egg)
        }
    } else if req.Method == "GET" && req.URL.Path == "/help" {
        fmt.Fprint(w, HelpText)

    } else {
        w.Header().Set("Content-Type", "text/plain")
        w.Write([]byte("This is an example server.\n"))
    }
}

func (h *WebHandler) StartServer() error {
    go func() {
        err := h.Server.ListenAndServeTLS("ssl/cert.pem", "ssl/key.pem")
        fmt.Println("Server error",  err)
    }()
    return nil
}

type WebHandler struct {
    Cfg *config.BotConfig
    Server *http.Server
    Speak chan plugin.PrivMsg
}


func NewWebHandler(Cfg *config.BotConfig, speak chan plugin.PrivMsg) WebHandler {
    mux := http.NewServeMux()

    srv := &http.Server{
        Addr: ":1443",
        Handler: mux,
    }
    wh := WebHandler{
        Cfg: Cfg,
        Server: srv,
        Speak: speak,
    }
    mux.HandleFunc("/", wh.handler)
    return wh
}

