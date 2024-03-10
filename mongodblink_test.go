package mongodblink_test

import (
	"log"
	"testing"

	mongodblink "github.com/lukas-czy/go-mongodb-link"
	"go.mongodb.org/mongo-driver/bson"
)

type testFile struct{
	Id int
	Name string
	Content string
}

const DB_URI = "mongodb://localhost:27017"
const TESTCOLLECTION = "testCollection"
const TESTDB = "testDB"

func TestAddGetAndRemoveFile(t *testing.T){
	newFile := testFile{
		Id:      0,
		Name:    "testFile",
		Content: "testContent",
	}
	client := mongodblink.New(DB_URI, func() {
		// This function is empty because it is not needed for the test.
	})
	if err := client.Connect(); err != nil {
		t.Fatalf("failed to connect to database. Error is %s", err)
	}
	if err := client.Add(&newFile, TESTCOLLECTION, TESTDB); err != nil {
		t.Fatalf("failed to add order. Error is %s", err)
	}
	dbOutput, err := client.GetAll(TESTCOLLECTION, TESTDB)
	if err != nil{
		t.Fatalf("failed to get orders. Error is %s", err)
	}
	orders, err := mongodblink.TransformInterfaces[testFile](dbOutput)
	if err != nil{
		t.Fatalf("failed to transform orders. Error is %s", err)
	
	}
	log.Printf("first order is %v", orders[0])
	if len(orders) <= 0{
		t.Fatalf("no orders received")
	}
	gotOrder := false
	for i := 0; i < len(orders); i++ {
		if orders[i].Id == newFile.Id{
			gotOrder = true
			break
		}
	}
	if gotOrder == false {
		t.Fatalf("no order with the added id found")
	}
	err3 := client.Remove(bson.D{{Key: "id", Value: newFile.Id}}, TESTCOLLECTION, TESTDB)
	if err3 != nil{
		t.Fatalf("failed to remove order. Error is %s", err3)
	}
}