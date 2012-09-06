package mgomodel

import (
	"encoding/json"
	"errors"
	"github.com/mdennebaum/cheshire"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"reflect"
)

type Modeler interface {
	ID() bson.ObjectId
	Collection() string
}

type ValidatedModeler interface {
	Validators() []func(interface{}) bool
	RequiredFields() []string
}

type IndexedModeler interface {
	Indexes() []*mgo.Index
}

func JSON(inst Modeler) (string, error) {
	//marshal object to byte array
	jsonBytes, err := json.Marshal(inst)
	//throw up the error if its there
	if err != nil {
		return "", err
	}
	//cast byte array to string and return
	return string(jsonBytes), nil
}

func Valid(inst ValidatedModeler) bool {

	//loop over our required fields
	for _, requiredField := range inst.RequiredFields() {

		//grab a reflect value object for this instance
		obj := reflect.ValueOf(inst)

		//check if the field is nil in this instance or not
		if obj.FieldByName(requiredField).IsNil() {
			//return false if the required field is nil
			return false
		}
	}

	//loop over our user defined validator funcs
	for _, validator := range inst.Validators() {

		//run the validator
		if !validator(inst) {
			//if validator is false return false
			return false
		}
	}

	//everything is skippy
	return true

}

func Save(inst Modeler) error {

	//check if the id is set or not
	if inst.ID() == "" {
		//this is a new doc so lets just insert it
		return cheshire.Mongo().GetCollection(inst.Collection()).Insert(inst)
	}

	//this is a existing doc so lets update it
	return cheshire.Mongo().GetCollection(inst.Collection()).UpdateId(inst.ID(), inst)
}

func Delete(inst Modeler) error {

	//check if the id is set or not
	if inst.ID() != "" {
		//delete this obj
		return cheshire.Mongo().GetCollection(inst.Collection()).RemoveId(inst.ID())
	}

	return errors.New("Can't delete an unloaded object")
}

func Load(inst Modeler) error {
	return cheshire.Mongo().GetCollection(inst.Collection()).FindId(inst.ID()).One(inst)
}
