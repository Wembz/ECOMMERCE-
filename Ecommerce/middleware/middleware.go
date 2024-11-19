package middleware

import (
	"net/http"

 	"github.com/gin-gonic/gin"
	token "github.com/rodrigueghenda/Ecommerce/tokens"
)

//Checking if the token is valid/present or not
func Authentication()gin.HandlerFunc{
	return func(c *gin.Context){
		ClientToken := c.Request.Header.Get("token")

		//if user client token is empty we send error message
		if ClientToken == ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No autherization header provided"})
			c.Abort()
			return
		}
		claims, err := token.ValidateToken(ClientToken)
		// if err is not empty we print out error message 
		if err != ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next() 
	}
}