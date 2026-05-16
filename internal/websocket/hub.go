package websocket

import (
	"sync"
)

type RoomID = int

type Hub struct {
	clients         map[*Client]bool
	rooms           map[RoomID]map[*Client]bool
	userConnections map[int]int
	registerChan    chan *Client
	unregisterChan  chan *Client
	broadcastChan   chan Event
	mu              sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:         make(map[*Client]bool),
		rooms:           make(map[RoomID]map[*Client]bool),
		userConnections: make(map[int]int),
		registerChan:    make(chan *Client),
		unregisterChan:  make(chan *Client),
		broadcastChan:   make(chan Event),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.registerChan:
			hub.addClient(client)

		case client := <-hub.unregisterChan:
			hub.removeClient(client)

		case event := <-hub.broadcastChan:
			hub.handleBroadcast(event)
		}
	}
}

func (hub *Hub) Register(client *Client) {
	hub.registerChan <- client
}

func (hub *Hub) Unregister(client *Client) {
	hub.unregisterChan <- client
}

func (hub *Hub) Broadcast(event Event) {
	hub.broadcastChan <- event
}

func (hub *Hub) addClient(client *Client) {
	hub.mu.Lock()

	hub.clients[client] = true

	if hub.rooms[client.RoomID] == nil {
		hub.rooms[client.RoomID] = make(map[*Client]bool)
	}
	hub.rooms[client.RoomID][client] = true
	hub.mu.Unlock()

	hub.userConnections[client.UserID]++
	if hub.userConnections[client.UserID] == 1 {
		hub.handleBroadcast(
			Event{
				Type:   EventIsOnline,
				RoomID: client.RoomID,
				Payload: UserPresencePayload{
					UserID: client.UserID,
				},
			},
		)
	}
}

func (hub *Hub) removeClient(client *Client) {
	hub.mu.Lock()

	if _, ok := hub.clients[client]; !ok {
		hub.mu.Unlock()
		return
	}
	defer close(client.SendChan)

	delete(hub.clients, client)

	if roomClients, ok := hub.rooms[client.RoomID]; ok {
		delete(roomClients, client)

		if len(roomClients) == 0 {
			delete(hub.rooms, client.RoomID)
		}
	}
	hub.mu.Unlock()

	hub.userConnections[client.UserID]--
	if hub.userConnections[client.UserID] <= 0 {
		delete(hub.userConnections, client.UserID)
		hub.handleBroadcast(
			Event{
				Type:   EventIsOffline,
				RoomID: client.RoomID,
				Payload: UserPresencePayload{
					UserID: client.UserID,
				},
			},
		)
	}
}

func (hub *Hub) handleBroadcast(event Event) {
	hub.mu.RLock()
	roomClients := hub.rooms[event.RoomID]

	clients := make([]*Client, 0, len(roomClients))
	for client := range roomClients {
		clients = append(clients, client)
	}
	hub.mu.RUnlock()

	for _, client := range clients {
		select {
		case client.SendChan <- event:
		default:
			hub.Unregister(client)
		}
	}
}
