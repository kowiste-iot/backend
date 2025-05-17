// shared/websocket/app/hub.go
package app

import (
	"context"
	"sync"
)

type Hub struct {
	// Clients map with tenant isolation
	// map[tenantID]map[userID]*Client
	clients       map[string]map[string]*Client
	mu            sync.RWMutex
	register      chan *Client
	unregister    chan *Client
	subscriptions map[string]map[string][]string
	subscribeMu   sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:       make(map[string]map[string]*Client),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		subscriptions: make(map[string]map[string][]string),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[client.tenantID]; !ok {
				h.clients[client.tenantID] = make(map[string]*Client)
			}
			h.clients[client.tenantID][client.userID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.tenantID]; ok {
				if _, ok := h.clients[client.tenantID][client.userID]; ok {
					delete(h.clients[client.tenantID], client.userID)
					close(client.send)
					// If tenant has no more users, clean up
					if len(h.clients[client.tenantID]) == 0 {
						delete(h.clients, client.tenantID)
					}
				}
			}
			h.mu.Unlock()

		case <-ctx.Done():
			return
		}
	}
}

// SendToUser sends a message to a specific user in a tenant
func (h *Hub) SendToUser(tenantID, branch, userID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if tenantClients, ok := h.clients[tenantID+branch]; ok {
		if client, ok := tenantClients[userID]; ok {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(tenantClients, userID)
			}
		}
	}
}

// GetConnectedUsers returns all connected users for a tenant
func (h *Hub) GetConnectedUsers(tenantID string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var users []string
	if tenantClients, ok := h.clients[tenantID]; ok {
		for userID := range tenantClients {
			users = append(users, userID)
		}
	}
	return users
}

// IsUserConnected checks if a specific user is connected
func (h *Hub) IsUserConnected(tenantID, userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if tenantClients, ok := h.clients[tenantID]; ok {
		_, exists := tenantClients[userID]
		return exists
	}
	return false
}

func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// SubscribeToMeasures subscribes a user to specific measures
func (h *Hub) SubscribeToMeasures(tenantID, userID string, measureIDs []string) {
	h.subscribeMu.Lock()
	defer h.subscribeMu.Unlock()

	// Initialize tenant map if it doesn't exist
	if _, ok := h.subscriptions[tenantID]; !ok {
		h.subscriptions[tenantID] = make(map[string][]string)
	}

	// Subscribe to each measure
	for _, measureID := range measureIDs {
		if _, ok := h.subscriptions[tenantID][measureID]; !ok {
			h.subscriptions[tenantID][measureID] = []string{}
		}

		// Check if user is already subscribed to avoid duplicates
		alreadySubscribed := false
		for _, id := range h.subscriptions[tenantID][measureID] {
			if id == userID {
				alreadySubscribed = true
				break
			}
		}

		// Add user to subscription list if not already there
		if !alreadySubscribed {
			h.subscriptions[tenantID][measureID] = append(h.subscriptions[tenantID][measureID], userID)
		}
	}
}

// UnsubscribeFromMeasures unsubscribes a user from specific measures
func (h *Hub) UnsubscribeFromMeasures(tenantID, userID string, measureIDs []string) {
	h.subscribeMu.Lock()
	defer h.subscribeMu.Unlock()

	if _, ok := h.subscriptions[tenantID]; !ok {
		return
	}

	// If measureIDs is empty, unsubscribe from all
	if len(measureIDs) == 0 {
		for measureID, users := range h.subscriptions[tenantID] {
			// Remove user from the slice
			newUsers := make([]string, 0, len(users))
			for _, id := range users {
				if id != userID {
					newUsers = append(newUsers, id)
				}
			}

			// Update or cleanup
			if len(newUsers) == 0 {
				delete(h.subscriptions[tenantID], measureID)
			} else {
				h.subscriptions[tenantID][measureID] = newUsers
			}
		}

		// Clean up empty tenant
		if len(h.subscriptions[tenantID]) == 0 {
			delete(h.subscriptions, tenantID)
		}
		return
	}

	// Unsubscribe from specific measures
	for _, measureID := range measureIDs {
		if users, ok := h.subscriptions[tenantID][measureID]; ok {
			// Remove user from the slice
			newUsers := make([]string, 0, len(users))
			for _, id := range users {
				if id != userID {
					newUsers = append(newUsers, id)
				}
			}

			// Update or cleanup
			if len(newUsers) == 0 {
				delete(h.subscriptions[tenantID], measureID)
			} else {
				h.subscriptions[tenantID][measureID] = newUsers
			}
		}
	}

	// Clean up empty tenant
	if len(h.subscriptions[tenantID]) == 0 {
		delete(h.subscriptions, tenantID)
	}
}

// GetSubscribedUsers returns all users subscribed to a specific measure
func (h *Hub) GetSubscribedUsers(tenantID, branch, measureID string) []string {
	h.subscribeMu.RLock()
	defer h.subscribeMu.RUnlock()
	//TODO: get tenant and branch
	if tenantSubs, ok := h.subscriptions[tenantID+branch]; ok {
		if users, ok := tenantSubs[measureID]; ok {
			// Return a copy of the slice to prevent external modification
			result := make([]string, len(users))
			copy(result, users)
			return result
		}
	}

	return []string{}
}
