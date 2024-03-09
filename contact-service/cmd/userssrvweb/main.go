package main

import (
	"contact-service/contactmanaging"
	"contact-service/inmem"
	"contact-service/stdout"
)

func main() {
	logger := stdout.NewLogger()
	inmemLockCache := inmem.NewLockCache()
	inmemRepo := inmem.NewUserRepository()
	inmemLRUCacheRepo := inmem.NewLRUCacheRepository(inmemRepo, 5, logger)

	service := contactmanaging.NewService(inmemLRUCacheRepo, inmemLockCache, logger)

	contactmanaging.ServeHTTP(service)
}
