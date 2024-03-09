package inmem

import (
	"container/list"
	"context"
	"fmt"
	"sync"

	"contact-service/contact"
	"infrastructure/myerror"
)

type Repository interface {
	CreateContact(context.Context, contact.Contact) error
	GetContact(ctx context.Context, userID string, contactID string) (contact.Contact, error)
	DeleteContact(ctx context.Context, userID string, contactID string) error
	SearchContacts(ctx context.Context, filters contact.Filters) (contacts []contact.Contact, err error)
	UpdateContact(context.Context, contact.Contact) error
	IsPhoneExistsForUser(ctx context.Context, userID, phone string) (bool, error)
}

type Logger interface {
	Info(ctx context.Context, msg string, keyvals ...interface{})
	Error(ctx context.Context, err error, keyvals ...interface{})
	Warning(ctx context.Context, err error, keyvals ...interface{})
	Debug(ctx context.Context, msg string, keyvals ...interface{})
}

type CacheEntry struct {
	Key   string
	Value any
}

type lruCache struct {
	repo     Repository
	capacity int
	cache    map[string]*list.Element
	lruList  *list.List
	mutex    sync.RWMutex
	logger   Logger
}

func NewLRUCacheRepository(repo Repository, capacity int, logger Logger) *lruCache {
	return &lruCache{
		repo:     repo,
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		lruList:  list.New(),
		logger:   logger,
	}
}

func (l *lruCache) CreateContact(ctx context.Context, c contact.Contact) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if err := l.repo.CreateContact(ctx, c); err != nil {
		return myerror.Wrap(err, "lruCache.CreateContact")
	}

	entryKey := getCacheKey(c.UserID, c.ID)
	l.addCacheEntry(ctx, entryKey, c)

	return nil
}

func (l *lruCache) GetContact(ctx context.Context, userID string, contactID string) (contact.Contact, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	entryKey := getCacheKey(userID, contactID)

	// Check if the key exists in the cache, if so, move it to the front
	if elem, ok := l.cache[entryKey]; ok {
		l.lruList.MoveToFront(elem)
		c, ok := elem.Value.(*CacheEntry).Value.(contact.Contact)
		if !ok {
			return contact.Contact{}, myerror.NewInternalError("lruCache.GetContact: failed to cast to contact.Contact")
		}
		return c, nil
	}

	// If not found in the cache, fetch from the underlying repository
	c, err := l.repo.GetContact(ctx, userID, contactID)
	if err != nil {
		return contact.Contact{}, myerror.Wrap(err, "get")
	}

	l.addCacheEntry(ctx, entryKey, c)

	return c, nil
}

func (l *lruCache) addCacheEntry(ctx context.Context, key string, value any) {
	// If the key already exists, update its value and move it to the front
	if elem, ok := l.cache[key]; ok {
		elem.Value.(*CacheEntry).Value = value
		l.lruList.MoveToFront(elem)
		l.logger.Info(ctx, "lruCache.addCacheEntry: updated cache entry", "key", key, "value", value)
		return
	}

	// If the cache is full, remove the least recently used element
	if l.lruList.Len() == l.capacity {
		l.logger.Info(ctx, "lruCache.addCacheEntry: cache is full, removing least recently used element")
		elem := l.lruList.Back()
		delete(l.cache, elem.Value.(*CacheEntry).Key)
		l.lruList.Remove(elem)
	}

	l.logger.Info(ctx, "lruCache.addCacheEntry: adding new cache entry", "key", key, "value", value)

	// Add the new entry to the front of the list
	elem := l.lruList.PushFront(&CacheEntry{
		Key:   key,
		Value: value,
	})

	l.cache[key] = elem
}

func (l *lruCache) DeleteContact(ctx context.Context, userID string, contactID string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if err := l.repo.DeleteContact(ctx, userID, contactID); err != nil {
		return myerror.Wrap(err, "lruCache.DeleteContact")
	}

	contactKey := getCacheKey(userID, contactID)
	if elem, ok := l.cache[contactKey]; ok {
		delete(l.cache, contactKey)
		l.lruList.Remove(elem)
	}

	return nil
}

// SearchContacts is not cached
func (l *lruCache) SearchContacts(ctx context.Context, filters contact.Filters) ([]contact.Contact, error) {
	contacts, err := l.repo.SearchContacts(ctx, filters)
	if err != nil {
		return nil, myerror.Wrap(err, "lruCache.SearchContacts")
	}

	return contacts, nil
}

func (l *lruCache) UpdateContact(ctx context.Context, c contact.Contact) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if err := l.repo.UpdateContact(ctx, c); err != nil {
		return myerror.Wrap(err, "lruCache.CreateContact")
	}

	contactKey := getCacheKey(c.UserID, c.ID)
	l.addCacheEntry(ctx, contactKey, c)

	return nil
}

// IsPhoneExistsForUser is not cached
func (l *lruCache) IsPhoneExistsForUser(_ context.Context, userID, phone string) (bool, error) {
	IsPhoneExistsForUser, err := l.repo.IsPhoneExistsForUser(context.Background(), userID, phone)
	if err != nil {
		return false, myerror.Wrap(err, "lruCache.IsPhoneExistsForUser")
	}

	return IsPhoneExistsForUser, nil
}

func getCacheKey(userID, contactID string) string {
	return fmt.Sprintf("%s:%s", userID, contactID)
}
