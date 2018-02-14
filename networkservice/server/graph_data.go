package server

import (
	"github.com/Sovianum/turbonetwork/networkservice/repr"
	"github.com/Sovianum/turbonetwork/pb"
)

type GraphData struct {
	graph           []repr.RepresentationNode
	callOrder       []repr.RepresentationNode
	domainCallOrder []domainCall
}

type domainCall struct {
	server    pb.NodeServiceClient
	nodes     []repr.RepresentationNode
	portLinks []portLink
}

type portLink struct {
	sourcePort remotePort
	destPort   remotePort
}

type remotePort struct {
	server pb.NodeServiceClient
	portID *pb.PortIdentifier
}
