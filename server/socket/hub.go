package socket

// Hub maintains the set of active clients and enables messages to be broadcast to clients
type Hub struct {
    // Registered clients
    Clients map[*Client]bool

    // a channel to broadcast messages to clients
    Broadcast chan []byte

    // Register requests from clients
    register chan *Client

    // Unregister requests from clients
    unregister chan *Client
}

func NewSocketHub() *Hub {
    hub := &Hub{
        Clients:    make(map[*Client]bool),
        Broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
    go hub.run()
    return hub
}

func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.Clients[client] = true
        case client := <-h.unregister:
            if _, ok := h.Clients[client]; ok {
                delete(h.Clients, client)
                close(client.send)
            }
        case message := <-h.Broadcast:
            for client := range h.Clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.Clients, client)
                }
            }
        }
    }
}