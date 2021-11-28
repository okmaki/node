package adapter

import (
	"github.com/gorilla/websocket"
	"github.com/okmaki/node/log/internal/core"
	"github.com/okmaki/node/log/internal/port"
)

type Client struct {
	con    *websocket.Conn
	log    chan core.Log
	Filter core.LogFilter
}

func (c *Client) GetFilter() core.LogFilter {
	return c.Filter
}

func (c *Client) Send(log core.Log) {
	c.log <- log
}

func (c *Client) Disconnect() {
	c.con.Close()
	close(c.log)
}

type Hub struct {
	join     chan port.HubClient
	leave    chan port.HubClient
	log      chan core.Log
	onLog    func(core.Log)
	shutdown chan bool
}

func NewHubAdapter(onLog func(core.Log)) *Hub {
	return &Hub{
		join:     make(chan port.HubClient),
		leave:    make(chan port.HubClient),
		log:      make(chan core.Log),
		onLog:    onLog,
		shutdown: make(chan bool),
	}
}

func (h *Hub) Start() error {
	clients := make(map[port.HubClient]core.LogFilter)

	for {
		select {
		case client := <-h.join:
			clients[client] = client.GetFilter()
		case client := <-h.leave:
			delete(clients, client)
			client.Disconnect()
		case log := <-h.log:
			h.onLog(log)
			for client, filter := range clients {
				if !filter.ShouldFilter(log) {
					client.Send(log)
				}
			}
		case <-h.shutdown:
			for client := range clients {
				delete(clients, client)
				client.Disconnect()
			}
			return nil
		}
	}
}

func (h *Hub) Shutdown() {
	h.shutdown <- true
}

func (h *Hub) Join(client port.HubClient) {
	h.join <- client
}

func (h *Hub) Leave(client port.HubClient) {
	h.leave <- client
}

func (h *Hub) Notify(log core.Log) {
	h.log <- log
}
