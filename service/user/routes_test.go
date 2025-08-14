package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ducklawrence/go-ecom/types"
	"github.com/gorilla/mux"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		payload := types.RegisterUserPayLoad{
			FirstName:       "duck",
			LastName:        "test",
			Email:           "invalid.email",
			Password:        "12345678",
			ConfirmPassword: "12345687089",
		}

		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register",
			bytes.NewReader(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d",
				http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should correctly register the user", func(t *testing.T) {
		payload := types.RegisterUserPayLoad{
			FirstName:       "duck",
			LastName:        "test",
			Email:           "valid@gmail.com",
			Password:        "12345678",
			ConfirmPassword: "12345678",
		}

		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register",
			bytes.NewReader(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d",
				http.StatusBadRequest, rr.Code)
		}
	})
}

type mockUserStore struct {
}

func (m *mockUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserStore) GetUserByID(ctx context.Context, id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(ctx context.Context, user types.User) error {
	return nil
}
