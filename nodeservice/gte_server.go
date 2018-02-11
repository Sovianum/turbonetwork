package nodeservice

import (
	"fmt"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/nodeservice/adapters"
	"github.com/Sovianum/turbonetwork/pb"
	"golang.org/x/net/context"
	"runtime/debug"
)

// NewGTEServer constructs gteServer which implements NodeService interface
func NewGTEServer(factory adapters.NodeAdapterFactory) pb.NodeServiceServer {
	return &gteServer{
		nodeStorage: NewMapNodeStorage(),
		factory:     factory,
	}
}

type gteServer struct {
	nodeStorage NodeStorage
	factory     adapters.NodeAdapterFactory
}

func (s *gteServer) CreateNodes(c context.Context, r *pb.NodeCreateRequest) (resp *pb.NodeModifyResponse, e error) {
	defer func() {
		if r := recover(); r != nil {
			resp = getModifyErrResponse(fmt.Sprintf("%v, %s", r, debug.Stack()), internalError)
		}
	}()

	responseItems := make([]*pb.NodeModifyResponse_UnitResponse, len(r.Items))

	for i, item := range r.Items {
		adapter, err := s.factory.GetAdapter(item.NodeType)
		if err != nil {
			responseItems[i] = getModifyErrResponseItem(err.Error(), notFound)
			continue
		}

		node, nodeErr := adapter.Create(item.Data)
		if nodeErr != nil {
			responseItems[i] = getModifyErrResponseItem(nodeErr.Error(), internalError)
			continue
		}
		node.SetName(item.NodeName)

		id, idErr := s.nodeStorage.Add(adapters.NewTypedNode(node, item.NodeType))
		if idErr != nil {
			responseItems[i] = getModifyErrResponseItem(idErr.Error(), internalError)
			continue
		}

		responseItems[i] = getModifySuccessResponseItem(id)
	}

	return getModifySuccessResponse(responseItems), nil
}

func (s *gteServer) UpdateNodes(c context.Context, r *pb.NodeUpdateRequest) (resp *pb.NodeModifyResponse, e error) {
	defer func() {
		if r := recover(); r != nil {
			resp = getModifyErrResponse(fmt.Sprintf("%v, %s", r, debug.Stack()), internalError)
		}
	}()

	responseItems := make([]*pb.NodeModifyResponse_UnitResponse, len(r.Items))

	for i, item := range r.Items {
		node, nodeErr := s.nodeStorage.Get(item.Identifier)
		if nodeErr != nil {
			responseItems[i] = getModifyErrResponseItem(nodeErr.Error(), notFound)
			continue
		}

		adapter, err := s.factory.GetAdapter(item.Identifier.NodeType)
		if err != nil {
			responseItems[i] = getModifyErrResponseItem(err.Error(), notFound)
			continue
		}

		updateErr := adapter.Update(node.Node, item.Data)
		if updateErr != nil {
			responseItems[i] = getModifyErrResponseItem(updateErr.Error(), internalError)
			continue
		}
		responseItems[i] = getModifySuccessResponseItem(item.Identifier)
	}

	return getModifySuccessResponse(responseItems), nil
}

func (s *gteServer) DeleteNodes(c context.Context, ids *pb.Identifiers) (resp *pb.NodeModifyResponse, e error) {
	defer func() {
		if r := recover(); r != nil {
			resp = getModifyErrResponse(fmt.Sprintf("%v, %s", r, debug.Stack()), internalError)
		}
	}()

	responseItems := make([]*pb.NodeModifyResponse_UnitResponse, len(ids.Ids))

	for i, id := range ids.Ids {
		if err := s.nodeStorage.Drop(id); err != nil {
			responseItems[i] = getModifyErrResponseItem(err.Error(), notFound)
		} else {
			responseItems[i] = getModifySuccessResponseItem(id)
		}
	}

	return getModifySuccessResponse(responseItems), nil
}

