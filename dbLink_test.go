package db

import (
	"log"
	"testing"

	"git.thm.de/vs-ws-23/efridge/interfaces"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAddGetAndRemoveOrder(t *testing.T){
	newOrder := interfaces.Order{
		Id: 0,
		CustomerId: 1234,
		Jobs: []*interfaces.Job{
			{ProductId: 1234, Quantity: 1},
		},
		Status: interfaces.Order_OPEN,
	}
	client := New("mongodb://localhost:27017/?replicaSet=rsHQ", func() {})
	if err := client.Connect(); err != nil{
		t.Fatalf("failed to connect to database. Error is %s", err)
	}
	if err := client.Add(&newOrder, ORDERS, HQ_DB); err != nil{
		t.Fatalf("failed to add order. Error is %s", err)
	}
	dbOutput, err := client.GetAll(ORDERS, HQ_DB)
	if err != nil{
		t.Fatalf("failed to get orders. Error is %s", err)
	}
	orders, err := TransformInterfaces[interfaces.Order](dbOutput)
	if err != nil{
		t.Fatalf("failed to transform orders. Error is %s", err)
	
	}
	log.Printf("first order is %v", orders[0])
	if len(orders) <= 0{
		t.Fatalf("no orders received")
	}
	gotOrder := false
	for i := 0; i < len(orders); i++ {
		if orders[i].Id == newOrder.Id{
			gotOrder = true
			break
		}
	}
	if gotOrder == false {
		t.Fatalf("no order with the added id found")
	}
	err3 := client.Remove(bson.D{{Key: "id", Value: newOrder.Id}}, ORDERS, HQ_DB)
	if err3 != nil{
		t.Fatalf("failed to remove order. Error is %s", err3)
	}
}