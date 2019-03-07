package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ajiyakin/gohealthz/internal/pkg/storage"
)

func TestCreateWebsiteHealthy(t *testing.T) {
	// arrange
	httpGetRequestFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
		}, nil
	}
	responseRecorder := httptest.NewRecorder()
	requestBody := createWebsiteRequest{URL: "https://www.example.com"}
	requestBodyRaw, err := json.Marshal(requestBody)
	if err != nil {
		t.Errorf("unable to marshal request body: %v", err)
	}
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/website", bytes.NewReader(requestBodyRaw))
	if err != nil {
		t.Errorf("unable to create new HTTP request instance: %v", err)
	}
	database := storage.NewInMemoryDatabase()

	// action
	handlerFunc := NewWebsiteHandler(database)
	handlerFunc(responseRecorder, request)

	// acceptance
	response := responseRecorder.Result()
	if response.StatusCode != http.StatusCreated {
		t.Errorf("expected response code %d, got %d", http.StatusCreated, response.StatusCode)
	}
	actualRecord, err := database.Get()
	if err != nil {
		t.Errorf("unable to retrieve website records: %v", err)
	}
	if len(actualRecord) < 1 {
		t.Errorf("expected records to not-empty, got %d records", len(actualRecord))
	}
	for _, website := range actualRecord {
		if website.URL != "https://www.example.com" {
			t.Errorf("expected URL to be https://www.example.com, got %s", website.URL)
		}
		if website.Healthy != true {
			t.Errorf("expected healthy to be true, got %v", website.Healthy)
		}
	}
}

func TestCreateWebsiteUnHealthy(t *testing.T) {
	// arrange
	httpGetRequestFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	responseRecorder := httptest.NewRecorder()
	requestBody := createWebsiteRequest{URL: "https://www.example.com"}
	requestBodyRaw, err := json.Marshal(requestBody)
	if err != nil {
		t.Errorf("unable to marshal request body: %v", err)
	}
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/website", bytes.NewReader(requestBodyRaw))
	if err != nil {
		t.Errorf("unable to create new HTTP request instance: %v", err)
	}
	database := storage.NewInMemoryDatabase()

	// action
	handlerFunc := NewWebsiteHandler(database)
	handlerFunc(responseRecorder, request)

	// acceptance
	response := responseRecorder.Result()
	if response.StatusCode != http.StatusCreated {
		t.Errorf("expected response code %d, got %d", http.StatusCreated, response.StatusCode)
	}
	actualRecord, err := database.Get()
	if err != nil {
		t.Errorf("unable to retrieve website records: %v", err)
	}
	if len(actualRecord) < 1 {
		t.Errorf("expected records to not-empty, got %d records", len(actualRecord))
	}
	for _, website := range actualRecord {
		if website.URL != "https://www.example.com" {
			t.Errorf("expected URL to be https://www.example.com, got %s", website.URL)
		}
		if website.Healthy != false {
			t.Errorf("expected healthy to be true, got %v", website.Healthy)
		}
	}
}

func TestCreateWebsiteWithInvalidURL(t *testing.T) {
	// arrange
	webURLTests := []struct {
		testName string
		URL      string
	}{
		{
			"URL with no scheme",
			"example.com",
		}, {
			"URL with no scheme and host",
			"example",
		}, {
			"URL with only domain",
			".com",
		},
	}
	database := storage.NewInMemoryDatabase()
	handlerFunc := NewWebsiteHandler(database)

	for _, tt := range webURLTests {
		requestBody := createWebsiteRequest{URL: tt.URL}
		requestBodyRaw, err := json.Marshal(requestBody)
		if err != nil {
			t.Errorf("unable to marshal request body: %v", err)
		}
		request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/website", bytes.NewReader(requestBodyRaw))
		if err != nil {
			t.Errorf("unable to create new HTTP request instance: %v", err)
		}

		// action
		responseRecorder := httptest.NewRecorder()
		t.Run(tt.testName, func(t *testing.T) {
			handlerFunc(responseRecorder, request)
		})

		// acceptance
		response := responseRecorder.Result()
		if response.StatusCode != http.StatusBadRequest {
			t.Errorf("expected response code %d, got %d", http.StatusBadRequest, response.StatusCode)
		}
	}
}

