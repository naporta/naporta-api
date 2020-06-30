package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/naporta/naporta-api/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func listVendedores(w http.ResponseWriter, r *http.Request) {
	condominio := mux.Vars(r)["c"]
	vendedores, err := mongo.FindAll(condominio)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, vendedores)
}

func listTags(w http.ResponseWriter, r *http.Request) {
	tags, err := mongo.GetTags()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, tags)
}

func listProdutos(w http.ResponseWriter, r *http.Request) {
	tags, err := mongo.GetProdutos()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, tags)
}

func insertRawVendedor(w http.ResponseWriter, r *http.Request) {
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
	vendedor.Verificado = false
	vendedor.Assinante = false

	respondWithJson(
		w,
		http.StatusCreated,
		map[string]string{"result": "success"},
	)
}

func updateVendedor(w http.ResponseWriter, r *http.Request) {
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
}

func getVendedorByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	vendedor, err := mongo.FindByID(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Vendedor ID")
		return
	}

	result := bson.M{
		"id":         vendedor.ID,
		"condominio": vendedor.Condominio,
		"nome":       vendedor.Nome,
		"empresa":    vendedor.Empresa,
		"profissao":  vendedor.Profissao,
		"whastapp":   vendedor.Whatsapp,
		"facebook":   vendedor.Facebook,
		"instagram":  vendedor.Instagram,
		"bloco":      vendedor.Bloco,
		"apt":        vendedor.Apt,
		"pagamento":  vendedor.Pagamento,
	}

	respondWithJson(
		w,
		http.StatusOK,
		result,
	)
}

func deleteVendedor(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var vendedor db.Vendedor
	if err := json.NewDecoder(r.Body).Decode(&vendedor); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	res, err := mongo.Delete(vendedor.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if res.DeletedCount < 1 {
		respondWithJson(w, http.StatusOK, map[string]string{"result": "no op"})
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
