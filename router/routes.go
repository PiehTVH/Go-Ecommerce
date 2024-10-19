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
