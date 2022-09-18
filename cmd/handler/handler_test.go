package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danielcesario/entry/cmd/handler"
	"github.com/danielcesario/entry/internal/entry"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (mock *MockService) CreateEntry(entryParam entry.Entry) (*entry.Entry, error) {
	arg := mock.Mock.Called(entryParam)
	result, _ := arg.Get(0).(entry.Entry)
	return &result, arg.Error(1)
}

func (mock *MockService) ListEntries(start, count int) ([]entry.Entry, error) {
	arg := mock.Mock.Called(start, count)
	result, _ := arg.Get(0).([]entry.Entry)
	return result, arg.Error(1)
}

func (mock *MockService) GetEntry(id int) (*entry.Entry, error) {
	arg := mock.Mock.Called(id)
	result, _ := arg.Get(0).(entry.Entry)
	return &result, arg.Error(1)
}

func (mock *MockService) DeleteEntry(id int) error {
	arg := mock.Mock.Called(id)
	return arg.Error(0)
}

func (mock *MockService) UpdateEntry(id int, entryParam entry.Entry) (*entry.Entry, error) {
	arg := mock.Mock.Called(id, entryParam)
	result, _ := arg.Get(0).(entry.Entry)
	return &result, arg.Error(1)
}

func TestCreateEntry(t *testing.T) {
	t.Run("Create Entry with Success", func(t *testing.T) {
		// Given: The service create and return an expected entry
		expectedCreatedUser := &entry.Entry{
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

func TestListEntries(t *testing.T) {
	service := new(MockService)

	t.Run("List Entries with Success", func(t *testing.T) {
		// Given: The service return an expected list of entries
		var entries []entry.Entry
		entries = append(entries, entry.Entry{
			ID:     1,
			UserId: 1,
			Weight: 105.6,
			Date:   "2022-05-10 00:30:00",
		})
		service.On("ListEntries", 0, 10).Once().Return(entries, nil)

		// And: The handler received a valid entry
		handler := handler.NewHandler(service)
		req := httptest.NewRequest(http.MethodGet, "/entries?count=10&start=0", nil)
		res := httptest.NewRecorder()

		// When: The List Entry Handler was called
		handler.HandleListEntries(res, req)

		// Then: Return the expected status code
		assert.Equal(t, http.StatusOK, res.Code)

		// And: The response body was correct
		expected := `[{"id":1,"user_id":1,"weight":105.6,"date":"2022-05-10 00:30:00"}]`
		assert.Equal(t, expected, res.Body.String())
	})

	t.Run("List Entries with error on call service", func(t *testing.T) {
		// Given: The service return an error
		service.On("ListEntries", 0, 10).Once().Return(nil, errors.New("Error on List Entries"))

		// And: The handler received a valid entry
		handler := handler.NewHandler(service)
		req := httptest.NewRequest(http.MethodGet, "/entries?count=10&start=0", nil)
		res := httptest.NewRecorder()

		// When: The List Entry Handler was called
		handler.HandleListEntries(res, req)

		// Then: Return the expected status code
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func TestGetEntry(t *testing.T) {
	t.Run("Get an Entry with Success", func(t *testing.T) {
		// Given: The service return an expected entry
		var entry = &entry.Entry{
			ID:     1,
			UserId: 1,
			Weight: 105.6,
			Date:   "2022-05-10 00:30:00",
		}
		service := new(MockService)
		service.On("GetEntry", mock.Anything).Once().Return(*entry, nil)

		// And: The request and response were valid
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/entries/{id}", nil)
		params := map[string]string{"id": "1"}
		req = mux.SetURLVars(req, params)
		handler := handler.NewHandler(service)

		// When: The Get Entry Handler was called
		handler.HandleGetEntry(res, req)

		// Then: Return the expected status code
		assert.Equal(t, http.StatusOK, res.Code)

		// And: The response body was correct
		expected := `{"id":1,"user_id":1,"weight":105.6,"date":"2022-05-10 00:30:00"}`
		assert.Equal(t, expected, res.Body.String())
	})

	t.Run("Not Found Error on Get an Entry", func(t *testing.T) {
		// Given: The service return a nil entry
		service := new(MockService)
		service.On("GetEntry", mock.Anything).Return(&entry.Entry{}, nil)

		// And: The request and response were valid
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/entries/{id}", nil)
		params := map[string]string{"id": "1"}
		req = mux.SetURLVars(req, params)
		handler := handler.NewHandler(service)

		// When: The Get Entry Handler was called
		handler.HandleGetEntry(res, req)

		// Then: Return the expected status code
		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("Internal Server Error on Get an Entry", func(t *testing.T) {
		// Given: The service return an error
		service := new(MockService)
		service.On("GetEntry", mock.Anything).Return(&entry.Entry{}, errors.New("Internal Error"))

		// And: The request and response were valid
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/entries/{id}", nil)
		params := map[string]string{"id": "1"}
		req = mux.SetURLVars(req, params)
		handler := handler.NewHandler(service)

		// When: The Get Entry Handler was called
		handler.HandleGetEntry(res, req)

		// Then: Return the expected status code
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func TestDeleteEntry(t *testing.T) {
	t.Run("Delete an Entry with Success", func(t *testing.T) {
		// Given: The service not return error
		service := new(MockService)
		service.On("DeleteEntry", mock.Anything).Once().Return(nil)

		// And: The request and response were valid
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/entries/{id}", nil)
		params := map[string]string{"id": "1"}
		req = mux.SetURLVars(req, params)
		handler := handler.NewHandler(service)

		// When: The Get Entry Handler was called
		handler.HandleDeleteEntry(res, req)

		// Then: Return the expected status code
		assert.Equal(t, http.StatusNoContent, res.Code)
	})

	t.Run("Internal Server Error on Delete an Entry", func(t *testing.T) {
		// Given: The service return an error
		service := new(MockService)
		service.On("DeleteEntry", mock.Anything).Return(errors.New("Internal Error"))

		// And: The request and response were valid
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/entries/{id}", nil)
		params := map[string]string{"id": "1"}
		req = mux.SetURLVars(req, params)
		handler := handler.NewHandler(service)

		// When: The Get Entry Handler was called
		handler.HandleDeleteEntry(res, req)

		// Then: Return the expected status code
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func TestUpdateEntry(t *testing.T) {
	t.Run("Update an Entry with Success", func(t *testing.T) {
		// Given: An expected updated entry
		var expectedUpdatedEntry = &entry.Entry{
			ID:     1,
			UserId: 1,
			Weight: 100.0,
			Date:   "2022-01-10 00:30:00",
		}

		// And: An expected Entry to update
		var expectedEntryToUpdate = entry.Entry{Weight: 100, Date: "2022-01-10 00:30:00"}

		// And: The service receive and return the expected values
		service := new(MockService)
		service.On("UpdateEntry", 1, expectedEntryToUpdate).Once().Return(*expectedUpdatedEntry, nil)

		// And: The request and response were valid
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/entries/{id}", strings.NewReader(`{"weight": 100.0, "date": "2022-01-10 00:30:00"}`))
		params := map[string]string{"id": "1"}
		req = mux.SetURLVars(req, params)
		handler := handler.NewHandler(service)

		// When: The Get Entry Handler was called
		handler.HandleUpdateEntry(res, req)

		// Then: Return the expected status code
		assert.Equal(t, http.StatusOK, res.Code)

		// And: The response body was correct
		expected := `{"id":1,"user_id":1,"weight":100,"date":"2022-01-10 00:30:00"}`
		assert.Equal(t, expected, res.Body.String())
	})

	t.Run("Error on Update with Invalid Json", func(t *testing.T) {
		// Given: The service wasn't called
		service := new(MockService)
		service.On("UpdateEntry", mock.Anything, mock.Anything).Times(0)

		// And: The request had an invalid value
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/entries/{id}", strings.NewReader(`{"weight": "one hundred", "date": "2022-01-10 00:30:00"}`))
		params := map[string]string{"id": "1"}
		req = mux.SetURLVars(req, params)
		handler := handler.NewHandler(service)

		// When: The Get Entry Handler was called
		handler.HandleUpdateEntry(res, req)

		// Then: Return the expected status code
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Error on Update with Internal Server Error", func(t *testing.T) {
		// Given: An expected Entry to update
		var expectedEntryToUpdate = entry.Entry{Weight: 100, Date: "2022-01-10 00:30:00"}

		// And: The service receive the expected values and return an error
		service := new(MockService)
		service.On("UpdateEntry", 1, expectedEntryToUpdate).Once().Return(nil, errors.New("Error on Update"))

		// And: The request and response were valid
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/entries/{id}", strings.NewReader(`{"weight": 100.0, "date": "2022-01-10 00:30:00"}`))
		params := map[string]string{"id": "1"}
		req = mux.SetURLVars(req, params)
		handler := handler.NewHandler(service)

		// When: The Get Entry Handler was called
		handler.HandleUpdateEntry(res, req)

		// Then: Return the expected status code
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("Error on Update with Not Found Entry", func(t *testing.T) {
		// Given: An expected Entry to update
		var expectedEntryToUpdate = entry.Entry{Weight: 100, Date: "2022-01-10 00:30:00"}

		// And: The service receive the expected values and return a nil entry
		service := new(MockService)
		service.On("UpdateEntry", 1, expectedEntryToUpdate).Once().Return(nil, nil)

		// And: The request and response were valid
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/entries/{id}", strings.NewReader(`{"weight": 100.0, "date": "2022-01-10 00:30:00"}`))
		params := map[string]string{"id": "1"}
		req = mux.SetURLVars(req, params)
		handler := handler.NewHandler(service)

		// When: The Get Entry Handler was called
		handler.HandleUpdateEntry(res, req)

		// Then: Return the expected status code
		assert.Equal(t, http.StatusNotFound, res.Code)
	})
}
