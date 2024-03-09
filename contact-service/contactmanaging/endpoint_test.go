package contactmanaging

import (
	"contact-service/inmem"
	"contact-service/stdout"
	"context"
	"testing"
)

func Test_endpointCreateContact(t *testing.T) {
	type args struct {
		request createContactRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid phone number",
			args: args{
				request: createContactRequest{
					UserID:    "123",
					Phone:     "123abc",
					FirstName: "John",
					LastName:  "Doe",
					Address:   "123 Main St",
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				request: createContactRequest{
					UserID:    "123",
					Phone:     "0546455401",
					FirstName: "John",
					LastName:  "Doe",
					Address:   "123 Main St",
				},
			},
			wantErr: false,
		},
		{
			name: "duplicated phone number for user",
			args: args{
				request: createContactRequest{
					UserID:    "123",
					Phone:     "0546455401",
					FirstName: "Not John",
					LastName:  "Doe",
					Address:   "123 Main St",
				},
			},
			wantErr: true,
		},
	}

	logger := stdout.NewLogger()
	inmemLockCache := inmem.NewLockCache()
	inmemRepo := inmem.NewUserRepository()
	inmemLRUCacheRepo := inmem.NewLRUCacheRepository(inmemRepo, 5, logger)

	s := service{
		repo:      inmemLRUCacheRepo,
		lockCache: inmemLockCache,
		logger:    logger,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := endpointCreateContact(context.Background(), s, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("endpointCreateContact() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
