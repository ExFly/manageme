package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/globalsign/mgo/bson"
)

func main() {
	data := make([]interface{}, 0)
	data = append(data, map[string]int{
		"a": 1,
		"b": 2,
	})

	jsons, err := json.MarshalIndent(data[0], "", "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsons))

	fmt.Println(T{"nt", bson.NewObjectId()})

	jsons, err = json.Marshal(T{Name: "test"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(jsons)

}

type T struct {
	Name string `json:"name"`
	//Objid string `json:"objid"`
	objid bson.ObjectId `json:"objid"`
}
