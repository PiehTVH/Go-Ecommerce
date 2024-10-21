// package main

// import (
// 	"log"
// 	"os"

// 	"github.com/PiehTVH/go-ecommerce/database"
// 	"github.com/joho/godotenv"
// )

// func init() {
// 	if _, err := os.Stat(".env"); err == nil {
// 		log.Println("Loading .env file")
// 		err := godotenv.Load()

// 		if err != nil {
// 			log.Println("Error loading .env file")
// 		}
// 		log.Println("Loaded .env file successfully")
// 	}
// 	database.ConnectDB()
// }

// func main() {
// 	// run the routes
// 	// router.Clie
// }

package main

import (
	"fmt"
)

func hienCoXinhGaiKhong() {
	fmt.Println("Hien rat xinh gai")
}
func main() {
	slice := []string{99: ""}
	fmt.Println(len(slice))
}
