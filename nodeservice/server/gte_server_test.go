package server

import (
	"fmt"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"github.com/Sovianum/turbonetwork/nodeservice/server/factories"
	mocks2 "github.com/Sovianum/turbonetwork/nodeservice/server/factories/mocks"
	"github.com/Sovianum/turbonetwork/nodeservice/server/mocks"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type GTEServerTestSuite struct {
	suite.Suite
	server *gteServer

	storageMock        *mocks.NodeStorageMock
	constructorFactory *mocks2.ConstructorFactoryMock
	getterFactory      *mocks2.StateGetterFactoryMock
	updaterFactory     *mocks2.UpdaterFactoryMock
}

func (s *GTEServerTestSuite) SetupTest() {
	s.server = NewGTEServer().(*gteServer)

	s.storageMock = mocks.NewNodeStorageMock()
	s.constructorFactory = mocks2.NewConstructorFactoryMock()
	s.getterFactory = mocks2.NewStateGetterFactoryMock()

	s.server.nodeStorage = s.storageMock
}

func (s *GTEServerTestSuite) TestCreateNodes_Success() {
	req, _ := GetCreateRequest(
		[]string{"node"},
		[]string{factories.PressureLossNodeType},
		[]map[string]float64{
			{"sigma": 1},
		},
	)
	response, err := s.server.CreateNodes(nil, req)

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))

	m := s.server.nodeStorage.(*mapNodeStorage).nodeMap
	s.Equal(1, len(m))
}

func (s *GTEServerTestSuite) TestCreateNodes_NotFound() {
	req, _ := GetCreateRequest(
		[]string{"node"},
		[]string{"notExist"},
		[]map[string]float64{{}},
	)

	response, err := s.server.CreateNodes(nil, req)
	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))

	s.EqualValues(NotFound, response.Items[0].Base.Status)

	m := s.server.nodeStorage.(*mapNodeStorage).nodeMap
	s.Equal(0, len(m))
}

func (s *GTEServerTestSuite) TestCreateNodes_ConstructorErr() {
	req, _ := GetCreateRequest(
		[]string{"node"},
		[]string{factories.PressureLossNodeType},
		[]map[string]float64{
			{"sigma": 1},
		},
	)

	e := fmt.Errorf("constructor fail")

	s.constructorFactory.ExpectResponse(
		func(data *pb.RequestData) (graph.Node, error) {
			return nil, e
		}, nil,
	)

	s.server.constructorFactory = s.constructorFactory
	resp, err := s.server.CreateNodes(nil, req)

	s.Nil(err)
	s.Require().Equal(1, len(resp.Items))

	s.EqualValues(InternalError, resp.Items[0].Base.Status)
	s.Equal(e.Error(), resp.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestCreateNodes_Panic() {
	req, _ := GetCreateRequest(
		[]string{"node"},
		[]string{factories.PressureLossNodeType},
		[]map[string]float64{
			{"sigma": 1},
		},
	)

	msg := "panic constructor"

	s.constructorFactory.ExpectResponse(
		func(data *pb.RequestData) (graph.Node, error) {
			panic(msg)
		}, nil,
	)

	s.server.constructorFactory = s.constructorFactory
	resp, err := s.server.CreateNodes(nil, req)

	s.Nil(err)
	s.Require().Equal(0, len(resp.Items))

	s.EqualValues(InternalError, resp.Base.Status)
	s.True(strings.HasPrefix(resp.Base.Description, msg))
}

func (s *GTEServerTestSuite) TestCreateNodes_StorageAddError() {
	req, _ := GetCreateRequest(
		[]string{"node"},
		[]string{factories.PressureLossNodeType},
		[]map[string]float64{
			{"sigma": 1},
		},
	)

	msg := "storage err"
	s.storageMock.ExpectAddResponse(
		nil, fmt.Errorf(msg),
	)
	s.server.nodeStorage = s.storageMock

	resp, err := s.server.CreateNodes(nil, req)

	s.Nil(err)
	s.Require().Equal(1, len(resp.Items))

	s.EqualValues(InternalError, resp.Items[0].Base.Status)
	s.Equal(msg, resp.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestDeleteNodes_Success() {
	req, _ := GetCreateRequest(
		[]string{"node"},
		[]string{factories.PressureLossNodeType},
		[]map[string]float64{
			{"sigma": 1},
		},
	)
	resp, _ := s.server.CreateNodes(nil, req)

	response, err := s.server.DeleteNodes(nil, s.getIdentifiers(resp.Items))
	s.Require().Nil(err)

	m := s.server.nodeStorage.(*mapNodeStorage).nodeMap
	s.Equal(0, len(m))

	s.Require().Equal(1, len(response.Items))
	s.Require().EqualValues(1, response.Items[0].Identifiers[0].Id)
}

func (s *GTEServerTestSuite) TestDelete_NotFound() {
	response, err := s.server.DeleteNodes(nil, &pb.Identifiers{
		Ids: []*pb.NodeIdentifier{
			{NodeType: factories.PressureLossNodeType, Id: 100},
		},
	})
	s.Require().Nil(err)

	s.Equal(1, len(response.Items))

	// it is safe to delete even non-existing item
	s.EqualValues(OK, response.Items[0].Base.Status)
}

func (s *GTEServerTestSuite) TestGetNodes_Success() {
	req, _ := GetCreateRequest(
		[]string{"node"},
		[]string{factories.PressureLossNodeType},
		[]map[string]float64{
			{"sigma": 1},
		},
	)
	resp, _ := s.server.CreateNodes(nil, req)

	response, err := s.server.GetNodes(nil, s.getIdentifiers(resp.Items))
	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))
	s.Require().EqualValues(OK, response.Items[0].Base.Status)

	s.Require().Equal(1, len(response.Items))
}

func (s *GTEServerTestSuite) TestGetNodes_GetterNotFound() {
	e := fmt.Errorf("err not found")
	s.getterFactory.ExpectResponse(nil, e)
	s.server.stateGetterFactory = s.getterFactory

	response, err := s.server.GetNodes(nil, &pb.Identifiers{
		Ids: []*pb.NodeIdentifier{
			{
				Id:       1,
				NodeType: "type",
			},
		},
	})

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))

	s.EqualValues(NotFound, response.Items[0].Base.Status)
	s.Equal(e.Error(), response.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestGetNodes_GetterError() {
	e := fmt.Errorf("error")

	s.getterFactory.ExpectResponse(func(node *factories.TypedNode) (*pb.NodeState, error) {
		return nil, e
	}, nil)
	s.server.stateGetterFactory = s.getterFactory

	s.storageMock.ExpectGetResponse(&factories.TypedNode{}, nil)
	s.server.nodeStorage = s.storageMock

	response, err := s.server.GetNodes(nil, &pb.Identifiers{
		Ids: []*pb.NodeIdentifier{
			{
				Id:       1,
				NodeType: "type",
			},
		},
	})

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))

	s.EqualValues(InternalError, response.Items[0].Base.Status)
	s.Equal(e.Error(), response.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestGetNodes_StorageError() {
	s.getterFactory.ExpectResponse(func(node *factories.TypedNode) (*pb.NodeState, error) {
		return nil, nil
	}, nil)
	s.server.stateGetterFactory = s.getterFactory

	e := fmt.Errorf("storage error")
	s.storageMock.ExpectGetResponse(nil, e)
	s.server.nodeStorage = s.storageMock

	response, err := s.server.GetNodes(nil, &pb.Identifiers{
		Ids: []*pb.NodeIdentifier{
			{
				Id:       1,
				NodeType: "type",
			},
		},
	})

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))

	s.EqualValues(NotFound, response.Items[0].Base.Status)
	s.Equal(e.Error(), response.Items[0].Base.Description)
}

