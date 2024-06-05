package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/aprole/ip-whitelist/pb"
	"github.com/aprole/ip-whitelist/utils"
	"github.com/oschwald/geoip2-golang"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type GrpcIPWhitelistServer struct {
	pb.UnimplementedIPWhitelistServer
	db *geoip2.Reader
}

func NewGrpcIPWhitelistServer(db *geoip2.Reader) *GrpcIPWhitelistServer {
	return &GrpcIPWhitelistServer{db: db}
}

func (server *GrpcIPWhitelistServer) CheckIP(ctx context.Context, req *pb.CheckIPRequest) (*pb.CheckIPResponse, error) {
	ip := net.ParseIP(req.Ip)
	if ip == nil {
		err := fmt.Errorf("invalid ip address: %s", req.Ip)
		log.Println(err)
		return nil, err
	}

	accepted, record, err := utils.CheckIP(ip, req.AllowedCountries, server.db)
	if err != nil {
		err := fmt.Errorf("server error: %s", err.Error())
		log.Println(err)
		return nil, err
	}

	utils.LogResult("gRPC", accepted, ip, record)

	return &pb.CheckIPResponse{
		Accepted:       &wrapperspb.BoolValue{Value: accepted},
		Ip:             req.Ip,
		CountryIsoCode: record.Country.IsoCode,
		CountryName:    record.Country.Names["en"],
	}, nil
}
