LoginServer for Golang (With BoltDB & JSON Web Token Authentication)
==================================


## Requirements
- Go 1.5 or later.
- [jwt-GO]
- [BoltDB]
- RSA Private Key with PEM encoded (2048 bits)

## Guide for generate required RSA Private Key
1. Download Putty
2. Open PuttyGen (Putty Key Generator)
3. Choose RSA Key Type
4. Press Generate Button
5. Menu => Conversions => Export Open SSH Key => Save the private key file

[jwt-GO]: https://github.com/dgrijalva/jwt-go
[BoltDB]: https://github.com/boltdb/bolt