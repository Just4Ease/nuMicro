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
	FindById(id string) map[string]interface{}
	FindOne(fields, projection map[string]interface{}) map[string]interface{}
	FindMany(fields, projection, sort map[string]interface{}, limit, skip int64) []map[string]interface{}
	UpdateById(id string, payload map[string]interface{}) (map[string]interface{}, error)
	UpdateOne(fields, payload map[string]interface{}) (map[string]interface{}, error)
	UpdateMany(fields, payload map[string]interface{}) error
	DeleteById(id string) error
	DeleteOne(fields map[string]interface{}) error
	DeleteMany(fields map[string]interface{}) error
}
