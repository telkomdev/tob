package utils

import (
	"errors"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v5"
)

// Alg represent jwt algorithm
type Alg string

const (
	// HS256 const
	HS256 Alg = "HS256"

	// RS256 const
	RS256 Alg = "RS256"
)

// Claim model
type Claim struct {
	Issuer    string
	Audience  string
	Subject   string
	ExpiredAt int64
	IssuedAt  int64
	User      struct {
		ID       string
		FullName string
		Email    string
	}
	Alg Alg
}

// JwtService represent jwt service
type JwtService interface {
	Generate(payload *Claim, expired time.Duration) (string, error)
	Validate(alg Alg, tokenString string) (*Claim, error)
}

// JWT implementation from JwtService
type JWT struct {
	hmacKey string
}

// NewJWT constructor
func NewJWT(hmacKey string) *JWT {
	return &JWT{
		hmacKey: hmacKey,
	}
}

// Generate token
func (r *JWT) Generate(payload *Claim, expired time.Duration) (string, error) {
	now := time.Now()
	exp := now.Add(expired)

	var key interface{}
	var token = new(jwtgo.Token)

	token = jwtgo.New(jwtgo.SigningMethodHS256)
	key = []byte(r.hmacKey)

	claims := jwtgo.MapClaims{
		"iss": "tob",
		"exp": exp.Unix(),
		"iat": now.Unix(),
		"sub": payload.Subject,
		"aud": "97b33193-43ff-4e58-9124-b3a9b9f72c34",
	}

	if payload.User.ID != "" {
		claims["id"] = payload.User.ID
	}

	if payload.User.Email != "" {
		claims["email"] = payload.User.Email
	}

	if payload.User.FullName != "" {
		claims["fullName"] = payload.User.FullName
	}

	token.Claims = claims

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err

	}

	return tokenString, nil
}

// Validate token
func (r *JWT) Validate(alg Alg, tokenString string) (*Claim, error) {
	tokenParse, err := jwtgo.Parse(tokenString, func(token *jwtgo.Token) (interface{}, error) {
		return []byte(r.hmacKey), nil
	})

	var errToken error
	switch {
	case errors.Is(err, jwtgo.ErrTokenMalformed):
		errToken = errors.New("error jwt: token in invalid format")
	case errors.Is(err, jwtgo.ErrTokenExpired) || errors.Is(err, jwtgo.ErrTokenNotValidYet):
		errToken = errors.New("error jwt: token expired or not active yet")
	default:
		errToken = err
	}

	if errToken != nil {
		return nil, errToken
	}

	if !tokenParse.Valid {
		return nil, errors.New("jwt error")
	}

	mapClaims, ok := tokenParse.Claims.(jwtgo.MapClaims)
	if !ok {
		return nil, errors.New("jwt error: cannot parse token to map claims")
	}

	var tokenClaim Claim
	tokenClaim.Issuer, _ = mapClaims["iss"].(string)
	tokenClaim.Audience, _ = mapClaims["aud"].(string)
	tokenClaim.IssuedAt, _ = mapClaims["iat"].(int64)
	tokenClaim.ExpiredAt, _ = mapClaims["exp"].(int64)
	tokenClaim.Subject, _ = mapClaims["sub"].(string)
	tokenClaim.User.ID, _ = mapClaims["id"].(string)
	tokenClaim.User.Email, _ = mapClaims["email"].(string)
	tokenClaim.User.FullName, _ = mapClaims["fullName"].(string)

	return &tokenClaim, nil
}
