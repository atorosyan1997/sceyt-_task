package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mbndr/figlet4go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"my-bank-service/internal/config"
	"my-bank-service/internal/data"
	"my-bank-service/internal/handler"
	"my-bank-service/internal/reposytory"
	"my-bank-service/internal/service"
	data2 "my-bank-service/internal/validation"
	"my-bank-service/pkg/logging"
	"my-bank-service/pkg/session"
	//_ "my-bank-service/docs"
)

var sf *session.SessionFactory

// Run initializes whole application
func Run(address string, port string) {
	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
		figlet4go.ColorGreen,
	}
	renderStr, _ := ascii.RenderOpts("API-Server!", options)
	fmt.Print(renderStr)

	logConfig := config.GetLogConfiguration()
	logging.Init(logConfig)
	logger := logging.GetLogger()
	logger.Info("logger initialized")

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	sessionRef := sf.GetSession()

	con := config.NewConfigurations(logger)

	// userRepository contains all the methods that interact with DB to perform CURD operations for user.
	userRepository := reposytory.NewUserRepository(sessionRef, logger)
	/*user :=&data.User{
		ID: "asdasdaa",
		Email:"asdasdas",
		Username: "username",
		FirstName: "Arman",
		LastName: "Torosyan",
		Password: "12345678",
		TokenHash: "asdasd",
	}
	userRepository.Create(user)*/

	/*user ,_:= userRepository.GetUserByUserName("admin")
	fmt.Println(user)*/

	authRepository := reposytory.NewAuthRepository(sessionRef, logger)

	authRepository.FetchAuth(&data.AuthDetails{
		UserId:   "asdadas",
		AuthUuid: "uwehrskadsaj",
	})

	// validation contains all the methods that are need to validate the user json in request
	validator := data2.NewValidation()

	// authService contains all methods that help in authorizing a user request
	authService := service.NewAuthService(logger, con, authRepository)

	// AuthHandler encapsulates all the services related to user
	authHandler := handlers.NewAuthHandler(logger, con, validator, userRepository, authService)

	authHandler.Routes(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := router.Run(fmt.Sprintf("%s:%v", address, port))
	if err != nil {
		logger.Error(err)
	}

}

func init() {
	var err error
	sf, err = session.NewSessionFactory(config.Driver)
	if err != nil {
		log.Panic(err)
	}
}
