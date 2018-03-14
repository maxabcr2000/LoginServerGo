package main

import (
	"os"
	"fmt"
	repository "github.com/maxabcr2000/LoginServerGo/pkg/repository"
	cmd "github.com/maxabcr2000/LoginServerGo/cmd"
)

const (
	DBName = "LoginServer.db"
	UsersBucketName = "Users"
	SignKeyVarName = "TOKEN_SIGN_KEY"
)

func main(){
	signKey, ok := os.LookupEnv(SignKeyVarName)
	if !ok {
		fmt.Printf("Failed to get TOKEN_SIGN_KEY from environmental settings.")
		return
	}

	repo,err:= repository.CreateBoltRepository(DBName);
	if err!=nil{
		panic(err)
	}

	repo.CreateBucket(UsersBucketName)
	if err!=nil{
		panic(err)
	}

	defer repo.Close()

	server,err:=cmd.NewLoginServer(cmd.WithRepository(repo), cmd.WithSignKey(signKey))
	if err!=nil{
		panic(err)
	}

	server.Start()
}