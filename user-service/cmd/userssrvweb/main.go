package main

import (
	"user-service/contactmanaging"
	"user-service/inmem"
	"user-service/stdout"
)

func main() {
	inmemRepo := inmem.NewUserRepository()
	inmemLockCache := inmem.NewLockCache()
	logger := stdout.NewLogger()

	service := contactmanaging.NewService(inmemRepo, inmemLockCache, logger)

	contactmanaging.ListenHTTP(service)
}
