package server

import "github.com/Sovianum/turbonetwork/nodeservice/pb"

func GetStateSuccessResponse(items []*pb.StateResponse_UnitResponse) *pb.StateResponse {
	return &pb.StateResponse{
		Base:  GetBaseSuccessResponseItem(),
		Items: items,
	}
}

func GetStateErrResponse(msg string, status int32) *pb.StateResponse {
	return &pb.StateResponse{
		Base: GetBaseErrResponseItem(msg, status),
	}
}

func GetStateSuccessResponseItem(id *pb.NodeIdentifier, state *pb.NodeState) *pb.StateResponse_UnitResponse {
	return &pb.StateResponse_UnitResponse{
		Base:       GetBaseSuccessResponseItem(),
		Identifier: id,
		State:      state,
	}
}

func GetStateErrResponseItem(msg string, status int32) *pb.StateResponse_UnitResponse {
	return &pb.StateResponse_UnitResponse{
		Base: GetBaseErrResponseItem(msg, status),
	}
}

func GetModifyErrResponse(msg string, status int32) *pb.ModifyResponse {
	return &pb.ModifyResponse{
		Base: GetBaseErrResponseItem(msg, status),
	}
}

func GetModifySuccessResponse(items []*pb.ModifyResponse_UnitResponse) *pb.ModifyResponse {
	return &pb.ModifyResponse{
		Base:  GetBaseSuccessResponseItem(),
		Items: items,
	}
}

func GetModifyErrResponseItem(msg string, status int32) *pb.ModifyResponse_UnitResponse {
	return &pb.ModifyResponse_UnitResponse{
		Base: GetBaseErrResponseItem(msg, status),
	}
}

func GetModifySuccessResponseItem(id ...*pb.NodeIdentifier) *pb.ModifyResponse_UnitResponse {
	return &pb.ModifyResponse_UnitResponse{
		Identifiers: id,
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
