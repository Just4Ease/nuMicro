package mongo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var dataStore *DataStore
var SampleData []map[string]interface{}

func init() {
	dataset, _ := os.Open("./dataset.json")
	defer dataset.Close()
	rawJSON, _ := ioutil.ReadAll(dataset)
	_ = json.Unmarshal(rawJSON, &SampleData)
}

func TestMongoDataStore(t *testing.T) {
	t.Run("Connect Database", func(t *testing.T) {
		dataStore = New("mongodb://localhost:27017/test", "test", "test")
	})

	t.Run("Delete Many", func(t *testing.T) {
		fields := make(map[string]interface{})
		fields["email"] = map[string]interface{}{
			"$exists": true,
		}
		if err := dataStore.DeleteMany(fields); err != nil {
			fmt.Println(err.Error(), " Error Deleting Existing Test Data")
			t.Fail()
		}

	})

	t.Run("Create many", func(t *testing.T) {
		var i []interface{}
		b, _ := json.Marshal(SampleData)
		_ = json.Unmarshal(b, &i)
		err := dataStore.SaveMany(i)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}
	})

	t.Run("Create Entry", func(t *testing.T) {
		if dataStore != nil {
			entry := map[string]interface{}{
				"_id":        "5e787f774a80487e0dd42488",
				"firstName":  "Justice",
				"lastName":   "Nefe",
				"guid":       "fa04ae4b-40f2-4018-911d-a5ae110b28e1",
				"isActive":   false,
				"balance":    "$1,983.33",
				"picture":    "http://placehold.it/32x32",
				"age":        25,
				"company":    "Neofortis",
				"email":      "justicenefe@gmail.com",
				"phone":      "07056031137",
				"address":    "908 Cortelyou Road, Salvo, Idaho, 3587",
				"about":      "lorem ispiym",
				"registered": time.Now(),
				"latitude":   "22.42297",
				"longitude":  "132.341514",
			}
			output, err := dataStore.Save(entry)

			if err != nil {
				t.Fail()
			}

			if output["firstName"] != entry["firstName"] {
				t.Fail()
			}
		}
	})

	t.Run("Find One by Id", func(t *testing.T) {
		output := dataStore.FindById("5e787f774a80487e0dd42488", nil)
		if output == nil {
			t.Fail()
		}
	})

	t.Run("Find One by fields", func(t *testing.T) {
		filters := make(map[string]interface{})
		filters["_id"] = "5e787f774a80487e0dd42488"
		filters["firstName"] = "Justice"
		output := dataStore.FindOne(filters, nil)
		if output == nil {
			t.Fail()
		}
	})

	t.Run("Find Many", func(t *testing.T) {
		filters := make(map[string]interface{})
		filters["firstName"] = map[string]interface{}{
			"$exists": true,
		}
		output := dataStore.FindMany(filters, nil, nil, 10, 0)
		if output == nil {
			t.Fail()
		}
		if len(output) != 10 {
			t.Fail()
		}
		fmt.Println(len(output), " Output.")
	})

	t.Run("Update one by Id", func(t *testing.T) {
		payload := map[string]interface{}{
			"firstName": "Emmanuel",
			"email":     "emmanuel@email.com",
		}
		if err := dataStore.UpdateById("5e787f774a80487e0dd42488", payload); err != nil {
			fmt.Println(err.Error(), " Error updating record.")
			t.Fail()
		}

		output := dataStore.FindById("5e787f774a80487e0dd42488", nil)
		if output == nil {
			t.Fail()
		}

		if output["firstName"] != payload["firstName"] {
			t.Fail()
		}
	})

	t.Run("Update one by field", func(t *testing.T) {
		payload := map[string]interface{}{
			"lastName": "Joseph",
			"phone":    "080924311111",
		}
		if err := dataStore.UpdateOne(map[string]interface{}{
			"_id":      "5e787f77f31d7514045a4f44",
			"lastName": "Hood",
		}, payload); err != nil {
			fmt.Println(err.Error(), " Error updating record by field")
			t.Fail()
		}

		output := dataStore.FindById("5e787f77f31d7514045a4f44", nil)
		if output == nil {
			t.Fail()
		}

		if output["lastName"] != payload["lastName"] {
			t.Fail()
		}
	})

	t.Run("Update many", func(t *testing.T) {
		payload := map[string]interface{}{
			"isActive": false,
		}
		if err := dataStore.UpdateMany(map[string]interface{}{
			"isActive": true,
		}, payload); err != nil {
			fmt.Println(err.Error(), " Error updating record by field")
			t.Fail()
		}

		output := dataStore.FindMany(map[string]interface{}{
			"isActive": true,
		}, nil, nil, 10, 0)
		if output == nil {
			t.Fail()
		}

		if len(output) >= 1 {
			t.Fail()
		}
	})

	t.Run("Delete one by Id", func(t *testing.T) {
		if err := dataStore.DeleteById("5e787f774a80487e0dd42488"); err != nil {
			fmt.Println(err.Error(), " Error deleting")
			t.Fail()
		}
	})

	t.Run("Delete one by field", func(t *testing.T) {
		if err := dataStore.DeleteOne(map[string]interface{}{
			"_id": "5e787f774a80487e0dd42488",
		}); err != nil {
			fmt.Println(err.Error(), " Error deleting")
			t.Fail()
		}
	})
}
