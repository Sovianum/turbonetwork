package nodeservice

import (
	"github.com/Sovianum/turbonetwork/nodeservice/adapters"
	"github.com/Sovianum/turbonetwork/pb"
)

var nodeDescriptionList = []*pb.NodeDescription{
	{
		NodeType: adapters.PressureLossNodeType,
	},
}
