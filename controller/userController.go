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

func EditAddress(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(400, gin.H{
			"message": "Token is required",
		})
		return
	}

	// Verified user
	_, _, err := helper.VerifyToken(token)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	var editAddress types.AddressData

	defer c.Request.Body.Close()

	// binding the request body to address
	reqErr := c.ShouldBindJSON(&editAddress)
	if reqErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": reqErr.Error()})
		return
	}

	var userCollection *mongo.Collection = database.GetCollection(database.DB, constant.UsersCollection)

	var dbUser types.User

	// checking if user exists
	emailExists := userCollection.FindOne(c, bson.M{"email": editAddress.Email}).Decode(&dbUser)
	if emailExists != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "email not found"})
		return
	}

	// updating the address
	_, updateErr := userCollection.UpdateOne(c, bson.M{"email": editAddress.Email}, bson.M{"$set": bson.M{"address": editAddress.Address}})
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success"})
}

func UpdateUser(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(400, gin.H{
			"message": "Token is required",
		})
		return
	}

	// verified user
	_, _, err := helper.VerifyToken(token)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	var updatePassword types.UpdatePassword

	defer c.Request.Body.Close()

	// binding the request body to updatePassword
	reqErr := c.ShouldBindJSON(&updatePassword)
	if reqErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": reqErr.Error()})
		return
	}

	var userCollection *mongo.Collection = database.GetCollection(database.DB, constant.UsersCollection)
	var dbUser types.User

	// checking if user exists
	emailExists := userCollection.FindOne(c, bson.M{"email": updatePassword.Email}).Decode(&dbUser)
	if emailExists != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "email not found"})
		return
	}

	// checking the password
	if !helper.ComparePassword(dbUser.Password, updatePassword.OldPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "password not matched"})
		return
	}

	// updating the password
	_, updateErr := userCollection.UpdateOne(c, bson.M{"email": updatePassword.Email}, bson.M{"$set": bson.M{"password": helper.EncryptPassword(updatePassword.NewPassword)}})
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success"})
}

func EditName(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.JSON(400, gin.H{
			"message": "Token is required",
		})
		return
	}

	// verified user
	_, _, err := helper.VerifyToken(token)

	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	var editName types.NameData

	defer c.Request.Body.Close()

	// binding the request body to address
	reqErr := c.ShouldBindJSON(&editName)
	if reqErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": reqErr.Error()})
		return
	}

	var userCollection *mongo.Collection = database.GetCollection(database.DB, constant.UsersCollection)

	var dbUser types.User

	// checking if user exists
	emailExists := userCollection.FindOne(c, bson.M{"email": editName.Email}).Decode(&dbUser)

	if emailExists != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "email not found"})
		return
	}

	// updating the name
	_, updateErr := userCollection.UpdateOne(c, bson.M{"email": editName.Email}, bson.M{"$set": bson.M{"name": editName.Name}})

	if updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success"})
}

func AddToFavorite(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.JSON(400, gin.H{
			"message": "Token is required",
		})
		return
	}

	// verified user
	_, _, err := helper.VerifyToken(token)

	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	var req struct {
		Email     string `json:"email" bson:"email"`
		ProductId string `json:"productId" bson:"productId"`
	}

	defer c.Request.Body.Close()

	// binding the request body to address
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}

	var userCollection *mongo.Collection = database.GetCollection(database.DB, constant.UsersCollection)

	var dbUser types.User

	// checking if user exists
	emailExists := userCollection.FindOne(c, bson.M{"email": req.Email}).Decode(&dbUser)

	if emailExists != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "email not found"})
		return
	}

	// updating the favorite
	_, updateErr := userCollection.UpdateOne(c, bson.M{"email": req.Email}, bson.M{"$push": bson.M{"favourite": req.ProductId}})
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Added to Favorite"})
}

func RemoveFromFavorite(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.JSON(400, gin.H{
			"message": "Token is required",
		})
		return
	}

	// verified user
	_, _, err := helper.VerifyToken(token)

	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	var req struct {
		Email     string `json:"email" bson:"email"`
		ProductId string `json:"productId" bson:"productId"`
	}

	defer c.Request.Body.Close()

	// binding the request body to address
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}

	var userCollection *mongo.Collection = database.GetCollection(database.DB, constant.UsersCollection)

	var dbUser types.User

	// checking if user exists
	emailExists := userCollection.FindOne(c, bson.M{"email": req.Email}).Decode(&dbUser)

	if emailExists != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "email not found"})
		return
	}

	// updating the favourite
	_, updateErr := userCollection.UpdateOne(c, bson.M{"email": req.Email}, bson.M{"$pull": bson.M{"favourite": req.ProductId}})
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Removed from Favorite"})
}

