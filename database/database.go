package database

type DataStore struct {
	IsConnected bool
	Collection  string
	Connection  interface{}
	Database    interface{}
}

type Database interface {
	New(databaseURL, database, collection string) DataStore
	Save(payload map[string]interface{}) (map[string]interface{}, error)
	Update(id string, payload map[string]interface{}) (map[string]interface{}, error)
	UpdateByField(key, value string, payload map[string]interface{}) error
	GetById(id string) map[string]interface{}
	GetByField(name string, value string) map[string]interface{}
	FindByFields(filters map[string]interface{}) []map[string]interface{}
	DeleteByField(name string, value string) error
	DeleteById(id string) error
}
