package nodeservice

import (
	"fmt"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/nodeservice/adapters"
	"github.com/Sovianum/turbonetwork/nodeservice/mocks"
	"github.com/Sovianum/turbonetwork/pb"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type GTEServerTestSuite struct {
	suite.Suite
	server *gteServer

	storage *mocks.NodeStorageMock
	factory *mocks.NodeAdapterFactoryMock
}

func (s *GTEServerTestSuite) SetupTest() {
	s.server = NewGTEServer(nil).(*gteServer)

	s.server.nodeStorage = mocks.NewNodeStorageMock()
	s.server.factory = mocks.NewNodeAdapterFactoryMock()

	s.storage = s.server.nodeStorage.(*mocks.NodeStorageMock)
	s.factory = s.server.factory.(*mocks.NodeAdapterFactoryMock)
}

func (s *GTEServerTestSuite) TestCreateNodes_Success() {
	s.factory.ExpectResponse(
		mocks.NodeAdapterMock{
			CreateFunc: func(data *pb.RequestData) (graph.Node, error) {
				return graph.NewTestNode(0, 0, true, func() error {
					return nil
				}), nil
			},
		}, nil,
	)

	ids := s.getNodeIdentifiers(1)
	s.storage.ExpectAddResponse(ids.Ids[0], nil)

	req := s.getValidCreateRequest()
	response, err := s.server.CreateNodes(nil, req)

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))
	s.EqualValues(ok, response.Items[0].Base.Status)
	s.EqualValues(1, response.Items[0].Identifiers[0].Id)
}

func (s *GTEServerTestSuite) TestCreateNodes_ConstructorNotFound() {
	e := fmt.Errorf("err not found")
	s.factory.ExpectResponse(nil, e)

	req := s.getValidCreateRequest()
	response, err := s.server.CreateNodes(nil, req)

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))
	s.EqualValues(notFound, response.Items[0].Base.Status)
	s.Equal(e.Error(), response.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestCreateNodes_ConstructorErr() {
	e := fmt.Errorf("err constructor failed")
	s.factory.ExpectResponse(
		mocks.NodeAdapterMock{
			CreateFunc: func(data *pb.RequestData) (graph.Node, error) {
				return nil, e
			},
		}, nil,
	)

	req := s.getValidCreateRequest()
	response, err := s.server.CreateNodes(nil, req)

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))
	s.EqualValues(internalError, response.Items[0].Base.Status)
	s.Equal(e.Error(), response.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestCreateNodes_StorageAddError() {
	s.factory.ExpectResponse(
		mocks.NodeAdapterMock{
			CreateFunc: func(data *pb.RequestData) (graph.Node, error) {
				return graph.NewTestNode(0, 0, true, func() error {
					return nil
				}), nil
			},
		}, nil,
	)

	e := fmt.Errorf("storage error")
	s.storage.ExpectAddResponse(nil, e)

	req := s.getValidCreateRequest()
	response, err := s.server.CreateNodes(nil, req)

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))
	s.EqualValues(internalError, response.Items[0].Base.Status)
	s.EqualValues(e.Error(), response.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestCreateNodes_Panic() {
	msg := "panic msg"
	s.factory.ExpectResponse(
		mocks.NodeAdapterMock{
			CreateFunc: func(data *pb.RequestData) (graph.Node, error) {
				panic(msg)
			},
		}, nil,
	)

	req := s.getValidCreateRequest()
	response, err := s.server.CreateNodes(nil, req)

	s.Require().Nil(err)
	s.Require().Equal(0, len(response.Items))
	s.EqualValues(internalError, response.Base.Status)
	s.True(strings.HasPrefix(response.Base.Description, msg))
}

func (s *GTEServerTestSuite) TestUpdateNodes_Success() {
	s.factory.ExpectResponse(
		mocks.NodeAdapterMock{
			UpdateFunc: func(node graph.Node, data *pb.RequestData) error {
				return nil
			},
		}, nil,
	)

	s.storage.ExpectGetResponse(&adapters.TypedNode{
		NodeType: "test",
		Node: graph.NewTestNode(0, 0, true, func() error {
			return nil
		}),
	}, nil)

	req := s.getValidUpdateRequest()
	response, err := s.server.UpdateNodes(nil, req)

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))
	s.EqualValues(ok, response.Items[0].Base.Status)
	s.EqualValues(1, response.Items[0].Identifiers[0].Id)
}

