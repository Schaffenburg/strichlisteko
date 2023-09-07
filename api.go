package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/mux"
	"net/http"
)

type UserAction string

const (
	UserBuy              UserAction = "buy"
	UserDeposit                     = "deposit"
	UserWithdraw                    = "withdraw"
	UserListTransactions            = "listtransaction"
	UserUpdate                      = "update"
)

type UserAPIRequest struct {
	Action UserAction `json:"action"`

	Product int `json:"product,omitempty"`
	Amount  int `json:"amount,omitempty"`

	Transaction int `json:"transaction,omitempty"`

	User User `json:"user,omitempty"`
}

func handleUserAPIList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr, ok := vars["id"]
	if !ok {
		errAPI(404, "User Not Found", w, r)
		return
	}

	id, err := strconv.ParseInt(idstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse userid '%s': %s", idstr, err)
		errAPI(404, "User Not Found", w, r)
		return
	}

	user, err := getUser(int(id))
	if err != nil {
		log.Printf("failed to get userid '%s': %s", idstr, err)
		errAPI(404, "User Not Found", w, r)
		return
	}

	apireq := new(UserAPIRequest)

	dec := json.NewDecoder(r.Body)
	dec.Decode(&apireq)

	switch apireq.Action {
	case UserBuy:
		prod, err := getProduct(apireq.Product)
		if err != nil {
			errAPI(404, "Product not found", w, r)
			return
		}

		user.Balance -= prod.Price
		prod.Stock -= 1

		err = setUser(user)
		if err != nil {
			log.Printf("Failed to update user after updating balance in transaction: %s", err)

			errAPI(500, "database access failed", w, r)
			return
		}

		err = setProduct(prod)
		if err != nil {
			log.Printf("Failed to update user after updating stock in transaction: %s", err)

			errAPI(500, "database access failed", w, r)
			return
		}

		infoAPI(200, fmt.Sprintf(
			"bought %s for %dct.",
			prod.Name, prod.Price),
			w, r)

		// add transaction
		err = addTransaction(-prod.Price, prod.Name, user.ID)
		if err != nil {
			log.Printf("Failed to add transaction to log: %s", err)
		}

	case UserDeposit:
		user.Balance += int(apireq.Amount)

		err = setUser(user)
		if err != nil {
			log.Printf("failed to update user after adjusting balance")
			errAPI(500, "database failure", w, r)
			return
		}

		err = addTransaction(apireq.Amount, "deposit", user.ID)
		if err != nil {
			log.Printf("Failed to add transaction to log: %s", err)
		}

		infoAPI(200, fmt.Sprintf("deposited %dct. new balace %dct.",
			apireq.Amount, user.Balance), w, r)

	case UserWithdraw:
		user.Balance -= int(apireq.Amount)

		err = setUser(user)
		if err != nil {
			log.Printf("failed to update user after adjusting balance")
			errAPI(500, "database failure", w, r)
			return
		}

		err = addTransaction(-apireq.Amount, "withdraw", user.ID)
		if err != nil {
			log.Printf("Failed to add transaction to log: %s", err)
		}

		infoAPI(200, fmt.Sprintf("withdrew %dct. new balace %dct.",
			apireq.Amount, user.Balance), w, r)

	case UserUpdate:
		if apireq.User.Username == "" {
			log.Printf("Failed to decode user json: %s", err)

			errAPI(400, "failed to decode user", w, r)
			return
		}

		// unchangable
		apireq.User.ID = user.ID

		err = setUser(apireq.User)
		if err != nil {
			log.Printf("Failed to update user: %s", err)
			errAPI(500, "database failure", w, r)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		enc := json.NewEncoder(w)
		err = enc.Encode(&UserAPIResponse{
			Info: "updated user",
			User: apireq.User,
		})

		if err != nil {
			log.Printf("Failed to encode response: %s", err)
			errAPI(500, "failed to encode response", w, r)
		}

	default:
		errAPI(400, "Invalid Action!", w, r)
	}
}

type InfoResponse struct {
	Info string `json:"info"`
}

func infoAPI(status int, text string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	err := enc.Encode(&InfoResponse{
		Info: text,
	})

	if err != nil {
		fmt.Fprintf(w, "Error encoding: '%s'", text)
		return
	}
}

func handleUserGetAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr, ok := vars["id"]
	if !ok {
		errAPI(404, "User Not Found", w, r)
		return
	}

	id, err := strconv.ParseInt(idstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse userid '%s': %s", idstr, err)
		errAPI(404, "User Not Found", w, r)
		return
	}

	user, err := getUser(int(id))
	if err != nil {
		log.Printf("failed to get userid '%s': %s", idstr, err)
		errAPI(404, "User Not Found", w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	enc := json.NewEncoder(w)
	err = enc.Encode(user)
	if err != nil {
		log.Printf("Failed to encode user %d: %s", id, err)
		errAPI(500, "failed to encode response", w, r)
	}
}

func handleUsersAPI(w http.ResponseWriter, r *http.Request) {

	users, err := getUsers(true)
	if err != nil {
		log.Printf("failed to get users '%s'", err)
		errAPI(404, "User Not Found", w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	enc := json.NewEncoder(w)
	err = enc.Encode(users)
	if err != nil {
		log.Printf("Failed to encode users %s", err)
		errAPI(500, "failed to encode response", w, r)
	}
}

func handleTransactionAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr, ok := vars["id"]
	if !ok {
		errAPI(404, "User Not Found", w, r)
		return
	}

	id, err := strconv.ParseInt(idstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse userid '%s': %s", idstr, err)
		errAPI(404, "User Not Found", w, r)
		return
	}

	_, err = getUser(int(id))
	if err != nil {
		log.Printf("failed to get userid '%s': %s", idstr, err)
		errAPI(404, "User Not Found", w, r)
		return
	}

	acts, err := getTransactions(int(id))
	if err != nil {
		log.Printf("failed to get transactions: %s", err)
		errAPI(500, "Database Failure", w, r)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(acts)
	if err != nil {
		errAPI(500, "error encoding response", w, r)
	}
}

type TransactionAction string

const (
	TransactionUndo TransactionAction = "undo"
)

type TransactionAPIRequest struct {
	Action TransactionAction `json:"action"`

	ID int `json:"ID"`
}

func handleTransactionAPIPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr, ok := vars["id"]
	if !ok {
		errAPI(404, "User Not Found", w, r)
		return
	}

	id, err := strconv.ParseInt(idstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse userid '%s': %s", idstr, err)
		errAPI(404, "User Not Found", w, r)
		return
	}

	user, err := getUser(int(id))
	if err != nil {
		log.Printf("failed to get userid '%s': %s", idstr, err)
		errAPI(404, "User Not Found", w, r)
		return
	}

	apireq := new(TransactionAPIRequest)

	dec := json.NewDecoder(r.Body)
	dec.Decode(&apireq)

	switch apireq.Action {
	default:
		errAPI(400, "Invalid Action!", w, r)
		return

	case TransactionUndo:
		trans, err := getTransaction(apireq.ID)
		if err != nil {
			log.Printf("failed to get trasaction %d (user %d): %s", apireq.ID, id, err)
			errAPI(404, "Transaction not found", w, r)
			return
		}

		// check if transaction is allready undone
		if trans.Undone {
			errAPI(400, "Transaction allready done", w, r)
			return
		}

		// undo unstock when transaction is product
	cond:
		if trans.Product != "deposit" && trans.Product != "withdrawal" {
			prod, err := getProductByName(trans.Product)
			if err != nil {
				log.Printf("Failed to get product '%s': %s", trans.Product, err)
				goto cond
			}

			prod.Stock += 1

			err = setProduct(prod)
			if err != nil {
				log.Printf("Faield to update product after updating stock:: %s", err)
			}
		}

		user.Balance += -trans.Value

		err = setUser(user)
		if err != nil {
			log.Printf("Failed to update user after updating balance in transaction: %s", err)

			errAPI(500, "database access failed", w, r)
			return
		}

		err = setTransactionUndo(trans.ID, true)
		if err != nil {
			log.Printf("Failed to update transaction after updating undo-status: %s", err)

			errAPI(500, "database access failed", w, r)
			return
		}

		infoAPI(200, "undone", w, r)
	}
}

func handleProductsAPIList(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts()
	if err != nil {
		log.Printf("failed to get products: %s", err)
		errAPI(500, "General Failure", w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	enc := json.NewEncoder(w)
	err = enc.Encode(products)
	if err != nil {
		log.Printf("Failed to encode response: %s", err)
		errAPI(500, "failed to encode response", w, r)
		return
	}
}

type UserAPIResponse struct {
	User User   `json:"user"`
	Info string `json:"info"`
}

func handleUserAPINew(w http.ResponseWriter, r *http.Request) {
	var user User

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&user)
	if err != nil {
		log.Printf("Failed to decode user json: %s", err)

		errAPI(400, "failed to decode user!", w, r)
		return
	}

	log.Printf("image: %s", user.Image)

	user.ID, err = addUser(user.Username, user.Image, user.Balance, user.Active)
	if err != nil {
		log.Printf("Failed to create user '%s': %s", user.Username, err)

		errAPI(500, "database failure", w, r)
		return
	}

	log.Printf("[POST] %s: created user: '%s' %d", r.URL.Path, user.Username, user.ID)

	enc := json.NewEncoder(w)
	err = enc.Encode(&UserAPIResponse{
		User: user,
		Info: "created user.",
	})

	if err != nil {
		log.Printf("Failed to encode response: %s", err)

		errAPI(500, "failed to encode response", w, r)
	}
}

type ProductAPIResponse struct {
	Product Product `json:"product"`
	Info    string  `json:"string"`
}

func handleProductsAPINew(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var prod Product
	err := dec.Decode(&prod)
	if err != nil {
		log.Printf("Failed to decode product json: %s", err)

		errAPI(400, "failed to decode product", w, r)
		return
	}

	id, err := addProduct(prod)
	if err != nil {
		log.Printf("Failed to add product: %s", err)

		errAPI(500, "database failure", w, r)
		return
	}

	prod.ID = int(id)

	enc := json.NewEncoder(w)
	err = enc.Encode(&ProductAPIResponse{
		Info:    "created product",
		Product: prod,
	})
	return
}

func handleUserListAPI(w http.ResponseWriter, r *http.Request) {
	users, err := getUsers(true)
	if err != nil {
		errAPI(500, "Database Failure", w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	enc := json.NewEncoder(w)
	err = enc.Encode(users)
	if err != nil {
		errAPI(500, "failed to encode response", w, r)

		log.Printf("Faield to encode response: %s", err)
		return
	}
}

type ProductAction string

const (
	ProductStock  ProductAction = "stock"
	ProductDelete ProductAction = "delete"
)

type ProductAPIRequest struct {
	Action ProductAction `json:"action"`

	Amount int `json:"amount,omitempty"`
}

func handleProductsAPI(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts()
	if err != nil {
		log.Printf("failed to get products: %s", err)
		errAPI(500, "General Failure", w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	enc := json.NewEncoder(w)
	err = enc.Encode(products)
	if err != nil {
		log.Printf("Failed to encode response: %s", err)
		errAPI(500, "failed to encode response", w, r)
		return
	}
}

func handleProductAPI(w http.ResponseWriter, r *http.Request) {
	req := new(ProductAPIRequest)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(req)
	if err != nil {
		errAPI(400, "invalid request", w, r)
		return
	}

	vars := mux.Vars(r)
	prodstr, ok := vars["product"]
	if !ok {
		errAPI(404, "Product not found", w, r)
		return
	}

	prodid, err := strconv.ParseInt(prodstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse product %s", err)
		errAPI(404, "Product not found", w, r)
		return
	}

	prod, err := getProduct(int(prodid))
	if err != nil {
		errAPI(404, "Product not found", w, r)
		return
	}

	switch req.Action {
	case ProductStock:
		prod.Stock += int(req.Amount)

		err = setProduct(prod)
		if err != nil {
			log.Printf("Failed to update user after updating stock: %s", err)

			errAPI(500, "database access failed", w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		enc := json.NewEncoder(w)
		err = enc.Encode(&ProductAPIResponse{
			Product: prod,
			Info:    "updated stock.",
		})
		if err != nil {
			errAPI(500, "failed to encode response", w, r)

			log.Printf("Faield to encode response: %s", err)
			return
		}

	case ProductDelete:
		err = delProduct(int(prodid))
		if err != nil {
			log.Printf("failed to delete prod %d: %s", prodid, err)

			errAPI(500, "Database Failure.", w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		enc := json.NewEncoder(w)
		err = enc.Encode(&ProductAPIResponse{
			Product: prod,
			Info:    "deleted product.",
		})
		if err != nil {
			errAPI(500, "failed to encode response", w, r)

			log.Printf("Faield to encode response: %s", err)
			return
		}

	default:
		errAPI(404, "invalid action", w, r)
	}
}
