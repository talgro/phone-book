package contactmanaging

import (
	"context"
	"fmt"
	"infrastructure/myerror"
	"time"

	"github.com/google/uuid"

	"user-service/contact"
)

type Repository interface {
	CreateContact(context.Context, contact.Contact) error
	GetContact(ctx context.Context, userID string, contactID string) (contact.Contact, error)
	UpdateContact(context.Context, contact.Contact) error
	IsPhoneExistsForUser(ctx context.Context, userID, phone string) (bool, error)
}

type LockCache interface {
	Lock(context.Context, string) (bool, error)
	Unlock(context.Context, string) error
}

type Logger interface {
	Info(ctx context.Context, msg string, keyvals ...interface{})
	Error(ctx context.Context, err error, keyvals ...interface{})
	Warning(ctx context.Context, err error, keyvals ...interface{})
	Debug(ctx context.Context, msg string, keyvals ...interface{})
}

type service struct {
	repo   Repository
	cache  LockCache
	logger Logger
}

func NewService(repo Repository, locker LockCache, logger Logger) *service {
	return &service{
		repo:   repo,
		cache:  locker,
		logger: logger,
	}
}

func (s service) CreateContact(ctx context.Context, c contact.Contact) (string, error) {
	lockKey := fmt.Sprintf("create:%s:%s", c.UserID, c.Phone)
	if err := s.lock(ctx, lockKey); err != nil {
		return "", myerror.Wrap(err, "service.CreateContact")
	}
	defer func() {
		if err := s.cache.Unlock(ctx, lockKey); err != nil {
			err = myerror.Wrap(err, "service.CreateContact")
			s.logger.Warning(ctx, err)
		}
	}()

	c.ID = uuid.New().String()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	if err := s.validateContactToCreate(ctx, c); err != nil {
		return "", myerror.Wrap(err, "service.CreateContact")
	}

	if err := s.repo.CreateContact(ctx, c); err != nil {
		return "", myerror.Wrap(err, "service.CreateContact")
	}

	return c.ID, nil
}

func (s service) UpdateContact(ctx context.Context, c contact.Contact) error {
	lockKey := fmt.Sprintf("update:%s:%s", c.UserID, c.Phone)
	if err := s.lock(ctx, lockKey); err != nil {
		return myerror.Wrap(err, "service.UpdateContact")
	}
	defer func() {
		if err := s.cache.Unlock(ctx, lockKey); err != nil {
			err = myerror.Wrap(err, "service.UpdateContact")
			s.logger.Warning(ctx, err)
		}
	}()

	contactPrevState, err := s.repo.GetContact(ctx, c.UserID, c.ID)
	if err != nil {
		return myerror.Wrap(err, "service.UpdateContact")
	}

	if !contactPrevState.UpdatedAt.Equal(c.UpdatedAt) {
		return myerror.NewBadRequestError("service.UpdateContact: contact has changed since the last read")
	}

	c = s.updateContactFields(contactPrevState, c)

	if err := s.repo.UpdateContact(ctx, c); err != nil {
		return myerror.Wrap(err, "service.UpdateContact")
	}

	return nil
}

func (s service) GetContact(ctx context.Context, userID, contactID string) (contact.Contact, error) {
	c, err := s.repo.GetContact(ctx, userID, contactID)
	if err != nil {
		return contact.Contact{}, myerror.Wrap(err, "service.GetContact")
	}

	return c, nil
}

func (s service) updateContactFields(contactPrevState contact.Contact, c contact.Contact) contact.Contact {
	contactPrevState.FirstName = c.FirstName
	contactPrevState.LastName = c.LastName
	contactPrevState.Address = c.Address
	contactPrevState.Phone = c.Phone
	contactPrevState.UpdatedAt = time.Now()

	return contactPrevState
}

func (s service) lock(ctx context.Context, key string) error {
	lockSuccess, err := s.cache.Lock(ctx, key)
	if err != nil {
		return myerror.Wrap(err, "lock")
	}
	if !lockSuccess {
		return myerror.NewInternalError("lock: failed to lock contact with key %s", key)
	}

	return nil
}

func (s service) validateContactToCreate(ctx context.Context, c contact.Contact) error {
	isPhoneExistsForUser, err := s.repo.IsPhoneExistsForUser(ctx, c.UserID, c.Phone)
	if err != nil {
		return myerror.NewInternalError("validateContactToCreate: %w", err)
	}
	if isPhoneExistsForUser {
		return myerror.NewBadRequestError("validateContactToCreate: contact with phone %s already exists for user %s", c.Phone, c.UserID)
	}

	return nil
}
