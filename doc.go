/*
mgomodel is a GO package that lets you interact with mongodb in a familiar model based way. It 
builds on top of the awsome mgo mongo driver http://labix.org/mgo.

mgomodel exposes the familiar save, delete, update, find, load methods one would expect with 
a normal ORM type setup. There is also a configurable mechanism for building in custom data validation.


Comming soon:

a utility for managing associated mongo indexes.


Example:
	import (
		"labix.org/v2/mgo/bson"
		"github.com/mdennebaum/mgomodel"
		"errors"
	)
	
	
	type AddressInfo struct{
		Address1 string
		Address2 string
		City string
		State string
		Zip string
	}

	type Person struct{
		Id bson.ObjectId "_id,omitempty"
		Name string
		Address *AddressInfo
		User string
	}

	//satisfy the modler interface
	func (this *Person) Collection() string {
		return "people"
	}

	//satisfy the modler interface
	func (this *Person) ID() bson.ObjectId {
		return this.Id
	}

	//satisfy the validatedmodler interface
	func (this *Person) RequiredFields() []string {
		return []string{"Name", "User"}
	}
	
	func NewPerson(name,user) *Person{
		return &Person {
			Name: name,
			User: user,
			Address: &AddressInfo{} 
		}
	}

	func main(){
		//init the mongo connection
		mgomodel.Mongo().Connect("localhost")

		//set the default database
		mgomodel.Mongo().SetDB("test")
		
		//create an empty user instance to hold our user document
		user := User{}

		//find the user with the username
		mgomodel.Mongo().Collection("users").find(bson.M{"Username":"mattdennebaum"}).One(&user)

		//create a new person object
		person := NewPerson("Matt Dennebaum",user.ID().String())
		
		//validate the person
		if mgomodel.Valid(person) {
			//save the person
			err := mgomodel.Save(person)
		}
	}
*/
package mgomodel