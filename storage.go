package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

func handleNewProductSubmit(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	log.Printf("Body: %s", body)

	dec := json.NewDecoder(bytes.NewReader(body))

	var prod Product
	err := dec.Decode(&prod)
	if err != nil {
		log.Printf("Failed to decode user json: %s", err)

		w.Header().Set("Content-Type", "application/json")
		//Add cors header (update allowed locs lateron)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"invalid argument"}`)
		return
	}

	id, err := addProduct(prod)
	if err != nil {
		log.Printf("Failed to add product: %s", err)

		w.Header().Set("Content-Type", "application/json")
		//Add cors header (update allowed locs lateron)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(500)
		fmt.Fprint(w, `{"error":"database failure"}`)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	//Add cors header (update allowed locs lateron)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	fmt.Fprintf(w, `{"info":"success, created with id %d"}`, id)
	return
}

func handleStorage(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts()
	if err != nil {
		log.Printf("failed to get products: %s", err)
		errPage(500, "Database Failure", w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	err = stockmgrTemplate.Execute(w, products)
	if err != nil {
		log.Printf("Failed to execute stockmgrTemplate for user %s\n", err)
	}
}

func handleProductSetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	prodstr, ok := vars["product"]
	if !ok {
		errPage(404, "Product not found", w, r)
		return
	}

	prodid, err := strconv.ParseInt(prodstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse product %s", err)
		errPage(404, "Product not found", w, r)
		return
	}

	prod, err := getProduct(int(prodid))
	if err != nil {
		log.Printf("Failed to get product %s: %s", prodstr, err)
		errPage(404, "Product not found", w, r)
		return

	}

	imgstr, ok := vars["image"]
	if !ok {
		errPage(400, "Invalid Request", w, r)
		return
	}

	var img = new(uuid.UUID)
	err = img.UnmarshalText([]byte(imgstr))
	if err != nil {
		log.Printf("Failed to unmarshal uuid '%s': %s", imgstr, err)
	}

	prod.Image = *img

	err = setProduct(prod)
	if err != nil {
		errPage(400, "Database Failure.", w, r)
		return
	}

	// redirect
	w.Header().Set("Location", "/storage")
	w.WriteHeader(307)
}

func handleProductDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	prodstr, ok := vars["product"]
	if !ok {
		errPage(404, "Product not found", w, r)
		return
	}

	prodid, err := strconv.ParseInt(prodstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse product %s", err)
		errPage(404, "Product not found", w, r)
		return
	}

	err = delProduct(int(prodid))
	if err != nil {
		log.Printf("failed to delete prod %d: %s", prodid, err)
		errPage(500, "Database Failure.", w, r)
		return
	}

	// redirect
	w.Header().Set("Location", "/storage")
	w.WriteHeader(307)
}

func handleProductStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	prodstr, ok := vars["product"]
	if !ok {
		errPage(404, "Product not found", w, r)
		return
	}

	prodid, err := strconv.ParseInt(prodstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse product %s", err)
		errPage(404, "Product not found", w, r)
		return
	}

	prod, err := getProduct(int(prodid))
	if err != nil {
		errPage(404, "Product not found", w, r)
		return

	}

	amtstr, ok := vars["amount"]
	if !ok {
		errPage(404, "Product not found", w, r)
		return
	}

	amt, err := strconv.ParseInt(amtstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse product %s", err)
		errPage(404, "Product not found", w, r)
		return
	}

	prod.Stock += int(amt)

	err = setProduct(prod)
	if err != nil {
		log.Printf("Failed to update user after updating stock in transaction: %s", err)

		errPage(500, "database access failed", w, r)
		return
	}

	// redicret back
	w.Header().Set("Location", "/storage")
	w.WriteHeader(307)
}
