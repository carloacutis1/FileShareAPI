package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/routes"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/storage"
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
