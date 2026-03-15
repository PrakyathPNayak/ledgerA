package handler

import (
	"ledgerA/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// RouterDependencies groups all handlers required by router.
type RouterDependencies struct {
	AuthHandler             *AuthHandler
	UserHandler             *UserHandler
	AccountHandler          *AccountHandler
	CategoryHandler         *CategoryHandler
	TransactionHandler      *TransactionHandler
	QuickTransactionHandler *QuickTransactionHandler
	StatsHandler            *StatsHandler
	AuditHandler            *AuditHandler
	AllowedOrigins          []string
}

// SetupRouter creates gin engine and registers all routes.
func SetupRouter(deps RouterDependencies) *gin.Engine {
	router := gin.New()
	router.Use(middleware.LoggerMiddleware(), middleware.RecoveryMiddleware())

	config := cors.DefaultConfig()
	if len(deps.AllowedOrigins) > 0 {
		config.AllowOrigins = deps.AllowedOrigins
	} else {
		config.AllowAllOrigins = true
	}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"}
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/auth/sync", deps.AuthHandler.Sync)

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/users/me", deps.UserHandler.Me)
			protected.PATCH("/users/me", deps.UserHandler.UpdateMe)

			protected.GET("/accounts", deps.AccountHandler.List)
			protected.POST("/accounts", deps.AccountHandler.Create)
			protected.PATCH("/accounts/:id", deps.AccountHandler.Update)
			protected.DELETE("/accounts/:id", deps.AccountHandler.Delete)

			protected.GET("/categories", deps.CategoryHandler.List)
			protected.POST("/categories", deps.CategoryHandler.Create)
			protected.PATCH("/categories/:id", deps.CategoryHandler.Update)
			protected.DELETE("/categories/:id", deps.CategoryHandler.Delete)
			protected.POST("/categories/:id/subcategories", deps.CategoryHandler.CreateSubcategory)
			protected.PATCH("/subcategories/:id", deps.CategoryHandler.UpdateSubcategory)
			protected.DELETE("/subcategories/:id", deps.CategoryHandler.DeleteSubcategory)

			protected.GET("/transactions", deps.TransactionHandler.List)
			protected.POST("/transactions", deps.TransactionHandler.Create)
			protected.GET("/transactions/:id", deps.TransactionHandler.Get)
			protected.PATCH("/transactions/:id", deps.TransactionHandler.Update)
			protected.DELETE("/transactions/:id", deps.TransactionHandler.Delete)

			protected.GET("/quick-transactions", deps.QuickTransactionHandler.List)
			protected.POST("/quick-transactions", deps.QuickTransactionHandler.Create)
			protected.PATCH("/quick-transactions/:id", deps.QuickTransactionHandler.Update)
			protected.DELETE("/quick-transactions/:id", deps.QuickTransactionHandler.Delete)
			protected.POST("/quick-transactions/:id/execute", deps.QuickTransactionHandler.Execute)
			protected.PATCH("/quick-transactions/reorder", deps.QuickTransactionHandler.Reorder)

			protected.GET("/stats/summary", deps.StatsHandler.Summary)
			protected.GET("/stats/export/pdf", deps.StatsHandler.ExportPDF)
			protected.GET("/stats/compare", deps.StatsHandler.Compare)

			protected.GET("/audit", deps.AuditHandler.List)
		}
	}

	return router
}
