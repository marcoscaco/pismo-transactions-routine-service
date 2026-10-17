package main

import "github.com/marcoscaco/pismo-transactions-routine-service/internal/adapter/httpapi"

const serverAddress = ":8080"

func main() {
	router := httpapi.NewRouter(httpapi.Handlers{
		Health: httpapi.NewHealthHandler(),
	})

	if err := router.Run(serverAddress); err != nil {
		panic(err)
	}
}
