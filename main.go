package main

import (
	"fmt"
	"log"
	"os"
	"project_absensi/handler"
	"project_absensi/handler/auth"
	"project_absensi/user"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/meilisearch/meilisearch-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
	}

	dbName := os.Getenv("DATABASE_NAME")
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")

	fmt.Println("Connecting to Meilisearch")
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://103.52.115.184:7700",
		APIKey: os.Getenv("MEILI_MASTER_KEY"),
	})

	fmt.Println("Connecting to Redis")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println("Connecting to the database...")
	fmt.Println(client, redisClient)

	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local&net_write_timeout=6000&multiStatements=true"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	err = db.AutoMigrate(&user.User{})
	if err != nil {
		log.Fatal("Error loading migrate database.")
	}

	fmt.Println("Database connected")

	fmt.Println("Initiating Repository & Services...")

	//REPOSITORY
	userRepository := user.NewRepository(db)

	//SERVICE
	userService := user.NewService(userRepository)
	authService := auth.NewService()

	//Handler
	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()
	api := router.Group("/api/v1")
	router.GET("/hallo", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	//ENPOINT
	api.POST("/user/register", userHandler.RegisterUser)
	api.POST("/user/login", userHandler.LoginUser)
	api.POST("/users", userHandler.GetAllUsers)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start the server!")
	}

	fmt.Println("Server is running on port 8080")
}
