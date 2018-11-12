package httpfileserver

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/kataras/iris"
	rootRouter "github.com/mmichaelb/gosharexserver/internal/gosharexserver/router"
	"github.com/mmichaelb/gosharexserver/internal/pkg/storage"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// upload is the endpoint of the upload path where new files are uploaded to by their ShareX clients.
func (router *Router) upload(ctx iris.Context) {
	var err error
	// parse multipart form file and if something goes wrong return an internal server error response code
	if err = ctx.Request().ParseMultipartForm(maximumMemoryBytes); err != nil {
		if err == http.ErrNotMultipart {
			rootRouter.HttpError(ctx, http.StatusBadRequest, "no multipart request")
			return
		}
		ctx.Application().Logger().Warnf("could not parse multipart form of file upload: %s", err.Error())
		rootRouter.InternalServerError(ctx)
		return
	}
	var file multipart.File
	// parse filename and mime type from multipart header
	var multipartFileHeader *multipart.FileHeader
	if file, multipartFileHeader, err = ctx.Request().FormFile(multipartFormName); err != nil {
		ctx.Application().Logger().Warnf("could not resolve file details of file upload: %s", err.Error())
		rootRouter.InternalServerError(ctx)
		return
	}
	// instantiate new entry from the given values
	fileName := multipartFileHeader.Filename
	mimeType := multipartFileHeader.Header.Get("Content-Type")
	uid := ctx.Values().Get(uuidContextValueKey).(uuid.UUID)
	entry := &storage.Entry{
		Author:      uid,
		Filename:    fileName,
		ContentType: mimeType,
		UploadDate:  time.Now(),
	}
	var fileWriter io.WriteCloser
	// store entry
	if fileWriter, err = router.Storage.Store(entry); err != nil {
		ctx.Application().Logger().Warnf("could not store new file entry: %s", err.Error())
		rootRouter.InternalServerError(ctx)
		return
	}
	// write file data to the returned writer
	defer func() {
		if err := fileWriter.Close(); err != nil {
			ctx.Application().Logger().Warnf("there was an error while closing the file writer: %s", err.Error())
		}
	}()
	total, err := writeFile(file, fileWriter)
	if err != nil {
		ctx.Application().Logger().Warnf("could not write file data to new entry: %s", err.Error())
		rootRouter.InternalServerError(ctx)
		return
	}
	ctx.Application().Logger().Printf("created entry %v (%d bytes)", entry.ID, total)
	// send json response
	response := Response{entry.CallReference, entry.DeleteReference}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		ctx.Application().Logger().Warnf("could not create json response message: %s", err.Error())
		rootRouter.InternalServerError(ctx)
		return
	}
	// set content type header to application/json
	ctx.ResponseWriter().Header().Set("Content-Type", "application/json")
	// write the above created json message to the client
	if _, err = ctx.ResponseWriter().Write([]byte(jsonResponse)); err != nil {
		fmt.Printf("there was an error while sending back the json response: %e", err)
	}
}

// writeFile writes the received uploaded data to the provided writer by the stored entry
func writeFile(file multipart.File, fileWriter io.WriteCloser) (int64, error) {
	// count total byte amount
	var total int64
	// do not stop iterating until no more bytes are available
	for {
		buffer := make([]byte, maximumBufferSize)
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
