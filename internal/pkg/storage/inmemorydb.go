package storage

import "errors"

var (
	// ErrNotFound an error (string) indicates that the records is not found
	ErrNotFound = errors.New("not found")
)

// InMemoryDatabase storage within memory
type InMemoryDatabase struct {
	webs map[string]Website
}

// NewInMemoryDatabase creates instances for database sotored within memeory
func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		webs: make(map[string]Website, 0),
	}
}

// Get retrieve all websites within database
func (database InMemoryDatabase) Get() ([]Website, error) {
	var websites []Website
	for _, web := range database.webs {
		websites = append(websites, web)
	}
	return websites, nil
}

// GetByID retrieve a website based on its ID
func (database InMemoryDatabase) GetByID(websiteID string) (Website, error) {
	for _, web := range database.webs {
		if web.ID == websiteID {
			return web, nil
		}
	}
	return Website{}, ErrNotFound
}

// Save store URL to in-memory database
func (database *InMemoryDatabase) Save(web Website) error {
	database.webs[web.ID] = web
	return nil
}

// Delete remove website URL from database
func (database *InMemoryDatabase) Delete(websiteID string) error {
	delete(database.webs, websiteID)
	return nil
}
