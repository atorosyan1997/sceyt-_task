package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sceyt_task/internal/cache"
	"sceyt_task/internal/config"
	"sceyt_task/internal/data"
	"sceyt_task/internal/repository"
	"sceyt_task/internal/validation"
	"sceyt_task/pkg/logging"
	"strings"
)

var (
	ErrUserNotFound = fmt.Sprintf("No user account exists with given email. Please sign in first")
	PgNoRowsMsg     = "no rows in result set"
)

// UserKey is used as a key for storing the User object in context at middleware
type UserKey struct{}

// UserHandler wraps instances needed to perform operations on user object
type UserHandler struct {
	logger    logging.Logger
	validator *validation.Validation
	repo      repository.UserRepository
	userCache cache.UserCache
}

// NewUserHandler returns a new UserHandler instance
func NewUserHandler(l logging.Logger, v *validation.Validation, r repository.UserRepository, cache cache.UserCache) *UserHandler {
	return &UserHandler{
		logger:    l,
		validator: v,
		repo:      r,
		userCache: cache,
	}
}

func (u *UserHandler) Routes(engine *gin.Engine) {
	user := engine.Group(config.GroupPath)
	{
		user.Use(u.MiddlewareValidateUser)
		user.POST(config.AddPath, u.Add)
		user.POST(config.UpdatePath, u.Update)
		user.POST(config.SearchPath, u.Search)
		user.DELETE(config.DeletePath, u.Delete)

	}

}

// GenericResponse is the format of our response
type GenericResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Errors []string `json:"errors"`
}

type SearchResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username" validate:"required"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type AddResponse struct {
	Username string `json:"username" validate:"required"`
}

// Add adds a user to the database
// @Summary add
// @Tags user
// @Description add user
// @ID user-add
// @Accept json
// @Produce json
// @Param input body swagger.UserAddUpdate true "add user"
// @Success 200 {integer} integer 1
// @Failure 400,404,500 {integer} integer 2
// @Router /add/ [post]
func (u *UserHandler) Add(ctx *gin.Context) {
	ctx.Set("Content-Type", "application/json")
	reqUser := ctx.Request.Context().Value(UserKey{}).(data.User)

	err := u.repo.Create(&reqUser)
	if err != nil {
		u.logger.Error("error while adding user", "error", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		_ = data.ToJSON(&GenericResponse{Status: false, Message: "error while adding user"}, ctx.Writer)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
	_ = data.ToJSON(&GenericResponse{
		Status:  true,
		Message: "user added successfully",
		Data:    &AddResponse{Username: reqUser.Username},
	}, ctx.Writer)
}

// Update updates user
// @Summary update
// @Tags user
// @Description update user
// @ID user-update
// @Accept json
// @Produce json
// @Param input body swagger.UserAddUpdate true "update user"
// @Success 200 {integer} integer 1
// @Failure 400,404,500 {integer} integer 2
// @Router /update/ [post]
func (u *UserHandler) Update(ctx *gin.Context) {
	ctx.Set("Content-Type", "application/json")
	reqUser := ctx.Request.Context().Value(UserKey{}).(data.User)

	_ = u.userCache.Del(reqUser.Username)
	err := u.repo.Update(&reqUser)
	if err != nil {
		u.logger.Error("error while updating user", "error", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		_ = data.ToJSON(&GenericResponse{Status: false, Message: "error while updating user"}, ctx.Writer)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
	_ = data.ToJSON(&GenericResponse{
		Status:  true,
		Message: "user updated successfully",
		Data:    &AddResponse{Username: reqUser.Username},
	}, ctx.Writer)
}

// Delete deletes user
// @Summary delete
// @Tags user
// @Description delete user
// @ID user-delete
// @Accept json
// @Produce json
// @Param input body swagger.UserSearchDelete true "delete user"
// @Success 200 {integer} integer 1
// @Failure 400,404,500 {integer} integer 2
// @Router /delete/ [delete]
func (u *UserHandler) Delete(ctx *gin.Context) {
	ctx.Set("Content-Type", "application/json")
	reqUser := ctx.Request.Context().Value(UserKey{}).(data.User)

	_ = u.userCache.Del(reqUser.Username)
	err := u.repo.Delete(reqUser.Username)
	if err != nil {
		u.logger.Error("error when deleting user", "error", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		_ = data.ToJSON(&GenericResponse{Status: false, Message: "error when deleting user"}, ctx.Writer)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
	_ = data.ToJSON(&GenericResponse{
		Status:  true,
		Message: "user deleted successfully",
		Data:    &AddResponse{Username: reqUser.Username},
	}, ctx.Writer)
}

// Search get user by username
// @Summary Search
// @Tags user
// @Description user search
// @ID user-search
// @Accept json
// @Produce json
// @Param input body swagger.UserSearchDelete true "user search"
// @Success 200 {integer} integer 1
// @Failure 400,404,500 {integer} integer 2
// @Router /search/ [post]
func (u *UserHandler) Search(ctx *gin.Context) {
	ctx.Set("Content-Type", "application/json")

	reqUser := ctx.Request.Context().Value(UserKey{}).(data.User)

	user, err := u.userCache.Get(reqUser.Username)
	if err != nil {
		u.logger.Error("error fetching the user", "error", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		_ = data.ToJSON(&GenericResponse{Status: false, Message: "An error occured while getting the user from Redis. Please try again later."}, ctx.Writer)
		return
	}
	if user == nil {
		user, err = u.repo.GetUserByUserName(reqUser.Username)
		if err != nil {
			u.logger.Error("error fetching the user", "error", err)
			errMsg := err.Error()
			if strings.Contains(errMsg, PgNoRowsMsg) {
				ctx.AbortWithStatus(http.StatusBadRequest)
				_ = data.ToJSON(&GenericResponse{Status: false, Message: ErrUserNotFound}, ctx.Writer)
			} else {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				_ = data.ToJSON(&GenericResponse{Status: false, Message: "Unable to retrieve user from database.Please try again later"}, ctx.Writer)
			}
			return
		}
		err = u.userCache.Set(user.Username, user)
		if err != nil {
			u.logger.Error("error while adding user to Redis", "error", err)
		}
	}

	ctx.AbortWithStatus(http.StatusOK)
	_ = data.ToJSON(&GenericResponse{
		Status:  true,
		Message: "user found successfully",
		Data:    &SearchResponse{ID: user.ID, Username: user.Username, FirstName: user.FirstName, LastName: user.LastName},
	}, ctx.Writer)
}
