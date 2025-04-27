package shortner

import (
	"linkshortener/internal/database"

	"github.com/gin-gonic/gin"
)

type Module struct {
	controller *Controller
	db         database.Database
	cache      database.Cache
}

func NewRepository(ginEngine *gin.Engine, db database.Database, cache database.Cache) *Module {
	service := NewService(db, cache)
	controller := NewController(service)
	module := &Module{
		controller: controller,
	}

	module.registerRoutes(ginEngine)

	return module
}

func (m *Module) registerRoutes(ginEngine *gin.Engine) {
	siteGroup := ginEngine.Group("/")
	{
		siteGroup.GET("/:shortCode", m.controller.getUrl)
		siteGroup.POST("/url", m.controller.createUrl)
	}
}
