package controller

import (
	"context"
	"os"

	"github.com/PiehTVH/go-ecommerce/constant"
	"github.com/PiehTVH/go-ecommerce/database"
	"github.com/PiehTVH/go-ecommerce/helper"
	"github.com/PiehTVH/go-ecommerce/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary List all products
// @Description List all products
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object}  string
// @Router /v1/ecommerce/products [get]
func ListProductsController(c *gin.Context) {

	var productCollection *mongo.Collection = database.GetCollection(database.DB, constant.ProductCollection)

	var products []types.Product

	cursor, err := productCollection.Find(context.Background(), bson.D{})

	if err != nil {
		c.JSON(400, gin.H{
			"message": constant.BadRequestMessage,
		})
		return
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var product types.Product
		cursor.Decode(&product)
		if product.Stock > 0 {
			products = append(products, product)
		}
	}

	c.JSON(200, gin.H{
		"products": products,
	})
}

// @Summary List all categories
// @Description List all categories
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object}  string
// @Router /v1/ecommerce/list-category [get]
func ListCategoryController(c *gin.Context) {
	var categoryCollection *mongo.Collection = database.GetCollection(database.DB, constant.CategoryCollection)

	var categories []types.Category

	cursor, err := categoryCollection.Find(context.Background(), bson.D{})

	if err != nil {
		c.JSON(400, gin.H{
			"message": constant.BadRequestMessage,
		})
		return
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var category types.Category
		cursor.Decode(&category)
		categories = append(categories, category)
	}

	c.JSON(200, gin.H{
		"categories": categories,
	})
}

// @Summary List single product by id
// @Description List single product by id
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object}  string
// @Router /v1/ecommerce/product/{id} [get]
func ListSingleProductController(c *gin.Context) {
	Id := c.Param("id")

	var productCollection *mongo.Collection = database.GetCollection(database.DB, constant.ProductCollection)

	var product types.Product

	err := productCollection.FindOne(context.Background(), bson.D{{"id", Id}}).Decode(&product)

	if err != nil {
		c.JSON(400, gin.H{
			"message": constant.BadRequestMessage,
		})
		return
	}

	c.JSON(200, gin.H{
		"product": product,
	})
}

// @Summary Get product link
// @Description Get product link
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object}  string
// @Router /v1/ecommerce/product-link/:id [get]
func GetProductLink(c *gin.Context) {
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

	Id := c.Param("id")
	frontend := os.Getenv("frontEndUrl")
	link := frontend + "/products/" + Id

	c.JSON(200, gin.H{
		"link": link,
	})

}

// @Summary Give rating
// @Description Give rating
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Token"
// @Param id path string true "Product ID"
// @Param rating body types.Rating true "Rating"
// @Success 200 {object}  string
// @Router /v1/ecommerce/rating/{id} [post]
func GiveRating(c *gin.Context) {
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

	Id := c.Param("id")

	var rating types.Rating

	err = c.ShouldBindJSON(&rating)

	if err != nil {
		c.JSON(400, gin.H{
			"message": constant.BadRequestMessage,
		})
		return
	}

	var productCollection *mongo.Collection = database.GetCollection(database.DB, constant.ProductCollection)

	var product types.Product

	err = productCollection.FindOne(context.Background(), bson.D{{"id", Id}}).Decode(&product)

	if err != nil {
		c.JSON(400, gin.H{
			"message": constant.BadRequestMessage,
		})
		return
	}

	product.Rating = (product.Rating*float64(product.NumRating) + rating.Rating) / float64(product.NumRating+1)
	product.NumRating = product.NumRating + 1

	_, updateErr := productCollection.UpdateOne(c, bson.M{"id": Id}, bson.M{"$set": bson.M{"rating": product.Rating, "numrating": product.NumRating}})
	if updateErr != nil {
		c.JSON(400, gin.H{
			"message": constant.BadRequestMessage,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Rating added",
	})
}

// @Summary Comment on product
// @Description Comment on product
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Token"
// @Param id path string true "Product ID"
// @Success 200 {object}  string
// @Router /v1/ecommerce/comment/{id} [post]
func CommentOnProduct(c *gin.Context) {
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

	Id := c.Param("id")

	var comment types.Comment

	err = c.ShouldBindJSON(&comment)

	if err != nil {
		c.JSON(400, gin.H{
			"message": constant.BadRequestMessage,
		})
		return
	}

	var productCollection *mongo.Collection = database.GetCollection(database.DB, constant.ProductCollection)

	var product types.Product

	err = productCollection.FindOne(context.Background(), bson.D{{"id", Id}}).Decode(&product)

	if err != nil {
		c.JSON(400, gin.H{
			"message": constant.BadRequestMessage,
		})
		return
	}

	output := append(product.Comments, comment)

	_, updateErr := productCollection.UpdateOne(c, bson.M{"id": Id}, bson.M{"$set": bson.M{"comments": output}})
	if updateErr != nil {
		c.JSON(400, gin.H{
			"message": constant.BadRequestMessage,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Comment added",
	})
}

// @Summary Search product
// @Description Search product
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Token"
// @Param search body string true "Search"
// @Param limit body int true "Limit"
// @Param page body int true "Page"
// @Param offset body int true "Offset"
// @Success 200 {object}  string
// @Router /v1/ecommerce/search [post]
func SearchProductController(c *gin.Context) {
	var reqSearch struct {
		Search string `json:"search"`
		Limit  int    `json:"limit"`
		Page   int    `json:"page"`
		Offset int    `json:"offset"`
	}

	err := c.ShouldBindJSON(&reqSearch)
	if err != nil {
		c.JSON(400, gin.H{
			"message": constant.BadRequestMessage,
		})
		return
	}

	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(400, gin.H{
			"message": "Token is required",
		})
		return
	}
	// verified user
	_, _, errToken := helper.VerifyToken(token)
	if errToken != nil {
		c.JSON(400, gin.H{
			"message": errToken.Error(),
		})
		return
	}

	skip := reqSearch.Limit * (reqSearch.Page - 1)
	if reqSearch.Offset > 0 {
		skip = reqSearch.Offset
	}

	var productCollection *mongo.Collection = database.GetCollection(database.DB, constant.ProductCollection)

	var products []types.Product

	findOptions := options.Find()
	findOptions.SetLimit(int64(reqSearch.Limit))
	findOptions.SetSkip(int64(skip))

	searchFilter := bson.M{}
	if len(reqSearch.Search) > 3 {
		searchFilter["$or"] = []bson.M{
			{"name": bson.M{"$regex": reqSearch.Search, "$options": "i"}},
			{"description": bson.M{"$regex": reqSearch.Search, "$options": "i"}},
		}
	}

	cursor, err := productCollection.Find(context.Background(), searchFilter, findOptions)

	if err != nil {
		c.JSON(400, gin.H{
			"message": constant.BadRequestMessage,
		})
		return
	}

	count, err := productCollection.CountDocuments(context.Background(), searchFilter)

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var product types.Product
		cursor.Decode(&product)
		if product.Stock > 0 {
			products = append(products, product)
		}
	}

	c.JSON(200, gin.H{
		"products": products,
		"total":    count,
	})

}
