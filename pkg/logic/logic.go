package logic

import (
	"context"
	"dbcache/models"
	"dbcache/repo"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Service interface {
	GetItem(w http.ResponseWriter, r *http.Request)
	GetTransport(w http.ResponseWriter, r *http.Request)
	GetItemFromCache(w http.ResponseWriter, r *http.Request)
	GetTransportedItems(w http.ResponseWriter, r *http.Request)

	CreateAlot(w http.ResponseWriter, r *http.Request)
}

type ItemService struct {
	wrapper repo.Wrapper
}

func NewItemService(wr repo.Wrapper) Service {
	return &ItemService{wrapper: wr}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (i *ItemService) GetItem(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	item, err := i.wrapper.GetItem(ctx, mux.Vars(r)["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	b, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func (i *ItemService) GetTransport(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	item, err := i.wrapper.GetTransport(ctx, mux.Vars(r)["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	b, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func (i *ItemService) GetTransportedItems(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	c, err := i.wrapper.GetTransportItems(ctx, mux.Vars(r)["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	b, err := json.Marshal(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func (i *ItemService) GetItemFromCache(w http.ResponseWriter, r *http.Request) {
	cache := i.wrapper.GetCache()
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	b, err := cache.Get(ctx, "b:1")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(b)
}

func (i *ItemService) CreateAlot(w http.ResponseWriter, r *http.Request) {
	var createReq models.CreateAlotItems

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&createReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if len(createReq.TransportId) == 0 || createReq.Length == 0 {
		http.Error(w, "Invalid request payload", http.StatusInternalServerError)
		return
	}

	db := i.wrapper.GetRepository()
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	if err := db.CreateAlotItems(ctx, createReq); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}
