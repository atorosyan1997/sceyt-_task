package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"sceyt_task/internal/data"
	"strings"
)

// MiddlewareValidateUser validates the user in the request
func (u *UserHandler) MiddlewareValidateUser(ctx *gin.Context) {
	ctx.Set("Content-Type", "application/json")

	u.logger.Debug("user json received")

	user := &data.User{}
	err := ctx.ShouldBindBodyWith(user, binding.JSON)
	if err != nil {
		u.logger.Error("deserialization of user json failed", "error", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		_ = data.ToJSON(&GenericResponse{Status: false, Message: err.Error()}, ctx.Writer)
		return
	}
	// validate the user
	errs := u.validator.Validate(user)
	if len(errs) != 0 {
		u.logger.Error("validation of user json failed", "error", errs)
		ctx.AbortWithStatus(http.StatusBadRequest)
		_ = data.ToJSON(&GenericResponse{Status: false, Message: strings.Join(errs.Errors(), ",")}, ctx.Writer)
		return
	}

	// add the user to the context
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), UserKey{}, *user))

	// call the next handler
	ctx.Next()

}
