package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	"hechi-go/graph"
	"hechi-go/graph/generated"
	"hechi-go/internal/auth"
	"hechi-go/internal/config"
	"hechi-go/internal/repository"
)

func main() {
	cfg := config.Load()

	db, err := repository.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	authService := auth.NewService(cfg.JWTSecret, cfg.JWTExpiryHours)

	resolver := &graph.Resolver{
		UserRepo:        repository.NewUserRepository(db),
		AccountRepo:     repository.NewAccountRepository(db),
		TransactionRepo: repository.NewTransactionRepository(db),
		Auth:            authService,
	}

	gqlSrv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: resolver}),
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler)
	r.Use(auth.Middleware(authService))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"hechi-go","version":"0.1.0"}`))
	})
	r.Handle("/graphql", gqlSrv)
	r.Handle("/playground", playground.Handler("Hechi GraphQL", "/graphql"))

	log.Printf("server listening on :%s", cfg.Port)
	log.Printf("GraphQL playground at http://localhost:%s/playground", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
