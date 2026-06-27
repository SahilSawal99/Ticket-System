package repository

import (
	"fmt"
	"sync"

	"github.com/sahilsawal99/ticket-system/internal/model"
)

type Store struct {
	users         map[string]model.User
	tickets       map[string]model.Ticket
	mu            sync.RWMutex
	ticketCounter int
	userCounter   int
}

func NewStore() *Store {
	return &Store{
		users:   make(map[string]model.User),
		tickets: make(map[string]model.Ticket),
	}
}

var DB = NewStore()

func (s *Store) IsUsernameTaken(username string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, u := range s.users {
		if u.Username == username {
			return true
		}
	}
	return false
}

func (s *Store) GetUserByUsername(username string) (model.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, u := range s.users {
		if u.Username == username {
			return u, true
		}
	}
	return model.User{}, false
}

func (s *Store) GetUserByID(id string) (model.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	return user, ok
}

func (s *Store) AddUser(user model.User) model.User {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.userCounter++
	user.ID = fmt.Sprintf("U%d", s.userCounter)
	s.users[user.ID] = user
	return user
}

func (s *Store) AddTicket(ticket model.Ticket) model.Ticket {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ticketCounter++
	ticket.ID = fmt.Sprintf("T%d", s.ticketCounter)
	s.tickets[ticket.ID] = ticket
	return ticket
}

func (s *Store) GetTicketsByUser(userID string) []model.Ticket {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tickets []model.Ticket
	for _, t := range s.tickets {
		if t.UserID == userID {
			tickets = append(tickets, t)
		}
	}
	return tickets
}

func (s *Store) GetTicketByID(id string) (model.Ticket, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ticket, ok := s.tickets[id]
	return ticket, ok
}

func (s *Store) UpdateTicket(ticket model.Ticket) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tickets[ticket.ID] = ticket
}
