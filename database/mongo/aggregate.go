package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *DataStore) Aggregate(pipeline interface{}, output interface{}, allowDiskUse bool) error {
	collection := d.Connection.(*mongo.Collection)
	opts := options.Aggregate()
	if allowDiskUse {
		opts.SetAllowDiskUse(true)
	}
	C, err := collection.Aggregate(nil, pipeline, opts)
	if err != nil {
		return err
	}
	if err := C.All(nil, output); err != nil {
		return err
	}

	return nil
}
