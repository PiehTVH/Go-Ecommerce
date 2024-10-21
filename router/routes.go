package router

import (
	"net/http"

	"github.com/PiehTVH/go-ecommerce/constant"
	"github.com/PiehTVH/go-ecommerce/controller"
)

// health check service
var healthCheckRoutes = Routes{
	Route{"Health check", http.MethodGet, constant.HealthCheckRoute, controller.HealthCheck},
}

var userRoutes = Routes{

	// Register User
	Route{"Register User", http.MethodPost, constant.UserRegisterRoute, controller.RegisterUser},
	Route{"Login User", http.MethodPost, constant.UserLoginRoute, controller.UserLogin},
	Route{"Sign Out", http.MethodPost, constant.UserLogoutRoute, controller.SignOut},
}

var productGlobalRoutes = Routes{
	Route{"List Product", http.MethodGet, constant.ListProductRoute, controller.ListProductsController},
	Route{"Search Product", http.MethodPost, constant.SearchProductRoute, controller.SearchProductController},
	Route{"List Category", http.MethodGet, constant.ListCategoryRoute, controller.ListCategoryController},
	Route{"List Single Product", http.MethodGet, constant.ListSingleProductRoute, controller.ListSingleProductController},
}
