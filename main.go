package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"

	_ "github.com/lib/pq"
	"github.com/vamshikrishna209/bank/api"
	db "github.com/vamshikrishna209/bank/db/sqlc"
	"github.com/vamshikrishna209/bank/db/util"
	_ "github.com/vamshikrishna209/bank/doc/statik"
	grpcapi "github.com/vamshikrishna209/bank/grpc_api"
	"github.com/vamshikrishna209/bank/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connection, err := sql.Open(config.DbDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	ctx := context.Background()

	store := db.NewStore(connection)
	go runGatewayServer(config, store, ctx)
	runGrpcServer(config, store, ctx)

}

func runGrpcServer(config util.Config, store *db.Store, ctx context.Context) {
	server, err := grpcapi.NewServer(config, store, ctx)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	log.Printf("start grpc ar %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Println("failed to serve: ", err)
	}
}

func runGinServer(config util.Config, store *db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Server failed to start", err)
		panic("Server failed to start")
	}
}

func runGatewayServer(config util.Config, store *db.Store, ctx context.Context) {
	server, err := grpcapi.NewServer(config, store, ctx)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register bank handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// fs := http.FileServer(http.Dir("./doc/swagger"))
	// mux.Handle("/swagger/", http.StripPrefix(
	// 	"/swagger/", fs,
	// ))

	staticFS, err := fs.New()
	if err != nil {
		log.Println("sttaik failed: ", err)
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(staticFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	log.Printf("start gateway at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Println("failed to serve: ", err)
	}
}
