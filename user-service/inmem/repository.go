package inmem

import (
	"context"
	"fmt"
	"infrastructure/myerror"
	"sync"

	"user-service/contact"
)

type contactRepository struct {
	mu       sync.RWMutex
	contacts map[string]contact.Contact
}

func NewUserRepository() *contactRepository {
	return &contactRepository{
		contacts: make(map[string]contact.Contact),
	}
}

func (r *contactRepository) GetContact(_ context.Context, userID string, contactID string) (contact.Contact, error) {
	contactKey := getContactKey(userID, contactID)
	if c, ok := r.contacts[contactKey]; ok {
		return c, nil
	}

	return contact.Contact{}, myerror.NewNotFoundError("inmem.GetContact: contact with ID %s not found for user %s", contactID, userID)
}

func (r *contactRepository) UpdateContact(_ context.Context, c contact.Contact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	contactKey := getContactKey(c.UserID, c.ID)
	r.contacts[contactKey] = c
	return nil
}

func (r *contactRepository) CreateContact(_ context.Context, c contact.Contact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	contactKey := getContactKey(c.UserID, c.ID)
	r.contacts[contactKey] = c
	return nil
}

func (r *contactRepository) IsPhoneExistsForUser(_ context.Context, userID, phone string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, c := range r.contacts {
		if c.UserID == userID && c.Phone == phone {
			return true, nil
		}
	}

	return false, nil
}

func getContactKey(userID, contactID string) string {
	return fmt.Sprintf("%s:%s", userID, contactID)
}
