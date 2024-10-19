package router

import (
	"log"
	"net/http"
	"os"

	"github.com/PiehTVH/go-ecommerce/docs"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(*gin.Context)
}

type routes struct {
	router *gin.Engine
}

type Routes []Route

/*
*	Function for grouping user routes
 */
func (r routes) EcommerceUser(rg *gin.RouterGroup) {
	orderRouteGrouping := rg.Group("/ecommerce")
	orderRouteGrouping.Use(CORSMiddleware())
	for _, route := range userRoutes {
		switch route.Method {
		case "GET":
			orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			orderRouteGrouping.POST(route.Pattern, route.HandlerFunc)
		case "OPTIONS":
			orderRouteGrouping.OPTIONS(route.Pattern, route.HandlerFunc)
		case "PUT":
			orderRouteGrouping.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			orderRouteGrouping.DELETE(route.Pattern, route.HandlerFunc)
		default:
			orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(200, gin.H{

					"result": "Specify a valid http method with this route.",
				})
			})
		}

	}
}

func ClientRoutes() {
	r := routes{
		router: gin.Default(),
	}

	v1 := r.router.Group(os.Getenv("API_VERSION"))
	r.EcommerceUser(v1)

	// Swagger docs
	docs.SwaggerInfo.Title = "Elegance API"
	docs.SwaggerInfo.Description = "A robust and scalable backend system built using Go and the Gin framework, designed to support a comprehensive eCommerce platform."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "https://ethnicelegance.onrender.com"
	docs.SwaggerInfo.BasePath = "/v1/ecommerce"

	if err := r.router.Run(":" + os.Getenv("PORT")); err != nil {
		log.Printf("Failed to run server: %v", err)
	}

}

// Middlewares
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
		}
	}
}
