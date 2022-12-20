package services

import (
	"go-todolist/entity"
	"go-todolist/model"
	"go-todolist/request"
	"log"

	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
)

// UserService is a contract about some user service can do
type UserService interface {
	// VerifyCredential is verify user credential
	VerifyCredential(email string, password string) interface{}

	// CreateUser is insert user to db and return user model to caller function
	CreateUser(user request.RegisterRequest) model.User

	// FindByEmail(email string, password string) model.User
}

// Create a new authService with the given userEntity.
type userService struct {
	userEntity entity.UserEntity
}

// NewUserService is creates a new instance of UserService with the given userEntity.
func NewUserService(userEntity entity.UserEntity) UserService {
	return &userService{userEntity: userEntity}
}

// VerifyCredential is verify user credential and return user model to caller function
func (s *userService) VerifyCredential(email string, password string) interface{} {
	// Verify user credential and return user model to caller function
	res := s.userEntity.VerifyCredential(email)

	// if res is user model then return user model to caller function
	if v, ok := res.(model.User); ok {

		// compare password with hashed password and return true if password is matched or return false if password is not matched
		comparedPassword := comparePassword(v.Password, []byte(password))

		// if email is matched and password is matched then return user model to caller function
		if v.Email == email && comparedPassword {
			// return user model to caller function
			return res
		}

		// return false if email is not matched or password is not matched
		return false
	}

	// return false if res is not user model
	return false
}

// CreateUser is insert user to db and return user model to caller function
func (s *userService) CreateUser(user request.RegisterRequest) model.User {
	// create user model
	userToCreate := model.User{}

	// fill user model with data from request model and return error if any error occur during mapping process or return nil if no error occur during mapping process and return user model to caller function to use it
	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}

	// insert user to db and return user model to caller function
	res := s.userEntity.InsertUser(userToCreate)

	// return user model to caller function
	return res
}

// func (s *userService) FindByEmail(email string) model.User {
// if err = gorm.Db
// }

// comparePassword is compare password with hashed password and return true if password is matched or return false if password is not matched
func comparePassword(hashedPwd string, plainPassword []byte) bool {
	// convert hashed password to byte array
	byteHash := []byte(hashedPwd)

	// compare password with hashed password
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}

	// return true if password is matched or return false if password is not matched
	return true
}
