package api

import (
	"net/http"
	"time"

	"airclean-tracker/backend/internal/auth"
	"airclean-tracker/backend/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	cfg  config.Config
	db   *gorm.DB
	auth *auth.Middleware
}

func New(cfg config.Config, db *gorm.DB, authMW *auth.Middleware) *Server {
	return &Server{cfg: cfg, db: db, auth: authMW}
}

func (s *Server) Router() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     s.cfg.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	api := r.Group("/api")
	api.Use(s.auth.RequireAuth())
	api.GET("/me", s.me)
	api.GET("/dashboard", s.dashboard)
	api.GET("/aircons", s.listAircons)
	api.POST("/aircons", s.auth.RequireRole("admin", "team"), s.createAircon)
	api.GET("/aircons/:id", s.getAircon)
	api.PUT("/aircons/:id", s.auth.RequireRole("admin", "team"), s.updateAircon)
	api.DELETE("/aircons/:id", s.auth.RequireRole("admin"), s.deleteAircon)
	api.GET("/aircons/:id/cleaning-records", s.listCleaningRecords)
	api.POST("/aircons/:id/cleaning-records", s.auth.RequireRole("admin", "team"), s.createCleaningRecord)
	api.PUT("/cleaning-records/:id", s.auth.RequireRole("admin", "team"), s.updateCleaningRecord)
	api.DELETE("/cleaning-records/:id", s.auth.RequireRole("admin"), s.deleteCleaningRecord)
	api.GET("/plans", s.listPlans)
	api.POST("/plans/bulk-update", s.auth.RequireRole("admin", "team"), s.bulkUpdatePlans)
	api.POST("/import/excel", s.auth.RequireRole("admin"), s.importExcel)
	api.GET("/users", s.auth.RequireRole("admin"), s.listUsers)
	api.PUT("/users/:id/role", s.auth.RequireRole("admin"), s.updateUserRole)
	return r
}
