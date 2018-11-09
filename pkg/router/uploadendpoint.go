package router

import (
	"encoding/json"
	"github.com/mmichaelb/gosharexserver/pkg/storage"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

const (
	maximumMemoryBytes = 1 << 20 // 1 MB maximum in memory
	receiveBufferSize  = 1 << 20 // 1 MB maximum in memory here too
	multipartFormName  = "file"
)

// handleUpload is the endpoint which handles new file upload requests.
func (shareXRouter *ShareXRouter) handleUpload(writer http.ResponseWriter, request *http.Request) {
	authorized, uid :=shareXRouter.checkAuthorization(request, writer)
	if !authorized {
		return
	}
	var err error
	// parse multipart form file and if something goes wrong return an internal server error response code
	if err = request.ParseMultipartForm(maximumMemoryBytes); err != nil {
		shareXRouter.sendInternalError(writer, "parsing multiform of file upload", err)
		return
	}
	var file multipart.File
	// parse filename and mime type from multipart header
	var multipartFileHeader *multipart.FileHeader
	if file, multipartFileHeader, err = request.FormFile(multipartFormName); err != nil {
		shareXRouter.sendInternalError(writer, "resolving file details of file upload", err)
		return
	}
	// instantiate new entry from the given values
	fileName := multipartFileHeader.Filename
	mimeType := multipartFileHeader.Header.Get(contentTypeHeader)
	entry := &storage.Entry{
		Author:      uid,
		Filename:    fileName,
		ContentType: mimeType,
		UploadDate:  time.Now(),
	}
	var fileWriter io.WriteCloser
	// store entry
	if fileWriter, err = shareXRouter.Storage.Store(entry); err != nil {
		shareXRouter.sendInternalError(writer, "storing new file entry", err)
		return
	}
	// write file data to the returned writer
	defer func() {
		if err := fileWriter.Close(); err != nil {
			log.Printf("There was an error while closing the file writer, %T: %+v", err, err)
		}
	}()
	total, err := writeFile(file, fileWriter)
	if err != nil {
		shareXRouter.sendInternalError(writer, "writing file data to new entry", err)
		return
	}
	log.Printf("Created entry %v (%v bytes)\n", entry.ID, total)
	// send json response
	response := Response{entry.CallReference, entry.DeleteReference}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		shareXRouter.sendInternalError(writer, "creation of the json response message", err)
		return
	}
	// set content type header to application/json
	writer.Header().Set("Content-Type", "application/json")
	// write the above created json message to the client
	writer.Write([]byte(jsonResponse))
}

// writeFile writes the received uploaded data to the provided writer by the stored entry
func writeFile(file multipart.File, fileWriter io.WriteCloser) (int64, error) {
	// count total byte amount
	var total int64
	// do not stop iterating until no more bytes are available
	for {
		buffer := make([]byte, receiveBufferSize)
		bytesRead, err := file.Read(buffer)
		total += int64(bytesRead)
		if bytesRead == 0 {
			break
		} else if err != nil {
			return -1, err
		} else {
			fileWriter.Write(buffer[:bytesRead])
		}
	}
	return total, nil
}

// Response holds all required data to respond to an upload.
type Response struct {
	CallReference   string `json:"call_reference"`
	DeleteReference string `json:"delete_reference"`
}
