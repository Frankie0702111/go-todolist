package services

import (
	"fmt"
	"go-todolist/entity"
	"reflect"

	"go-todolist/utils/log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT Service is a contract of what a JWT Service should be able to do.
type JWTService interface {
	// Generate a new token
	GenerateToken(userID uint64, t time.Time) string

	// Validate the token
	ValidateToken(token string) (*jwt.Token, error)

	// Refresh the token
	RefreshToken(token string) string

	// User logout and remove the token from redis
	Logout(authHeader string) bool

	// Authorize JWT for middleware
	AuthJWT(authHeader string) string

	GoogleGenerateToken(data interface{}) string
}

// jwtCustomClaim is a struct that contains the custom claims for the JWT
type jwtCustomClaim struct {
	// The userId is the only required field
	UserID uint64 `json:"user_id"`

	// This is a registered JWT claim (StandardClaims are deprecated)
	jwt.RegisteredClaims
}

// jwtService is a struct that implements the JWTService interface
type jwtService struct {
	// Secret key used to sign the token
	secretKey string

	// Who creates the token
	issuer string

	// conntection to redis
	redisEntity entity.RedisEntity
	userEntity  entity.UserEntity
}

//NewJWTService method is creates a new instance of JWTService
func NewJWTService(redisEntity entity.RedisEntity, userEntity entity.UserEntity) JWTService {
	return &jwtService{
		// Call the getSecretKey function to get the secret key
		secretKey: getSecretKey(),

		// who creates the token
		issuer: "gojwt",

		// connection: rdb,
		redisEntity: redisEntity,
		userEntity:  userEntity,
	}
}

// getSecretKey Create get the secret key from the environment variable
func getSecretKey() string {
	// Get the secret key from the environment variable
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		// If the environment variable is empty, use a default value
		secretKey = "learnGolangJWTToken"
	}
	return secretKey
}

// getUserDataByToken Get user data by token
func getUserDataByToken(authHeader string, claims jwt.Claims, secretKey string) (interface{}, error) {
	token, err := jwt.ParseWithClaims(authHeader, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	return token, err
}

// GetTokenTTL Get token TTL from .env file
func GetTokenTTL() int {
	// Get the secret key from the environment variable
	stringTTL := os.Getenv("JWT_TTL")
	if stringTTL == "" {
		// If the environment variable is empty, use a default value
		stringTTL = "900"
	}

	intTTL, _ := strconv.Atoi(stringTTL)

	return intTTL
}

// GenerateToken Create a new token object, specifying signing method and the claims
func (s *jwtService) GenerateToken(userID uint64, t time.Time) string {
	jwtTTL := GetTokenTTL()
	// Create the Claims struct with the required claims for the JWT
	claims := &jwtCustomClaim{
		// userId is the only required field
		userID,
		jwt.RegisteredClaims{
			// 1 day expiration
			ExpiresAt: jwt.NewNumericDate(t),
			// when the token was issued/created (now)
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// Who creates the token
			Issuer: s.issuer,
		},
	}

	// Sign the token with our secret
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with an expiration time
	token, err := generateToken.SignedString([]byte(s.secretKey))
	if err != nil {
		// If there is an error, return empty string
		log.Error("Failed to process request : Signature failed")
		return ""
	}

	_, setRedisErr := s.redisEntity.Set("token"+strconv.FormatUint(userID, 10), token, time.Duration(jwtTTL)*time.Second)
	if setRedisErr != nil {
		log.Error("Failed to set the token in redis : " + setRedisErr.Error())
	}

	// Return the token to the user, along with an expiration time
	return token
}

// ValidateToken validates the token and returns the claims
func (s *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	// Parse the token
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			// Return an error if the signing method isn't HMAC
			log.Errorf("Unexpected signing method", t_.Header["alg"])
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		// Return the key
		return []byte(s.secretKey), nil
	})
}

// RefreshToken refresh the token
func (s *jwtService) RefreshToken(authHeader string) string {
	jwtTTL := GetTokenTTL()
	claims := &jwtCustomClaim{}
	_, err := getUserDataByToken(authHeader, claims, s.secretKey)
	if err != nil {
		log.Error("Failed to get user data (RefreshToken) : " + err.Error())
		return ""
	}

	// int to time.Duration
	token := s.GenerateToken(claims.UserID, time.Now().Add(time.Duration(jwtTTL)*time.Second))

	return token
}

// Logout User logout and remove token from redis
func (s *jwtService) Logout(authHeader string) bool {
	claims := &jwtCustomClaim{}
	_, erro := getUserDataByToken(authHeader, claims, s.secretKey)
	if erro != nil {
		log.Error("Failed to get user data (logout) : " + erro.Error())
		return false
	}

	get, err := s.redisEntity.Del("token" + strconv.Itoa(int(claims.UserID)))
	if get == int64(1) {
		return true
	}
	if err != nil {
		log.Error("Failed to delete the token in redis : " + err.Error())
		return false
	}

	return false
}

// AuthJWT Get the redis token and return the middleware
func (s *jwtService) AuthJWT(authHeader string) string {
	claims := &jwtCustomClaim{}
	_, erro := getUserDataByToken(authHeader, claims, s.secretKey)
	if erro != nil {
		log.Error("Failed to get user data (AuthJWT) : " + erro.Error())
		return ""
	}

	get, err := s.redisEntity.Get("token" + strconv.Itoa(int(claims.UserID)))
	if err != nil {
		log.Error("Failed to set the token in redis (AuthJWT) : " + err.Error())
		return ""
	}

	return fmt.Sprintf("%v", get)
}

func (s *jwtService) GoogleGenerateToken(data interface{}) string {
	jwtTTL := GetTokenTTL()
	// Test how to get interface values
	googleInfo := reflect.ValueOf(data)

	findByEmail := s.userEntity.FindByEmail(googleInfo.FieldByName("Email").String())
	if findByEmail.ID == 0 {
		return ""
	}

	fmt.Println("GoogleGenerateToken ID : ", googleInfo.FieldByName("Id").Uint())

	token := s.GenerateToken(googleInfo.FieldByName("Id").Uint(), time.Now().Add(time.Duration(jwtTTL)*time.Second))
	fmt.Println(googleInfo.FieldByName("Id"))

	return token
}
