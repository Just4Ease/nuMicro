package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/Just4Ease/nuMicro/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataStore database.DataStore

/**
 * New
 * This initialises a new MongoDB DataStore
 * param: string databaseURl
 * param: string databaseName
 * param: string collection
 * return: *DataStore
 */
func New(databaseURl string, databaseName string, collection string) *DataStore {
	clientOptions := options.Client().ApplyURI(databaseURl)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}

	db := client.Database(databaseName)

	c := db.Collection(collection)

	dataStore := DataStore{
		IsConnected: true,
		Database:    db,
		Collection:  collection,
		Connection:  c,
	}

	return &dataStore
}

/**
 * Save
 * Save is used to save a record in the DataStore
 */
func (d *DataStore) Save(payload interface{}, out interface{}) error {
	collection := d.Connection.(*mongo.Collection)
	_result_, err := collection.InsertOne(nil, payload)

	if err != nil {
		return err
	}

	err = collection.FindOne(nil, bson.M{"_id": _result_.InsertedID}).Decode(out)
	if err != nil {
		return err
	}
	return nil
}

/**
 * SaveMany
 * SaveMany is used to bulk insert into the DataStore
 *
 * param: []interface{} payload
 * return: error
 */
func (d *DataStore) SaveMany(payload []interface{}) error {
	collection := d.Connection.(*mongo.Collection)
	_, err := collection.InsertMany(nil, payload)

	if err != nil {
		return err
	}
	return nil
}

/**
 * FindById
 * find a single record by id in the DataStore
 * returns nil if record isn't found.
 *
 * param: interface{}            id
 * param: map[string]interface{} projection
 * return: map[string]interface{}
 */
func (d *DataStore) FindById(id interface{}, projection map[string]interface{}, result interface{}) error {
	collection := d.Connection.(*mongo.Collection)
	ops := options.FindOne()
	if projection != nil {
		ops.Projection = projection
	}
	if err := collection.FindOne(nil, bson.M{"_id": id}, ops).Decode(result); err != nil {
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
func (d *DataStore) FindOne(fields, projection map[string]interface{}, result interface{}) error {
	collection := d.Connection.(*mongo.Collection)
	ops := options.FindOne()
	ops.Projection = projection
	if err := collection.FindOne(nil, fields, ops).Decode(result); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return errors.New("document not found")
		}
		return err
	}
	return nil
}

func (d *DataStore) FindMany(fields, projection, sort map[string]interface{}, limit, skip int64, results interface{}) error {
	collection := d.Connection.(*mongo.Collection)
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
	cursor, err := collection.Find(nil, fields, ops)
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
 * Updates a single record by id in the DataStore
 *
 * param: interface{} id
 * param: interface{} payload
 * return: error
 */
func (d *DataStore) UpdateById(id interface{}, payload interface{}) error {
	collection := d.Connection.(*mongo.Collection)
	var u map[string]interface{}
	opts := options.FindOneAndUpdate()
	up := true
	opts.Upsert = &up
	if err := collection.FindOneAndUpdate(nil, bson.M{"_id": id}, bson.M{
		"$set": payload,
	}).Decode(&u); err != nil {
		return err
	}
	return nil
}

/**
 * UpdateOne
 *
 * Updates one item in the DataStore using fields as the criteria.
 *
 * param: map[string]interface{} fields
 * param: interface{}            payload
 * return: error
 */
func (d *DataStore) UpdateOne(fields map[string]interface{}, payload interface{}) error {
	collection := d.Connection.(*mongo.Collection)
	var u map[string]interface{}
	if err := collection.FindOneAndUpdate(nil, fields, bson.M{
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
func (d *DataStore) UpdateMany(fields map[string]interface{}, payload interface{}) error {
	// TODO: Update Many.
	collection := d.Connection.(*mongo.Collection)
	if _, err := collection.UpdateMany(nil, fields, bson.M{
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
func (d *DataStore) DeleteById(id interface{}) error {
	collection := d.Connection.(*mongo.Collection)
	var u map[string]interface{}
	if e := collection.FindOneAndDelete(nil, bson.M{
		"_id": id,
	}).Decode(&u); e != nil {
		return e
	}

	return nil
}

/**
 * DeleteOne
 * Deletes one item from the DataStore using fields a hash map to properly filter what is to be deleted.
 *
 * param: map[string]interface{} fields
 * return: error
 */
func (d *DataStore) DeleteOne(fields map[string]interface{}) error {
	collection := d.Connection.(*mongo.Collection)
	_, err := collection.DeleteOne(nil, fields)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Delete Many items from the DataStore
 *
 * param: map[string]interface{} fields
 * return: error
 */
func (d *DataStore) DeleteMany(fields map[string]interface{}) error {
	collection := d.Connection.(*mongo.Collection)
	_, err := collection.DeleteMany(nil, fields)
	if err != nil {
		return err
	}

	return nil
}
