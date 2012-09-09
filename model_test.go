package mgomodel

import (
	"labix.org/v2/mgo/bson"
	"time"
	"errors"
	"testing"
)

type User struct{
	Id bson.ObjectId "_id,omitempty"
	Username string
	Password string
	Email string
	CreatedAt time.Time
	Roles []string
}

//satisfy the modler interface
func (this *User) Collection() string {
	return "users"
}

//satisfy the modler interface
func (this *User) ID() bson.ObjectId {
	return this.Id
}

//satisfy the validatedmodler interface
func (this *User) RequiredFields() []string {
	return []string{"Username", "Password", "Email"}
}

//satisfy the defaultedmodler interface
func (this *User) DefaultValues() map[string]interface{} {
	return map[string]interface{}{
		"Roles": []string{"member"},
		"CreatedAt": time.Now(),
	}
}

func NewUser(username string,password string,email string) *User{
	return &User{
		Username: username,
		Password: password, //WARNING: in real life you would hash this right?!!
		Email:    email,
	}
}

//satisfy the validatedmodler interface
func (this *User) Validators() []func(interface{})error {
	
	//a username validator
	var usernameValidator = func(inst interface{})error {
		//make sure go knows its a User type
		user := inst.(*User)

		//check if we have a username
		if user.Username != "" {
			//check if username is at least 5 chars long
			if len(user.Username) >= 5 {
				//return true cause we are valid
				return nil
			}else{
				return errors.New("username is less then 5 chars long")
			}
		}
		//this shouldnt be possible as we allready checked it with the RequiredFields but what the hell
		return errors.New("username is empty")
	}

	//return an array of the validators
	return []func(inst interface{})error{usernameValidator}
}

//TestCreateUser tests validating and creating
func TestCreateUser(t *testing.T){
	//init the mongo connection
	Mongo().Connect("localhost")

	//close mongo once were done
	defer Mongo().Session.Close()

	//set the default database
	Mongo().SetDB("test")

	//create a new user object
	user := NewUser("mattdennebaum","s0m#Pa$2wD", "matt@quantumsp.in")

	//check to be sure the instance is valid
	if Valid(user) == nil {
		//save the valid user
		err := Save(user)
		//check if we had any issues saving
		if err != nil{
			//output the error			
			t.Error(err.Error())
		}
	}
}

//TestRequiredFieldUser validating required field validation
func TestRequiredFieldUser(t *testing.T){
	//init the mongo connection
	Mongo().Connect("localhost")

	//close mongo once were done
	defer Mongo().Session.Close()

	//set the default database
	Mongo().SetDB("test")

	//create a new user object
	user := NewUser("mattdennebaum","s0m#Pa$2wD","")

	//check to be sure the instance is valid
	if Valid(user) == nil{
		t.Error("user is valid but it shouldnt be")
		t.Fail()
	}
}