package grpcapi

import (
	"context"
	"fmt"

	logger "bitbucket.org/ninjafactory/ucc-federal-go-gin-base-module/v2/log"
	logwriter "bitbucket.org/ninjafactory/ucc-federal-go-gin-base-module/v2/log/logwriter"
	utils "bitbucket.org/ninjafactory/ucc-federal-go-gin-base-module/v2/utils"
	db "github.com/vamshikrishna209/bank/db/sqlc"
	"github.com/vamshikrishna209/bank/db/util"
	"github.com/vamshikrishna209/bank/pb"
	"github.com/vamshikrishna209/bank/token"
)

type Server struct {
	pb.UnimplementedBankServer
	config     util.Config
	store      *db.Store
	tokenMaker token.Maker
	log        *logger.Logger
}

func NewServer(config util.Config, s *db.Store, ctx context.Context) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	consoleLogger := logwriter.NewConsoleWriter(logger.HostParams{})
	lMux := logger.NewDefaultLogMux(consoleLogger)
	log := logger.NewLogger(ctx, &logger.Config{
		HostParams: logger.HostParams{
			Version:     "v1",
			ServiceName: "bank",
		},
		LogLevel: utils.GetEnvInt("LOG_LEVEL", 7),
	}, "bank", lMux, nil)
	fmt.Println("test")
	log.Debug(context.Background(), "Hitesh is a good boy", "")
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("cannot create token maker")
	}
	server := &Server{
		config:     config,
		store:      s,
		tokenMaker: tokenMaker,
		log:        log,
	}
	return server, nil
}
