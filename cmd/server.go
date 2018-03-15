package server

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"time"
	"github.com/rs/cors"
	authentication "github.com/uniontsai/httpmiddlewarego"
	jwt "github.com/dgrijalva/jwt-go"
)

const(
	BOLT_BUCKET_NAME_USERS = "Users"
	PARM_USER_ACC = "account"
	PARM_USER_PASS= "password"
)

type ServerDependency func(*LoginServer) error

type Repository interface {
	SaveMessage(bucketName, key, value string) error
	ReadMessage(bucketName, key string) (string, error)
}

type LoginServer struct{
	signKey string
	repo Repository
}

type User struct {
    Account string `json:"account"`
	Password string `json:"password"`
	Name string `json:"name"`
	Email string `json:"email"`
}

func (server *LoginServer) handleLogin(w http.ResponseWriter, req *http.Request){
	fmt.Println("listJSON Endpoint: ", req.RemoteAddr)

	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "We only support application/json format in POST.", http.StatusUnsupportedMediaType)
		return
	}

	if req.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	user:= &User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		panic(err)
	}

	savedUser:= &User{}
	savedData,err := server.repo.ReadMessage(BOLT_BUCKET_NAME_USERS, user.Account)
	err = json.Unmarshal([]byte(savedData), savedUser)
	if err != nil {
		panic(err)
	}

	if user.Password!=savedUser.Password{
		http.Error(w, "Wrong account / password.", http.StatusUnauthorized)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		Subject: user.Account,
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(1* time.Hour).Unix(),
	})
	
	fmt.Println("signKey:", server.signKey)

	privateKey, err:=jwt.ParseRSAPrivateKeyFromPEM([]byte(server.signKey))
	if err != nil {
		panic(err)
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, tokenString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (server *LoginServer) handleCreateUser(w http.ResponseWriter, req *http.Request){
	fmt.Println("listJSON Endpoint: ", req.RemoteAddr)

	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "We only support application/json format in POST.", http.StatusUnsupportedMediaType)
		return
	}

	if req.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	user:= &User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		panic(err)
	}
	
	existedUser, err := server.repo.ReadMessage(BOLT_BUCKET_NAME_USERS, user.Account)
	if existedUser !=""{
		http.Error(w, fmt.Sprintf("User Account: %s is already used!", user.Account), http.StatusInternalServerError)
		return
	}

	err = server.repo.SaveMessage(BOLT_BUCKET_NAME_USERS, user.Account, string(body))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w,"User creation has completed.")
}

func WithRepository(repo Repository) ServerDependency{
	return func(server *LoginServer) error{
		server.repo = repo
		return nil
	}
}

func WithSignKey(signKey string) ServerDependency{
	return func(server *LoginServer) error{
		server.signKey = signKey
		return nil
	}
}

func NewLoginServer(deps ...ServerDependency) (*LoginServer, error){
	server:=&LoginServer{}
	for _,dep := range deps{
		dep(server)
	}

	return server,nil
}

func (server *LoginServer) Start(){
	access := cors.AllowAll().Handler
	mux := http.NewServeMux()
	mux.HandleFunc("/login", authentication.PostOnly(server.handleLogin))
	mux.HandleFunc("/createUser", authentication.PostOnly(server.handleCreateUser))

	err:=http.ListenAndServe(":8889", access(mux))
	if err!=nil{
		fmt.Println("ListenAndServe Error: ", err)
	}
}
