package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	app := pocketbase.New()

	isDev := os.Getenv("ENV") != "production"
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isDev,
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/health", func(e *core.RequestEvent) error {
			return e.JSON(http.StatusOK, map[string]string{
				"status":  "ok",
				"service": "hechi-go",
				"version": "0.1.0",
			})
		})

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
