package mongo

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *mongoStore) Aggregate(pipeline interface{}, output interface{}, allowDiskUse bool) error {
	opts := options.Aggregate()
	if allowDiskUse {
		opts.SetAllowDiskUse(true)
	}
	C, err := d.Collection.Aggregate(nil, pipeline, opts)
	if err != nil {
		return err
	}
	if err := C.All(nil, output); err != nil {
		return err
	}

	return nil
}
