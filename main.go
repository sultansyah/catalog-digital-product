package main

import (
	"catalog-digital-product/internal/category"
	"catalog-digital-product/internal/config"
	"catalog-digital-product/internal/middleware"
	"catalog-digital-product/internal/product"
	"catalog-digital-product/internal/store"
	"catalog-digital-product/internal/token"
	"catalog-digital-product/internal/user"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	router.Static("/static/images", "./public/images")

	api := router.Group("/api/v1")

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

	productRepository := product.NewProductRepository()
	productService := product.NewProductService(productRepository, db)
	productHandler := product.NewProductHandler(productService)

	templatesDir := "./internal/templates"
	router.LoadHTMLFiles(
		filepath.Join(templatesDir, "index.html"),
		filepath.Join(templatesDir, "product.html"),
		filepath.Join(templatesDir, "about.html"),
		filepath.Join(templatesDir, "home.html"),
		filepath.Join(templatesDir, "layouts_home", "header.html"),
		filepath.Join(templatesDir, "layouts_home", "footer.html"),
	)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":    "Katalog Digital Produk",
			"htmlfile": "home.html",
		})
	})

	router.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":    "Tentang Kami - Katalog Digital Produk",
			"htmlfile": "about.html",
		})
	})

	router.GET("/product/:slug", func(c *gin.Context) {
		var input product.SlugProductInput

		if err := c.ShouldBindUri(&input); err != nil {
			input.Slug = ""
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":    "Tentang Kami - Katalog Digital Produk",
			"id":       input.Slug,
			"htmlfile": "product.html",
		})
	})

	api.POST("/auth/login", userHandler.Login)
	api.POST("/auth/password", middleware.AuthMiddleware(tokenService), userHandler.UpdatePassword)

	api.POST("/categories", middleware.AuthMiddleware(tokenService), categoryHandler.Create)
	api.GET("/categories", categoryHandler.GetAll)
	api.GET("/categories/:id", categoryHandler.Get)
	api.DELETE("/categories/:id", middleware.AuthMiddleware(tokenService), categoryHandler.Delete)
	api.PUT("/categories/:id", middleware.AuthMiddleware(tokenService), categoryHandler.Update)

	api.PUT("/store", middleware.AuthMiddleware(tokenService), storeHandler.Update)
	api.GET("/store", storeHandler.GetStore)

	api.GET("/products", productHandler.GetAll)
	api.GET("/products/id/:id", productHandler.Get)
	api.GET("/products/slug/:slug", productHandler.GetBySlug)
	api.POST("/products", middleware.AuthMiddleware(tokenService), productHandler.Insert)
	api.PUT("/products/:id", middleware.AuthMiddleware(tokenService), productHandler.Update)
	api.DELETE("/products/:id", middleware.AuthMiddleware(tokenService), productHandler.Delete)

	api.POST("/products/:id/images", middleware.AuthMiddleware(tokenService), productHandler.InsertImage)
	api.PUT("/products/:id/images/:imageId", middleware.AuthMiddleware(tokenService), productHandler.SetLogoImage)
	api.DELETE("/products/:id/images/:imageId", middleware.AuthMiddleware(tokenService), productHandler.DeleteImage)

	if err := router.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
