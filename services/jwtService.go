package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT Service is a contract of what a JWT Service should be able to do.
type JWTService interface {
	// Generate a new token
	GenerateToken(userID string, t time.Time) string

	// Validate the token
	ValidateToken(token string) (*jwt.Token, error)

	// Refresh the token
	RefreshToken(token string) string
}

// jwtCustomClaim is a struct that contains the custom claims for the JWT
type jwtCustomClaim struct {
	// The userId is the only required field
	UserID string `json:"user_id"`

	// This is a registered JWT claim (StandardClaims are deprecated)
	jwt.RegisteredClaims
}

// jwtService is a struct that implements the JWTService interface
type jwtService struct {
	// Secret key used to sign the token
	secretKey string

	// Who creates the token
	issuer string
}

//NewJWTService method is creates a new instance of JWTService
func NewJWTService() JWTService {
	return &jwtService{
		// Call the getSecretKey function to get the secret key
		secretKey: getSecretKey(),

		// who creates the token
		issuer: "gojwt",
	}
}

// Create get the secret key from the environment variable
func getSecretKey() string {
	// Get the secret key from the environment variable
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		// If the environment variable is empty, use a default value
		secretKey = "learnGolangJWTToken"
	}
	return secretKey
}

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

// Create a new token object, specifying signing method and the claims
func (s *jwtService) GenerateToken(userID string, t time.Time) string {
	// Create the Claims struct with the required claims for the JWT
	claims := &jwtCustomClaim{
		// userId is the only required field
		userID,
		jwt.RegisteredClaims{
			// 1 day expiration
			ExpiresAt: jwt.NewNumericDate(t),
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
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
		return ""
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
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		// Return the key
		return []byte(s.secretKey), nil
	})
}

func (s *jwtService) RefreshToken(authHeader string) string {
	jwtTTL := GetTokenTTL()
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(authHeader, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return ""
	}

	userID := fmt.Sprintf("%d", claims["user_id"])
	// int to time.Duration
	token := s.GenerateToken(userID, time.Now().Add(time.Duration(jwtTTL)*time.Second))

	return token
}
