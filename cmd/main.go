package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/jonathannavas/gocourse_user/internal/user"
	"github.com/jonathannavas/gocourse_user/pkg/bootstrap"
)

func main() {

	router := mux.NewRouter()
	_ = godotenv.Load()

	logs := bootstrap.InitLogger()
	db, err := bootstrap.DBConnection()

	if err != nil {
		logs.Fatal(err)
	}

	pagLimitDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimitDef == "" {
		logs.Fatal("paginator limit default is required")
	}

	userRepository := user.NewRepo(logs, db)
	userService := user.NewService(logs, userRepository)
	userEndpoints := user.MakeEndpoints(userService, user.Config{LimitPageDef: pagLimitDef})
	usersById := "/users/{id}"

	router.HandleFunc("/users", userEndpoints.Create).Methods("POST")
	router.HandleFunc(usersById, userEndpoints.Get).Methods("GET")
	router.HandleFunc("/users", userEndpoints.GetAll).Methods("GET")
	router.HandleFunc(usersById, userEndpoints.Update).Methods("PATCH")
	router.HandleFunc(usersById, userEndpoints.Delete).Methods("DELETE")

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		// Handler:      http.TimeoutHandler(router, time.Second*3, "Timeout!!"),
		Handler:      router,
		Addr:         address,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
