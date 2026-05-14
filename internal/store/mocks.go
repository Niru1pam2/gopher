package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
}

func (m *MockUserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	// Simulate a successful creation
	return nil
}

func (m *MockUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	// Return a dummy user or nil.
	// For testing, you might want to return a specific user based on the email passed.
	return &User{}, nil
}

func (m *MockUserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	// Return a dummy user or nil
	return &User{}, nil
}

func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error {
	// Simulate successful creation and invitation
	return nil
}

func (m *MockUserStore) Activate(ctx context.Context, token string) error {
	// Simulate successful activation
	return nil
}

func (m *MockUserStore) Delete(ctx context.Context, id int64) error {
	// Simulate successful deletion
	return nil
}
