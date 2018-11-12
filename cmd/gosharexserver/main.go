package main

import (
	"context"
	"flag"
	"github.com/kataras/iris"
	"github.com/mmichaelb/gosharexserver/internal/gosharexserver/router/httpfileserver"
	"github.com/mmichaelb/gosharexserver/internal/pkg/config"
	"github.com/mmichaelb/gosharexserver/internal/pkg/storage"
	"github.com/mmichaelb/gosharexserver/internal/pkg/storage/storages"
	"github.com/mmichaelb/gosharexserver/internal/pkg/user"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"log"
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
	app := iris.New()
	var fileStorage storage.FileStorage
	session := connectToMongoDB()
	// use MongoStorage per default
	database := session.DB(viper.GetString("mongodb.db"))
	fileStorage = &storages.MongoStorage{
		Database:        database,
		GridFSPrefix:    viper.GetString("mongodb.gridfs_prefix"),
		GridFSChunkSize: viper.GetInt("mongodb.gridfs_chunk_size"),
	}
	// initialization via interface method Initialize of the file storage instance
	log.Println("Initializing file storage...")
	if err := fileStorage.Initialize(); err != nil {
		log.Fatalf("There was an error while initializing the storage: %v", err)
	}
	// initiate user manager
	userManager := &user.Manager{
		Database: database,
	}
	if err := userManager.InitializeCollection(); err != nil {
		log.Fatalf("There was an error while initializing the user manager: %v", err)
	}
	log.Println("Done with storage initialization! Continuing with the binding of the Iris web server initialization...")
	// bind ShareXRouter to previously initialized mux muxRouter
	fileHttpRouter := &httpfileserver.Router{
		Storage:     fileStorage,
		UserManager: userManager,
	}
	fileHttpRouter.BindToIris(app.Party("/"))
	irisCfg := &iris.Configuration{
		DisableStartupLog:true,
		DisableInterruptHandler:true,
	}
	// check if a reverse proxy is used
	if reverseProxyHeader := viper.GetString("webserver.reverse_proxy_header"); reverseProxyHeader != "" {
		irisCfg.RemoteAddrHeaders = map[string]bool{
			reverseProxyHeader: true,
		}
	}
	webserverAddress := viper.GetString("webserver.address")
	log.Printf("Running ShareX server in background and listening for connections on %s. "+
		"Press CTRL-C to stop the application.\n", strconv.Quote(webserverAddress))
	app.Logger().SetLevel("debug")
	go func() {
		// run http server in background
		if err := app.Run(iris.Addr(webserverAddress), iris.WithConfiguration(*irisCfg)); err != nil && err != iris.ErrServerClosed {
			panic(err)
		}
	}()
	// wait for stop signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Println("Shutting down ShareX server and MongoDB connection...")
	if err := app.Shutdown(context.Background()); err != nil {
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
