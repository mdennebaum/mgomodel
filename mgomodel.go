package mgomodel

import (
	"encoding/json"
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"reflect"
	"log"
)

//Modeler is the basic interface type for model instances
type Modeler interface {
	ID() bson.ObjectId
	Collection() string
}

//ValidatedModeler defines an interface for models that support validation funtionality
type ValidatedModeler interface {
	Validators() []func(interface{})error
	RequiredFields() []string
}

//DefaultedModeler defines an interface for models that support default value funtionality
type DefaultedModeler interface {
	DefaultValues() map[string]interface{}
}

//IndexedModeler defines an interface for models that support indexing funtionality.
//indexes are not added automatically. You must use the index manager utility. 
type IndexedModeler interface {
	Indexes() []*mgo.Index
}

//Json outputs a model instance as a json packet
func Json(inst Modeler) (string, error) {
	//marshal object to byte array
	jsonBytes, err := json.Marshal(inst)
	//throw up the error if its there
	if err != nil {
		return "", err
	}
	//cast byte array to string and return
	return string(jsonBytes), nil
}

//Valid validates a model based on the RequiredFields and Validators methods
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

//setDefaults will set the default values defined in the DefaultValues method
func setDefaults(inst Modeler){
	for key, val := range inst.(DefaultedModeler).DefaultValues() {
		reflect.ValueOf(inst).Elem().FieldByName(key).Set(reflect.ValueOf(val))
	}
}

//Save will save this instance. If there is an Id field thats not set 
//we assume this is a new model and we create a new objectid and insert
//the new document. If there is an Id then we update the document at that 
//Id. If its an insert we will first set all the default values. 
func Save(inst Modeler) error {

	//check if the id is set or not
	if inst.ID() == "" {
		//try to get the id field
		field := reflect.ValueOf(inst).Elem().FieldByName("Id")

		//make sure this model has an id field
		if field != reflect.ValueOf(nil){
			//set the Id field so we have it later
			field.Set(reflect.ValueOf(bson.NewObjectId()))
		}

		//if DefaulyValues are defined then set them
		if reflect.ValueOf(inst).MethodByName("DefaultValues").IsValid() {
			//set the defaults
			setDefaults(inst)
		}

		//this is a new doc so lets just insert it
		return Mongo().Collection(inst.Collection()).Insert(inst)
	}

	//this is a existing doc so lets update it
	return Mongo().Collection(inst.Collection()).UpdateId(inst.ID(), inst)
}

//Delete will remove the current document from the collection
func Delete(inst Modeler) error {

	//check if the id is set or not
	if inst.ID() != "" {
		//delete this obj
		return Mongo().Collection(inst.Collection()).RemoveId(inst.ID())
	}

	//return an error since we cant delete something thats not created
	return errors.New("Can't delete an unloaded object")
}

//Load a document. You must first create an empty Model and set the Id
func Load(inst Modeler) error {
	return Mongo().Collection(inst.Collection()).FindId(inst.ID()).One(inst)
}