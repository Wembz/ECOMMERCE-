package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rodrigueghenda/Ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid code"})
			c.Abort()
			return
		}

		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		var addresses models.Address
		addresses.Address_ID = primitive.NewObjectID()

		if err = c.BindJSON(&addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		match_filter := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}

		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal server error")
			return
		}

		var addressinfo []bson.M
		if err = pointcursor.All(ctx, &addressinfo); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Error processing data")
			return
		}

		var size int32
		for _, address_no := range addressinfo {
			count := address_no["count"]
			size = count.(int32)
		}

		if size < 2 {
			filter := bson.D{{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{{Key: "address", Value: addresses}}}}

			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, "Error updating data")
				return
			}
		} else {
			c.IndentedJSON(http.StatusBadRequest, "Not Allowed")
			return
		}

		c.IndentedJSON(http.StatusOK, "Address added successfully")
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "address.0.house_name", Value: editaddress.House},
			{Key: "address.0.street_name", Value: editaddress.Street},
			{Key: "address.0.city_name", Value: editaddress.City},
			{Key: "address.0.pin_code", Value: editaddress.Pincode},
		}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong")
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully updated the home address")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "address.1.house_name", Value: editaddress.House},
			{Key: "address.1.street_name", Value: editaddress.Street},
			{Key: "address.1.city_name", Value: editaddress.City},
			{Key: "address.1.pin_code", Value: editaddress.Pincode},
		}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong")
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully updated the work address")
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Search Index"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		addresses := make([]models.Address, 0) // Empty slice to remove all addresses
		filter := bson.D{{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "address", Value: addresses}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Error deleting address")
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully deleted address")
	}
}
