package main

import (
	"context"
	"log"

	"github.com/carloacutis1/FileShareAPI/internal/database"
	"github.com/carloacutis1/FileShareAPI/internal/routes"
	"github.com/carloacutis1/FileShareAPI/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()

	blobStore, cleanup, err := storage.NewBlobStoreFromEnv(context.Background())
	if err != nil {
		log.Fatalf("blob storage: %v", err)
	}
	defer cleanup()

	router := gin.Default()
	routes.SetupRoutes(router, blobStore)

	log.Fatal(router.Run(":8080"))
}
