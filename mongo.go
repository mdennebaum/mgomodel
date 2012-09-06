package mgomodel

import (
	"labix.org/v2/mgo"
	"sync"
)

type mongoDB struct {
	Session *mgo.Session
	DB      *mgo.Database
}

//singleton instance
var mongo_instance *mongoDB
var mongo_once sync.Once

//Mongo is an accessor for the mongo singleton
func Mongo() *mongoDB {
	mongo_once.Do(func() {
		mongo_instance = new(mongoDB)
	})

	return mongo_instance
}

//Connect connects use to the mongo server
func (this *mongoDB) Connect(server string) {
	var err error

	this.Session, err = mgo.Dial(server)

	if err != nil {
		panic(err)
	}
}

//SetDB is an accessor to allow you to set the target DB
func (this *mongoDB) SetDB(name string) {
	this.DB = this.Session.DB(name)
}

//Collection is an accessor to get back a pointer to the collection object
func (this *mongoDB) Collection(name string) *mgo.Collection {
	return this.DB.C(name)
}