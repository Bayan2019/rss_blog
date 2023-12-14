package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	// feed, err := urlToFeed("https://wagslane.dev/index.xml")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(feed)

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

	db := database.New(conn)

	apiCfg := apiConfig{
		DB: db,
	}

	go startScraping(db, 10, time.Minute)

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
	v1_router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	v1_router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1_router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1_router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1_router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1_router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

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
