package main

import (
	"user-service/contactmanaging"
	"user-service/inmem"
	"user-service/stdout"
)

func main() {
	inmemRepo := inmem.NewUserRepository()
	inmemLockCache := inmem.NewLockCache()
	inmemLRUCacheRepo := inmem.NewLRUCacheRepository(inmemRepo, 5)
	logger := stdout.NewLogger()

	service := contactmanaging.NewService(inmemLRUCacheRepo, inmemLockCache, logger)

	contactmanaging.ListenHTTP(service)
}
