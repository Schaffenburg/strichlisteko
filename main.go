package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
}

func main() {
	r := mux.NewRouter()

	r.Use(accessControlMiddleware)
	r.HandleFunc("/", handleUsersList).Methods("GET")
	r.HandleFunc("/inactive", handleUsersListInactive).Methods("GET")

	r.HandleFunc("/confirm", handleConfirmPage).Methods("GET")

	// lager
	r.HandleFunc("/storage/{product}/setimg/{image}", handleProductSetImage).Methods("GET")
	r.HandleFunc("/storage/{product}/delete", handleProductDelete).Methods("GET")
	r.HandleFunc("/storage/{product}/stock/{amount}", handleProductStock).Methods("GET")

	r.HandleFunc("/storage/new", handleNewProductSubmit).Methods("POST")

	r.HandleFunc("/storage", handleStorage).Methods("GET")

	// images
	r.HandleFunc("/img/{image}", handleImage).Methods("GET")
	r.HandleFunc("/img/upload", handlePutImage).Methods("PUT")

	// user releated stffs
	// new user page
	r.HandleFunc("/user/new", handleNewUser).Methods("GET")
	r.HandleFunc("/user/new/imgselector", handleNewUserImgSelector).Methods("GET")

	r.HandleFunc("/user/new", handleNewUserSubmit).Methods("POST")

	//API START
	r.HandleFunc("/api/img/{image}", handleImage).Methods("GET")
	r.HandleFunc("/api/img/new", handlePutImageAPI).Methods("POST")

	r.HandleFunc("/api/storage", handleProductsAPIList).Methods("GET")
	r.HandleFunc("/api/storage/new", handleProductsAPINew).Methods("POST")
	r.HandleFunc("/api/storage/new", PreflightRoot).Methods("OPTIONS")
	r.HandleFunc("/api/storage/{product}", handleProductAPI).Methods("POST")
	r.HandleFunc("/api/storage/{product}", PreflightRoot).Methods("OPTIONS")

	r.HandleFunc("/api/products", handleProductsAPI).Methods("GET")
	r.HandleFunc("/api/products/new", handleNewProductSubmit).Methods("POST")

	r.HandleFunc("/api/user/new", handleUserAPINew).Methods("POST")
	r.HandleFunc("/api/users", handleUserListAPI).Methods("GET")

	r.HandleFunc("/api/user/{id}", handleUserGetAPI).Methods("GET")
	r.HandleFunc("/api/user/{id}", handleUserAPIList).Methods("POST")
	r.HandleFunc("/api/user/{id}", PreflightRoot).Methods("OPTIONS")

	r.HandleFunc("/api/user/{id}/transactions", handleTransactionAPI).Methods("GET")
	r.HandleFunc("/api/user/{id}/transactions", handleTransactionAPIPost).Methods("POST")

	//API END

	r.HandleFunc("/user/{id}", handleUserPage).Methods("GET")
	r.HandleFunc("/user/{id}/buy/{product}", handleBuy).Methods("GET")

	r.HandleFunc("/user/{id}/transactions/{transaction}/undo", handleTransactionUndo).Methods("GET")

	r.HandleFunc("/user/{id}/transactions", handleTransactionsPage).Methods("GET")
	r.HandleFunc("/user/{id}/settings", handleEditUser).Methods("GET")
	r.HandleFunc("/user/{id}/delete", handleDeleteAsk).Methods("GET")
	r.HandleFunc("/user/{id}/settings", handleEditUserSubmit).Methods("POST")

	r.HandleFunc("/user/{id}/wallet", handleWallet).Methods("GET")
	r.HandleFunc("/user/{id}/wallet/deposit/{amount}", handleWalletDeposit).Methods("GET")
	r.HandleFunc("/user/{id}/wallet/withdraw/{amount}", handleWalletWithdraw).Methods("GET")
	r.PathPrefix("/s/").
		Handler(http.StripPrefix("/s/",
			AddPrefix("html/static/",
				http.FileServer(http.FS(staticFS)),
			),
		))

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// access control and  CORS middleware
func accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func PreflightRoot(respWriter http.ResponseWriter, request *http.Request) {
	respWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	respWriter.Header().Set("Access-Control-Allow-Headers", "*")
}
