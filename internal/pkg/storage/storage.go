package storage

// Database interface to do database operations
type Database interface {
	// Get retrieve all stored websites within database
	Get() ([]Website, error)
	// GetByID retrieve a website based on its ID
	GetByID(websiteID string) (Website, error)
	// Save store new website URL into database
	Save(web Website) error
	// Delete remove URL from database based on its ID
	Delete(websiteID string) error
}

// Website models that holds URL address of the website
type Website struct {
	ID      string
	URL     string
	Healthy bool
}
