package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/naporta/naporta-api/db"
	"github.com/naporta/naporta-api/telegram"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mongo = db.Connection{}

func init() {

	cfg, err := loadConfig()
	if err != nil {
		log.Printf("Config file err:%v", err)
		return
	}

	mongo = db.Connection{
		User:     cfg.MongoUser,
		Password: cfg.MongoPassword,
		Server:   cfg.MongoServer,
		Database: cfg.MongoDB,
	}

	err = mongo.Connect()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to the database")

	go telegram.Start(cfg.TelegramToken, mongo)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/vendedor", func(w http.ResponseWriter, r *http.Request) {
		vendedores, err := mongo.FindAll()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJson(w, http.StatusOK, vendedores)
	}).Methods("GET")

	r.HandleFunc("/vendedor", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var vendedor db.Vendedor
		if err := json.NewDecoder(r.Body).Decode(&vendedor); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		res, err := mongo.Insert(vendedor)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		vendedor.ID = res.InsertedID.(primitive.ObjectID)
		respondWithJson(w, http.StatusCreated, vendedor)
	}).Methods("POST")

	r.HandleFunc("/vendedor", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var vendedor db.Vendedor

		if err := json.NewDecoder(r.Body).Decode(&vendedor); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		log.Println(vendedor)
		res, err := mongo.Update(vendedor)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		log.Printf("%+v\n", res)
		if res.MatchedCount < 1 {
			respondWithJson(w, http.StatusOK, map[string]string{"status": "no op"})
			return
		}

		respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
	}).Methods("PUT")

	r.HandleFunc("/vendedor/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		movie, err := mongo.FindByID(params["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid Vendedor ID")
			return
		}
		respondWithJson(w, http.StatusOK, movie)
	}).Methods("GET")

	r.HandleFunc("/vendedor", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var vendedor db.Vendedor
		if err := json.NewDecoder(r.Body).Decode(&vendedor); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		res, err := mongo.Delete(vendedor)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if res.DeletedCount < 1 {
			respondWithJson(w, http.StatusOK, map[string]string{"result": "no op"})
			return
		}
		respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
	}).Methods("DELETE")

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}
