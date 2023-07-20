package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/vamshikrishna209/bank/db/sqlc"
	"github.com/vamshikrishna209/bank/db/util"
	"github.com/vamshikrishna209/bank/token"
)

type Server struct {
	config     util.Config
	store      *db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config util.Config, s *db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("cannot create token maker")
	}
	server := &Server{
		config:     config,
		store:      s,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// router := gin.Default()

	// router.POST("/login", server.loginUser)
	server.setupRouter()
	// authRoutes := router.Group("/").Use(authMiddleware(tokenMaker))

	// router.POST("/accounts", server.createAccount)
	// router.GET("/accounts/:id", server.getAccount)
	// router.GET("/accounts", server.listAccounts)
	// router.POST("/transfer", server.transferAmount)
	// router.POST("/user", server.createUser)

	// server.router = router
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/token/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfers", server.transferAmount)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	fmt.Println("error", err)
	return gin.H{"error": err.Error()}
}