func TestCreateWebsiteWithInvalidRequestBody(t *testing.T) {
	// arrange
	responseRecorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/website", bytes.NewReader([]byte(`...`)))
	if err != nil {
		t.Errorf("unable to create new HTTP request instance: %v", err)
	}
	database := storage.NewInMemoryDatabase()

	// action
	handlerFunc := NewWebsiteHandler(database)
	handlerFunc(responseRecorder, request)

	// acceptance
	response := responseRecorder.Result()
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("expected response code %d, got %d", http.StatusBadRequest, response.StatusCode)
	}
}

func TestGetListOfWebsitesWithEmptyRecords(t *testing.T) {
	// arrange
	request, err := http.NewRequest(http.MethodGet, "http://localhost:8080/website", nil)
	if err != nil {
		t.Errorf("unable to create new HTTP request instance: %v", err)
	}
	database := storage.NewInMemoryDatabase()
	responseRecorder := httptest.NewRecorder()

	// action
	handlerFunc := NewWebsiteHandler(database)
	handlerFunc(responseRecorder, request)

	// acceptance
	response := responseRecorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected response code %d, got %d", http.StatusOK, response.StatusCode)
	}
	var responseBody []getWebsitesResponse
	if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		t.Errorf("unable to decode response body: %v", err)
	}
	if len(responseBody) > 0 {
		t.Errorf("expected empty response: got %#v", responseBody)
	}
}

func TestGetListOfWebsitesWithNonEmptyRecords(t *testing.T) {
	// arrange
	request, err := http.NewRequest(http.MethodGet, "http://localhost:8080/website", nil)
	if err != nil {
		t.Errorf("unable to create new HTTP request instance: %v", err)
	}
	database := storage.NewInMemoryDatabase()
	err = database.Save(storage.Website{
		ID:      "1234",
		URL:     "https://example.com",
		Healthy: true,
	})
	if err != nil {
		t.Errorf("unable to save to database: %v", err)
	}
	responseRecorder := httptest.NewRecorder()

	// action
	handlerFunc := NewWebsiteHandler(database)
	handlerFunc(responseRecorder, request)

	// acceptance
	response := responseRecorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected response code %d, got %d", http.StatusOK, response.StatusCode)
	}
	var responseBody []getWebsitesResponse
	if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		t.Errorf("unable to decode response body: %v", err)
	}
	if len(responseBody) < 1 {
		t.Errorf("expected non-empty response: got %#v", responseBody)
	}
}

func TestDeleteWebsiteWithRecordExists(t *testing.T) {
	// arrange
	formValues := url.Values{
		"website_id": []string{"1234"},
	}
	request, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/website?"+formValues.Encode(), nil)
	if err != nil {
		t.Errorf("unable to create new HTTP request instance: %v", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	database := storage.NewInMemoryDatabase()
	err = database.Save(storage.Website{
		ID:      "1234",
		URL:     "https://example.com",
		Healthy: true,
	})
	if err != nil {
		t.Errorf("unable to save to database: %v", err)
	}
	responseRecorder := httptest.NewRecorder()

	// action
	handlerFunc := NewWebsiteHandler(database)
	handlerFunc(responseRecorder, request)

	// acceptance
	response := responseRecorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected response code %d, got %d", http.StatusOK, response.StatusCode)
	}
	deletedRecord, err := database.GetByID("123")
	if err != storage.ErrNotFound {
		t.Errorf("expected record not found after delete: %v", err)
	}
	expectedDeletedRecord := storage.Website{}
	if deletedRecord != expectedDeletedRecord {
		t.Errorf("expected no record after delete. got: %v", deletedRecord)
	}
}

func TestDeleteWebsiteWithRecordNotExists(t *testing.T) {
	// arrange
	formValues := url.Values{
		"website_id": []string{"1234"},
	}
	request, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/website?"+formValues.Encode(), nil)
	if err != nil {
		t.Errorf("unable to create new HTTP request instance: %v", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	database := storage.NewInMemoryDatabase()
	responseRecorder := httptest.NewRecorder()

	// action
	handlerFunc := NewWebsiteHandler(database)
	handlerFunc(responseRecorder, request)

	// acceptance
	response := responseRecorder.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected response code %d, got %d", http.StatusOK, response.StatusCode)
	}
	deletedRecord, err := database.GetByID("123")
	if err != storage.ErrNotFound {
		t.Errorf("expected record not found after delete: %v", err)
	}
	expectedDeletedRecord := storage.Website{}
	if deletedRecord != expectedDeletedRecord {
		t.Errorf("expected no record after delete. got: %v", deletedRecord)
	}
}
