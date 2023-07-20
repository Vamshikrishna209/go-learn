package grpcapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayAgentHeader = "grpcgateway-user-agent"
	userAgentHeader        = "user-agent"
	xForwardedForHeader    = "x-forwarded-for"
)

type MetaData struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMeta(ctx context.Context) *MetaData {
	mtd := &MetaData{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("md : %+v\n", md)
		if userAgents := md.Get(grpcGatewayAgentHeader); len(userAgents) > 0 {
			mtd.UserAgent = userAgents[0]
		}
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtd.UserAgent = userAgents[0]
		}
		if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
			mtd.ClientIP = clientIPs[0]
		}
	}
	if p, ok := peer.FromContext(ctx); ok {
		mtd.ClientIP = p.Addr.String()
	}
	return mtd
}
