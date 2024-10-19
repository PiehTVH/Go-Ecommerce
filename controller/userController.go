package controller

import (
	"net/http"
	"time"

	"github.com/PiehTVH/go-ecommerce/constant"
	"github.com/PiehTVH/go-ecommerce/database"
	"github.com/PiehTVH/go-ecommerce/helper"
	"github.com/PiehTVH/go-ecommerce/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type JwtClaim struct {
	UserId   primitive.ObjectID
	Email    string
	UserType string
	jwt.StandardClaims
}

// @Summary		User Signup
// @Description	user can signup by giving their details
// @Tags			User
// @Accept			json
// @Produce		    json
// @Success		200	{object}	string
// @Failure		500	{object}	string
// @Router			/v1/ecommerce/signup [post]
func RegisterUser(c *gin.Context) {
	var userClient types.UserClient
	var dbUser types.User

	defer c.Request.Body.Close()
	// binding the request body to userClient
	reqErr := c.ShouldBindJSON(&userClient)
	if reqErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": reqErr.Error()})
	}

	// checking the payload
	err := helper.CheckUserValidation(userClient)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	var userCollection *mongo.Collection = database.GetCollection(database.DB, constant.UsersCollection)

	// checking if email is unique
	emailExists := userCollection.FindOne(c, bson.M{"email": userClient.Email}).Decode(&dbUser)
	if emailExists == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "email already exists"})
		return
	}

	// creating the user object
	dbUser = types.User{
		Name:      userClient.Name,
		Email:     userClient.Email,
		Phone:     userClient.Phone,
		Password:  helper.EncryptPassword(userClient.Password),
		UserType:  "user",
		IsBlocked: false,
		Address:   "",
		Favourite: []string{},
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		Id:        primitive.NewObjectID(),
	}

	_, insertErr := userCollection.InsertOne(c, dbUser)
	if insertErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": insertErr.Error()})
		return
	}

	// jwt token
	token, err := helper.GenerateToken(dbUser.Id.Hex(), dbUser.Email, dbUser.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}

	dbUser.Password = ""

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success", "data": dbUser, "token": token})
}

// Login
func UserLogin(c *gin.Context) {
	var loginReq types.Login

	defer c.Request.Body.Close()

	// binding the request body to lognReq
	reqErr := c.ShouldBindJSON(&loginReq)
	if reqErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": reqErr.Error()})
		return
	}

	var userCollection *mongo.Collection = database.GetCollection(database.DB, constant.UsersCollection)
	var dbUser types.User

	// checking if email exists
	emailExists := userCollection.FindOne(c, bson.M{"email": loginReq.Email}).Decode(&dbUser)
	if emailExists != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "email not found"})
		return
	}

	// checking the password
	if !helper.ComparePassword(dbUser.Password, loginReq.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "password not matched"})
		return
	}

	// jwt token
	token, err := helper.GenerateToken(dbUser.Id.Hex(), dbUser.Email, dbUser.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}

	dbUser.Password = ""

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success", "data": dbUser, "token": token})
}

func SignOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success"})
}

func AddAddress(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(400, gin.H{
			"message": "Token is required",
		})
		return
	}

	// verified token
	_, _, err := helper.VerifyToken(token)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	var addAddress types.AddressData
	defer c.Request.Body.Close()

	// binding the request body to address
	reqErr := c.ShouldBindJSON(&addAddress)
	if reqErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": reqErr.Error()})
		return
	}

	var userCollection *mongo.Collection = database.GetCollection(database.DB, constant.UsersCollection)
	var dbUser types.User

	// checking if user exists
	emailExists := userCollection.FindOne(c, bson.M{"email": addAddress.Email}).Decode(&dbUser)
	if emailExists != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "email not found"})
		return
	}

	// updating the address
	_, updateErr := userCollection.UpdateOne(c,
		bson.M{"email": addAddress.Email},
		bson.M{"$set": bson.M{"address": addAddress.Address}})
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success"})
}
