package main

import (
	"context"
	"embed"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/romakorinenko/task-manager/configs"
	_ "github.com/romakorinenko/task-manager/docs"
	"github.com/romakorinenko/task-manager/internal/controller"
	"github.com/romakorinenko/task-manager/internal/dbpool"
	"github.com/romakorinenko/task-manager/internal/repository"
	"github.com/romakorinenko/task-manager/internal/server"
	"github.com/romakorinenko/task-manager/internal/service"
)

// @title Task Manager API
// @version 1.0
// @description This is a sample API for Task Manager.
// @host localhost:8080
// @BasePath /

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	cfg := configs.MustLoadConfig()

	dbPool, err := dbpool.NewDBPool(context.Background(), cfg.DB)
	if err != nil {
		log.Fatalln("cannot create dbPool", err)
	}

	MigrateData(dbPool)

	userController := controller.NewUserController(
		service.NewUserService(repository.NewUserRepo(dbPool)),
		service.NewTaskService(repository.NewTaskRepo(dbPool)),
	)
	taskController := controller.NewTaskController(
		service.NewTaskService(repository.NewTaskRepo(dbPool)),
		service.NewUserService(repository.NewUserRepo(dbPool)),
	)
	server.RegisterServerAndHandlers(userController, taskController, cfg.Server.Port)
}

func MigrateData(dbPool *pgxpool.Pool) {
	db := stdlib.OpenDBFromPool(dbPool)
	err := db.Ping()
	if err != nil {
		log.Fatalln("cannot ping database", err)
	}

	goose.SetBaseFS(embedMigrations)
	dialectErr := goose.SetDialect("postgres")
	if dialectErr != nil {
		log.Fatalln("cannot set postgres dialect", dialectErr)
	}
	if migrationsErr := goose.Up(db, "migrations"); migrationsErr != nil {
		log.Fatalln("cannot migrate data", migrationsErr)
	}
}

// swagger - в первую очередь
// поправить мидлвейр - во вторую очередь
// разгрузить хендлеры переносом логики в сервис - в третью очередь
// транзакции - в четвертую очередь
// валидация - в пятую очередь
// логирование - в шестую очередь
// тесты - в седьмую очередь
// ci/cd - в восьмую очередь
// бейджи - в девятую очередь
