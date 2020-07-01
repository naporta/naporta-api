package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/naporta/naporta-api/db"
	"github.com/naporta/naporta-api/telegram"
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

	go telegram.Start(cfg.TelegramToken, cfg.Admin, mongo)
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/vendedor", listVendedores).Queries("c", "{c}", "cat", "{cat}").Methods("GET")
	r.HandleFunc("/vendedor", listVendedores).Methods("GET")
	r.HandleFunc("/tags", listTags).Methods("GET")
	r.HandleFunc("/produtos", listProdutos).Methods("GET")
	r.HandleFunc("/categorias", listCategorias).Methods("GET")
	r.HandleFunc("/vendedor", insertRawVendedor).Methods("POST")
	//r.HandleFunc("/vendedor", updateVendedor).Methods("PUT")
	r.HandleFunc("/vendedor/{id}", getVendedorByID).Methods("GET")
	//r.HandleFunc("/vendedor", deleteVendedor).Methods("DELETE")

	if err := http.ListenAndServe("0.0.0.0:3000", r); err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}
