package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// GetImoveis handles GET requests to /imoveis
func GetImoveis(w http.ResponseWriter, r *http.Request) {

	var imoveis []Imovel
	val, err := rdb.GetEx(ctx, "imovel", 10*time.Second).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		err = json.Unmarshal([]byte(val), &imoveis)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(imoveis)
		return
	}

	rows, err := pgPool.Query(ctx, "SELECT id, endereco, preco, area FROM imovel")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var imovel Imovel
		err := rows.Scan(&imovel.ID, &imovel.Endereco, &imovel.Preco, &imovel.Area)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		imoveis = append(imoveis, imovel)
	}

	jsonResponse, err := json.Marshal(imoveis)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err2 := rdb.SetEX(ctx, "imovel", jsonResponse, 10*time.Second).Err()
	if err2 != nil {
		log.Printf("Error setting value in Redis: %v\n", err2)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(imoveis)
}

// GetImovelByID handles GET requests to /imoveis/{imovelId}
func GetImovelByID(w http.ResponseWriter, r *http.Request) {
	imovelId := chi.URLParam(r, "imovelId")
	id, err := strconv.Atoi(imovelId)
	if err != nil {
		http.Error(w, "Invalid imovel ID", http.StatusBadRequest)
		return
	}

	var imovel Imovel
	val, err := rdb.GetEx(ctx, "imovel#"+imovelId, 10*time.Second).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		err = json.Unmarshal([]byte(val), &imovel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(imovel)
		return
	}

	err = pgPool.QueryRow(ctx, "SELECT id, endereco, preco, area FROM imovel WHERE id=$1", id).Scan(&imovel.ID, &imovel.Endereco, &imovel.Preco, &imovel.Area)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(imovel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err2 := rdb.SetEX(ctx, "imovel#"+imovelId, jsonResponse, 10*time.Second).Err()
	if err2 != nil {
		log.Printf("Error setting value in Redis: %v\n", err2)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
