package server

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"github.com/Sovianum/turbonetwork/nodeservice/server/factories"
	"golang.org/x/net/context"
)

func NewGTEServer() pb.NodeServiceServer {
	return &gteServer{
		nodeStorage:        NewMapNodeStorage(),
		constructorFactory: factories.NewConstructorFactory(),
		updaterFactory:     factories.NewUpdaterFactory(),
		stateGetterFactory: factories.NewStateGetterFactory(),
		portGetterFactory:  factories.NewPortGetterFactory(),
	}
}

type gteServer struct {
	nodeStorage NodeStorage

	constructorFactory factories.ConstructorFactory
	updaterFactory     factories.UpdaterFactory
	stateGetterFactory factories.StateGetterFactory
	portGetterFactory  factories.PortGetterFactory
}

func (s *gteServer) CreateNodes(c context.Context, r *pb.CreateRequest) (*pb.ModifyResponse, error) {
	responseItems := make([]*pb.ModifyResponse_UnitResponse, len(r.Items))

	for i, item := range r.Items {
		constructor, err := s.constructorFactory.GetConstructor(item.NodeType)

		if err != nil {
			responseItems[i] = s.getModifyErrResponseItem(err, NotFound)
			continue
		}

		node, nodeErr := constructor(item.Data)
		if nodeErr != nil {
			responseItems[i] = s.getModifyErrResponseItem(nodeErr, NotFound)
			continue
		}
		node.SetName(item.NodeName)

		id, idErr := s.nodeStorage.Add(factories.NewTypedNode(node, item.NodeType))
		if idErr != nil {
			responseItems[i] = s.getModifyErrResponseItem(idErr, NotFound)
			continue
		}

		item := s.getModifySuccessResponseItem()
		item.Identifier = id
		responseItems[i] = item
	}

	return &pb.ModifyResponse{Items: responseItems}, nil
}

func (s *gteServer) UpdateNodes(c context.Context, r *pb.UpdateRequest) (*pb.ModifyResponse, error) {
	responseItems := make([]*pb.ModifyResponse_UnitResponse, len(r.Items))

	for i, item := range r.Items {
		updater, err := s.updaterFactory.GetUpdater(item.Identifier.NodeType)
		if err != nil {
			responseItems[i] = s.getModifyErrResponseItem(err, NotFound)
			continue
		}

		updateErr := updater(item.Data)
		if updateErr != nil {
			responseItems[i] = s.getModifyErrResponseItem(updateErr, NotFound)
			continue
		}
		responseItems[i] = s.getModifySuccessResponseItem()
	}

	return &pb.ModifyResponse{Items: responseItems}, nil
}

func (s *gteServer) DeleteNodes(c context.Context, ids *pb.Identifiers) (*pb.ModifyResponse, error) {
	responseItems := make([]*pb.ModifyResponse_UnitResponse, len(ids.Ids))

	for i, id := range ids.Ids {
		if err := s.nodeStorage.Drop(id); err != nil {
			responseItems[i] = s.getModifyErrResponseItem(err, NotFound)
		} else {
			item := s.getModifySuccessResponseItem()
			item.Identifier = id
			responseItems[i] = item
		}
	}

	return &pb.ModifyResponse{Items: responseItems}, nil
}

func (s *gteServer) GetNodes(c context.Context, r *pb.Identifiers) (*pb.StateResponse, error) {
	responseItems := make([]*pb.StateResponse_UnitResponse, len(r.Ids))

	for i, item := range r.Ids {
		stateGetter, err := s.stateGetterFactory.GetStateGetter(item.NodeType)
		if err != nil {
			responseItems[i] = s.getStateErrResponseItem(err, NotFound)
			continue
		}

		node, nodeErr := s.nodeStorage.Get(item)
		if nodeErr != nil {
			responseItems[i] = s.getStateErrResponseItem(err, NotFound)
			continue
		}

		state, stateErr := stateGetter(node)
		if stateErr != nil {
			responseItems[i] = s.getStateErrResponseItem(stateErr, NotFound)
			continue
		}

		responseItems[i] = &pb.StateResponse_UnitResponse{
			Identifier: item,
			Base:       s.getBaseSuccessResponseItem(),
			State:      state,
		}
	}

	return &pb.StateResponse{Items: responseItems}, nil
}

