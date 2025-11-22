package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/catarium/avito_test_task/internal/db"
	"github.com/catarium/avito_test_task/internal/handlers/teams"
	"github.com/catarium/avito_test_task/internal/middlewares"
)

func main() {

	var err error

	log.Println("connecting to db")
	database, err := db.Connect()

	if err != nil {
		log.Fatal(err.Error())
	}
	defer database.Close()

	log.Println("creating tables")
	err = db.CreateTables(database)

	if err != nil {
		log.Fatal(err.Error())
	}

	mux := http.NewServeMux()
	mux.Handle("/team/", http.StripPrefix("/team", teams.CreateTeamRouter(database)))

	loggedMux := middlewares.LoggingMiddleware(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")),
		Handler: loggedMux,
	}

	// log.Println("starting gorutine")
	// go func() {
	// 	<-ctx.Done()
	// 	server.Shutdown(ctx)
	// }()
	err = server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err.Error())
	}
}
