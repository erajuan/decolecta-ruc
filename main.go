package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

var CTX = context.Background()

var redisClient *redis.Client

var pgsqlClient *pgxpool.Pool

type counter struct{ r *redis.Client }

type repo struct{ db *pgxpool.Pool }

func main() {
	Init()
	app := fiber.New(fiber.Config{
		ServerHeader:  "api.decolecta.com",
		CaseSensitive: true,
		Immutable:     true,
		Prefork:       false,
		ReadTimeout:   time.Second * 30,
	})

	app.Use(logger.New())
	app.Get("/", Home)
	app.Get("/ruc", GetCompanyHandler)
	app.Get("/ruc/full", GetCompanyAdvanceHandler)
	app.Get("/ruc/:numero", Pro5RucHandler)
	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
	log.Fatal(app.Listen(":3000"))
}

func Init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: "",
		DB:       1,
	})
	pgsqlClient, _ = pgxpool.New(CTX, os.Getenv("DATABASE_URI"))

	// Load ubigeos
	ubigeos, err := os.Open("ubigeos.json")
	if err == nil {
		defer ubigeos.Close()
		byteValue, _ := io.ReadAll(ubigeos)
		json.Unmarshal(byteValue, &DEPARTAMENTOS)
	} else {
		panic(err)
	}
}
