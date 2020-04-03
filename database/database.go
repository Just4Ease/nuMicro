package database

type DataStore struct {
	IsConnected bool
	Collection  string
	Connection  interface{}
	Database    interface{}
}

type Database interface {
	New(databaseURL, database, collection string) DataStore
	Save(payload interface{}, out interface{}) error
	FindById(id interface{}, projection map[string]interface{}, result interface{}) error
	FindOne(fields, projection map[string]interface{}, result interface{}) error
	FindMany(fields, projection, sort map[string]interface{}, limit, skip int64, results interface{}) error
	UpdateById(id interface{}, payload interface{}) error
	UpdateOne(fields map[string]interface{}, payload interface{})
	UpdateMany(fields, payload map[string]interface{}) error
	DeleteById(id string) error
	DeleteOne(fields map[string]interface{}) error
	DeleteMany(fields map[string]interface{}) error
}