func (s *GTEServerTestSuite) TestGetNodes_Panic() {
	msg := "panic msg"
	s.getterFactory.ExpectResponse(func(node *factories.TypedNode) (*pb.NodeState, error) {
		panic(msg)
	}, nil)
	s.server.stateGetterFactory = s.getterFactory

	s.storageMock.ExpectGetResponse(&factories.TypedNode{}, nil)
	s.server.nodeStorage = s.storageMock

	response, err := s.server.GetNodes(nil, &pb.Identifiers{
		Ids: []*pb.NodeIdentifier{
			{
				Id:       1,
				NodeType: "type",
			},
		},
	})

	s.Require().Nil(err)
	s.Require().Equal(0, len(response.Items))

	s.EqualValues(InternalError, response.Base.Status)
	s.True(strings.HasPrefix(response.Base.Description, msg))
}

func (s *GTEServerTestSuite) TestProcess_Success() {
	s.storageMock.ExpectGetResponse(
		&factories.TypedNode{
			NodeType: "test",
			Node: graph.NewTestNode(0, 0, true, func() error {
				return nil
			}),
		}, nil,
	)
	s.server.nodeStorage = s.storageMock

	r, err := s.server.Process(nil, &pb.Identifiers{
		Ids: []*pb.NodeIdentifier{
			{
				Id:       1,
				NodeType: "test",
			},
		},
	})
	s.Require().Nil(err)

	s.Require().Equal(1, len(r.Items))
	s.EqualValues(OK, r.Items[0].Base.Status)
}

func (s *GTEServerTestSuite) TestProcess_ProcessError() {
	req, _ := GetCreateRequest(
		[]string{"node"},
		[]string{factories.PressureLossNodeType},
		[]map[string]float64{
			{"sigma": 1},
		},
	)
	resp, _ := s.server.CreateNodes(nil, req)

	r, err := s.server.Process(nil, s.getIdentifiers(resp.Items))
	s.Require().Nil(err)

	s.Require().Equal(1, len(r.Items))
	s.EqualValues(OK, r.Base.Status)
	s.EqualValues(InternalError, r.Items[0].Base.Status)
}

func (s *GTEServerTestSuite) TestProcess_StorageError() {
	e := fmt.Errorf("storage error")
	s.storageMock.ExpectGetResponse(
		nil, e,
	)
	s.server.nodeStorage = s.storageMock

	r, err := s.server.Process(nil, &pb.Identifiers{
		Ids: []*pb.NodeIdentifier{
			{
				Id:       1,
				NodeType: "test",
			},
		},
	})
	s.Require().Nil(err)

	s.Require().Equal(1, len(r.Items))
	s.EqualValues(NotFound, r.Items[0].Base.Status)
}

// returns only zeros's identifiers
func (s *GTEServerTestSuite) getIdentifiers(items []*pb.ModifyResponse_UnitResponse) *pb.Identifiers {
	result := make([]*pb.NodeIdentifier, len(items))
	for i, item := range items {
		result[i] = item.Identifiers[0]
	}
	return &pb.Identifiers{Ids: result}
}

func TestGTEServerTestSuite(t *testing.T) {
	suite.Run(t, new(GTEServerTestSuite))
}
