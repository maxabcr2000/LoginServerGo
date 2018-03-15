package main

import (
	// "os"
	// "fmt"
	repository "github.com/maxabcr2000/LoginServerGo/pkg/repository"
)

const (
	DBName = "LoginServer.db"
	SignKeyVarName = "TOKEN_SIGN_KEY"
)

func main(){
	// signKey, ok := os.LookupEnv(SignKeyVarName)
	// if !ok {
	// 	fmt.Printf("Failed to get TOKEN_SIGN_KEY from environmental settings.")
	// 	return
	// }
	signKey:=`-----BEGIN RSA PRIVATE KEY-----
MIIEoAIBAAKCAQEAhj36cVEAhI0qZeKpJBdLdwYIVzf1xPoMHcx7A+KJStWQYjWp
a8Oe3o4SFypDH91T/sfNRmvFsWdRoq8ytVE9cSzk9g51zUhegtd4OP+QDQ1P1IMX
j6QdWpp3vaxTcKORoHWgrovacahlV7+T/LhBzWRmJomx6vs/0Ar01mIHEmVF70hM
xmto7XXum3xGMoN+JIJgnuWr/nLIe8H5bKxBf2q9maMQ+lBe+Z/kuWo0gK7spoZg
TZ8PFc5nE/aLdNnEzHXccSm3o+miSslxNxLQVOy/vhEsy/dl5NS1P+41UmzHE385
sq4tbL7MM4hf/CFvVsp7cCdt+q4CmK5m+oldpwIBJQKCAQBh9eBE2j6l6njikKTw
0rqHSZdiPZeralveZMh4dOC6EXcy6ONqUHrZwZ5/m2FqR5BSIxlBOeNKGvZh9XgH
xalgbN5uXZQ72t0vC+/yPfN/JWPElwNa+zgD2IDkWyghw3gbJWdqnWgNwBKZ+oC/
VgaBH8Ap9ccFazxnYfDvd/dSm96FZkbW9tzn5B53bPcufvUVhA0p0wxxohbsJOzP
Aqt6Lj3QXPr35/B1q2xMQUBfYz2zSBx+oxjfycMIuSMXZL7Bs3BmSIQh3x8EN4ON
G6eetDAoDCHAiYDAkwUVzuuA/5Wd9Vrm1oVsU57gioR7b56OnxoW9knZPQA6TWX2
qaOVAoGBAMfhjbUShPREhyVhIHzYYAFfRaZvmuSbtLod5hwwnPtdsjszi019P1IT
XPRx1tWAkWWRDVi6a4jZvS1PSKiFKJeKXmjXWARr+K/x3wFx0f0g7GKd8kGPVcG6
sTxkwY+fPBSJJhOUKYQzn1M2StuzqAWEeeA/s4BoR4MZMkq6o9WDAoGBAKvuoDag
SE9rJ/TPlw5TT5U1+Yuzm1eEpwtmGwp/MYqNQECnX3UM5iutTotHX+DcNrdFt2VE
efJeul3oIfU37jrQWnFGisVCadrOZXkgU9EOqp/fwnEh0j+UBKrGxT0L5WpAG+Jw
aCtF3EnvdjcHnhVHwm30irfHhK6X4qN8GYINAoGAW9ZV3ZLjJB98XUhNMnFA3gkt
1tlVAUCfJRSn/x1BNTjjL/UWfYyStwH6RsyFTVa/WC3jiaHCtH+3yLW5mYlKAHAB
3SSr5lsfPBUvHFbj8NfGjizk7bCPL39Kg4g9Qf0NxD8DqCF62+Bd9czWSUS0+527
dN7/clls9wuTpctgCCcCgYBTpIVLAd39H1+L2pyDgnm+Ne8FsVJ9mnrHRm4FG0iI
l75kpHOguirJI+EFfKsvHwXcnnTlKDtTUK2xJNI/8bIqczLq/7khsZtH2gfaD7oS
rTBpc8ZgjP5/y3fkYLL7GotIcjcUNq83pXIWIXfHvoRBtDUS5atEpkCMSePg1JbJ
pQKBgD/T8eYduojbQBe8cEnAm0YNrV7mvtGVyCEgbuipMN3tgQh1QM0KCW4DxExZ
ZKZOzMJXJOXcSqRrQpRTHNCCaPeQYUZ0pGNDGfLfdX/Wn+OcmLcBQeFQA2LHdKKh
ayfgA9uN2HtMyqn1EdYuc7ry5Fsq5F30yfzLDmQEu+bVHXP+
-----END RSA PRIVATE KEY-----`

	repo,err:= repository.CreateBoltRepository(DBName)
	if err!=nil{
		panic(err)
	}

	defer repo.Close()

	server,err:=NewLoginServer(WithRepository(repo), WithSignKey(signKey))
	if err!=nil{
		panic(err)
	}

	server.Start()
}