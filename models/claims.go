package models

import (
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID string `json:"userID"` // Can incldue public information like user id
	Roles  string `json:"roles"`  // TODO: only have roels in teh database and retrive via the userid this would be better becuase else if a users roles update he has to grep a new jwt
	jwt.StandardClaims
}
