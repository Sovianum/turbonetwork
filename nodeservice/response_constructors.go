package nodeservice

import "github.com/Sovianum/turbonetwork/pb"

func getStateSuccessResponse(items []*pb.NodeStateResponse_UnitResponse) *pb.NodeStateResponse {
	return &pb.NodeStateResponse{
		Base:  getBaseSuccessResponseItem(),
		Items: items,
	}
}

func getStateErrResponse(msg string, status int32) *pb.NodeStateResponse {
	return &pb.NodeStateResponse{
		Base: getBaseErrResponseItem(msg, status),
	}
}

func getStateSuccessResponseItem(id *pb.NodeIdentifier, state *pb.NodeState) *pb.NodeStateResponse_UnitResponse {
	return &pb.NodeStateResponse_UnitResponse{
		Base:       getBaseSuccessResponseItem(),
		Identifier: id,
		State:      state,
	}
}

func getStateErrResponseItem(msg string, status int32) *pb.NodeStateResponse_UnitResponse {
	return &pb.NodeStateResponse_UnitResponse{
		Base: getBaseErrResponseItem(msg, status),
	}
}

func getModifyErrResponse(msg string, status int32) *pb.NodeModifyResponse {
	return &pb.NodeModifyResponse{
		Base: getBaseErrResponseItem(msg, status),
	}
}

func getModifySuccessResponse(items []*pb.NodeModifyResponse_UnitResponse) *pb.NodeModifyResponse {
	return &pb.NodeModifyResponse{
		Base:  getBaseSuccessResponseItem(),
		Items: items,
	}
}

func getModifyErrResponseItem(msg string, status int32) *pb.NodeModifyResponse_UnitResponse {
	return &pb.NodeModifyResponse_UnitResponse{
		Base: getBaseErrResponseItem(msg, status),
	}
}

func getModifySuccessResponseItem(ids ...*pb.NodeIdentifier) *pb.NodeModifyResponse_UnitResponse {
	return &pb.NodeModifyResponse_UnitResponse{
		Identifiers: ids,
		Base:        getBaseSuccessResponseItem(),
	}
}

func getBaseErrResponseItem(msg string, status int32) *pb.BaseResponse {
	return &pb.BaseResponse{
		Status:      status,
		Description: msg,
	}
}

func getBaseSuccessResponseItem() *pb.BaseResponse {
	return &pb.BaseResponse{
		Status:      ok,
		Description: "ok",
	}
}
