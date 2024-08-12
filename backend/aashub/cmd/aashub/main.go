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
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
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
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Initialize repositories
	verificationRepo := &repositories.VerificationRepository{DB: database}
	mailVerificationRepo := &repositories.EmailVerificationRepository{VerificationRepository: verificationRepo}
	userRepo := &repositories.UserRepository{DB: database, VerificationRepository: mailVerificationRepo}

	// Initialize handlers
	userHandler := &api.UserHandler{Repo: userRepo}
	verificationHandler := &api.VerificationHandler{VerificationRepository: verificationRepo}

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", Helloworld)
		}
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
	r.Run(":9000")
}
