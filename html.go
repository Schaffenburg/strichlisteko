package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"strconv"

	"github.com/gorilla/mux"
	"net/http"
)

//go:embed html/static
var staticFS embed.FS

// templates
//
//go:embed html/user.html
var userMenu string
var userMenuTemplate *template.Template

type userMenuTemplateArgs struct {
	User     User
	Products []Product
}

type usertransactionTemplateArgs struct {
	User         User
	Transactions []Transaction
}

//go:embed html/userlist.html
var userlist string
var userlistTemplate *template.Template

//go:embed html/transactionlist.html
var transactionlist string
var transactionlistTemplate *template.Template

//go:embed html/stockmgr.html
var stockmgr string
var stockmgrTemplate *template.Template

//go:embed html/error.html
var errpage string
var errpageTemplate *template.Template

//go:embed html/wallet.html
var wallet string
var walletTemplate *template.Template

//go:embed html/newuser_imgselector.html
var newuser_imgselector string
var newuser_imgselectorTemplate *template.Template

//go:embed html/edituser.html
var edituser string
var edituserTemplate *template.Template

//go:embed html/newuser.html
var newuserpage []byte

//go:embed html/confirm.html
var confirmpage []byte

//go:embed html/static/404.png
var image404 []byte

//go:embed html/static/400.png
var image400 []byte

//go:embed html/static/500.png
var image500 []byte

func init() {
	userMenuTemplate = template.New("usermenu")
	var err error
	userMenuTemplate, err = userMenuTemplate.Parse(userMenu)
	if err != nil {
		log.Fatalf("Failed to parse userMenuTemplate: %s", err)
	}

	userlistTemplate = template.New("userlist")
	userlistTemplate, err = userlistTemplate.Parse(userlist)
	if err != nil {
		log.Fatalf("Failed to parse userlistTemplate: %s", err)
	}

	transactionlistTemplate = template.New("transactionlist")
	transactionlistTemplate, err = transactionlistTemplate.Parse(transactionlist)
	if err != nil {
		log.Fatalf("Failed to parse transactionlistTemplate: %s", err)
	}

	stockmgrTemplate = template.New("stockmgr")
	stockmgrTemplate, err = stockmgrTemplate.Parse(stockmgr)
	if err != nil {
		log.Fatalf("Failed to parse stockmgrTemplate: %s", err)
	}

	errpageTemplate = template.New("errpage")
	errpageTemplate, err = errpageTemplate.Parse(errpage)
	if err != nil {
		log.Fatalf("Failed to parse errpage: %s", err)
	}

	walletTemplate = template.New("wallet")
	walletTemplate, err = walletTemplate.Parse(wallet)
	if err != nil {
		log.Fatalf("Failed to parse walletTemplate: %s", err)
	}

	newuser_imgselectorTemplate = template.New("imgselector")
	newuser_imgselectorTemplate, err = newuser_imgselectorTemplate.Parse(newuser_imgselector)
	if err != nil {
		log.Fatalf("Failed to parse newuser_imgselectorTemplate: %s", err)
	}

	edituserTemplate = template.New("edituser")
	edituserTemplate, err = edituserTemplate.Parse(edituser)
	if err != nil {
		log.Fatalf("Failed to parse edituserTemplate: %s", err)
	}
}

func handleUserPage(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("failed to get userid '%s': %s", idstr, err)
		errPage(404, "User Not Found", w, r)
		return
	}

	products, err := getProducts()
	if err != nil {
		log.Printf("failed to get products: %s", err)
		errPage(500, "General Failure", w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	err = userMenuTemplate.Execute(w, userMenuTemplateArgs{user, products})
	if err != nil {
		log.Printf("Failed to execute userMenuTemplate for user '%s': %s\n", idstr, err)
	}
}

func handleUsersList(w http.ResponseWriter, r *http.Request) {
	users, err := getUsers(true)
	if err != nil {
		errPage(500, "General Failure", w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	err = userlistTemplate.Execute(w, users)
	if err != nil {
		log.Printf("Failed to execute userlistTemplate:: %s\n", err)
	}
}

func handleUsersListInactive(w http.ResponseWriter, r *http.Request) {
	users, err := getUsers(false)
	if err != nil {
		errPage(500, "General Failure", w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	err = userlistTemplate.Execute(w, users)
	if err != nil {
		log.Printf("Failed to execute userlistTemplate:: %s\n", err)
	}
}

func handleTransactionsPage(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("failed to get userid '%s': %s", idstr, err)
		errPage(404, "User Not Found", w, r)
		return
	}

	acts, err := getTransactions(int(id))
	if err != nil {
		log.Printf("failed to get transactions: %s", err)
		errPage(500, "General Failure", w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	err = transactionlistTemplate.Execute(w, usertransactionTemplateArgs{user, acts})
	if err != nil {
		log.Printf("Failed to execute userMenuTemplate for user '%s': %s\n", idstr, err)
	}
}

func handleNewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(newuserpage)
}

func handleNewUserSubmit(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	log.Printf("Body: %s", body)

	dec := json.NewDecoder(bytes.NewReader(body))

	var user User
	err := dec.Decode(&user)
	if err != nil {
		log.Printf("Failed to decode user json: %s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"invalid argument"}`)
		return
	}

	log.Printf("image: %s", user.Image)

	userid, err := addUser(user.Username, user.Image, 0, user.Active)
	if err != nil {
		log.Printf("Failed to create user '%s': %s", user.Username, err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		fmt.Fprint(w, `{"error":"database failure"}`)
		return
	}

	log.Printf("[POST] %s: created user: '%s' %d", r.URL.Path, user.Username, userid)

	w.Header().Set("Location", "/user/"+strconv.FormatInt(int64(userid), 10))
	w.WriteHeader(201)
	fmt.Fprintf(w, `{"info":"created account, redirecting.","id":%d}`, userid)
}

func handleNewUserImgSelector(w http.ResponseWriter, r *http.Request) {
	imgs, err := getImages()
	if err != nil {
		log.Printf("failed to get images: %s", err)
		errPage(500, "Database Failure", w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	err = newuser_imgselectorTemplate.Execute(w, imgs)
	if err != nil {
		log.Printf("Failed to execute newuser_imgselectorTemplate: %s\n", err)
	}
}

func handleEditUser(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("failed to get userid '%s': %s", idstr, err)
		errPage(404, "User Not Found", w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)

	err = edituserTemplate.Execute(w, user)
	if err != nil {
		log.Printf("error editusertemplate: %s", err)
	}
}

func handleEditUserSubmit(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("failed to get userid '%s': %s", idstr, err)
		errPage(404, "User Not Found", w, r)
		return
	}

	body, _ := io.ReadAll(r.Body)
	log.Printf("Body: %s", body)

	dec := json.NewDecoder(bytes.NewReader(body))

	var updatuser User
	err = dec.Decode(&updatuser)
	if err != nil {
		log.Printf("Failed to decode user json: %s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"invalid argument"}`)
		return
	}

	user.Username = updatuser.Username
	user.Image = updatuser.Image
	user.Active = updatuser.Active

	err = setUser(user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		fmt.Fprint(w, `{"error":"database failure"}`)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	fmt.Fprintf(w, `{"info":"updated.","id":%d}`, user.ID)
}

func handleDeleteAsk(w http.ResponseWriter, r *http.Request) {
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

	err = deleteUser(int(id))
	if err != nil {
		log.Printf("Failed to delete user %d: %s", id, err)
		errPage(500, "Database Failure.", w, r)
		return
	}

	errPage(200, "deleted", w, r)
}

func handleConfirmPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(confirmpage)
}
