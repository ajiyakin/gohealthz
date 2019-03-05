package storage

import (
	"testing"
)

func TestGetWebsiteByIDSuccess(t *testing.T) {
	// arrange
	db := NewInMemoryDatabase()
	err := db.Save(Website{
		ID:      "123",
		URL:     "http://example.com",
		Healthy: true})
	if err != nil {
		t.Errorf("unable to save website to database: %v", err)
	}

	// action
	actual, err := db.GetByID("123")
	if err != nil {
		t.Errorf("unable to get website by ID: %v", err)
	}

	// acceptance
	expected := Website{
		ID:      "123",
		URL:     "http://example.com",
		Healthy: true,
	}
	if actual != expected {
		t.Errorf("expected %#v got %#v", expected, actual)
	}
}

func TestGetWebsiteAllSuccess(t *testing.T) {
	// arrange
	db := NewInMemoryDatabase()
	err := db.Save(Website{
		ID:      "123",
		URL:     "http://one.example.com",
		Healthy: true,
	})
	if err != nil {
		t.Errorf("unable to save website to database: %v", err)
	}
	err = db.Save(Website{
		ID:      "456",
		URL:     "http://two.example.com",
		Healthy: true,
	})
	if err != nil {
		t.Errorf("unable to save website to database: %v", err)
	}

	// action
	actuals, err := db.Get()
	if err != nil {
		t.Errorf("unable to get all websites record from database: %v", err)
	}

	// acceptance
	expecteds := []Website{
		{
			ID:      "123",
			URL:     "http://one.example.com",
			Healthy: true,
		}, {
			ID:      "456",
			URL:     "http://two.example.com",
			Healthy: true,
		},
	}
	for index, expected := range expecteds {
		if expected != actuals[index] {
			t.Errorf("expected %#v got %#v", expected, actuals)
		}
	}
}

func TestSaveWebsiteToDatabaseSuccess(t *testing.T) {
	// arrange
	db := NewInMemoryDatabase()

	// action
	err := db.Save(Website{
		ID:      "123",
		URL:     "https://www.example.com",
		Healthy: true,
	})
	if err != nil {
		t.Errorf("unable to store to database successfully: %v", err)
	}

	// acceptance
	actual, err := db.GetByID("123")
	if err != nil {
		t.Errorf("unable to get website by ID: %v", err)
	}
	expected := Website{
		ID:      "123",
		URL:     "https://www.example.com",
		Healthy: true,
	}
	if actual != expected {
		t.Errorf("expected %#v got %#v", expected, actual)
	}
}

func TestDeleteWebsiteFromDatabaseSuccess(t *testing.T) {
	// arrange
	db := NewInMemoryDatabase()
	err := db.Save(Website{
		ID:      "123",
		URL:     "http://example.com",
		Healthy: true,
	})
	if err != nil {
		t.Errorf("unable to save website: %v", err)
	}

	// action
	if err = db.Delete("123"); err != nil {
		t.Errorf("unabel to delete website: %v", err)
	}

	// acceptance
	_, err = db.GetByID("123")
	if err != ErrNotFound {
		t.Errorf("expected error not found, got %v", err)
	}
}
