package plugin

func NewPrivMsg(Channel, Message string) PrivMsg {
    return PrivMsg{
        Channel: Channel,
        Message: Message,
    }
}

type PrivMsg struct {
    Channel,
    Message string
}
