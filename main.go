package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2"
)

var (
	AdminGroup = "Admin"
)
var validator *auth0.JWTValidator

// LoadPublicKey loads a public key from PEM/DER-encoded data.
func LoadPublicKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	// Try to load SubjectPublicKeyInfo
	pub, err0 := x509.ParsePKIXPublicKey(input)
	if err0 == nil {
		return pub, nil
	}

	cert, err1 := x509.ParseCertificate(input)
	if err1 == nil {
		return cert.PublicKey, nil
	}

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s' and '%s'", err0, err1)
}

func init() {
	//Creates a configuration with the Auth0 information
	data, err := ioutil.ReadFile("./dev-dgoly5h6.pem")
	if err != nil {
		panic("Impossible to read key form disk")
	}

	secret, err := LoadPublicKey(data)
	if err != nil {
		panic("Invalid provided key")
	}
	secretProvider := auth0.NewKeyProvider(secret)
	configuration := auth0.NewConfiguration(secretProvider, []string{"casseur_flutter"}, "https://dev-dgoly5h6.eu.auth0.com/", jose.RS256)
	validator = auth0.NewValidator(configuration, nil)
}

func Auth0Groups(wantedGroups ...string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		tok, err := validator.ValidateRequest(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			log.Println("Invalid token:", err)
			return
		}

		claims := map[string]interface{}{}
		err = validator.Claims(c.Request, tok, &claims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			c.Abort()
			log.Println("Invalid claims:", err)
			return
		}

		fmt.Println(claims)
		//metadata, okMetadata := claims["app_metadata"].(map[string]interface{})
		//authorization, okAuthorization := metadata["authorization"].(map[string]interface{})
		//groups, hasGroups := authorization["groups"].([]interface{})
		//fmt.Println(okMetadata, okAuthorization, hasGroups)
		//fmt.Println(groups)
		//if !okMetadata || !okAuthorization || !hasGroups {
		//	c.JSON(http.StatusUnauthorized, gin.H{"error": "need more privileges"})
		//	c.Abort()
		//	log.Println("Need more provileges")
		//	return
		//}
		c.Next()
	})
}

func main() {
	r := gin.Default()
	r.GET("/", Auth0Groups(AdminGroup), func(c *gin.Context) {
		c.JSON(200, gin.H{"test": "ok"})
	})
	r.Run(":9090")
}
