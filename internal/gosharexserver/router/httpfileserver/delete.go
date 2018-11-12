package httpfileserver

import (
	"fmt"
	"github.com/kataras/iris"
	rootRouter "github.com/mmichaelb/gosharexserver/internal/gosharexserver/router"
	"strconv"
)

// delete is the endpoint which handles incoming delete requests of entries.
func (router *Router) delete(ctx iris.Context) {
	//get the delete reference
	deleteReference := ctx.Params().Get(deleteReferenceVar)
	//delete the entry
	err := router.Storage.Delete(deleteReference)
	if err != nil {
		ctx.Application().Logger().Warnf("could not delete entry with call reference %s: %s", deleteReference, err)
		rootRouter.InternalServerError(ctx)
		return
	}
	if _, err := ctx.JSON(iris.Map{"message": fmt.Sprintf("entry with reference %s successfully deleted.", strconv.Quote(deleteReference))}); err != nil {
		ctx.Application().Logger().Warnf("could not send deletion response to %s: %s", ctx.RemoteAddr(), err.Error())
		rootRouter.InternalServerError(ctx)
		return
	}
}
