package entity

import (
	"go-todolist/model"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//UserEntity is contract what UserEntity can do to db
type UserEntity interface {
	//InsertUser is insert user to db
	InsertUser(user model.User) model.User

	//VerifyCredential is verify user login
	VerifyCredential(email string) interface{}

	// FindByEmail(email string) model.User
}

// userConnection is a struct that implements connection to db with gorm
type userConnection struct {
	// connection to db with gorm
	connection *gorm.DB
}

// NewUserEntity is creates a new instance of UserEntity with gorm connection instance as parameter and return UserEntity interface instance to use
func NewUserEntity(db *gorm.DB) UserEntity {
	return &userConnection{
		connection: db,
	}
}

// InsertUser is insert user to db and return user model to caller function
func (db *userConnection) InsertUser(user model.User) model.User {
	// hash password
	user.Password = hashAndSalt([]byte(user.Password))
	db.connection.Save(&user)
	return user
}

// VerifyCredential is verify user credential and return user model to caller function if credential is correct or return nil if credential is incorrect
func (db *userConnection) VerifyCredential(email string) interface{} {
	var user model.User
	res := db.connection.Where("email = ?", email).Take(&user)

	// GORM debug model, show raw sql
	// panic(db.connection.Debug().Where("email = ?", email).Take(&user))

	// show user model struct
	// fmt.Printf("%+v\n", user)

	if res.Error == nil {
		return user
	}

	return nil
}

// func (db *userConnection) FindByEmail(email string) model.User {
// 	var user model.User
// 	// db.connection.Where("email = ? AND password = ?", email, password).Take(&user)
// 	db.connection.Where("email = ?", email).Take(&user)
// 	return user
// }

// hashAndSalt is hash password and return hashed password
func hashAndSalt(pwd []byte) string {
	// hash password
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)

	if err != nil {
		log.Println(err)

		// panic if failed to hash password
		panic("Failed to hash a password")
	}

	// return hashed password
	return string(hash)
}
