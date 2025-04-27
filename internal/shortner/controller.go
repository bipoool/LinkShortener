package shortner

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

type CreateShortUrlRequestBody struct {
	OriginalUrl string `json:"original-url"`
}

func NewController(service *Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) getUrl(ctx *gin.Context) {
	shortCode := ctx.Param("shortCode")
	originalUrl, err := c.service.getOriginalUrl(shortCode)

	if err != nil {
		slog.Error("Error in getting url", "Error", err)
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if len(originalUrl) == 0 {
		slog.Error("URL not found", "ShortCode", shortCode)
		ctx.JSON(404, gin.H{"error": "Shortcode Not Found"})
		return
	}

	slog.Debug("Redirecting to", "URL", originalUrl, "Shortcode", shortCode)
	ctx.Redirect(http.StatusMovedPermanently, originalUrl)
}

func (c *Controller) createUrl(ctx *gin.Context) {
	var reqBody CreateShortUrlRequestBody
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		slog.Error("Error in parsing request body", "Error", err)
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	shortCode, err := c.service.createShortUrl(reqBody.OriginalUrl)
	if err != nil {
		slog.Error("Error in creating url", "Error", err)
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"short-url": "https://shtln.xyz/" + shortCode})
}
