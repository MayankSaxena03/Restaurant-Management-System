package helpers

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/MayankSaxena03/Restaurant-Management-System/database"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userCollection = database.OpenCollection(database.Client, "users")
var SECRETKEY = os.Getenv("SECRETKEY")

type SignedDetails struct {
	UserID   primitive.ObjectID
	Username string
	Email    string
	jwt.StandardClaims
}

func GenerateAllTokens(userID primitive.ObjectID, email string, username string) (string, string, error) {
	claims := &SignedDetails{
		UserID:   userID,
		Username: username,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		UserID:   userID,
		Username: username,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRETKEY))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRETKEY))
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func UpdateAllTokens(ctx context.Context, userID primitive.ObjectID, token string, refreshToken string) error {
	query := bson.M{
		"_id": userID,
	}
	update := bson.M{
		"$set": bson.M{
			"token":        token,
			"refreshToken": refreshToken,
			"updatedOn":    time.Now(),
		},
	}
	_, err := userCollection.UpdateOne(ctx, query, update)
	return err
}

func ValidateToken(signedToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRETKEY), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, errors.New("Token is invalid")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("Token has expired")
	}

	return claims, nil
}
