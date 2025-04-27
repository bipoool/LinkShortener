package server

import (
	"context"
	"fmt"
	"linkshortener/internal/config"
	"linkshortener/internal/database"
	"linkshortener/internal/logger"
	"linkshortener/internal/middleware"
	"linkshortener/internal/shortner"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LinkShortenerServer struct {
	ginEngine  *gin.Engine
	httpServer *http.Server
	db         database.Database
	cache      database.Cache
}

func Init() {

	// InitLink Shortener API server
	server := InitialiseLinkShortenerServerAndDependencies()

	// Initialize global manager with common dependencies
	shortner.NewRepository(server.ginEngine, server.db, server.cache)

	server.runServer()

	// server.shutdown()
}

func InitialiseLinkShortenerServerAndDependencies() *LinkShortenerServer {
	slog.SetDefault(logger.New())
	// Initialize router, Middleware & HttpServer
	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(middleware.SlogMiddleware())

	httpServer := &http.Server{
		Addr:    ":" + config.Config.Server.Port,
		Handler: ginEngine,
	}

	// Initialize DB
	db := database.NewSqlDatabase(
		config.Config.Sql.Host,
		config.Config.Sql.Port,
		config.Config.Sql.User,
		config.Config.Sql.Password,
		config.Config.Sql.Db,
		config.Config.Sql.SslMode,
	)

	// Initialize cache
	cache := database.NewCache(config.Config.Cache.Host, config.Config.Cache.Port)

	server := &LinkShortenerServer{
		ginEngine:  ginEngine,
		httpServer: httpServer,
		db:         db,
		cache:      cache,
	}

	// Print Banner
	server.printBanner()

	// Print configuration
	server.printConfiguration()

	return server
}

func (server *LinkShortenerServer) runServer() {
	// Start server
	if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Listen", "Error", err.Error())
	}
}

func (server *LinkShortenerServer) printBanner() {
	fmt.Println(`
██╗     ██╗███╗   ██╗██╗  ██╗    ███████╗██╗  ██╗ ██████╗ ██████╗ ████████╗███████╗███╗   ██╗███████╗██████╗ 
██║     ██║████╗  ██║██║ ██╔╝    ██╔════╝██║  ██║██╔═══██╗██╔══██╗╚══██╔══╝██╔════╝████╗  ██║██╔════╝██╔══██╗
██║     ██║██╔██╗ ██║█████╔╝     ███████╗███████║██║   ██║██████╔╝   ██║   █████╗  ██╔██╗ ██║█████╗  ██████╔╝
██║     ██║██║╚██╗██║██╔═██╗     ╚════██║██╔══██║██║   ██║██╔══██╗   ██║   ██╔══╝  ██║╚██╗██║██╔══╝  ██╔══██╗
███████╗██║██║ ╚████║██║  ██╗    ███████║██║  ██║╚██████╔╝██║  ██║   ██║   ███████╗██║ ╚████║███████╗██║  ██║
╚══════╝╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝    ╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚══════╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝
	`)
}

func (server *LinkShortenerServer) printConfiguration() {
	// TODO update version from version file
	slog.Info("Starting Link Shortener", slog.String("version", "1.0.0"))
	slog.Info("Sql", slog.String("DB", config.Config.Sql.Db), slog.String("User", config.Config.Sql.User))
	slog.Info("On", slog.String("Port", config.Config.Server.Port))
	slog.Info("Log", slog.String("Level", config.Config.Server.LogLevel))
}

func (server *LinkShortenerServer) shutdown() {
	slog.Info("Shutting down the server")
	// Give server 10s to gracefully shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "Error", err)
	}
	server.db.Close()
}
