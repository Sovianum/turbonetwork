package server

import "github.com/Sovianum/turbonetwork/pb"

func GetStateSuccessResponse(items []*pb.NodeStateResponse_UnitResponse) *pb.NodeStateResponse {
	return &pb.NodeStateResponse{
		Base:  GetBaseSuccessResponseItem(),
		Items: items,
	}
}

func GetStateErrResponse(msg string, status int32) *pb.NodeStateResponse {
	return &pb.NodeStateResponse{
		Base: GetBaseErrResponseItem(msg, status),
	}
}

func GetStateSuccessResponseItem(id *pb.NodeIdentifier, state *pb.NodeState) *pb.NodeStateResponse_UnitResponse {
	return &pb.NodeStateResponse_UnitResponse{
		Base:       GetBaseSuccessResponseItem(),
		Identifier: id,
		State:      state,
	}
}

func GetStateErrResponseItem(msg string, status int32) *pb.NodeStateResponse_UnitResponse {
	return &pb.NodeStateResponse_UnitResponse{
		Base: GetBaseErrResponseItem(msg, status),
	}
}

func GetModifyErrResponse(msg string, status int32) *pb.NodeModifyResponse {
	return &pb.NodeModifyResponse{
		Base: GetBaseErrResponseItem(msg, status),
	}
}

func GetModifySuccessResponse(items []*pb.NodeModifyResponse_UnitResponse) *pb.NodeModifyResponse {
	return &pb.NodeModifyResponse{
		Base:  GetBaseSuccessResponseItem(),
		Items: items,
	}
}

func GetModifyErrResponseItem(msg string, status int32) *pb.NodeModifyResponse_UnitResponse {
	return &pb.NodeModifyResponse_UnitResponse{
		Base: GetBaseErrResponseItem(msg, status),
	}
}

func GetModifySuccessResponseItem(ids ...*pb.NodeIdentifier) *pb.NodeModifyResponse_UnitResponse {
	return &pb.NodeModifyResponse_UnitResponse{
		Identifiers: ids,
		Base:        GetBaseSuccessResponseItem(),
	}
}

func GetBaseErrResponseItem(msg string, status int32) *pb.BaseResponse {
	return &pb.BaseResponse{
		Status:      status,
		Description: msg,
	}
}

func GetBaseSuccessResponseItem() *pb.BaseResponse {
	return &pb.BaseResponse{
		Status:      OK,
		Description: "ok",
	}
}
