package token

import (
	"context"
	"log"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/rodrigueghenda/Ecommerce/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email string
	First_Name string
	Last_Name string
	Uid string
	jwt.StandardClaims
}

//User collection 
var UserData *mongo.Collection = database.UserData(database.Client, "Users")

var SECRET_KEY = os.Getenv("SECRET_KEY")
func TokenGenerator(email string, firstname string, lastname string, uid string)(signedtoken string, signedrefreshtoken string, err error){

	claims := &SignedDetails{
		Email: email,
		First_Name: firstname,
		Last_Name: lastname,
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour *time.Duration(24)).Unix(), 
		},
	}

	refreshclaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour *time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, refreshclaims).SignedString([]byte(SECRET_KEY))
	
	//Error handling
	if err != nil {
	log.Panic(err)
	return
	}

	//return value if line of code work
	return token, refreshtoken, err 
}

//
func ValidateToken(signedtoken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedtoken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return nil, msg // Return nil for claims in case of error
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid" // Fix typo
		return nil, msg
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is already expired"
		return nil, msg
	}
	return claims, msg
}

//creating a Mongodb query 
func UpdateAllTokens(signedtoken string, signedrefreshtoken string, userid string) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateobj primitive.D

	updateobj = append(updateobj, bson.E{Key: "token", Value: signedtoken})
	updateobj = append(updateobj, bson.E{Key: "refresh_token", Value: signedrefreshtoken})
	updateobj = append(updateobj, bson.E{Key: "updatedat", Value: time.Now()}) // Use time.Now()

	upsert := true

	filter := bson.M{"user_id": userid}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := UserData.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updateobj},
	}, &opt)

	if err != nil {
		log.Println("Error updating tokens:", err) // Handle error gracefully
		return
	}
}