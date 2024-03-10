package db

// dbs
var HQ_DB = "Headquarters"
var PLANT_USA_DB = "Plant_USA"
var PLANT_CHINA_DB = "Plant_China"
var SUPPORT_DB = "Support"

// collections
var PRODUCTS = "Products"
var PARTS = "Parts"
var SUPPLIERS = "Suppliers"
var ORDERS = "Orders"
var ORDERS_BROKEN = "Orders_Broken"
var KPI = "KPI"
var LOAD = "Load"
var TICKETS = "Tickets"
var STATE = "State"

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

func UpdateNamings(namings DatabaseNamings){
	HQ_DB = namings.HqDb
	PLANT_USA_DB = namings.PlantUsaDb
	PLANT_CHINA_DB = namings.PlantChinaDb
	SUPPORT_DB = namings.SupportDb
	PRODUCTS = namings.Products
	PARTS = namings.Parts
	SUPPLIERS = namings.Suppliers
	ORDERS = namings.Orders
	ORDERS_BROKEN = namings.OrdersBroken
	KPI = namings.Kpi
	LOAD = namings.Load
	TICKETS = namings.Tickets
	STATE = namings.State
}