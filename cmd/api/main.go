package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories"
	submitjob "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/submit_job"
	"github.com/tapiaw38/practiq-be/internal/adapters/web"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	"github.com/tapiaw38/practiq-be/internal/platform/config"
	"github.com/tapiaw38/practiq-be/internal/platform/database"
	"github.com/tapiaw38/practiq-be/internal/platform/revocation"
	"github.com/tapiaw38/practiq-be/internal/platform/storage"
	"github.com/tapiaw38/practiq-be/internal/usecases"
)

// staleSubmitJobAfter is longer than a submission can legitimately take: the
// goroutine that runs one gives up at five minutes. Anything older than this
// is not slow, it is gone.
const staleSubmitJobAfter = 10 * time.Minute

// startSubmitJobSweeper closes submissions whose process is no longer running.
//
// The work lives in a goroutine and its payload nowhere else, so a restart
// leaves jobs marked processing that nothing will ever finish. The client
// polls them forever, and the student waits on a spinner instead of
// resubmitting. This runs once at boot — the restart case — and then on a
// timer, for a goroutine that died without saying so.
func startSubmitJobSweeper(repo submitjob.Repository) {
	sweep := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		closed, err := repo.FailStale(ctx, staleSubmitJobAfter)
		if err != nil {
			log.Printf("[submit_jobs] could not close interrupted submissions: %v", err)
			return
		}
		if closed > 0 {
			log.Printf("[submit_jobs] closed %d interrupted submission(s)", closed)
		}
	}

	sweep()
	go func() {
		ticker := time.NewTicker(staleSubmitJobAfter)
		defer ticker.Stop()
		for range ticker.C {
			sweep()
		}
	}()
}

func main() {
	loadConfig()

	cfg := config.GetConfigService()
	gin.SetMode(cfg.ServerConfig.GinMode)

	db, err := database.GetSQLClient()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	ds := datasources.CreateDatasources(db)
	reposFactory := repositories.NewFactory(ds)
	repos := reposFactory()
	integ := integrations.CreateIntegrations(
		cfg.ServerConfig.AuthAPIURL,
		cfg.ServerConfig.PaymentsURL,
		cfg.ServerConfig.PaymentsAPIKey,
	)
	imageStorage := storage.NewS3ImageStorage(cfg.S3Config)
	factory := appcontext.NewFactory(repos, integ, imageStorage)
	uc := usecases.NewUsecases(factory)

	app := gin.Default()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.ServerConfig.FrontendURL, "https://app.practiq.com.ar", "https://practiq.com.ar", "https://www.practiq.com.ar", "https://practiq-landing.onrender.com", "http://localhost:5174", "http://localhost:5173", "http://localhost:4321", "http://127.0.0.1:4321"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-School-ID"},
		AllowCredentials: true,
	}))

	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "practiq-be"})
	})

	revoked := revocation.NewChecker(integ.AuthAPI.GetTokenVersion, 60*time.Second)

	web.RegisterRoutes(app, uc, repos.SubmitJob, repos.UserProfile, repos.SiteContact, revoked)

	startSubmitJobSweeper(repos.SubmitJob)

	port := cfg.ServerConfig.Port
	log.Printf("practiq-be running on port %s", port)
	if err := app.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
