# db go package

this package abstracts the connection between go code and the mongo db to easily [read and write data](#functionality) from whereever needed.

The package takes care of the connection by using [environment variables](#environment-variables).

The package is setup to be limited to work with go compilted protofile structs. 
Click [here](#how-to-use-with-your-custom-structs) for further details 

## how to set up
First copy the package within your project root folder.

Then import the package by using ```import <yourRootFolder/>db```.

The public functions will create a connection the first time one of them is used. 
The established connection is then reused.

The connection will try to reconnect if the database will go offline.

If you want to add custom code that should be running when reconnecting to the database, you can use this function:
```
var RunWhenReconnected func()
```
here is an example on how to use
```
db.RunWhenReconnected = func(){
    log.Printf("reconnected to database")
}
```

## configuration and startup

The link is configured by giving an uri and using the New function.
You can create a new link by using this code snippet:
```
link := db.New(<your uri>, func(){
	//insert your custom functionality to be called after a reconnect
})
hqLink.Connect()
```

There is also a naming struct to help with database namings. The struct looks like this:
```
type DatabaseNamings struct {
	HqDb			string `yaml:"hqDb"`
	PlantUsaDb   	string `yaml:"plantUsaDb"`
	PlantChinaDb 	string `yaml:"plantChinaDb"`
	SupportDb 		string `yaml:"supportDb"`
	Orders    		string `yaml:"orders"`
	OrdersBroken	string `yaml:"ordersBroken"`
	Kpi       		string `yaml:"kpi"`
	Load      		string `yaml:"load"`
	Products  		string `yaml:"products"`
	Parts     		string `yaml:"parts"`
	Suppliers 		string `yaml:"suppliers"`
	Tickets			string `yaml:"Tickets"`
	State			string `yaml:"State"`
}
```
to use the functionality, it has to be added onto the db package with this line of code:
```
db.UpdateNamings(configs.DatabaseNamings)
```

## functionality

- New(uri string, runWhenConnectedMethod func()) Link -> returns a link to use all the other functions

### link related functions
- link.Connect(dbSettings Database) error -> establish connection
- link.TryReconnecting -> for manual reconnection calling (is also called from within the following functions)
- link.Disconnect() error -> for disconnecting

- link.Add(toAdd *T, collName string, dbName string) error -> Add a dbDocument
- link.UpdateById(toUpdate *T, id uint32, collName string, dbName string) -> update a dbDocument by using an [id](#specification-id)
- link.Remove(filter primitive.D, collName string, dbName string) error -> remove a dbDocument with a filter
- link.RemoveById(id uint32, collName string, dbName string) error -> remove a dbDocument by using an [id](#specification-id)
- link.GetAll(collName string, dbName string) ([]interface{}, error) -> returns all instances of the dbDocument as a list
- link.Get(filter primitive.D, collName string, dbName string) ([]interface{}, error) -> returns a list of instances of the dbDocument based on the given filter
- link.GetLast(collName string, dbName string) (*interface{}, error) -> get the last inserted instance of the collection
- link.WatchIncoming(collName string, dbName string, callback func(*interface{})) error -> add a handler for incoming database changes for the specified collection

### helpers
- TransformInterface[T interface{}](in interface{}) (out *T, err error) -> to transform the output interface into a custom struct. Gives an error if it cant be converted
- TransformInterfaces[T interface{}](in []interface{}) (out []*T, err error) -> transform the output interfaces in the same way as TransformInterface

direct access (should only be used if the functionality is not given with other functions):
- GetCollection(collName string) -> for more direct interactions with the database
- GetContext() -> a way to use the same context as the internal functions when using GetCollection()

### specification id
Everytime a function is using an id it refers to the id of the structs given within the dbDocument interface.
The user has to ensure the structs have an id or should not use the xxxById functions

## how to use with your custom structs
As mentioned in the introduction this package limits the interaction with the database to a specified set of structs.
To update or alter the list of structs the source code needs to be altered in two places.

### 1. update the dbDocument interface
This interface holds all the allowed struct names. All functions refer to this interface when accepting structs.
If you want to add "yourStruct" to the interfaces just add it like this:
```
type dbDocument interface {
	existingStruct | 
	yourStruct
}
```