package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danielcesario/controlepeso/cmd/handler"
	"github.com/danielcesario/controlepeso/internal/controlepeso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (mock *MockService) CreateEntry(entry controlepeso.Entry) (*controlepeso.Entry, error) {
	arg := mock.Mock.Called(entry)
	result, _ := arg.Get(0).(controlepeso.Entry)
	return &result, arg.Error(1)
}

func (mock *MockService) ListEntries(start, count int) ([]controlepeso.Entry, error) {
	arg := mock.Mock.Called(start, count)
	result, _ := arg.Get(0).([]controlepeso.Entry)
	return result, arg.Error(1)
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

		service := new(MockService)
		service.On("CreateEntry", mock.Anything).Once().Return(*expectedCreatedUser, nil)

		// And: The handler received a valid entry
		handler := handler.NewHandler(service)
		req := httptest.NewRequest(http.MethodPost, "/entries", strings.NewReader(`{"user_id": 1, "weight": 105.6, "date": "2022-05-10 00:30:00"}`))
		res := httptest.NewRecorder()

		// When: The Create Entry Handler was called
		handler.HandleCreateEntry(res, req)

		// Then: Return the expected status code
		assert.Equal(t, http.StatusCreated, res.Code)

		// And: The response body was correct
		expected := `{"id":1,"user_id":1,"weight":105.6,"date":"2022-05-10 00:30:00"}`
		assert.Equal(t, expected, res.Body.String())
	})

	t.Run("Error on decode entry from JSON", func(t *testing.T) {
		// Given: a Mock service
		service := new(MockService)
		service.On("CreateEntry", mock.Anything).Once().Return(nil, nil)

		// And: The handler received an invalid entry
		handler := handler.NewHandler(service)
		req := httptest.NewRequest(http.MethodPost, "/entries", strings.NewReader(`{"user_id": "teste", "weight": 105.6, "date": "2022-05-10 00:30:00"}`))
		res := httptest.NewRecorder()

		// When: The Create Entry Handler was called
		handler.HandleCreateEntry(res, req)

		// Then: Return the expected status code
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Error on create entry by Service", func(t *testing.T) {
		// Given: The CreateEntry Service return an error
		service := new(MockService)
		service.On("CreateEntry", mock.Anything).Once().Return(nil, errors.New("Error on Create Entry"))

		// And: The handler received a valid entry
		handler := handler.NewHandler(service)
		req := httptest.NewRequest(http.MethodPost, "/entries", strings.NewReader(`{"user_id": 1, "weight": 105.6, "date": "2022-05-10 00:30:00"}`))
		res := httptest.NewRecorder()

		// When: The Create Entry Handler was called
		handler.HandleCreateEntry(res, req)

		// Then: Return the expected status code
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}
