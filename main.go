package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Bayan2019/rss_blog/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't  connect to database")
	}

	queries := database.New(conn)
	// if err != nil {
	// 	log.Fatal("Can't create db connection", err)
	// }

	apiCfg := apiConfig{
		DB: queries,
	}

	app_router := chi.NewRouter()

	app_router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1_router := chi.NewRouter()

	v1_router.Get("/healthz", handlerReadiness)
	v1_router.Get("/err", handlerError)

	v1_router.Post("/users", apiCfg.handlerCreateUser)

	app_router.Mount("/v1", v1_router)

	srv := &http.Server{
		Handler: app_router,
		Addr:    ":" + portString,
	}

	fmt.Printf("Server starting on port %v\n", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
