package routes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func ConfigureServer(server *gin.Engine) error {
	server.UseRawPath = true
	server.UnescapePathValues = true
	origin := os.Getenv("CORS_ORIGIN")
	if origin == "" {
		origin = "http://localhost:4200"
	}
	parsed, err := url.Parse(origin)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") ||
		parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" || parsed.User != nil {
		return fmt.Errorf("CORS_ORIGIN must be a valid HTTP origin")
	}
	if err := server.SetTrustedProxies(nil); err != nil {
		return err
	}
	server.Use(func(ctx *gin.Context) {
		requestContext, cancel := context.WithTimeout(ctx.Request.Context(), 25*time.Second)
		defer cancel()
		ctx.Request = ctx.Request.WithContext(requestContext)
		requestOrigin := ctx.GetHeader("Origin")
		ctx.Writer.Header().Add("Vary", "Origin")
		if requestOrigin == origin {
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			ctx.Header("Access-Control-Allow-Headers", "Content-Type")
		}
		if ctx.Request.Method == http.MethodOptions {
			if requestOrigin != origin {
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			}
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 1024*1024)
		ctx.Next()
	})
	return nil
}

func RegisterDocsRoutes(server *gin.Engine) {
	server.GET("/swagger", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")
	})
	server.GET("/swagger/index.html", func(ctx *gin.Context) {
		content, err := os.ReadFile("./docs/index.html")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not load Swagger UI"})
			return
		}
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", content)
	})
	server.StaticFile("/swagger/openapi.yaml", "./docs/openapi.yaml")
	server.Static("/swagger/assets", "./docs/assets")
}
