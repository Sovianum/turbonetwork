package server

import (
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"fmt"
)

func GetUpdateRequest(ids []*pb.NodeIdentifier, args []map[string]float64) (*pb.UpdateRequest, error) {
	if len(ids) != len(args) {
		return nil, fmt.Errorf("length of arguments are not equal")
	}

	result := &pb.UpdateRequest{
		Items:make([]*pb.UpdateRequest_UnitRequest, len(ids)),
	}
	for i := range ids {
		result.Items[i] = &pb.UpdateRequest_UnitRequest{
			Identifier:ids[i],
			Data:&pb.RequestData{
				DKwargs:args[i],
			},
		}
	}
	return result, nil
}

func GetCreateRequest(nodeNames, nodeTypes []string, args []map[string]float64) (*pb.CreateRequest, error) {
	if len(nodeNames) != len(nodeTypes) || len(nodeNames) != len(args) {
		return nil, fmt.Errorf("length of arguments are not equal")
	}

	result := &pb.CreateRequest{
		Items:make([]*pb.CreateRequest_UnitRequest, len(nodeNames)),
	}
	for i := range nodeNames {
		result.Items[i] = &pb.CreateRequest_UnitRequest{
			NodeType:nodeTypes[i],
			NodeName:nodeNames[i],
			Data:&pb.RequestData{
				DKwargs:args[i],
			},
		}
	}
	return result, nil
}
