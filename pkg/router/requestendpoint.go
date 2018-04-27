package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mmichaelb/gosharexserver/pkg/storage"
	"net/http"
	"strconv"
	"strings"
)

const (
	callReferenceVar       = "callreference"
	dispositionHeader      = "Content-Disposition"
	dispositionValueFormat = "%v; filename=\"%v\""
)

// handleRequest is the endpoint which handles incoming file requests via link. It uses the var with the key stored in
// callReferenceVar to resolve the database entry.
func (shareXRouter *ShareXRouter) handleRequest(writer http.ResponseWriter, request *http.Request) {
	callReference, ok := mux.Vars(request)[callReferenceVar]
	if !ok {
		http.Error(writer, "400 the client sent a bad request", http.StatusBadRequest)
		return
	}
	// resolve the remote entry and check if it could be found
	entry, err := shareXRouter.Storage.Request(callReference)
	if err == storage.ErrEntryNotFound {
		http.NotFound(writer, request)
		return
	} else if err != nil {
		shareXRouter.sendInternalError(writer, fmt.Sprintf("requesting entry with call reference %v",
			strconv.Quote(callReference)), err)
		return
	}
	// make sure that the reader gets closed after sending the data
	defer entry.Reader.Close()
	// send disposition header
	var dispositionType string
	for _, entryMimeType := range shareXRouter.WhitelistedContentTypes {
		if strings.EqualFold(entryMimeType, entry.ContentType) {
			dispositionType = "inline"
		}
	}
	if len(dispositionType) == 0 {
		dispositionType = "attachment"
	}
	writer.Header().Set(dispositionHeader, fmt.Sprintf(dispositionValueFormat, dispositionType, entry.Filename))
	// set content type header
	writer.Header().Set(contentTypeHeader, entry.ContentType)
	// write file data from the opened reader to the remote client
	http.ServeContent(writer, request, "", entry.UploadDate, entry.Reader)
}
