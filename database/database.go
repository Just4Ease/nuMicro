package database

type Database interface {
	Save(payload interface{}, out interface{}) error
	SaveMany(payload []interface{}) error
	FindById(id interface{}, projection map[string]interface{}, result interface{}) error
	FindOne(fields, projection map[string]interface{}, result interface{}) error
	FindMany(fields, projection, sort map[string]interface{}, limit, skip int64, results interface{}) error
	UpdateById(id interface{}, payload interface{}) error
	UpdateOne(fields map[string]interface{}, payload interface{}) error
	UpdateMany(fields, payload map[string]interface{}) error
	DeleteById(id interface{}) error
	DeleteOne(fields map[string]interface{}) error
	DeleteMany(fields map[string]interface{}) error
	Aggregate(pipelines interface{}, result interface{}, allowDiskUse bool) error
	Count(fields interface{}) (int64, error)
}
