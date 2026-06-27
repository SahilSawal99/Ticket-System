package service

import (
	"errors"

	"github.com/sahilsawal99/ticket-system/internal/auth"
	"github.com/sahilsawal99/ticket-system/internal/model"
	"github.com/sahilsawal99/ticket-system/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameTaken           = errors.New("username taken")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrTicketNotFound          = errors.New("ticket not found")
	ErrForbidden               = errors.New("forbidden")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

const (
	StatusOpen       = "open"
	StatusInProgress = "in_progress"
	StatusClosed     = "closed"
)

func RegisterUser(username, password string) error {
	if repository.DB.IsUsernameTaken(username) {
		return ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	repository.DB.AddUser(model.User{Username: username, PasswordHash: string(hash)})
	return nil
}

func AuthenticateUser(username, password string) (string, error) {
	user, found := repository.DB.GetUserByUsername(username)
	if !found {
		return "", ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", ErrInvalidCredentials
	}

	return auth.GenerateToken(user.ID)
}

func CreateTicket(userID, title, description string) model.Ticket {
	ticket := model.Ticket{
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      StatusOpen,
	}
	return repository.DB.AddTicket(ticket)
}

func ListTickets(userID string) []model.Ticket {
	tickets := repository.DB.GetTicketsByUser(userID)
	if tickets == nil {
		return []model.Ticket{}
	}
	return tickets
}

func GetTicket(userID, id string) (model.Ticket, error) {
	ticket, found := repository.DB.GetTicketByID(id)
	if !found {
		return model.Ticket{}, ErrTicketNotFound
	}
	if ticket.UserID != userID {
		return model.Ticket{}, ErrForbidden
	}
	return ticket, nil
}

func UpdateTicketStatus(userID, id, status string) (model.Ticket, error) {
	ticket, err := GetTicket(userID, id)
	if err != nil {
		return model.Ticket{}, err
	}

	if ticket.Status == StatusClosed {
		return model.Ticket{}, errors.New("cannot reopen a closed ticket")
	}

	validTransition := false
	if ticket.Status == StatusOpen && (status == StatusInProgress || status == StatusClosed) {
		validTransition = true
	} else if ticket.Status == StatusInProgress && status == StatusClosed {
		validTransition = true
	} else if ticket.Status == status {
		validTransition = true
	}

	if !validTransition {
		return model.Ticket{}, ErrInvalidStatusTransition
	}

	ticket.Status = status
	repository.DB.UpdateTicket(ticket)
	return ticket, nil
}
