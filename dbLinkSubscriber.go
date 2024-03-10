package db

import (
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
This function is used to watch the database for incoming changes.
Specify the collection and database and give a function to handle the incoming data.
*/
func (link Link)WatchIncoming(collName string, dbName string, callback func(id int)) error {
	pipeline := mongo.Pipeline{bson.D{{Key: "$match", Value: bson.D{{Key: "operationType", Value: "insert"}}}}}

	for {
		coll, err := link.GetCollection(collName, dbName)
		if err != nil {
			log.Printf("Establishing order change stream failed with error:\n\n%v\n\n Retrying...", err.Error())
			continue
		}
		changeStream, err := coll.Watch(ctx, pipeline)
		if err != nil {
			log.Printf("Establishing order change stream failed with error:\n\n%v\n\n Retrying...", err.Error())
			time.Sleep(time.Second * 5)
			continue
		}

		for changeStream.Next(ctx) {
			id, _, err := getIdFromInsertDocument(changeStream)
			if err != nil {
				log.Print(err.Error())
				continue
			}
			callback(id)
		}
	}
}

// getIdFromInsertDocument fetches the id from a change event insert document for further processing.
func getIdFromInsertDocument(cs *mongo.ChangeStream) (id int, opType string, err error) {
	var changeEvent bson.M
	if err = cs.Decode(&changeEvent); err != nil {
		return id, opType, err
	}

	var docMap primitive.M
	opType = changeEvent["operationType"].(string)

	switch opType {
	case "delete":
		docMap = changeEvent["fullDocumentBeforeChange"].(primitive.M)
	case "insert", "replace", "update":
		docMap = changeEvent["fullDocument"].(primitive.M)
	}

	docId, ok := docMap["id"].(int64)

	if !ok {
		return id, opType, fmt.Errorf("database contains corrupted document")
	}

	id = int(docId)

	return
}