package server

import (
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"github.com/Sovianum/turbonetwork/nodeservice/server/factories"
)

var nodeDescriptionList = []*pb.NodeDescription{
	{
		Type:factories.PressureLossNodeType,
		Description:factories.PressureLossNodeType,
	},
}
