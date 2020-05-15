package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Just4Ease/nuMicro/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoStore struct {
	IsConnected    bool
	CollectionName string
	Collection     *mongo.Collection
	Database       *mongo.Database
}

/**
 * New
 * This initialises a new MongoDB mongoStore
 * param: string databaseURl
 * param: string databaseName
 * param: string collection
 * return: *mongoStore
 */
func New(databaseURl string, databaseName string, collection string) (database.Database, error) {
	clientOptions := options.Client().ApplyURI(databaseURl)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	if err = client.Ping(context.TODO(), nil); err != nil {
		//log.Fatal(err)
		return nil, err
	}

	db := client.Database(databaseName)

	mongoStore := mongoStore{
		IsConnected:    true,
		CollectionName: collection,
		Collection:     db.Collection(collection),
		Database:       db,
	}

	return &mongoStore, nil
}

//func (d *mongoStore) pingDBClient()  {
//
//}
/**
 * Save
 * Save is used to save a record in the mongoStore
 */
func (d *mongoStore) Save(payload interface{}, out interface{}) error {
	_result_, err := d.Collection.InsertOne(context.Background(), payload)

	if err != nil {
		return err
	}

	err = d.Collection.FindOne(nil, bson.M{"_id": _result_.InsertedID}).Decode(out)
	if err != nil {
		return err
	}
	return nil
}

/**
 * SaveMany
 * SaveMany is used to bulk insert into the mongoStore
 *
 * param: []interface{} payload
 * return: error
 */
func (d *mongoStore) SaveMany(payload []interface{}) error {
	_, err := d.Collection.InsertMany(nil, payload)

	if err != nil {
		return err
	}
	return nil
}

/**
 * FindById
 * find a single record by id in the mongoStore
 * returns nil if record isn't found.
 *
 * param: interface{}            id
 * param: map[string]interface{} projection
 * return: map[string]interface{}
 */
func (d *mongoStore) FindById(id interface{}, projection map[string]interface{}, result interface{}) error {
	ops := options.FindOne()
	if projection != nil {
		ops.Projection = projection
	}
	if err := d.Collection.FindOne(nil, bson.M{"_id": id}, ops).Decode(result); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return errors.New("document not found")
		}
		return err
	}
	return nil
}

/**
 * Find One by
 */
func (d *mongoStore) FindOne(fields, projection map[string]interface{}, result interface{}) error {
	ops := options.FindOne()
	ops.Projection = projection
	if err := d.Collection.FindOne(nil, fields, ops).Decode(result); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return errors.New("document not found")
		}
		return err
	}
	return nil
}

func (d *mongoStore) FindMany(fields, projection, sort map[string]interface{}, limit, skip int64, results interface{}) error {
	ops := options.Find()
	if limit > 0 {
		ops.Limit = &limit
	}
	if skip > 0 {
		ops.Skip = &skip
	}
	if projection != nil {
		ops.Projection = projection
	}
	if sort != nil {
		ops.Sort = sort
	}
	cursor, err := d.Collection.Find(nil, fields, ops)
	if err != nil {
		return err
	}

	var output []map[string]interface{}
	for cursor.Next(nil) {
		var item map[string]interface{}
		_ = cursor.Decode(&item)
		output = append(output, item)
	}

	if b, e := json.Marshal(output); e == nil {
		_ = json.Unmarshal(b, &results)
	} else {
		return e
	}
	return nil
}

/**
 * UpdateById
 * Updates a single record by id in the mongoStore
 *
 * param: interface{} id
 * param: interface{} payload
 * return: error
 */
func (d *mongoStore) UpdateById(id interface{}, payload interface{}) error {
	var u map[string]interface{}
	opts := options.FindOneAndUpdate()
	up := true
	opts.Upsert = &up
	if err := d.Collection.FindOneAndUpdate(nil, bson.M{"_id": id}, bson.M{
		"$set": payload,
	}).Decode(&u); err != nil {
		return err
	}
	return nil
}

/**
 * UpdateOne
 *
 * Updates one item in the mongoStore using fields as the criteria.
 *
 * param: map[string]interface{} fields
 * param: interface{}            payload
 * return: error
 */
func (d *mongoStore) UpdateOne(fields map[string]interface{}, payload interface{}) error {
	var u map[string]interface{}
	if err := d.Collection.FindOneAndUpdate(nil, fields, bson.M{
		"$set": payload,
	}).Decode(&u); err != nil {
		return err
	}
	return nil
}

/**
 * UpdateMany
 * Updates many items in the collection
 * `fields` this is the search criteria
 * `payload` this is the update payload.
 *
 * param: map[string]interface{} fields
 * param: interface{}            payload
 * return: error
 */
func (d *mongoStore) UpdateMany(fields, payload map[string]interface{}) error {
	if _, err := d.Collection.UpdateMany(nil, fields, bson.M{
		"$set": payload,
	}); err != nil {
		return err
	}
	return nil
}

/**
 * DeleteById
 * Deletes a single record by id
 * where ID can be a string or whatever.
 * param: interface{} id
 * return: error
 */
func (d *mongoStore) DeleteById(id interface{}) error {
	var u map[string]interface{}
	if e := d.Collection.FindOneAndDelete(nil, bson.M{
		"_id": id,
	}).Decode(&u); e != nil {
		return e
	}

	return nil
}

/**
 * DeleteOne
 * Deletes one item from the mongoStore using fields a hash map to properly filter what is to be deleted.
 *
 * param: map[string]interface{} fields
 * return: error
 */
func (d *mongoStore) DeleteOne(fields map[string]interface{}) error {
	_, err := d.Collection.DeleteOne(nil, fields)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Delete Many items from the mongoStore
 *
 * param: map[string]interface{} fields
 * return: error
 */
func (d *mongoStore) DeleteMany(fields map[string]interface{}) error {
	_, err := d.Collection.DeleteMany(nil, fields)
	if err != nil {
		return err
	}

	return nil
}

func (d *mongoStore) Count(fields interface{}) (int64, error) {
	return d.Collection.CountDocuments(nil, fields)
}
