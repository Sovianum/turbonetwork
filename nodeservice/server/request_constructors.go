package server

import (
	"github.com/Sovianum/turbonetwork/pb"
	"fmt"
)

func GetUpdateRequest(ids []*pb.NodeIdentifier, args []map[string]float64) (*pb.NodeUpdateRequest, error) {
	if len(ids) != len(args) {
		return nil, fmt.Errorf("length of arguments are not equal")
	}

	result := &pb.NodeUpdateRequest{
		Items:make([]*pb.NodeUpdateRequest_UnitRequest, len(ids)),
	}
	for i := range ids {
		result.Items[i] = &pb.NodeUpdateRequest_UnitRequest{
			Identifier:ids[i],
			Data:&pb.RequestData{
				DKwargs:args[i],
			},
		}
	}
	return result, nil
}

func GetCreateRequest(nodeNames, nodeTypes []string, args []map[string]float64) (*pb.NodeCreateRequest, error) {
	if len(nodeNames) != len(nodeTypes) || len(nodeNames) != len(args) {
		return nil, fmt.Errorf("length of arguments are not equal")
	}

	result := &pb.NodeCreateRequest{
		Items:make([]*pb.NodeCreateRequest_UnitRequest, len(nodeNames)),
	}
	for i := range nodeNames {
		result.Items[i] = &pb.NodeCreateRequest_UnitRequest{
			NodeType:nodeTypes[i],
			NodeName:nodeNames[i],
			Data:&pb.RequestData{
				DKwargs:args[i],
			},
		}
	}
	return result, nil
}