func ListFavorite(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.JSON(400, gin.H{
			"message": "Token is required",
		})
		return
	}

	// verified user
	_, _, err := helper.VerifyToken(token)

	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	var req struct {
		Email string `json:"email" bson:"email"`
	}

	defer c.Request.Body.Close()

	// binding the request body to address
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}

	var userCollection *mongo.Collection = database.GetCollection(database.DB, constant.UsersCollection)

	var dbUser types.User

	// checking if user exists
	emailExists := userCollection.FindOne(c, bson.M{"email": req.Email}).Decode(&dbUser)

	if emailExists != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "email not found"})
		return
	}

	var allFav []types.Product
	var productCollection *mongo.Collection = database.GetCollection(database.DB, constant.ProductCollection)
	for _, v := range dbUser.Favourite {
		var singleProduct types.Product
		productCollection.FindOne(c, bson.M{"id": v}).Decode(&singleProduct)
		allFav = append(allFav, singleProduct)
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success", "data": allFav})
}

func AddToCart(c *gin.Context) {
	var addToCart types.AddToCart
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.JSON(400, gin.H{
			"message": "Token is required",
		})
		return
	}

	defer c.Request.Body.Close()

	// binding the request body to address
	if err := c.ShouldBindJSON(&addToCart); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}

	var cartCollection *mongo.Collection = database.GetCollection(database.DB, constant.CartItemCollection)
	var productCollection *mongo.Collection = database.GetCollection(database.DB, constant.ProductCollection)

	var dbCart types.CartItem

	// if user already has products on the cart
	emailExists := cartCollection.FindOne(c, bson.M{"email": addToCart.Email}).Decode(&dbCart)
	if emailExists != nil {
		var product types.Product
		err := productCollection.FindOne(c, bson.M{"id": addToCart.ProductID}).Decode(&product)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "product not found"})
			return
		}

		// creating the cart object
		dbCart = types.CartItem{
			Email:     addToCart.Email,
			NumItems:  1,
			ChekedOut: false,
			Total:     float64(product.Price * addToCart.Quantity),
			Products: []types.ProductInCart{
				{
					ProductID: addToCart.ProductID,
					Quantity:  addToCart.Quantity,
				},
			},
		}

		_, insertErr := cartCollection.InsertOne(c, dbCart)
		if insertErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": insertErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"error": false, "message": "success"})
		return
	}

	// if user does not have products on the cart
	var product types.Product
	err := productCollection.FindOne(c, bson.M{"id": addToCart.ProductID}).Decode(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "product not found"})
		return
	}

	// getting the cart for the user
	var cart types.CartItem
	cartCollection.FindOne(c, bson.M{"email": addToCart.Email}).Decode(&cart)

	dbCart.Total = cart.Total
	dbCart = types.CartItem{
		Email:     addToCart.Email,
		NumItems:  dbCart.NumItems + 1,
		ChekedOut: false,
		Total:     dbCart.Total + float64(product.Price*addToCart.Quantity),
		Products: append(dbCart.Products, types.ProductInCart{
			ProductID: addToCart.ProductID,
			Quantity:  addToCart.Quantity,
		}),
	}

	_, updateErr := cartCollection.UpdateOne(c, bson.M{"email": addToCart.Email}, bson.M{"$set": dbCart})
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success"})
}

func RemoveFromCart(c *gin.Context) {
	var addToCart struct {
		Email     string `json:"email" bson:"email"`
		ProductId string `json:"productId" bson:"productId"`
	}
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(400, gin.H{
			"message": "Token is required",
		})
		return
	}

	defer c.Request.Body.Close()

	// binding the request body to address
	if err := c.ShouldBindJSON(&addToCart); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}

	var cartCollection *mongo.Collection = database.GetCollection(database.DB, constant.CartItemCollection)
	var productCollection *mongo.Collection = database.GetCollection(database.DB, constant.ProductCollection)

	var dbCart types.CartItem
	var product types.Product

	// get the product
	err := productCollection.FindOne(c, bson.M{"id": addToCart.ProductId}).Decode(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "product not found"})
		return
	}

	// if user already has products on the cart
	emailExists := cartCollection.FindOne(c, bson.M{"email": addToCart.Email}).Decode(&dbCart)
	if emailExists != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "cart not found"})
		return
	}

	var numberofItems int

	// total minus for removed product
	for _, v := range dbCart.Products {
		if v.ProductID == addToCart.ProductId {
			numberofItems = v.Quantity
		}
	}

	dbCart.NumItems = dbCart.NumItems - 1
	dbCart.Total = dbCart.Total - float64(product.Price*numberofItems)

	_, updateErr := cartCollection.UpdateOne(
		c,
		bson.M{"email": addToCart.Email},
		bson.M{"$pull": bson.M{"products": bson.M{"product_id": addToCart.ProductId}},
			"$set": bson.M{"total": dbCart.Total}})
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success"})
}
