package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

var (
	pgPool *pgxpool.Pool
	rdb    *redis.Client
	ctx    = context.Background()
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to PostgreSQL
	databaseUrl := os.Getenv("DATABASE_URL")
	pgPool, err = pgxpool.Connect(ctx, databaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pgPool.Close()

	// Connect to Redis
	redisUrl := os.Getenv("REDIS_URL")
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatalf("Unable to parse Redis URL: %v\n", err)
	}
	rdb = redis.NewClient(opt)
	defer rdb.Close()

	// Set up the router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Define routes
	r.Get("/api/v1/imoveis", GetImoveis)
	r.Get("/api/v1/imoveis/{imovelId}", GetImovelByID)

	log.Println("Starting server on :8081...")
	http.ListenAndServe(":8081", r)
}
