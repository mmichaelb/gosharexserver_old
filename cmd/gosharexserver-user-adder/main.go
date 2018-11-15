package main

import (
	"flag"
	"fmt"
	"github.com/mmichaelb/gosharexserver/internal/pkg/user"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	hostname       = flag.String("hostname", "localhost", "hostname of the MongoDB server to connect to")
	connectTimeout = flag.Duration("connect_timeout", time.Second*4, "connect timeout to wait until cancelling the connection process to the MongoDB server")
	authDb         = flag.String("auth_db", "", "database to authenticate with")
	authUser       = flag.String("auth_user", "", "user to authenticate with")
	authPassword   = flag.String("auth_password", "", "password to authenticate with")
	dbName         = flag.String("db_name", "gosharexserver", "database name")

	username        = flag.String("username", "", "username of the new user (required)")
	password        = flag.String("password", "", "password of the new user (required)")
	authTokenLength = flag.Int("auth_token_length", 20, "length of the authorization token to generate")
)

func main() {
	log.SetOutput(os.Stderr)
	flag.Parse()
	validateFlag("username", username)
	validateFlag("password", password)
	// set length of each generated authorization token
	viper.SetDefault("webserver.authorization_token_length", *authTokenLength)
	info := &mgo.DialInfo{
		Addrs:    []string{*hostname},
		Timeout:  *connectTimeout,
		Database: *authDb,
		Username: *authUser,
		Password: *authPassword,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		log.Fatalf("Could not dial MongoDB server: %v", err)
	}
	database := session.DB(*dbName)
	userManager := &user.Manager{
		Database: database,
	}
	if err := userManager.InitializeCollection(); err != nil {
		log.Fatalf("Could not initialize user collection: %v", err)
	}
	newUser := userManager.GetNewUserInstance(*username)
	uuid, err := newUser.CreateNewEntry([]byte(*password))
	if err != nil {
		log.Fatalf("Could not create new user entry: %v", err)
	}
	token, err := newUser.RegenerateAuthorizationToken()
	if err != nil {
		log.Fatalf("Could not regenerate authorization token: %v", err)
	}
	fmt.Println("------------------------------------------------------------------------")
	fmt.Printf("| %-8s -> %-56s |\n", "username", *username)
	fmt.Printf("| %-8s -> %-56s |\n", "uuid", uuid.String())
	fmt.Printf("| %-8s -> %-56s |\n", "token", token.String())
	fmt.Println("------------------------------------------------------------------------")
}

func validateFlag(name string, flagValue *string) {
	if *flagValue == "" {
		flag.PrintDefaults()
		log.Fatalf("required flagValue %s not set", strconv.Quote(name))
	}
}
