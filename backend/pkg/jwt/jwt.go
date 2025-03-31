package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// MyClaims Custom Claims struct embedding jwt.StandardClaims
// The jwt package's jwt.StandardClaims only contains official fields.
// Here, we need to record an additional UserID field, so we define a custom struct.
// If more information needs to be stored, you can add them to this struct.
type MyClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// Define the Secret used for encryption
var mySecret = []byte("news")

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return mySecret, nil
}

// TokenExpireDuration Define the expiration time for JWT
const TokenExpireDuration = time.Hour * 24

// AccessTokenExpireDuration Define the expiration time for the access token
const AccessTokenExpireDuration = time.Hour * 24      // access_token expiration time
const RefreshTokenExpireDuration = time.Hour * 24 * 7 // refresh_token expiration time

// GenToken Generate JWT - creates an access_token and refresh_token
func GenToken(userID uint64, username string) (aToken, rToken string, err error) {
	// Create our custom claims
	c := MyClaims{
		userID,   // Custom field
		username, // Custom field
		jwt.StandardClaims{ // 7 official fields specified by JWT
			ExpiresAt: time.Now().Add(AccessTokenExpireDuration).Unix(), // Expiration time
			Issuer:    "news",                                           // Issuer
		},
	}
	// Encrypt and obtain the complete encoded token string
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)

	// The refresh token does not need to store any custom data
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(RefreshTokenExpireDuration).Unix(), // Expiration time
		Issuer:    "news",                                            // Issuer
	}).SignedString(mySecret)
	// Sign with the specified secret and obtain the complete encoded token string
	return
}

// ParseToken Parse JWT
func ParseToken(tokenString string) (claims *MyClaims, err error) {
	// Parse the token
	var token *jwt.Token
	claims = new(MyClaims)
	token, err = jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return
	}
	if !token.Valid { // Validate token
		err = errors.New("invalid token")
	}
	return
}

// RefreshToken Refresh the AccessToken
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// If refresh token is invalid, return immediately
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return
	}

	// Parse claims from the old access token (parse the payload data)
	var claims MyClaims
	_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)
	v, _ := err.(*jwt.ValidationError)

	// If the access token is expired and the refresh token is still valid, create a new access token
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.UserID, claims.Username)
	}
	return
}
