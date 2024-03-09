package contactmanaging

import (
	"context"
	"infrastructure/myerror"
	"strings"
	"time"
	"user-service/contact"
)

type Service interface {
	CreateContact(ctx context.Context, c contact.Contact) (string, error)
	UpdateContact(ctx context.Context, c contact.Contact) error
}

// Create

type createContactRequest struct {
	UserID    string
	Phone     string
	FirstName string
	LastName  string
	Address   string
}

func (r createContactRequest) Validate() error {
	var errorMessages []string

	if r.UserID == "" {
		errorMessages = append(errorMessages, "userID is required")
	}

	if r.Phone == "" {
		errorMessages = append(errorMessages, "phone is required")
	}

	if r.FirstName == "" {
		errorMessages = append(errorMessages, "firstName is required")
	}

	if r.LastName == "" {
		errorMessages = append(errorMessages, "lastName is required")
	}

	if r.Address == "" {
		errorMessages = append(errorMessages, "address is required")
	}

	if len(errorMessages) > 0 {
		return myerror.NewBadRequestError("invalid request: %s", strings.Join(errorMessages, ", "))
	}

	return nil
}

func (r createContactRequest) ToContact() contact.Contact {
	return contact.Contact{
		UserID:    r.UserID,
		Phone:     r.Phone,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Address:   r.Address,
	}
}

type createContactResponse struct {
	ID string
}

func endpointCreateContact(ctx context.Context, s Service, request createContactRequest) (createContactResponse, error) {
	if err := request.Validate(); err != nil {
		return createContactResponse{}, myerror.Wrap(err, "endpointCreateContact")
	}

	id, err := s.CreateContact(ctx, request.ToContact())
	if err != nil {
		return createContactResponse{}, myerror.Wrap(err, "endpointCreateContact")
	}

	return createContactResponse{
		ID: id,
	}, nil
}

// Update

type updateContactRequest struct {
	UserID          string
	ContactID       string
	Phone           string
	FirstName       string
	LastName        string
	Address         string
	UpdateAtVersion time.Time
}

func (r updateContactRequest) Validate() error {
	var errorMessages []string

	if r.UserID == "" {
		errorMessages = append(errorMessages, "userID is required")
	}

	if r.ContactID == "" {
		errorMessages = append(errorMessages, "contactID is required")
	}

	if r.Phone == "" {
		errorMessages = append(errorMessages, "phone is required")
	}

	if r.FirstName == "" {
		errorMessages = append(errorMessages, "firstName is required")
	}

	if r.LastName == "" {
		errorMessages = append(errorMessages, "lastName is required")
	}

	if r.Address == "" {
		errorMessages = append(errorMessages, "address is required")
	}

	if len(errorMessages) > 0 {
		return myerror.NewBadRequestError("invalid request: %s", strings.Join(errorMessages, ", "))
	}

	return nil
}

func (r updateContactRequest) ToContact() contact.Contact {
	return contact.Contact{
		UserID:    r.UserID,
		ID:        r.ContactID,
		Phone:     r.Phone,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Address:   r.Address,
		UpdatedAt: r.UpdateAtVersion,
	}
}

func endpointUpdateContact(ctx context.Context, s Service, request updateContactRequest) error {
	if err := request.Validate(); err != nil {
		return myerror.Wrap(err, "endpointUpdateContact")
	}

	if err := s.UpdateContact(ctx, request.ToContact()); err != nil {
		return myerror.Wrap(err, "endpointUpdateContact")
	}

	return nil
}
