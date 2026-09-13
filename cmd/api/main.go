package main

import (
	"context"
	"log"
	"restapi/internal/config"
	"restapi/internal/handler"
	"restapi/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	cfg := config.Load()

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	taskRepo := repository.NewTaskRepository(db)
	taskHandler := handler.NewTaskHandler(taskRepo)

	router.GET("/api/tasks", taskHandler.GetTasks)
	router.GET("/api/tasks/:id", taskHandler.GetTask)
	router.POST("/api/tasks", taskHandler.CreateTask)
	router.PATCH("/api/tasks/:id", taskHandler.UpdateTask)
	router.DELETE("api/tasks/:id", taskHandler.DeleteTask)

	router.Run(":" + cfg.Port)
}
