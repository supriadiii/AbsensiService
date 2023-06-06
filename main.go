package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"project_absensi/auth"
	"project_absensi/handler"
	"project_absensi/helper"
	"project_absensi/user"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
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
		Host:   "http://127.0.0.1:7700",
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

	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local&net_write_timeout=6000"
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
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Next()
	})

	api := router.Group("/api/v1")
	api.GET("/hallo", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	//ENPOINT
	api.POST("/user/register", userHandler.RegisterUser)
	api.POST("/users", autMiddleware(authService, userService), userHandler.GetAllUsers)
	api.POST("/user/login", userHandler.Login)
	api.POST("/user/nim", userHandler.CheckNimAvailability)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start the server!")
	}

	fmt.Println("Server is running on port 8080")
}

func autMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHandler := c.GetHeader(("Authorization"))
		if !strings.Contains(authHandler, "Baerer") {
			respone := helper.ResponseFormatter("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, respone)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHandler, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			respone := helper.ResponseFormatter("Unauthorized1", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, respone)
			return
		}
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			respone := helper.ResponseFormatter("Unauthorized2", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, respone)
			return
		}
		userID := uint(claim["user_id"].(float64))
		user, err := userService.GetUserByID(userID)
		if err != nil {
			respone := helper.ResponseFormatter("Unauthorized3", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, respone)
			return
		}
		c.Set("CurrentUser", user)
	}
}
