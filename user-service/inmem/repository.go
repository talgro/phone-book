package inmem

import (
	"context"
	"fmt"
	"infrastructure/myerror"
	"sync"

	"user-service/contact"
)

type userRepository struct {
	mu       sync.RWMutex
	contacts map[string]contact.Contact
}

func (r *userRepository) GetContact(_ context.Context, ID string) (contact.Contact, error) {
	for _, c := range r.contacts {
		if c.ID == ID {
			return c, nil
		}
	}

	return contact.Contact{}, myerror.NewNotFoundError("inmem.GetContact: contact with ID %s not found", ID)
}

func (r *userRepository) UpdateContact(ctx context.Context, c contact.Contact) error {
	//TODO implement me
	panic("implement me")
}

func NewUserRepository() *userRepository {
	return &userRepository{
		contacts: make(map[string]contact.Contact),
	}
}

func (r *userRepository) CreateContact(_ context.Context, c contact.Contact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	contactKey := getContactKey(c.UserID, c.Phone)
	r.contacts[contactKey] = c
	return nil
}

func (r *userRepository) IsPhoneExistsForUser(_ context.Context, userID, phone string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	contactKey := getContactKey(userID, phone)

	_, ok := r.contacts[contactKey]
	return ok, nil
}

func getContactKey(userID, phone string) string {
	return fmt.Sprintf("%s:%s", userID, phone)
}
