package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/aprole/ip-whitelist/grpcserver"
	"github.com/aprole/ip-whitelist/handler"
	"github.com/aprole/ip-whitelist/pb"
	"github.com/gorilla/mux"
	"github.com/oschwald/geoip2-golang"
)

func main() {
	logFile, err := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stderr, logFile)
	log.SetOutput(mw)

	db, err := geoip2.Open("/usr/share/GeoIP/GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/check-ip", func(w http.ResponseWriter, r *http.Request) {
		handler.Handler(w, r, db)
	}).Methods("POST")

	httpPort := 8080
	go func() {
		log.Printf("HTTP server started at :%d\n", httpPort)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), r))
	}()

	grpcPort := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterIPWhitelistServer(grpcServer, grpcserver.NewGrpcIPWhitelistServer(db))
	reflection.Register(grpcServer)

	log.Printf("gRPC server started at :%d\n", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