func (s *gteServer) Process(c context.Context, r *pb.Identifiers) (resp *pb.BaseResponseList, error error) {
	defer func() {
		if r := recover(); r != nil {
			resp = &pb.BaseResponseList{
				Base:s.getBaseSuccessResponseItem(),
			}
			error = nil
			return
		}
	}()

	responseItems := make([]*pb.BaseResponse, len(r.Ids))

	for i, item := range r.Ids {
		node, nodeErr := s.nodeStorage.Get(item)
		if nodeErr != nil {
			responseItems[i] = s.getBaseErrResponseItem(nodeErr, NotFound)
			continue
		}
		err := node.Node.Process()

		if err != nil {
			responseItems[i] = s.getBaseErrResponseItem(err, InternalError)
			continue
		}

		responseItems[i] = s.getBaseSuccessResponseItem()
	}

	return &pb.BaseResponseList{
		Base:s.getBaseSuccessResponseItem(),
		Items: responseItems,
	}, nil
}

func (s *gteServer) Link(c context.Context, r *pb.LinkRequest) (*pb.BaseResponseList, error) {
	portExtractor := func(portIdentifier *pb.PortIdentifier) (graph.Port, error) {
		node, nodeErr := s.nodeStorage.Get(portIdentifier.NodeIdentifier)
		if nodeErr != nil {
			return nil, nodeErr
		}

		getter, getterErr := s.portGetterFactory.GetPortGetter(portIdentifier.NodeIdentifier.NodeType)
		if getterErr != nil {
			return nil, getterErr
		}

		port, portErr := getter(node, portIdentifier.PortTag)
		if portErr != nil {
			return nil, portErr
		}

		return port, nil
	}

	responseItems := make([]*pb.BaseResponse, len(r.Items))
	for i, item := range r.Items {
		port1, portErr1 := portExtractor(item.Id1)
		if portErr1 != nil {
			responseItems[i] = s.getBaseErrResponseItem(portErr1, NotFound)
		}

		port2, portErr2 := portExtractor(item.Id2)
		if port2 != nil {
			responseItems[i] = s.getBaseErrResponseItem(portErr2, NotFound)
		}

		switch item.LinkType {
		case pb.LinkRequest_UnitRequest_WEAK_FIRST:
			graph.Link(graph.NewWeakPort(port1), port2)
		case pb.LinkRequest_UnitRequest_WEAK_SECOND:
			graph.Link(port1, graph.NewWeakPort(port2))
		case pb.LinkRequest_UnitRequest_WEAK_BOTH:
			graph.Link(graph.NewWeakPort(port1), graph.NewWeakPort(port2))
		default:
			graph.Link(port1, port2)
		}
		responseItems[i] = s.getBaseSuccessResponseItem()
	}

	return &pb.BaseResponseList{Items: responseItems}, nil
}

func (s *gteServer) GetDescription(context.Context, *pb.Empty) (*pb.ServiceDescription, error) {
	return &pb.ServiceDescription{
		Description: "gte_service",
		Nodes:       nodeDescriptionList,
	}, nil
}

func (s *gteServer) getStateErrResponseItem(err error, status int32) *pb.StateResponse_UnitResponse {
	return &pb.StateResponse_UnitResponse{
		Base: s.getBaseErrResponseItem(err, status),
	}
}

func (s *gteServer) getModifyErrResponseItem(err error, status int32) *pb.ModifyResponse_UnitResponse {
	return &pb.ModifyResponse_UnitResponse{
		Base: s.getBaseErrResponseItem(err, status),
	}
}

func (s *gteServer) getBaseErrResponseItem(err error, status int32) *pb.BaseResponse {
	return &pb.BaseResponse{
		Status:      status,
		Description: err.Error(),
	}
}

func (s *gteServer) getModifySuccessResponseItem() *pb.ModifyResponse_UnitResponse {
	return &pb.ModifyResponse_UnitResponse{
		Base: s.getBaseSuccessResponseItem(),
	}
}

func (s *gteServer) getBaseSuccessResponseItem() *pb.BaseResponse {
	return &pb.BaseResponse{
		Status:      OK,
		Description: "ok",
	}
}
