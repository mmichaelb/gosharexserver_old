package httpfileserver

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/google/uuid"
	"github.com/kataras/iris"
	rootRouter "github.com/mmichaelb/gosharexserver/internal/gosharexserver/router"
	"github.com/mmichaelb/gosharexserver/internal/pkg/storage"
	"github.com/mmichaelb/gosharexserver/internal/pkg/storage/storages"
	"github.com/mmichaelb/gosharexserver/internal/pkg/user"
	"strconv"
	"strings"
)

const (
	uuidContextValueKey    = "user-uuid"
	maximumMemoryBytes     = 1 << 20 // 1 MB maximum in memory
	maximumBufferSize      = 1 << 20 // 1 MB maximum in memory here too
	multipartFormName      = "file"
	deleteReferenceVar     = "deleteReference"
	callReferenceVar       = "callreference"
	dispositionHeader      = "Content-Disposition"
	dispositionValueFormat = "%v; filename=\"%v\""
)

// Router contains the basic endpoints to upload, delete and view/request files.
type Router struct {
	// Storage is an implementation of the Storage interface which is used by the ShareX router.
	Storage storage.FileStorage
	// UserManager is used for authentication and identification purposes.
	UserManager *user.Manager
}

// BindToIris binds all endpoints to the given iris party.
func (router *Router) BindToIris(party iris.Party) {
	// register macro to filter
	party.Macros().Get("string").RegisterFunc("reference", router.referenceMacro)
	// register upload and delete endpoints
	updatesParty := party.Party("/updates", router.authorizationTokenHandler)
	updatesParty.Post("/upload", router.upload)
	updatesParty.Get(fmt.Sprintf("/delete/{%s:string reference(%d)}", deleteReferenceVar, storages.DeleteReferenceLength), router.delete)
	// register normal request endpoint
	party.Get(fmt.Sprintf("/{%s:string reference(%d)}", callReferenceVar, storages.CallReferenceLength), router.request)
}

// authorizationTokenHandler checks if the requesting client is authorized to interact with the given route.
func (router *Router) authorizationTokenHandler(ctx iris.Context) {
	rawToken := strings.TrimSpace(ctx.Request().Header.Get("Authorization"))
	if rawToken == "" {
		rootRouter.Unauthorized(ctx)
		return
	}
	parsedToken := user.AuthorizationTokenFromString(rawToken)
	if parsedToken == nil {
		rootRouter.Unauthorized(ctx)
		return
	}
	var uid uuid.UUID
	var err error
	if uid, err = router.UserManager.CheckAuthorizationToken(parsedToken); err == mgo.ErrNotFound {
		rootRouter.Unauthorized(ctx)
		return
	} else if err != nil {
		ctx.Application().Logger().Warnf("could not check user authorization token (%s): %s", strconv.Quote(rawToken), err.Error()) // TODO better error logging
		rootRouter.InternalServerError(ctx)
		return
	}
	// set context value
	ctx.Values().Set(uuidContextValueKey, uid)
	ctx.Next()
}

// referenceMacro returns the macro to validate if the given param is a valid reference.
func (router *Router) referenceMacro(length int) func(string) bool {
	return func(paramValue string) (passed bool) {
		if len(paramValue) != length {
			return
		}
		if !storages.ReferenceRegex.MatchString(paramValue) {
			return
		}
		return true
	}
}
