package mongo

import (
	"context"
	"fmt"
	"log"

	"github.com/Just4Ease/nuMicro/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataStore database.DataStore

/**
Save
*/
func (d DataStore) Save(payload interface{}) (map[string]interface{}, error) {
	collection := d.Connection.(*mongo.Collection)
	_result_, err := collection.InsertOne(nil, payload)

	if err != nil {
		return nil, err
	}

	var response map[string]interface{}

	err = collection.FindOne(nil, bson.M{"_id": _result_.InsertedID}).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

/**
Update One
*/
func (d DataStore) Update(id string, payload interface{}) error {
	collection := d.Connection.(*mongo.Collection)
	Id, _ := primitive.ObjectIDFromHex(id)
	var u interface{}
	if err := collection.FindOneAndUpdate(nil, Id, bson.M{
		"$set": payload,
	}).Decode(u); err != nil {
		return err
	}
	return nil
}

func (d *DataStore) UpdateByField(key, value string, payload interface{}) error {
	collection := d.Connection.(*mongo.Collection)
	err := collection.FindOneAndUpdate(nil, bson.M{key: value}, bson.M{
		"$set": payload,
	}).Decode(&payload)
	if err != nil {
		return err
	}
	return nil
}

/**
Get By Id
*/
func (d DataStore) GetById(id string) map[string]interface{} {
	collection := d.Connection.(*mongo.Collection)
	Id, _ := primitive.ObjectIDFromHex(id)
	var result map[string]interface{}
	err := collection.FindOne(nil, bson.M{"_id": Id}).Decode(&result)
	if err != nil {
		return nil
	}
	return result
}

/**
Get One By Field Name
*/
func (d DataStore) GetByField(name string, value string) map[string]interface{} {
	collection := d.Connection.(*mongo.Collection)
	var result map[string]interface{}
	err := collection.FindOne(nil, bson.M{name: value}).Decode(&result)
	if err != nil {
		// Log error here.
		return nil
	}

	return result
}

func (d DataStore) FindByFields(filters map[string]interface{}) []map[string]interface{} {
	collection := d.Connection.(*mongo.Collection)
	cursor, err := collection.Find(nil, filters)
	result := make([]map[string]interface{}, 0)
	if err != nil {
		return result
	}

	for cursor.Next(nil) {
		var item map[string]interface{}
		_ = cursor.Decode(&item)
		result = append(result, item)
	}

	return result
}

func (d DataStore) DeleteByField(name string, value string) error {
	collection := d.Connection.(*mongo.Collection)
	_, err := collection.DeleteOne(nil, bson.M{name: value})
	if err != nil {
		return err
	}

	return nil
}

func (d DataStore) DeleteById(id string) error {
	collection := d.Connection.(*mongo.Collection)
	output, err := collection.DeleteOne(nil, bson.M{"_id": id})
	if err != nil {
		return err
	}

	fmt.Println(output, " Deleted document.")
	return nil
}

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
