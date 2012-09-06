package mgomodel

import (
	"labix.org/v2/mgo/bson"
	"time"
	"regexp"
	"crypto/sha1"
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
	return []string{"username", "password"}
}

func NewUser(username,password,email) *User{
	var hashedpw string
	
	if password != "" {
		hasher := sha1.New()
		hashedpw = string(hasher.Sum([]byte(password)))
	}

	return &User{
		Username: username,
		Password: hashedpw,
		Email:    email,	
	}
}

//satisfy the validatedmodler interface
func (this *User) Validators() []func(interface{})(bool,error) {
	
	//an email validator
	var emailValidator = func(conf *Config) (bool,error) {
		//check if we have an email
		if this.Email != "" {

			//compile a regex to match against email addresses
			exp, err := regexp.Compile("[a-zA-Z0-9.!#$%&'*+-/=?\^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)* ")
			
			//check to see if we had an issue compiling the regex
			if err != nil {
				//return false and tell everyone our regex is shit
				return false, errors.New("problem compiling regex")
			}
			//check if the email passes the regex
			if !exp.MatchString(this.Email) { 
				//woops looks like its a shit email... or our regex is too stupid (always possible)
				return false, errors.New("email doesn't seem to be valid")
			}
		}else{
			//since email address isnt required if its just not set we will return valid
			return true, nil
		}

	}

	//return an array of the validators
	return []func(conf *Config)bool{}{emailValidator} nil
}

//TestCreateUser tests validating and creating of a user
func TestCreateUser(*testing.T){
	//init the mongo connection
	Mongo().Connect("localhost")

	//close mongo once were done
	defer Mongo().Session.Close()

	//set the default database
	Mongo().SetDB("test")

	//create a new user object
	user := NewUser("mattdennebaum","s0m#Pa$2wD", "matt@quantumsp.in")
	
	//check to be sure the instance is valid
	if Valid(user) {
		//save the valid user
		err := Save(user)
		//check if we had any issues saving
		if err != nil{
			//output the error			
			t.Errorf(err)
		}
	}
}