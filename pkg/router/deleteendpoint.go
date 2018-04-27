package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	deleteReferenceVar = "deleteReference"
)

// handleDelete is the endpoint which handles incoming delete requests via link.
func (shareXRouter *ShareXRouter) handleDelete(writer http.ResponseWriter, request *http.Request) {
	//get the delete reference
	deleteReference, ok := mux.Vars(request)[deleteReferenceVar]
	if !ok {
		http.Error(writer, "400 Bad request", http.StatusBadRequest)
		return
	}
	//delete the entry
	err := shareXRouter.Storage.Delete(deleteReference)
	if err != nil {
		shareXRouter.sendInternalError(writer, fmt.Sprintf("deleting entry with call reference %v", strconv.Quote(deleteReference)), err)
		return
	}
	writer.WriteHeader(http.StatusOK)
}
