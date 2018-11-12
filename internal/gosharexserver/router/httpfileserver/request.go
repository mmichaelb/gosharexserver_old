package httpfileserver

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	rootRouter "github.com/mmichaelb/gosharexserver/internal/gosharexserver/router"
	"github.com/mmichaelb/gosharexserver/internal/pkg/storage"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strings"
)

// request is the endpoint which handles incoming file requests. It uses the var with the key stored in callReferenceVar
// to resolve the database entry.
func (router *Router) request(ctx iris.Context) {
	callReference := ctx.Params().Get(callReferenceVar)
	// resolve the remote entry and check if it could be found
	entry, err := router.Storage.Request(callReference)
	if err == storage.ErrEntryNotFound {
		rootRouter.HttpError(ctx, http.StatusNotFound, fmt.Sprintf("could not find entry %s", callReference))
		return
	} else if err != nil {
		ctx.Application().Logger().Warnf("could not request entry with call reference")
		rootRouter.InternalServerError(ctx)
		return
	}
	// make sure that the reader gets closed after sending the data
	defer func() {
		if err := entry.Reader.Close(); err != nil {
			ctx.Application().Logger().Warnf("there was an error while closing the entry reader: %s", err.Error())
		}
	}()
	if modified, err := ctx.CheckIfModifiedSince(entry.UploadDate); !modified && err == nil {
		ctx.WriteNotModified()
		return
	}
	// send disposition header
	var dispositionType string
	for _, entryMimeType := range viper.GetStringSlice("webserver.whitelisted_content_types") {
		if strings.EqualFold(entryMimeType, entry.ContentType) {
			dispositionType = "inline"
		}
	}
	if len(dispositionType) == 0 {
		dispositionType = "attachment"
	}
	ctx.Header(dispositionHeader, fmt.Sprintf(dispositionValueFormat, dispositionType, entry.Filename))
	// set content type header
	ctx.Header(context.ContentTypeHeaderKey, entry.ContentType)
	// set last modified header
	ctx.SetLastModified(entry.UploadDate)
	// write file data from the opened reader to the remote client
	buf := make([]byte, maximumBufferSize)
	for {
		bytesRead, err := entry.Reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				ctx.Application().Logger().Warnf("could not read entry %s: %s", callReference, err.Error())
				rootRouter.InternalServerError(ctx)
				return
			}
			break
		}
		if _, err = ctx.TryWriteGzip(buf[:bytesRead]); err != nil {
			ctx.Application().Logger().Warnf("could not write entry data of %s: %s", callReference, err.Error())
			rootRouter.InternalServerError(ctx)
			return
		}
	}
	ctx.StatusCode(http.StatusOK)
}
