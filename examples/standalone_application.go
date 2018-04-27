package main

import (
	"github.com/gorilla/mux"
	"github.com/mmichaelb/gosharexserver/pkg/router"
	"github.com/mmichaelb/gosharexserver/pkg/storage/storages"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

func main() {
	// initialize main gorilla/mux router
	mainRouter := mux.NewRouter()
	mainRouter.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Hello there, this is my custom ShareX server application."))
	})
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	// use MongoDB storage type
	fileStorage := &storages.MongoStorage{
		// MongoDB database
		Database: session.DB("gosharexserver"),
		// MongoDB GridFS prefix
		GridFSPrefix: "uploads",
		// MongoDB GridFS chunk size
		GridFSChunkSize: 255000,
	}
	if err := fileStorage.Initialize(); err != nil {
		log.Println("Could not initialize file storage!")
		panic(err)
	}
	// setup ShareX router
	shareXRouter := router.ShareXRouter{
		Storage:                 fileStorage,
		WhitelistedContentTypes: []string{"image/png", "image/jpeg"},
	}
	// add ShareX handler to main router
	shareXRouter.WrapHandler(mainRouter.PathPrefix("/sharex/").Subrouter())
	httpServer := http.Server{
		Handler: mainRouter,        // use the gorilla/mux router as the http handler
		Addr:    "localhost:10711", // bind to local loop-back interface on port 8080
	}
	// run server and log occurring errors
	log.Fatal(httpServer.ListenAndServe())
}
