package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"

	"log"
	"strconv"
)

func handleBuy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr, ok := vars["id"]
	if !ok {
		errPage(404, "User Not Found", w, r)
		return
	}

	id, err := strconv.ParseInt(idstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse userid '%s': %s", idstr, err)
		errPage(404, "User Not Found", w, r)
		return
	}

	user, err := getUser(int(id))
	if err != nil {
		errPage(404, "User Not Found", w, r)
		return
	}

	prodstr, ok := vars["product"]
	if !ok {
		errPage(404, "Product not found", w, r)
		return
	}

	prodid, err := strconv.ParseInt(prodstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse prodint '%s': %s", idstr, err)
		errPage(404, "Product not found", w, r)
		return
	}

	prod, err := getProduct(int(prodid))
	if err != nil {
		errPage(404, "Product not found", w, r)
		return

	}

	// do transaction thing
	user.Balance -= prod.Price
	prod.Stock -= 1

	err = setUser(user)
	if err != nil {
		log.Printf("Failed to update user after updating balance in transaction: %s", err)

		errPage(500, "database access failed", w, r)
		return
	}

	err = setProduct(prod)
	if err != nil {
		log.Printf("Failed to update user after updating stock in transaction: %s", err)

		errPage(500, "database access failed", w, r)
		return
	}

	// add transaction
	err = addTransaction(-prod.Price, prod.Name, user.ID)
	if err != nil {
		log.Printf("Failed to add transaction to log: %s", err)
	}

	// redirect to user page
	w.Header().Set("Location", "/user/"+strconv.FormatInt(int64(user.ID), 10))
	w.WriteHeader(307)
}

func handleTransactionUndo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr, ok := vars["id"]
	if !ok {
		errPage(404, "User Not Found", w, r)
		return
	}

	id, err := strconv.ParseInt(idstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse userid '%s': %s", idstr, err)
		errPage(404, "User Not Found", w, r)
		return
	}

	user, err := getUser(int(id))
	if err != nil {
		errPage(404, "User Not Found", w, r)
		return
	}

	transstr, ok := vars["transaction"]
	if !ok {
		errPage(404, "Product not found", w, r)
		return
	}

	transid, err := strconv.ParseInt(transstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse prodint '%s': %s", idstr, err)
		errPage(404, "Transaction not found", w, r)
		return
	}

	trans, err := getTransaction(int(transid))
	if err != nil {
		log.Printf("failed to get trasaction %d (user %d): %s", transid, id, err)
		errPage(404, "Transaction not found", w, r)
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

		errPage(500, "database access failed", w, r)
		return
	}

	err = setTransactionUndo(trans.ID, true)
	if err != nil {
		log.Printf("Failed to update transaction after updating undo-status: %s", err)

		errPage(500, "database access failed", w, r)
		return
	}

	// redirect
	w.Header().Set("Location", "/user/"+strconv.FormatInt(id, 10)+"/transactions")
	w.WriteHeader(307)
}

func handleWallet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr, ok := vars["id"]
	if !ok {
		errPage(404, "User Not Found", w, r)
		return
	}

	id, err := strconv.ParseInt(idstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse userid '%s': %s", idstr, err)
		errPage(404, "User Not Found", w, r)
		return
	}

	user, err := getUser(int(id))
	if err != nil {
		errPage(404, "User Not Found", w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)

	walletTemplate.Execute(w, user)
}

func handleWalletDeposit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr, ok := vars["id"]
	if !ok {
		errPage(404, "User Not Found", w, r)
		return
	}

	id, err := strconv.ParseInt(idstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse userid '%s': %s", idstr, err)
		errPage(404, "User Not Found", w, r)
		return
	}

	user, err := getUser(int(id))
	if err != nil {
		errPage(404, "User Not Found", w, r)
		return
	}

	amtstr, ok := vars["amount"]
	if !ok {
		errPage(404, "User Not Found", w, r)
		return
	}

	amt, err := strconv.ParseInt(amtstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse amount '%s': %s", idstr, err)
		errPage(400, "Invalid Amount", w, r)
		return
	}

	user.Balance += int(amt)

	err = setUser(user)
	if err != nil {
		log.Printf("failed to update user after adjusting balance")
		errPage(500, "database failure", w, r)
		return
	}

	// add transaction
	err = addTransaction(int(amt), "deposit", user.ID)
	if err != nil {
		log.Printf("Failed to add transaction to log: %s", err)
	}

	// redirect to user page
	w.Header().Set("Location", "/user/"+strconv.FormatInt(int64(user.ID), 10)+"/wallet")
	w.WriteHeader(307)
}

func handleWalletWithdraw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr, ok := vars["id"]
	if !ok {
		errPage(404, "User Not Found", w, r)
		return
	}

	id, err := strconv.ParseInt(idstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse userid '%s': %s", idstr, err)
		errPage(404, "User Not Found", w, r)
		return
	}

	user, err := getUser(int(id))
	if err != nil {
		errPage(404, "User Not Found", w, r)
		return
	}

	amtstr, ok := vars["amount"]
	if !ok {
		errPage(404, "User Not Found", w, r)
		return
	}

	amt, err := strconv.ParseInt(amtstr, 10, 32)
	if err != nil {
		log.Printf("failed to parse amount '%s': %s", idstr, err)
		errPage(400, "Invalid Amount", w, r)
		return
	}

	user.Balance -= int(amt)

	err = setUser(user)
	if err != nil {
		log.Printf("failed to update user after adjusting balance")
		errPage(500, "database failure", w, r)
		return
	}

	// add transaction
	err = addTransaction(-int(amt), "withdrawal", user.ID)
	if err != nil {
		log.Printf("Failed to add transaction to log: %s", err)
	}

	// redirect to user page
	w.Header().Set("Location", "/user/"+strconv.FormatInt(int64(user.ID), 10)+"/wallet")
	w.WriteHeader(307)
}
