package main

import (
	"catalog-digital-product/internal/category"
	"catalog-digital-product/internal/config"
	"catalog-digital-product/internal/middleware"
	"catalog-digital-product/internal/store"
	"catalog-digital-product/internal/token"
	"catalog-digital-product/internal/user"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	dbConfig := config.DBConfig{
		User:     dbUser,
		Password: dbPassword,
		Host:     dbHost,
		Port:     dbPort,
		Name:     dbName,
	}

	db, err := config.InitDb(dbConfig)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api/v1")
	api.Static("/static/store", "./public/images/store")

	tokenService := token.NewTokenService([]byte(jwtSecretKey))

	userRepository := user.NewUserRepository()
	userService := user.NewUserService(db, userRepository, tokenService)
	userHandler := user.NewUserHandler(userService)

	categoryRepository := category.NewCategoryRepository()
	categoryService := category.NewCategoryService(categoryRepository, db)
	categoryHandler := category.NewCategoryHandler(categoryService)

	storeRepository := store.NewStoreRepository()
	storeService := store.NewStoreService(storeRepository, db)
	storeHandler := store.NewStoreHandler(storeService)

	api.POST("/auth/login", userHandler.Login)
	api.POST("/auth/password", middleware.AuthMiddleware(tokenService), userHandler.UpdatePassword)

	api.POST("/categories", middleware.AuthMiddleware(tokenService), categoryHandler.Create)
	api.GET("/categories", categoryHandler.GetAll)
	api.GET("/categories/:id", categoryHandler.Get)
	api.DELETE("/categories/:id", middleware.AuthMiddleware(tokenService), categoryHandler.Delete)
	api.PUT("/categories/:id", middleware.AuthMiddleware(tokenService), categoryHandler.Update)

	api.PUT("/store", middleware.AuthMiddleware(tokenService), storeHandler.Update)
	api.GET("/store", storeHandler.GetStore)

	if err := router.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
