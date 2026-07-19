package main

import (
	"log"
	"net/http"
	"time"

	api "github.com/aas-hub-org/aashub/api/handler"
	"github.com/aas-hub-org/aashub/internal/database"
	repositories "github.com/aas-hub-org/aashub/internal/database/repositories"

	docs "github.com/aas-hub-org/aashub/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @BasePath /api/v1

// Health godoc
// @Summary Health Check
// @Description Responds with OK if the service is up and running
// @Tags health
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /health [get]
func Health(g *gin.Context) {
	g.JSON(http.StatusOK, "healthy")
}

func main() {
	r := gin.Default()

	// Configure CORS
	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	// Initialize database
	database, err := database.NewDB()
	if err != nil {
		log.Printf("Could not connect to the database: %v", err)
	}

	env_err := godotenv.Load("/workspace/backend/aashub/.env")
	if env_err != nil {
		log.Printf("Error loading .env file")
	}

	// Initialize repositories
	verificationRepo := &repositories.VerificationRepository{DB: database}
	mailVerificationRepo := &repositories.VerificationRepository{DB: database}
	userRepo := &repositories.UserRepository{DB: database, VerificationRepository: mailVerificationRepo}

	// Initialize handlers
	userHandler := &api.UserHandler{Repo: userRepo}
	verificationHandler := &api.VerificationHandler{VerificationRepository: verificationRepo}

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		ug := v1.Group("/users")
		{
			ug.POST("/register", gin.WrapF(userHandler.RegisterUser))
			ug.POST("/login", gin.WrapF(userHandler.LoginUser))
		}
		vg := v1.Group("/verify")
		{
			vg.GET("/", gin.WrapF(verificationHandler.VerifyUser))
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/health", Health)
	r.Run(":9000")
}
