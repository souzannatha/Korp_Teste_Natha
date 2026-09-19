package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/infra/db"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/internal/controller"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/internal/repository"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/internal/routes"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/internal/usecases"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("could not load .env")
	}
	dbConnection, err := db.ConnectDB()
	if err != nil {
		return err
	}
	defer dbConnection.Close()
	server := gin.Default()
	if err := routes.ConfigureServer(server); err != nil {
		return err
	}
	productRepository := repository.NewProductRepository(dbConnection)
	productUseCase := usecase.NewProductUseCase(productRepository)
	productController := controller.NewProductController(productUseCase)
	stockRepository := repository.NewStockRepository(dbConnection)
	stockUseCase := usecase.NewStockUseCase(stockRepository)
	stockController := controller.NewStockController(stockUseCase)
	routes.RegisterProductRoutes(server, &productController)
	routes.RegisterStockRoutes(server, &stockController)
	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	routes.RegisterDocsRoutes(server)
	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           server,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	signalContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serverErrors := make(chan error, 1)
	go func() { serverErrors <- httpServer.ListenAndServe() }()
	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-signalContext.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownContext); err != nil {
			httpServer.Close()
			return err
		}
		return nil
	}
}
