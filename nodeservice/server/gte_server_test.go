package server

import (
	"github.com/stretchr/testify/suite"
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"github.com/Sovianum/turbonetwork/nodeservice/server/factories"
	"testing"
)

type GTEServerTestSuite struct {
	suite.Suite
	server *gteServer
	ids *pb.Identifiers
}

func (s *GTEServerTestSuite) SetupTest() {
	s.server = NewGTEServer().(*gteServer)

	response, err := s.server.CreateNodes(nil, &pb.CreateRequest{
		Items:[]*pb.CreateRequest_UnitRequest{
			{
				NodeName:"node",
				NodeType:factories.PressureLossNodeType,
				Data:&pb.RequestData{
					DKwargs:map[string]float64{
						"sigma": 1,
					},
				},
			},
		},
	})

	s.Require().Nil(err)
	s.Require().Equal(1, len(response.Items))

	s.ids = s.getIdentifiers(response.Items)
}

func (s *GTEServerTestSuite) TestCreateNodes_Success() {
	m := s.server.nodeStorage.(*mapNodeStorage).nodeMap
	s.Equal(1, len(m))
}

func (s *GTEServerTestSuite) TestCreateNodes_NotFound() {
	response, err := s.server.CreateNodes(nil, &pb.CreateRequest{
		Items:[]*pb.CreateRequest_UnitRequest{
			{
				NodeName:"node",
				NodeType:"notExist",
				Data:nil,
			},
		},
	})
	s.Require().Nil(err)

	s.EqualValues(NotFound, response.Items[0].Base.Status)

	m := s.server.nodeStorage.(*mapNodeStorage).nodeMap
	s.Equal(1, len(m))
}

func (s *GTEServerTestSuite) TestDeleteNodes_Success() {
	response, err := s.server.DeleteNodes(nil, s.ids)

	s.Require().Nil(err)

	m := s.server.nodeStorage.(*mapNodeStorage).nodeMap
	s.Equal(0, len(m))

	s.Require().Equal(1, len(response.Items))
	s.Require().EqualValues(1, response.Items[0].Identifier.Id)
}

func (s *GTEServerTestSuite) TestDelete_NotFound() {
	response, err := s.server.DeleteNodes(nil, &pb.Identifiers{
		Ids:[]*pb.NodeIdentifier{
			{NodeType:factories.PressureLossNodeType, Id:100},
		},
	})
	s.Require().Nil(err)

	s.Equal(1, len(response.Items))

	// it is safe to delete even non-existing item
	s.EqualValues(OK, response.Items[0].Base.Status)
}

func (s *GTEServerTestSuite) TestGetNodes_Success() {
	response, err := s.server.GetNodes(nil, s.ids)
	s.Require().Nil(err)
	s.Require().EqualValues(OK, response.Items[0].Base.Status)

	s.Require().Equal(1, len(response.Items))
}

func (s *GTEServerTestSuite) TestProcess_Fail() {
	r, err := s.server.Process(nil, s.ids)
	s.Require().Nil(err)

	s.Equal(1, len(r.Items))
	s.EqualValues(OK, r.Base.Status)
	s.EqualValues(InternalError, r.Items[0].Status)
}

func (s *GTEServerTestSuite) getIdentifiers(items []*pb.ModifyResponse_UnitResponse) *pb.Identifiers {
	result := make([]*pb.NodeIdentifier, len(items))
	for i, item := range items {
		result[i] = item.Identifier
	}
	return &pb.Identifiers{Ids:result}
}

func TestGTEServerTestSuite(t *testing.T) {
	suite.Run(t, new(GTEServerTestSuite))
}
