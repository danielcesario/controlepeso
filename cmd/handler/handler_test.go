package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danielcesario/controlepeso/cmd/handler"
	"github.com/danielcesario/controlepeso/internal/controlepeso"
)

type MockService struct {
	CreateEntryFunc func(entry controlepeso.Entry) (*controlepeso.Entry, error)
	CreatedEntries  []controlepeso.Entry
}

func (mock *MockService) CreateEntry(entry controlepeso.Entry) (*controlepeso.Entry, error) {
	mock.CreatedEntries = append(mock.CreatedEntries, entry)
	return mock.CreateEntryFunc(entry)
}

func (mock *MockService) ListEntries(start, count int) ([]controlepeso.Entry, error) {
	return mock.CreatedEntries, nil
}

func TestCreateEntry(t *testing.T) {
	t.Run("Create Entry with Success", func(t *testing.T) {
		// Given: The service create and return an expected entry
		expectedCreatedUser := &controlepeso.Entry{
			ID:     1,
			UserId: 1,
			Weight: 105.6,
			Date:   "2022-05-10 00:30:00",
		}

		service := &MockService{
			CreateEntryFunc: func(entry controlepeso.Entry) (*controlepeso.Entry, error) {
				return expectedCreatedUser, nil
			},
		}

		// And: The handler received a valid entry
		handler := handler.NewHandler(service)
		req := httptest.NewRequest(http.MethodPost, "/entries", strings.NewReader(`{"user_id": 1, "weight": 105.6, "date": "2022-05-10 00:30:00"}`))
		res := httptest.NewRecorder()

		// When: The Create Entry Handler was called
		handler.HandleCreateEntry(res, req)

		// Then: Return the expected status code
		if status := res.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		// And: The response body was correct
		expected := `{"id":1,"user_id":1,"weight":105.6,"date":"2022-05-10 00:30:00"}`
		if res.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				res.Body.String(), expected)
		}
	})

	t.Run("Error on create entry with invalid JSON", func(t *testing.T) {
		// Given: a Mock service
		service := &MockService{
			CreateEntryFunc: func(entry controlepeso.Entry) (*controlepeso.Entry, error) {
				return nil, nil
			},
		}

		// And: The handler received an invalid entry
		handler := handler.NewHandler(service)
		req := httptest.NewRequest(http.MethodPost, "/entries", strings.NewReader(`{"user_id": "teste", "weight": 105.6, "date": "2022-05-10 00:30:00"}`))
		res := httptest.NewRecorder()

		// When: The Create Entry Handler was called
		handler.HandleCreateEntry(res, req)

		// Then: Return the expected status code
		if status := res.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("Error on create entry with invalid JSON", func(t *testing.T) {
		// Given: The CreateEntry Service return an error
		service := &MockService{
			CreateEntryFunc: func(entry controlepeso.Entry) (*controlepeso.Entry, error) {
				return nil, errors.New("Error on Create Entry")
			},
		}

		// And: The handler received a valid entry
		handler := handler.NewHandler(service)
		req := httptest.NewRequest(http.MethodPost, "/entries", strings.NewReader(`{"user_id": 1, "weight": 105.6, "date": "2022-05-10 00:30:00"}`))
		res := httptest.NewRecorder()

		// When: The Create Entry Handler was called
		handler.HandleCreateEntry(res, req)

		// Then: Return the expected status code
		if status := res.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
}
