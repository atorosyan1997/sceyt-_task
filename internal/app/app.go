package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mbndr/figlet4go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	_ "sceyt_task/docs"
	"sceyt_task/internal/cache"
	"sceyt_task/internal/config"
	"sceyt_task/internal/handler"
	"sceyt_task/internal/repository"
	"sceyt_task/internal/validation"
	"sceyt_task/pkg/logging"
	"sceyt_task/pkg/session"
)

var sf *session.SessionFactory

// Run initializes whole application
func Run(address string, port string) {
	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
		figlet4go.ColorGreen,
	}
	renderStr, _ := ascii.RenderOpts("User-Server!", options)
	fmt.Print(renderStr)

	logConfig := config.GetLogConfiguration()
	logging.Init(logConfig)
	logger := logging.GetLogger()
	logger.Info("logger initialized")

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	sessionRef := sf.GetSession()

	// userRepository contains all the methods that interact with DB to perform CURD operations for user.
	userRepository := repository.NewUserRepository(sessionRef, logger)

	// userCache contains all the methods that interact with redis cache
	userCache := cache.NewRedisCache(fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort), config.RedisDb, config.RedisExpires)

	// validation contains all the methods that are need to validate the user json in request
	validator := validation.NewValidation()

	// AuthHandler encapsulates all the services related to user
	authHandler := handler.NewUserHandler(logger, validator, userRepository, userCache)

	authHandler.Routes(router)
	router.GET(config.SwaggerPath, ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := router.Run(fmt.Sprintf("%s:%v", address, port))
	if err != nil {
		logger.Error(err)
	}

}

func init() {
	var err error
	sf, err = session.NewSessionFactory()
	if err != nil {
		log.Panic(err)
	}
}
