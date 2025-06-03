package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/bruguedes/gobid/internal/api"
	"github.com/bruguedes/gobid/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	gob.Register(uuid.UUID{}) // Registra o tipo uuid.UUID para serialização

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ctx := context.Background()
	dbParams := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("GOBID_DB_USER"),
		os.Getenv("GOBID_DB_PASSWORD"),
		os.Getenv("GOBID_DB_HOST"),
		os.Getenv("GOBID_DB_PORT"),
		os.Getenv("GOBID_DB_NAME"),
	)

	pool, err := pgxpool.New(ctx, dbParams)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	s := scs.New()
	s.Store = pgxstore.New(pool)
	s.Lifetime = 24 * time.Hour              // define o tempo de vida da sessão
	s.Cookie.HttpOnly = true                 // define o cookie como HttpOnly
	s.Cookie.SameSite = http.SameSiteLaxMode // define o SameSite do cookie

	api := api.API{
		Router:      chi.NewMux(),
		UserService: services.NewUserService(pool),
		Sessions:    s,
	}

	api.BindRoutes()
	fmt.Println("Server is running on port 8080	")

	if err := http.ListenAndServe(":8080", api.Router); err != nil {
		panic(err)

	}
}