func (s *gteServer) GetNodes(c context.Context, r *pb.NodeStateRequest) (resp *pb.NodeStateResponse, e error) {
	defer func() {
		if r := recover(); r != nil {
			resp = getStateErrResponse(fmt.Sprintf("%v, %s", r, debug.Stack()), internalError)
		}
	}()

	responseItems := make([]*pb.NodeStateResponse_UnitResponse, len(r.Items))
	for i, item := range r.Items {
		adapter, err := s.factory.GetAdapter(item.Identifier.NodeType)
		if err != nil {
			responseItems[i] = getStateErrResponseItem(err.Error(), notFound)
			continue
		}

		node, nodeErr := s.nodeStorage.Get(item.Identifier)
		if nodeErr != nil {
			responseItems[i] = getStateErrResponseItem(nodeErr.Error(), notFound)
			continue
		}

		state, stateErr := adapter.GetState(node.Node, item.RequiredFields)
		if stateErr != nil {
			responseItems[i] = getStateErrResponseItem(stateErr.Error(), internalError)
			continue
		}

		responseItems[i] = getStateSuccessResponseItem(item.Identifier, state)
	}

	return getStateSuccessResponse(responseItems), nil
}

func (s *gteServer) Process(c context.Context, r *pb.Identifiers) (resp *pb.NodeModifyResponse, error error) {
	defer func() {
		if r := recover(); r != nil {
			resp = getModifyErrResponse(fmt.Sprintf("%v, %s", r, debug.Stack()), internalError)
			return
		}
	}()

	responseItems := make([]*pb.NodeModifyResponse_UnitResponse, len(r.Ids))

	for i, item := range r.Ids {
		node, nodeErr := s.nodeStorage.Get(item)
		if nodeErr != nil {
			responseItems[i] = getModifyErrResponseItem(nodeErr.Error(), notFound)
			continue
		}

		err := node.Node.Process()
		if err != nil {
			responseItems[i] = getModifyErrResponseItem(err.Error(), internalError)
			continue
		}

		responseItems[i] = getModifySuccessResponseItem(item)
	}

	return getModifySuccessResponse(responseItems), nil
}

func (s *gteServer) Link(c context.Context, r *pb.LinkRequest) (resp *pb.NodeModifyResponse, e error) {
	defer func() {
		if r := recover(); r != nil {
			resp = getModifyErrResponse(fmt.Sprintf("%v, %s", r, debug.Stack()), internalError)
			return
		}
	}()

	portExtractor := func(portIdentifier *pb.PortIdentifier) (graph.Port, error) {
		node, nodeErr := s.nodeStorage.Get(portIdentifier.NodeIdentifier)
		if nodeErr != nil {
			return nil, nodeErr
		}

		adapter, err := s.factory.GetAdapter(portIdentifier.NodeIdentifier.NodeType)
		//getter, err := s.portGetterFactory.GetPortGetter(portIdentifier.NodeIdentifier.NodeType)
		if err != nil {
			return nil, err
		}

		port, portErr := adapter.GetPort(portIdentifier.PortTag, node.Node)
		if portErr != nil {
			return nil, portErr
		}

		return port, nil
	}

	responseItems := make([]*pb.NodeModifyResponse_UnitResponse, len(r.Items))
	for i, item := range r.Items {
		port1, portErr1 := portExtractor(item.Id1)
		if portErr1 != nil {
			responseItems[i] = getModifyErrResponseItem(portErr1.Error(), notFound)
			continue
		}

		port2, portErr2 := portExtractor(item.Id2)
		if portErr2 != nil {
			responseItems[i] = getModifyErrResponseItem(portErr2.Error(), notFound)
			continue
		}

		switch item.LinkType {
		case pb.LinkType_WEAK_FIRST:
			graph.Link(graph.NewWeakPort(port1), port2)
		case pb.LinkType_WEAK_SECOND:
			graph.Link(port1, graph.NewWeakPort(port2))
		case pb.LinkType_WEAK_BOTH:
			graph.Link(graph.NewWeakPort(port1), graph.NewWeakPort(port2))
		default:
			graph.Link(port1, port2)
		}
		responseItems[i] = getModifySuccessResponseItem(item.Id1.NodeIdentifier, item.Id2.NodeIdentifier)
	}

	return getModifySuccessResponse(responseItems), nil
}

func (s *gteServer) GetDescription(context.Context, *pb.Empty) (*pb.ServiceDescription, error) {
	return &pb.ServiceDescription{
		Description: "gte_service",
		Nodes:       nodeDescriptionList,
	}, nil
}
