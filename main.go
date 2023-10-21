package main

import (
	"doodocs-challenge/internal/delivery"
	"doodocs-challenge/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	au := usecase.NewArchiveUsecase()
	delivery.NewArchiveHandler(router, au)
	router.Run("localhost:7169")
}