func (s *GTEServerTestSuite) TestUpdateNodes_NodeNotFound() {
	s.factory.ExpectResponse(
		mocks.NodeAdapterMock{
			UpdateFunc: func(node graph.Node, data *pb.RequestData) error {
				return nil
			},
		}, nil,
	)

	e := fmt.Errorf("node not found")
	s.storage.ExpectGetResponse(nil, e)

	req := s.getValidUpdateRequest()
	response, err := s.server.UpdateNodes(nil, req)

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))
	s.EqualValues(notFound, response.Items[0].Base.Status)
	s.EqualValues(e.Error(), response.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestUpdateNodes_UpdaterNotFound() {
	e := fmt.Errorf("updater not found")
	s.factory.ExpectResponse(
		mocks.NodeAdapterMock{
			UpdateFunc: func(node graph.Node, data *pb.RequestData) error {
				return nil
			},
		}, e,
	)

	s.storage.ExpectGetResponse(&adapters.TypedNode{
		NodeType: "test",
		Node: graph.NewTestNode(0, 0, true, func() error {
			return nil
		}),
	}, nil)

	req := s.getValidUpdateRequest()
	response, err := s.server.UpdateNodes(nil, req)

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))
	s.EqualValues(notFound, response.Items[0].Base.Status)
	s.EqualValues(e.Error(), response.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestUpdateNodes_UpdaterFailed() {
	e := fmt.Errorf("updater failed")
	s.factory.ExpectResponse(
		mocks.NodeAdapterMock{
			UpdateFunc: func(node graph.Node, data *pb.RequestData) error {
				return e
			},
		}, nil,
	)

	s.storage.ExpectGetResponse(&adapters.TypedNode{
		NodeType: "test",
		Node: graph.NewTestNode(0, 0, true, func() error {
			return nil
		}),
	}, nil)

	req := s.getValidUpdateRequest()
	response, err := s.server.UpdateNodes(nil, req)

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))
	s.EqualValues(internalError, response.Items[0].Base.Status)
	s.EqualValues(e.Error(), response.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestDeleteNodes_Success() {
	s.storage.ExpectDropResponse(nil)

	ids := s.getNodeIdentifiers(1)
	response, err := s.server.DeleteNodes(nil, ids)

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))
	s.EqualValues(ok, response.Items[0].Base.Status)
	s.EqualValues(1, response.Items[0].Identifiers[0].Id)
}

func (s *GTEServerTestSuite) TestDelete_NotFound() {
	s.storage.ExpectDropResponse(fmt.Errorf("err"))

	ids := s.getNodeIdentifiers(1)
	response, err := s.server.DeleteNodes(nil, ids)

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))

	s.EqualValues(notFound, response.Items[0].Base.Status)
}

func (s *GTEServerTestSuite) TestGetNodes_Success() {
	s.factory.ExpectResponse(
		mocks.NodeAdapterMock{
			GetStateFunc: func(node graph.Node, requiredFields []string) (*pb.NodeState, error) {
				return &pb.NodeState{}, nil
			},
		}, nil,
	)

	s.storage.ExpectGetResponse(&adapters.TypedNode{
		NodeType: "test",
		Node: graph.NewTestNode(0, 0, true, func() error {
			return nil
		}),
	}, nil)

	response, err := s.server.GetNodesState(nil, s.getValidGetStateRequest())

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))
	s.EqualValues(ok, response.Items[0].Base.Status)
	s.EqualValues(1, response.Items[0].Identifier.Id)
}

