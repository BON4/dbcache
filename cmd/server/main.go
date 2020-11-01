package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	logic "dbcache/pkg/logic"
	rp "dbcache/repo"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		//Change to os.Getenv
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	status := rdb.Ping(context.Background())
	if status.Err() != nil {
		log.Fatal(status.Err(), "Fail to dial to redis")
	}

	defer rdb.Close()

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")

	if err != nil {
		log.Fatal(err)
		//panic(err)
	}

	psqlRepo := rp.NewPsqlRepository(db)
	redisRepo := rp.NewRedisCache(rdb)
	wrapper := rp.NewCacheWrapper(redisRepo, psqlRepo)

	service := logic.NewItemService(wrapper)

	r := mux.NewRouter()
	r.HandleFunc("/item/{id:[0-9]+}", service.GetItem).Methods("GET")
	r.HandleFunc("/transport/{id:[0-9]+}", service.GetTransport).Methods("GET")
	r.HandleFunc("/cacheitem/{id:[0-9]+}", service.GetItemFromCache).Methods("GET")
	r.HandleFunc("/transported_items/{id:[0-9]+}", service.GetTransportedItems).Methods("GET")

	r.HandleFunc("/items_create", service.CreateAlot).Methods("POST")
	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}
