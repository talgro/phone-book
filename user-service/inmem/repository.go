package inmem

import (
	"context"
	"fmt"
	"infrastructure/myerror"
	"sort"
	"sync"

	"user-service/contact"
)

type repository struct {
	mu       sync.RWMutex
	contacts map[string]contact.Contact
}

func NewUserRepository() *repository {
	return &repository{
		contacts: make(map[string]contact.Contact),
	}
}

func (r *repository) GetContact(_ context.Context, userID string, contactID string) (contact.Contact, error) {
	contactKey := getContactKey(userID, contactID)
	if c, ok := r.contacts[contactKey]; ok {
		return c, nil
	}

	return contact.Contact{}, myerror.NewNotFoundError("inmem.GetContact: contact with ID %s not found for user %s", contactID, userID)
}

func (r *repository) SearchContacts(_ context.Context, filters contact.Filters) ([]contact.Contact, error) {
	var userContacts []contact.Contact
	for _, c := range r.contacts {
		if c.UserID == filters.UserID {
			userContacts = append(userContacts, c)
		}
	}

	sort.Slice(userContacts, func(i, j int) bool {
		return userContacts[i].FirstName < userContacts[j].FirstName
	})

	var contacts []contact.Contact
	for i := filters.Offset; i < len(userContacts); i++ {
		if len(contacts) == filters.Limit {
			break
		}

		if (filters.Phone != "" && userContacts[i].Phone != filters.Phone) ||
			(filters.FirstName != "" && userContacts[i].FirstName != filters.FirstName) ||
			(filters.LastName != "" && userContacts[i].LastName != filters.LastName) ||
			(filters.Address != "" && userContacts[i].Address != filters.Address) {
			continue
		}

		contacts = append(contacts, userContacts[i])
	}

	return contacts, nil
}

func (r *repository) UpdateContact(_ context.Context, c contact.Contact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	contactKey := getContactKey(c.UserID, c.ID)
	r.contacts[contactKey] = c
	return nil
}

func (r *repository) DeleteContact(_ context.Context, userID string, contactID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	contactKey := getContactKey(userID, contactID)
	delete(r.contacts, contactKey)
	return nil
}

func (r *repository) CreateContact(_ context.Context, c contact.Contact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	contactKey := getContactKey(c.UserID, c.ID)
	r.contacts[contactKey] = c
	return nil
}

func (r *repository) IsPhoneExistsForUser(_ context.Context, userID, phone string) (bool, error) {
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