func (s *GTEServerTestSuite) TestGetNodes_GetterNotFound() {
	e := fmt.Errorf("getter not found")
	s.factory.ExpectResponse(
		mocks.NodeAdapterMock{
			GetStateFunc: func(node graph.Node, requiredFields []string) (*pb.NodeState, error) {
				return &pb.NodeState{}, nil
			},
		}, e,
	)

	response, err := s.server.GetNodesState(nil, s.getValidGetStateRequest())

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))
	s.EqualValues(notFound, response.Items[0].Base.Status)
	s.EqualValues(e.Error(), response.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestGetNodes_GetterError() {
	e := fmt.Errorf("getter failed")
	s.factory.ExpectResponse(
		mocks.NodeAdapterMock{
			GetStateFunc: func(node graph.Node, requiredFields []string) (*pb.NodeState, error) {
				return nil, e
			},
		}, nil,
	)

	s.storage.ExpectGetResponse(&adapters.TypedNode{
		NodeType: "test",
		Node: graph.NewTestNode(0, 0, true, func() error {
			return nil
		}),
	}, nil)

	response, err := s.server.GetNodesState(nil, s.getValidGetStateRequest())

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))
	s.EqualValues(internalError, response.Items[0].Base.Status)
	s.EqualValues(e.Error(), response.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestGetNodes_StorageError() {
	s.factory.ExpectResponse(
		mocks.NodeAdapterMock{
			GetStateFunc: func(node graph.Node, requiredFields []string) (*pb.NodeState, error) {
				return &pb.NodeState{}, nil
			},
		}, nil,
	)

	e := fmt.Errorf("getter failed")
	s.storage.ExpectGetResponse(nil, e)

	response, err := s.server.GetNodesState(nil, s.getValidGetStateRequest())

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))
	s.EqualValues(notFound, response.Items[0].Base.Status)
	s.EqualValues(e.Error(), response.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestGetNodes_Panic() {
	msg := "panic msg"
	s.factory.ExpectResponse(
		mocks.NodeAdapterMock{
			GetStateFunc: func(node graph.Node, requiredFields []string) (*pb.NodeState, error) {
				panic(msg)
			},
		}, nil,
	)

	s.storage.ExpectGetResponse(&adapters.TypedNode{
		NodeType: "test",
		Node: graph.NewTestNode(0, 0, true, func() error {
			return nil
		}),
	}, nil)

	response, err := s.server.GetNodesState(nil, s.getValidGetStateRequest())

	s.Require().Nil(err)
	s.Require().Equal(0, len(response.Items))
	s.EqualValues(internalError, response.Base.Status)
	s.True(strings.HasPrefix(response.Base.Description, msg))
}

func (s *GTEServerTestSuite) TestProcess_Success() {
	s.storage.ExpectGetResponse(
		&adapters.TypedNode{
			NodeType: "test",
			Node: graph.NewTestNode(0, 0, true, func() error {
				return nil
			}),
		}, nil,
	)

	ids := s.getNodeIdentifiers(1)
	r, err := s.server.Process(nil, ids)
	s.Require().Nil(err)

	s.Require().Equal(1, len(r.Items))
	s.EqualValues(ok, r.Items[0].Base.Status)
}

func (s *GTEServerTestSuite) TestProcess_ProcessError() {
	e := fmt.Errorf("process error")
	s.storage.ExpectGetResponse(
		&adapters.TypedNode{
			NodeType: "test",
			Node: graph.NewTestNode(0, 0, true, func() error {
				return e
			}),
		}, nil,
	)

	ids := s.getNodeIdentifiers(1)
	r, err := s.server.Process(nil, ids)
	s.Require().Nil(err)

	s.Require().Equal(1, len(r.Items))
	s.EqualValues(internalError, r.Items[0].Base.Status)
	s.EqualValues(e.Error(), r.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestProcess_NodeNotFound() {
	e := fmt.Errorf("err not found")
	s.storage.ExpectGetResponse(
		nil, e,
	)

	ids := s.getNodeIdentifiers(1)
	r, err := s.server.Process(nil, ids)
	s.Require().Nil(err)

	s.Require().Equal(1, len(r.Items))
	s.EqualValues(notFound, r.Items[0].Base.Status)
	s.EqualValues(e.Error(), r.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestLink_Success() {
	s.storage.ExpectGetResponse(
		&adapters.TypedNode{
			NodeType: "test",
			Node: graph.NewTestNode(0, 0, true, func() error {
				return nil
			}),
		}, nil,
	)

	s.storage.ExpectGetResponse(
		&adapters.TypedNode{
			NodeType: "test",
			Node: graph.NewTestNode(0, 0, true, func() error {
				return nil
			}),
		}, nil,
	)

	s.factory.ExpectResponse(
		mocks.NodeAdapterMock{
			GetPortFunc: func(tag string, node graph.Node) (graph.Port, error) {
				return graph.NewAttachedPort(node), nil
			},
		}, nil,
	)

	s.factory.ExpectResponse(
		mocks.NodeAdapterMock{
			GetPortFunc: func(tag string, node graph.Node) (graph.Port, error) {
				return graph.NewAttachedPort(node), nil
			},
		}, nil,
	)

	r, err := s.server.Link(nil, s.getValidLinkRequest())
	s.Require().Nil(err)

	s.Require().Equal(1, len(r.Items))
	s.EqualValues(ok, r.Items[0].Base.Status)
	s.EqualValues(1, r.Items[0].Identifiers[0].Id)
	s.EqualValues(2, r.Items[0].Identifiers[1].Id)
}

func (s *GTEServerTestSuite) TestLink_NodeNotFound() {
	e := fmt.Errorf("err not found")
	s.storage.ExpectGetResponse(nil, e)

	r, err := s.server.Link(nil, s.getValidLinkRequest())
	s.Require().Nil(err)

	s.Require().Equal(1, len(r.Items))
	s.EqualValues(notFound, r.Items[0].Base.Status)
	s.EqualValues(e.Error(), r.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestLink_GetterNotFound() {
	s.storage.ExpectGetResponse(
		&adapters.TypedNode{
			NodeType: "test",
			Node: graph.NewTestNode(0, 0, true, func() error {
				return nil
			}),
		}, nil,
	)

	e := fmt.Errorf("err not found")
	s.factory.ExpectResponse(nil, e)

	r, err := s.server.Link(nil, s.getValidLinkRequest())
	s.Require().Nil(err)

	s.Require().Equal(1, len(r.Items))
	s.EqualValues(notFound, r.Items[0].Base.Status)
	s.EqualValues(e.Error(), r.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestLink_PortGetterError() {
	s.storage.ExpectGetResponse(
		&adapters.TypedNode{
			NodeType: "test",
			Node: graph.NewTestNode(0, 0, true, func() error {
				return nil
			}),
		}, nil,
	)

	e := fmt.Errorf("err not found")
	s.factory.ExpectResponse(
		mocks.NodeAdapterMock{
			GetPortFunc: func(tag string, node graph.Node) (graph.Port, error) {
				return nil, e
			},
		}, nil,
	)

	r, err := s.server.Link(nil, s.getValidLinkRequest())
	s.Require().Nil(err)

	s.Require().Equal(1, len(r.Items))
	s.EqualValues(notFound, r.Items[0].Base.Status)
	s.EqualValues(e.Error(), r.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) getValidGetStateRequest() *pb.NodeStateRequest {
	ids := s.getNodeIdentifiers(1)
	result := &pb.NodeStateRequest{
		Items: make([]*pb.NodeStateRequest_UnitRequest, 0),
	}

	for _, id := range ids.Ids {
		result.Items = append(result.Items, &pb.NodeStateRequest_UnitRequest{
			Identifier: id,
		})
	}

	return result
}

func (s *GTEServerTestSuite) getValidLinkRequest() *pb.LinkRequest {
	nodeIds := s.getNodeIdentifiers(1, 2)
	return &pb.LinkRequest{
		Items: []*pb.LinkRequest_UnitRequest{
			{
				Id1: &pb.PortIdentifier{
					NodeIdentifier: nodeIds.Ids[0],
					PortTag:        "port",
				},
				Id2: &pb.PortIdentifier{
					NodeIdentifier: nodeIds.Ids[1],
					PortTag:        "port",
				},
			},
		},
	}
}

func (s *GTEServerTestSuite) getValidCreateRequest() *pb.NodeCreateRequest {
	req, _ := GetCreateRequest(
		[]string{"node"},
		[]string{"test"},
		[]map[string]float64{
			{},
		},
	)
	return req
}

func (s *GTEServerTestSuite) getValidUpdateRequest() *pb.NodeUpdateRequest {
	ids := s.getNodeIdentifiers(1)
	req, _ := GetUpdateRequest(ids.Ids, []map[string]float64{{}})
	return req
}

func (s *GTEServerTestSuite) getNodeIdentifiers(ids ...int32) *pb.NodeIdentifiers {
	result := &pb.NodeIdentifiers{Ids: make([]*pb.NodeIdentifier, len(ids))}
	for i, id := range ids {
		result.Ids[i] = &pb.NodeIdentifier{Id: id}
	}
	return result
}

func TestGTEServerTestSuite(t *testing.T) {
	suite.Run(t, new(GTEServerTestSuite))
}
