package router

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mmichaelb/gosharexserver/pkg/storage"
	"github.com/mmichaelb/gosharexserver/pkg/user"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const contentTypeHeader = "Content-Type"

// ShareXRouter represents the main router which serves/accepts files.
type ShareXRouter struct {
	// Storage is an implementation of the Storage interface which is used by the ShareX router.
	Storage storage.FileStorage
	// WhitelistedContentTypes is a slice of content types which will be displayed embed in the browser.
	WhitelistedContentTypes []string
	// UserManager is used for authentication and identification purposes.
	UserManager *user.Manager
}

// WrapHandler wraps the endpoints to the given mux.Router. At the moment this is bound to the usage of gorilla/mux in
// your dependency but in the future this should be generalized. //TODO
func (shareXRouter *ShareXRouter) WrapHandler(router *mux.Router) {
	// register endpoints
	router.Path("/upload").Methods(http.MethodPost).HandlerFunc(shareXRouter.handleUpload)
	router.Path(fmt.Sprintf("/delete/{%v}", deleteReferenceVar)).HandlerFunc(shareXRouter.handleDelete)
	router.Path(fmt.Sprintf("/{%v}", callReferenceVar)).HandlerFunc(shareXRouter.handleRequest)
}

// checkAuthorization returns whether the request is authorized or not.
func (shareXRouter *ShareXRouter) checkAuthorization(request *http.Request, writer http.ResponseWriter) (authorized bool, uid uuid.UUID) {
	rawToken := strings.TrimSpace(request.Header.Get("Authorization"))
	if rawToken == "" {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}
	parsedToken := user.AuthorizationTokenFromString(rawToken)
	if parsedToken == nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var err error
	if uid, err = shareXRouter.UserManager.CheckAuthorizationToken(parsedToken); err == mgo.ErrNotFound {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("there was an error while checking the authorization token (%s): %e", strconv.Quote(rawToken), err) // TODO better error logging
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	return true, uid
}

// sendInternalError generalizes the internal error method.
func (shareXRouter *ShareXRouter) sendInternalError(writer http.ResponseWriter, action string, err error) {
	http.Error(writer, "500 an internal error occurred", http.StatusInternalServerError)
	log.Printf("An error occurred while doing the action \"%v\", %T: %+v\n", action, err, err)
}

// Close stops and closes the ShareX router. It returns an error if something goes wrong.
func (shareXRouter *ShareXRouter) Close() error {
	return shareXRouter.Storage.Close()
}
