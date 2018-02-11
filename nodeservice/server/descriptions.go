package server

import (
	"github.com/Sovianum/turbonetwork/pb"
	"github.com/Sovianum/turbonetwork/nodeservice/server/adapters"
)

var nodeDescriptionList = []*pb.NodeDescription{
	{
		NodeType:        adapters.PressureLossNodeType,
	},
}
