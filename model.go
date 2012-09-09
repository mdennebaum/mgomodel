package mgomodel

import (
	"encoding/json"
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"reflect"
	"log"
)

type Modeler interface {
	ID() bson.ObjectId
	Collection() string
}

type ValidatedModeler interface {
	Validators() []func(interface{})error
	RequiredFields() []string
}

type DefaultedModeler interface {
	DefaultValues() map[string]interface{}
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

func Valid(inst ValidatedModeler) error {
	//loop over our required fields
	for _, requiredField := range inst.RequiredFields() {
		
		//grab the field
		field := reflect.ValueOf(inst).Elem().FieldByName(requiredField)

		//check if this is an invalid field
		if field == reflect.ValueOf(nil){
			//kill at compile time. no reason to let it become a runtime issue
			log.Fatalf("mgomodel.Valid Error: invalid requiredField defined: %s", requiredField)
		}

		//check if the type is nil
		switch field.Kind() {
			case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
				if field.Len() == 0 {
					return errors.New("required field isn't set: "+requiredField)
				}
			case reflect.Interface, reflect.Ptr:
				if field.IsNil() {
					return errors.New("required field isn't set: "+requiredField)
				}
		}

	}

	//loop over our user defined validator funcs
	for _, validator := range inst.Validators() {
		//run the validator
		err := validator(inst)
		
		//run the validator
		if err != nil{
			//print the error
			log.Println(err.Error())

			//if validator is false return false
			return err
		}
	}

	//everything is skippy
	return nil

}

func setDefaults(inst Modeler){
	for key, val := range inst.(DefaultedModeler).DefaultValues() {
		reflect.ValueOf(inst).Elem().FieldByName(key).Set(reflect.ValueOf(val))
	}
}

func Save(inst Modeler) error {

	//check if the id is set or not
	if inst.ID() == "" {
		
		//if DefaulyValues are defined then set them
		if reflect.ValueOf(inst).MethodByName("DefaultValues").IsValid() {
			setDefaults(inst)
		}

		//this is a new doc so lets just insert it
		return Mongo().Collection(inst.Collection()).Insert(inst)
	}

	//this is a existing doc so lets update it
	return Mongo().Collection(inst.Collection()).UpdateId(inst.ID(), inst)
}

func Delete(inst Modeler) error {

	//check if the id is set orMethodByName("")not
	if inst.ID() != "" {
		//delete this obj
		return Mongo().Collection(inst.Collection()).RemoveId(inst.ID())
	}

	return errors.New("Can't delete an unloaded object")
}

func Load(inst Modeler) error {
	return Mongo().Collection(inst.Collection()).FindId(inst.ID()).One(inst)
}
