package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

type Link struct{
	Uri string
	Client *mongo.Client
	Reconnecting bool
	isStarted bool
	RunWhenReconnected func()
}

func New(uri string, runWhenReconnectedMethod func()) Link {

	return Link {
		Uri: uri,
		Reconnecting: false,
		isStarted: false,
		RunWhenReconnected: runWhenReconnectedMethod,
	}
}

/*
This function connects the database to the given uri.
*/
func (link *Link)Connect() error {
	link.Reconnecting = false
	link.isStarted = false

	clientOptions := options.Client().ApplyURI(link.Uri)
	clientOptions.SetServerSelectionTimeout(1 * time.Second)
	newClient, err := mongo.Connect(ctx, clientOptions)
	link.Client = newClient
	if err != nil {
		return err
	}

	if err := link.Client.Ping(ctx, nil); err != nil {
		link.TryReconnecting()
		return err
	}
	log.Printf("database link connected to %s", link.Uri)
	return nil
}
/*
Internal function for handling reconnection.
This function should be handled by the database link itself.
The function is written to run as a goroutine.
*/
func (link *Link)reconnect(){
	if !link.isStarted{
		link.isStarted = true
	}else{
		return
	}
	log.Println("starting reconnect mode")
	link.Client.Disconnect(ctx)
	for{
		log.Println("reconnecting...")
		clientOptions := options.Client().ApplyURI(link.Uri)
		clientOptions.SetServerSelectionTimeout(1 * time.Second)
		newClient, err := mongo.Connect(ctx, clientOptions)
		if err != nil{
			time.Sleep(5 * time.Second)
			if newClient != nil{
				newClient.Disconnect(ctx)
			}
			continue
		}
		link.Client = newClient
		if err := link.Client.Ping(ctx, nil); err != nil {
			time.Sleep(5 * time.Second)
			link.Client.Disconnect(ctx)
			continue
		}
		log.Println("database link reconnected")
		go link.RunWhenReconnected()
		return
		
	}
}
/*
This function triggers the reconnection process of the database link.
Usually this is not needed because the database link tries to reconnect automatically.
*/
func (link *Link)TryReconnecting(){
	if err := link.Client.Ping(ctx, nil); err != nil {
		if link.Reconnecting{
			log.Println("database is already trying to reconnect")
			return
		}
		link.Reconnecting = true
		log.Println("database link is going into reconnect mode")
		go link.reconnect()
	}
	link.Reconnecting = false
}
/*
disconnect the database link
*/
func (link *Link)Disconnect() error {
	err := link.Client.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (link *Link)IsAlive() bool {
	if link.Client == nil {
		return false
	}
	if link.Reconnecting{
		return false
	}
	if err := link.Client.Ping(ctx, nil); err != nil {
		return false
	}
	return true

}
/*
Get a mongo db collection for a custom interaction with the mongo db.
Be advised that using the given functions is easier but this gives more freedom of interaction.
Usually this function is called from within the other functions.
*/
func (link *Link)GetCollection(collName string, dbName string) (*mongo.Collection, error) {
	//check if a client is already connected. If not, connect it
	if link.Client == nil {
		return nil, fmt.Errorf("database is not connected yet. Please connect with the Connect method")
	}
	if link.Reconnecting{
		return nil, fmt.Errorf("database is currently reconnecting. Please try again later")
	}
	if err := link.Client.Ping(ctx, nil); err != nil {
		link.TryReconnecting()
		return nil, err
	}
	//create the collection reference
	indexId := mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	coll := link.Client.Database(dbName).Collection(collName)
	_, err3 := coll.Indexes().CreateOne(ctx, indexId)
	if err3 != nil {
		return nil, err3
	}
	return coll, nil
}
/*
When using GetCollection this is your ticket to get the same context as all the other functions
*/
func (link *Link)GetContext() context.Context {
	return ctx
}

/*
Add a new dbDocument into the mongodb.

example call:
err := Add(toAdd, db.ORDERS, db.HQ_DB)
*/
func (link *Link)Add(toAdd interface{}, collName string, dbName string) error {
	coll, err1 := link.GetCollection(collName, dbName)
	if err1 != nil {
		return err1
	}
	_, err2 := coll.InsertOne(ctx, toAdd)
	if err2 != nil {
		return err2
	}
	return nil
}

/*
Update an existing document within the mongo db using a the id (not the internal mongodb _id, but the one from the struct).

example call: 
err := UpdateById(toUpdate, toUpdate.Id, db.ORDERS, db.HQ_DB)
*/
func (link *Link)UpdateById(toUpdate interface{}, id uint32, collName string, dbName string) error {
	coll, err1 := link.GetCollection(collName, dbName)
	if err1 != nil {
		return err1
	}
	filter := bson.D{{Key: "id", Value: id}}
	_, err2 := coll.ReplaceOne(ctx, filter, toUpdate)
	if err2 != nil {
		return err2
	}
	return nil
}

/*
Remove a dbDocument from the mongodb by using a filter.

example filter: 
bson.D{{Key: "id", Value: toRemove.Id}}
example call: 
err := Remove[interfaces.Order](bson.D{{Key: "id", Value: toRemove.Id, db.ORDERS, db.HQ_DB}})
*/
func (link *Link)Remove(filter primitive.D, collName string, dbName string) error {
	coll, err1 := link.GetCollection(collName, dbName)
	if err1 != nil {
		return err1
	}

	_, err2 := coll.DeleteMany(ctx, filter)
	if err2 != nil {
		return err2
	}
	return nil
}

/*
Remove a dbDocument from the mongodb by using a the struct id (not the internal _id from mongodb)
This is a simplified version of the Remove function.

To use the function the type that you want to remove needs to be specified.

example call: 
err := Remove[interfaces.Order](toRemove.id, db.ORDERS, db.HQ_DB})
*/
func (link *Link)RemoveById(id uint32, collName string, dbName string) error {
	coll, err1 := link.GetCollection(collName, dbName)
	if err1 != nil {
		return err1
	}
	filter := bson.D{{Key: "id", Value: id}}
	_, err2 := coll.DeleteMany(ctx, filter)
	if err2 != nil {
		return err2
	}
	return nil
}

/*
Get all instances of a dbDocument from the mongodb.
example call:
orders, err := GetAll[interfaces.Order](db.ORDERS, db.HQ_DB)
*/
func (link *Link)GetAll(collName string, dbName string) ([]interface{}, error) {
	coll, err1 := link.GetCollection(collName, dbName)
	if err1 != nil {
		return nil, err1
	}
	filter := bson.D{}
	cursor, err2 := coll.Find(ctx, filter)
	if err2 != nil {
		return nil, err2
	}
	var results []interface{}
	err3 := cursor.All(ctx, &results)
	if err3 != nil {
		return nil, err3
	}
	return results, nil
}

/*
Get all filtered instances of a dbDocument from the mongodb.

example filter: 
bson.D{{Key: "id", Value: toRemove.Id}}
example call: 
orders, err := Get[interfaces.Order](bson.D{{Key: "id", Value: toRemove.Id, db.ORDERS, db.HQ_DB}})
*/
func (link *Link)Get(filter primitive.D, collName string, dbName string) ([]interface{}, error) {
	coll, err1 := link.GetCollection(collName, dbName)
	if err1 != nil {
		return nil, err1
	}
	cursor, err2 := coll.Find(ctx, filter)
	if err2 != nil {
		return nil, err2
	}
	var results []interface{}
	err3 := cursor.All(ctx, &results)
	if err3 != nil {
		return nil, err3
	}
	return results, nil
}
/*
Get the last inserted document from the mongodb.
example call:
lastOrder, err := GetLast[interfaces.Order](db.ORDERS, db.HQ_DB)
*/
func (link *Link)GetLast(collName string, dbName string) (*interface{}, error) {
	coll, err1 := link.GetCollection(collName, dbName)
	if err1 != nil {
		return nil, err1
	}
	opts := options.FindOne().SetSort(bson.M{"$natural": -1})
	var lastrecord interface{}
	err2 := coll.FindOne(ctx, bson.M{}, opts).Decode(&lastrecord)
	if err2 != nil {
		if err2.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err2
	}
	return &lastrecord, nil
}
/*
helper function to transform a received interface into the correct structs
example usage:
examplePartInterface, err := hqLink.GetLast(db.PARTS, db.HQ_DB)
if err != nil {
	//handle error
}
examplePart, err := db.TransformInterface[interfaces.Part](examplePartInterface)
if err != nil {
	//handle error
}
*/
func TransformInterface[T interface{}](in interface{}) (out *T, err error){
	var result T
		bytes, err := bson.Marshal(in)
		if err != nil{
			return nil, err
		}
		if err := bson.Unmarshal(bytes, &result); err != nil{
			return nil, err
		}
		return &result, nil
}
/*
helper function to transform a list of received interfaces into the correct structs.
See TransformInterface for more information.
*/
func TransformInterfaces[T interface{}](in []interface{}) (out []*T, err error){
	output := make([]*T, len(in))
	for i, v := range in {
		out, err :=  TransformInterface[T](v)
		if err != nil{
			return nil, err
		}
		output[i] = out
	}
	return output, nil
}