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
	GetContact(context.Context, string) (contact.Contact, error)
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

	if err := s.validateContactToUpdate(ctx, c); err != nil {
		return myerror.Wrap(err, "service.UpdateContact")
	}

	c.UpdatedAt = time.Now()

	if err := s.repo.UpdateContact(ctx, c); err != nil {
		return myerror.Wrap(err, "service.UpdateContact")
	}

	return nil
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

func (s service) validateContactToUpdate(ctx context.Context, c contact.Contact) error {
	contactPrevState, err := s.repo.GetContact(ctx, c.ID)
	if err != nil {
		return myerror.Wrap(err, "validateContactToUpdate")
	}

	if contactPrevState.UserID != c.UserID {
		return myerror.NewForbiddenError("validateContactToUpdate")
	}

	if contactPrevState.UpdatedAt != c.UpdatedAt {
		return myerror.NewBadRequestError("validateContactToUpdate: contact has changed since the last read")
	}
	return nil
}