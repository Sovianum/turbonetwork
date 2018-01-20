package server

import (
	"fmt"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"github.com/Sovianum/turbonetwork/nodeservice/server/factories"
	"golang.org/x/net/context"
	"runtime/debug"
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

func (s *gteServer) CreateNodes(c context.Context, r *pb.CreateRequest) (resp *pb.ModifyResponse, e error) {
	defer func() {
		if r := recover(); r != nil {
			resp = GetModifyErrResponse(fmt.Sprintf("%v, %s", r, debug.Stack()), InternalError)
		}
	}()

	responseItems := make([]*pb.ModifyResponse_UnitResponse, len(r.Items))

	for i, item := range r.Items {
		constructor, err := s.constructorFactory.GetConstructor(item.NodeType)
		if err != nil {
			responseItems[i] = GetModifyErrResponseItem(err.Error(), NotFound)
			continue
		}

		node, nodeErr := constructor(item.Data)
		if nodeErr != nil {
			responseItems[i] = GetModifyErrResponseItem(nodeErr.Error(), InternalError)
			continue
		}
		node.SetName(item.NodeName)

		id, idErr := s.nodeStorage.Add(factories.NewTypedNode(node, item.NodeType))
		if idErr != nil {
			responseItems[i] = GetModifyErrResponseItem(idErr.Error(), InternalError)
			continue
		}

		responseItems[i] = GetModifySuccessResponseItem(id)
	}

	return GetModifySuccessResponse(responseItems), nil
}

func (s *gteServer) UpdateNodes(c context.Context, r *pb.UpdateRequest) (resp *pb.ModifyResponse, e error) {
	defer func() {
		if r := recover(); r != nil {
			resp = GetModifyErrResponse(fmt.Sprintf("%v, %s", r, debug.Stack()), InternalError)
		}
	}()

	responseItems := make([]*pb.ModifyResponse_UnitResponse, len(r.Items))

	for i, item := range r.Items {
		updater, err := s.updaterFactory.GetUpdater(item.Identifier.NodeType)
		if err != nil {
			responseItems[i] = GetModifyErrResponseItem(err.Error(), NotFound)
			continue
		}

		updateErr := updater(item.Data)
		if updateErr != nil {
			responseItems[i] = GetModifyErrResponseItem(updateErr.Error(), NotFound)
			continue
		}
		responseItems[i] = GetModifySuccessResponseItem(item.Identifier)
	}

	return GetModifySuccessResponse(responseItems), nil
}

func (s *gteServer) DeleteNodes(c context.Context, ids *pb.Identifiers) (resp *pb.ModifyResponse, e error) {
	defer func() {
		if r := recover(); r != nil {
			resp = GetModifyErrResponse(fmt.Sprintf("%v, %s", r, debug.Stack()), InternalError)
		}
	}()

	responseItems := make([]*pb.ModifyResponse_UnitResponse, len(ids.Ids))

	for i, id := range ids.Ids {
		if err := s.nodeStorage.Drop(id); err != nil {
			responseItems[i] = GetModifyErrResponseItem(err.Error(), NotFound)
		} else {
			responseItems[i] = GetModifySuccessResponseItem(id)
		}
	}

	return GetModifySuccessResponse(responseItems), nil
}

func (s *gteServer) GetNodes(c context.Context, r *pb.Identifiers) (resp *pb.StateResponse, e error) {
	defer func() {
		if r := recover(); r != nil {
			resp = GetStateErrResponse(fmt.Sprintf("%v, %s", r, debug.Stack()), InternalError)
		}
	}()

	responseItems := make([]*pb.StateResponse_UnitResponse, len(r.Ids))
	for i, item := range r.Ids {
		stateGetter, err := s.stateGetterFactory.GetStateGetter(item.NodeType)
		if err != nil {
			responseItems[i] = GetStateErrResponseItem(err.Error(), NotFound)
			continue
		}

		node, nodeErr := s.nodeStorage.Get(item)
		if nodeErr != nil {
			responseItems[i] = GetStateErrResponseItem(nodeErr.Error(), NotFound)
			continue
		}

		state, stateErr := stateGetter(node)
		if stateErr != nil {
			responseItems[i] = GetStateErrResponseItem(stateErr.Error(), InternalError)
			continue
		}

		responseItems[i] = GetStateSuccessResponseItem(item, state)
	}

	return GetStateSuccessResponse(responseItems), nil
}

func (s *gteServer) Process(c context.Context, r *pb.Identifiers) (resp *pb.ModifyResponse, error error) {
	defer func() {
		if r := recover(); r != nil {
			resp = GetModifyErrResponse(fmt.Sprintf("%v, %s", r, debug.Stack()), InternalError)
			return
		}
	}()

	responseItems := make([]*pb.ModifyResponse_UnitResponse, len(r.Ids))

	for i, item := range r.Ids {
		node, nodeErr := s.nodeStorage.Get(item)
		if nodeErr != nil {
			responseItems[i] = GetModifyErrResponseItem(nodeErr.Error(), NotFound)
			continue
		}

		err := node.Node.Process()
		if err != nil {
			responseItems[i] = GetModifyErrResponseItem(err.Error(), InternalError)
			continue
		}

		responseItems[i] = GetModifySuccessResponseItem(item)
	}

	return GetModifySuccessResponse(responseItems), nil
}

func (s *gteServer) Link(c context.Context, r *pb.LinkRequest) (resp *pb.ModifyResponse, e error) {
	defer func() {
		if r := recover(); r != nil {
			resp = GetModifyErrResponse(fmt.Sprintf("%v, %s", r, debug.Stack()), InternalError)
			return
		}
	}()

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

	responseItems := make([]*pb.ModifyResponse_UnitResponse, len(r.Items))
	for i, item := range r.Items {
		port1, portErr1 := portExtractor(item.Id1)
		if portErr1 != nil {
			responseItems[i] = GetModifyErrResponseItem(portErr1.Error(), NotFound)
		}

		port2, portErr2 := portExtractor(item.Id2)
		if port2 != nil {
			responseItems[i] = GetModifyErrResponseItem(portErr2.Error(), NotFound)
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
		responseItems[i] = GetModifySuccessResponseItem(item.Id1.NodeIdentifier, item.Id2.NodeIdentifier)
	}

	return GetModifySuccessResponse(responseItems), nil
}

func (s *gteServer) GetDescription(context.Context, *pb.Empty) (*pb.ServiceDescription, error) {
	return &pb.ServiceDescription{
		Description: "gte_service",
		Nodes:       nodeDescriptionList,
	}, nil
}
