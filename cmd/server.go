package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"time"
	"github.com/rs/cors"
	domain "github.com/maxabcr2000/LoginServerGo/pkg/domain"
	mw "github.com/uniontsai/httpmiddlewarego"
	jwt "github.com/dgrijalva/jwt-go"
)

type ServerDependency func(*LoginServer) error

type Repository interface {
	SaveUser(user *domain.User, key string) error
	ReadUser(key string) (*domain.User, error)
}

type LoginServer struct{
	signKey interface{}
	repo Repository
}

func (server *LoginServer) handleLogin(w http.ResponseWriter, req *http.Request){
	fmt.Println("listJSON Endpoint: ", req.RemoteAddr)

	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	if req.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return 
	}

	type CreateUserRequest struct{
		Account string
		Password string
		Email string
	}

	user:= &domain.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return 
	}

	savedUser,err := server.repo.ReadUser(user.Account)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return 
	}

	if user.Password!=savedUser.Password{
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		Subject: user.Account,
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(1* time.Hour).Unix(),
	})
	
	fmt.Println("signKey:", server.signKey)

	tokenString, err := token.SignedString(server.signKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return 
	}

	fmt.Fprint(w, tokenString)
}

func (server *LoginServer) handleCreateUser(w http.ResponseWriter, req *http.Request){
	fmt.Println("listJSON Endpoint: ", req.RemoteAddr)

	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	if req.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		panic(err)
	}

	user:= &domain.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		panic(err)
	}

	err = server.repo.SaveUser(user, user.Account)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
		privateKey, err:=jwt.ParseRSAPrivateKeyFromPEM([]byte(signKey))
		if err != nil {
			panic(err)
		}

		server.signKey = privateKey
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
	mux.HandleFunc("/login", mw.PostOnly(server.handleLogin))
	mux.HandleFunc("/createUser", mw.PostOnly(server.handleCreateUser))

	err:=http.ListenAndServe(":8889", access(mux))
	if err!=nil{
		fmt.Println("ListenAndServe Error: ", err)
	}
}
