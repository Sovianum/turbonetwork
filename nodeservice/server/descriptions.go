package server

import (
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"github.com/Sovianum/turbonetwork/nodeservice/server/adapters"
)

var nodeDescriptionList = []*pb.NodeDescription{
	{
		NodeType:        adapters.PressureLossNodeType,
	},
}
