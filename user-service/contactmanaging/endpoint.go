package contactmanaging

import (
	"context"
	"fmt"
	"infrastructure/myerror"
	"regexp"
	"strings"
	"time"
	"user-service/contact"
)

const (
	LimitMaxContacts = 10 // LimitMax is the maximum number of contacts that can be returned
)

type Service interface {
	CreateContact(ctx context.Context, c contact.Contact) (string, error)
	UpdateContact(ctx context.Context, c contact.Contact) error
	GetContact(ctx context.Context, userID, contactID string) (contact.Contact, error)
	SearchContacts(context.Context, contact.Filters) (contacts []contact.Contact, err error)
	DeleteContact(ctx context.Context, userID, contactID string) error
}

// Create

type createContactRequest struct {
	UserID    string
	Phone     string
	FirstName string
	LastName  string
	Address   string
}

func validatePhoneOnlyDigits(phone string) bool {
	regex := regexp.MustCompile(`^[0-9]+$`)
	return regex.MatchString(phone)
}

func (r createContactRequest) Validate() error {
	var errorMessages []string

	if r.UserID == "" {
		errorMessages = append(errorMessages, "userID is required")
	}

	if r.Phone == "" || !validatePhoneOnlyDigits(r.Phone) {
		errorMessages = append(errorMessages, "phone is required and must be digits only")
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

	if r.Phone == "" || !validatePhoneOnlyDigits(r.Phone) {
		errorMessages = append(errorMessages, "phone is required and must be digits only")
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

// Get

type getContactRequest struct {
	UserID    string
	ContactID string
}

func (r getContactRequest) Validate() error {
	var errorMessages []string

	if r.UserID == "" {
		errorMessages = append(errorMessages, "userID is required")
	}

	if r.ContactID == "" {
		errorMessages = append(errorMessages, "contactID is required")
	}

	if len(errorMessages) > 0 {
		return myerror.NewBadRequestError("invalid request: %s", strings.Join(errorMessages, ", "))
	}

	return nil
}

type getContactResponse struct {
	UserID    string
	ID        string
	Phone     string
	FirstName string
	LastName  string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func contactToGetContactResponse(c contact.Contact) getContactResponse {
	return getContactResponse{
		UserID:    c.UserID,
		ID:        c.ID,
		Phone:     c.Phone,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Address:   c.Address,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func endpointGetContact(ctx context.Context, s Service, request getContactRequest) (getContactResponse, error) {
	if err := request.Validate(); err != nil {
		return getContactResponse{}, myerror.Wrap(err, "endpointGetContact")
	}

	c, err := s.GetContact(ctx, request.UserID, request.ContactID)
	if err != nil {
		return getContactResponse{}, myerror.Wrap(err, "endpointGetContact")
	}

	return contactToGetContactResponse(c), nil
}

// Search

type searchContactsRequest struct {
	UserID    string
	Phone     string
	FirstName string
	LastName  string
	Address   string
	Limit     int
	Offset    int
}

func (r searchContactsRequest) Validate() error {
	var errorMessages []string

	if r.UserID == "" {
		errorMessages = append(errorMessages, "userID is required")
	}

	if r.Phone != "" && !validatePhoneOnlyDigits(r.Phone) {
		errorMessages = append(errorMessages, "phone must be digits only")
	}

	if r.Limit < 0 || r.Limit > LimitMaxContacts {
		errorMessages = append(errorMessages, fmt.Sprintf("limit must be a positive number smaller than or equal to %d", LimitMaxContacts))
	}

	if r.Offset < 0 {
		errorMessages = append(errorMessages, "offset must be greater than or equal to 0")
	}

	if len(errorMessages) > 0 {
		return myerror.NewBadRequestError("invalid request: %s", strings.Join(errorMessages, ", "))
	}

	return nil
}

func (r searchContactsRequest) ToFilters() contact.Filters {
	return contact.Filters{
		UserID:    r.UserID,
		Phone:     r.Phone,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Address:   r.Address,
		Limit:     r.Limit,
		Offset:    r.Offset,
	}
}

type searchContactsResponse struct {
	Contacts []getContactResponse
}

func endpointSearchContacts(ctx context.Context, s Service, request searchContactsRequest) (searchContactsResponse, error) {
	if err := request.Validate(); err != nil {
		return searchContactsResponse{}, myerror.Wrap(err, "endpointSearchContacts")
	}

	filters := request.ToFilters()
	if filters.Limit == 0 {
		filters.Limit = LimitMaxContacts
	}

	contacts, err := s.SearchContacts(ctx, filters)
	if err != nil {
		return searchContactsResponse{}, myerror.Wrap(err, "endpointSearchContacts")
	}

	contactsResponse := make([]getContactResponse, len(contacts))
	for i, c := range contacts {
		contactsResponse[i] = contactToGetContactResponse(c)
	}

	return searchContactsResponse{
		Contacts: contactsResponse,
	}, nil
}

// Delete

type deleteContactRequest struct {
	UserID    string
	ContactID string
}

func (r deleteContactRequest) Validate() error {
	var errorMessages []string

	if r.UserID == "" {
		errorMessages = append(errorMessages, "userID is required")
	}

	if r.ContactID == "" {
		errorMessages = append(errorMessages, "contactID is required")
	}

	if len(errorMessages) > 0 {
		return myerror.NewBadRequestError("invalid request: %s", strings.Join(errorMessages, ", "))
	}

	return nil
}

func endpointDeleteContact(ctx context.Context, s Service, request deleteContactRequest) error {
	if err := request.Validate(); err != nil {
		return myerror.Wrap(err, "endpointDeleteContact")
	}

	if err := s.DeleteContact(ctx, request.UserID, request.ContactID); err != nil {
		return myerror.Wrap(err, "endpointDeleteContact")
	}

	return nil
}
