package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/mmichaelb/gosharexserver/internal/gosharexserver"
	"github.com/mmichaelb/gosharexserver/internal/gosharexserver/config"
	"github.com/mmichaelb/gosharexserver/pkg/router"
	"github.com/mmichaelb/gosharexserver/pkg/storage"
	"github.com/mmichaelb/gosharexserver/pkg/storage/storages"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// general information about the application
var applicationName = "{application_name}"
var version = "{version}"
var branch = "{branch}"
var commit = "{commit}"
var author = "mmichaelb"

var configFilepath = flag.String(
	"config", "./config.toml", "The filepath to the configuration file used by the ShareX server.")

func main() {
	// parse flags
	flag.Parse()
	// main start process
	log.Printf("Starting %v %v (%v/%v) by %v...\n", applicationName, version, branch, commit, author)
	// load main configuration
	log.Printf("Loading configuration file from %s...\n", strconv.Quote(*configFilepath))
	if err := config.LoadMainConfig(*configFilepath); err != nil {
		log.Fatalf("Could not load configuration from file, %T: %v\n", err, err)
	}
	log.Printf("Successfully loaded %d configuration keys.\n", len(viper.AllKeys()))
	// setup default mux router
	muxRouter := mux.NewRouter()
	var fileStorage storage.FileStorage
	session := connectToMongoDB()
	// use MongoStorage per default
	fileStorage = &storages.MongoStorage{
		Database:        session.DB(viper.GetString("mongodb.db")),
		GridFSPrefix:    viper.GetString("mongodb.gridfs_prefix"),
		GridFSChunkSize: viper.GetInt("mongodb.gridfs_chunk_size"),
	}
	// initialization via interface method Initialize of the file storage instance
	log.Println("Initializing file storage...")
	if err := fileStorage.Initialize(); err != nil {
		log.Fatalf("There was an error while initializing the storage: %v\n", err)
	}
	log.Println("Done with storage initialization! Continuing with the binding of the ShareX muxRouter...")
	// bind ShareXRouter to previously initialized mux muxRouter
	shareXRouter := &router.ShareXRouter{
		Storage:                 fileStorage,
		WhitelistedContentTypes: viper.GetStringSlice("webserver.whitelisted_content_types"),
		AuthorizationToken:      viper.GetString("webserver.authorization_token"),
	}
	// bind ShareX server handler to existing mux muxRouter
	shareXRouter.WrapHandler(muxRouter.PathPrefix("/").Subrouter())
	var handler http.Handler
	// check if a reverse proxy is used
	if reverseProxyHeader := viper.GetString("webserver.reverse_proxy_header"); reverseProxyHeader != "" {
		handler = gosharexserver.WrapRouterToReverseProxyRouter(muxRouter, reverseProxyHeader)
	} else {
		handler = muxRouter
	}
	webserverAddress := viper.GetString("webserver.address")
	httpServer := http.Server{
		Addr:    webserverAddress,
		Handler: handler,
	}
	log.Printf("Running ShareX server in background and listening for connections on %s. "+
		"Press CTRL-C to stop the application.\n", strconv.Quote(webserverAddress))
	var closed bool
	go func() {
		// run http server in background
		if err := httpServer.ListenAndServe(); err != nil && !closed {
			panic(err)
		}
	}()
	// wait for stop signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Println("Shutting down ShareX server and MongoDB connection...")
	closed = true
	if err := httpServer.Close(); err != nil {
		log.Printf("There was an error while closing the ShareX server, %T: %v\n", err, err)
	}
	if err := fileStorage.Close(); err != nil {
		log.Printf("There was an error while closing the ShareX file storage, %T: %v\n", err, err)
	}
	session.Close()
	log.Println("Thank you for using the ShareX server. Bye!")
}

func connectToMongoDB() *mgo.Session {
	dialInfo := parseDialInfoFromConfig()
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB server: %v\n", err)
	}
	return session
}

// parseDialInfoFromConfig parses the dial information used to connect to the MongoDB server.
func parseDialInfoFromConfig() *mgo.DialInfo {
	return &mgo.DialInfo{
		Addrs:    []string{viper.GetString("mongodb.address")},
		Timeout:  viper.GetDuration("mongodb.connect_timeout"),
		Source:   viper.GetString("mongodb.auth_db"),
		Username: viper.GetString("mongodb.auth_user"),
		Password: viper.GetString("mongodb.auth_passwd"),
	}
}
