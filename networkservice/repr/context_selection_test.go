package repr

import (
	"fmt"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/pb"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ContextSelectorTestSuite struct {
	suite.Suite
	source1  RepresentationNode
	source2  RepresentationNode
	sink1    RepresentationNode
	sink2    RepresentationNode
	c1       RepresentationNode
	c2       RepresentationNode
	b        RepresentationNode
	t1       RepresentationNode
	t2       RepresentationNode
	selector *contextSelector
}

func (s *ContextSelectorTestSuite) SetupTest() {
	s.source1, _ = NewRepresentationNode(getSourceDescription(), nil)
	s.source1.SetName("source1")
	s.source2, _ = NewRepresentationNode(getSourceDescription(), nil)
	s.source2.SetName("source2")
	s.sink1, _ = NewRepresentationNode(getSinkDescription(), nil)
	s.sink1.SetName("sink1")
	s.sink2, _ = NewRepresentationNode(getSinkDescription(), nil)
	s.sink2.SetName("sink2")

	s.c1, _ = NewRepresentationNode(get1In2Out(), nil)
	s.c1.SetName("c1")
	s.c2, _ = NewRepresentationNode(get1In2Out(), nil)
	s.c2.SetName("c2")
	s.b, _ = NewRepresentationNode(getBipoleDescription(), nil)
	s.b.SetName("b")
	s.t1, _ = NewRepresentationNode(get2In1Out(), nil)
	s.t1.SetName("t1")
	s.t2, _ = NewRepresentationNode(get2In1Out(), nil)
	s.t2.SetName("t2")

	mustLink(s.source1, s.c1, outputTag, portATag)
	mustLink(s.c1, s.c2, portBTag, portATag)
	mustLink(s.c1, s.t1, portCTag, portCTag)
	mustLink(s.c2, s.b, portBTag, portATag)
	mustLink(s.c2, s.sink1, portCTag, inputTag)
	mustLink(s.b, s.t2, portBTag, portATag)
	mustLink(s.t2, s.t1, portBTag, portATag)
	mustLink(s.source2, s.t2, outputTag, portCTag)
	mustLink(s.t1, s.sink2, portBTag, inputTag)

	s.selector = newContextSelector([]RepresentationNode{
		s.source1, s.source2, s.sink1, s.sink2, s.c1, s.c2, s.b, s.t1, s.t2,
	})
}

func (s *ContextSelectorTestSuite) TestInitConnMatrix() {
	r, c := s.selector.connMatrix.Dims()
	s.Equal(18, r)
	s.Equal(18, c)

	for i := 0; i != r; i++ {
		for j := 0; j != c; j++ {
			s.Equal(
				int(pb.NodeDescription_AttachedPortDescription_NEUTRAL),
				s.selector.connMatrix.At(i, j), "%d %d", i, j,
			)
		}
	}
}

func (s *ContextSelectorTestSuite) TestMakePortIndex() {
	s.Equal(18, len(s.selector.portIndex))
	s.Require().Equal(9, len(s.selector.connConfigs))

	s.Equal(1, len(s.selector.connConfigs[0]))
	s.Equal(1, len(s.selector.connConfigs[1]))
	s.Equal(1, len(s.selector.connConfigs[2]))
	s.Equal(1, len(s.selector.connConfigs[3]))
	s.Equal(3, len(s.selector.connConfigs[4]))
	s.Equal(3, len(s.selector.connConfigs[5]))
	s.Equal(2, len(s.selector.connConfigs[6]))
	s.Equal(3, len(s.selector.connConfigs[7]))
	s.Equal(3, len(s.selector.connConfigs[8]))
}

func (s *ContextSelectorTestSuite) TestUpdateConnMatrix() {
	line := s.source1.GetConnectionLines()[0]
	s.selector.updateConnMatrix(line)

	for from, connType := range line {
		to := from.GetLinkPort()
		got := s.selector.connMatrix.At(
			s.selector.portIndex[from],
			s.selector.portIndex[to],
		)
		s.EqualValues(connType, got)
	}
}

func (s *ContextSelectorTestSuite) TestCheckGraphMatrix() {
	selectors := s.selector.findValidConfigurations()
	s.Require().EqualValues(1, len(selectors))
	selector := selectors[0]

	getMap := func(id int) map[graph.Port]connType {
		return s.selector.connConfigs[id][selector[id]]
	}

	source1Map := getMap(0)
	s.Require().Equal(1, len(source1Map))
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_OUTPUT,
		source1Map[mustExtract(s.source1, outputTag)],
	)

	source2Map := getMap(1)
	s.Require().Equal(1, len(source2Map))
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_OUTPUT,
		source2Map[mustExtract(s.source2, outputTag)],
	)

	sink1Map := getMap(2)
	s.Require().Equal(1, len(sink1Map))
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_INPUT,
		sink1Map[mustExtract(s.sink1, inputTag)],
	)

	sink2Map := getMap(3)
	s.Require().Equal(1, len(sink2Map))
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_INPUT,
		sink2Map[mustExtract(s.sink2, inputTag)],
	)

	c1Map := getMap(4)
	s.Require().Equal(3, len(c1Map))
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_INPUT,
		c1Map[mustExtract(s.c1, portATag)],
	)
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_OUTPUT,
		c1Map[mustExtract(s.c1, portBTag)],
	)
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_OUTPUT,
		c1Map[mustExtract(s.c1, portCTag)],
	)

	c2Map := getMap(5)
	s.Require().Equal(3, len(c2Map))
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_INPUT,
		c2Map[mustExtract(s.c2, portATag)],
	)
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_OUTPUT,
		c2Map[mustExtract(s.c2, portBTag)],
	)
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_OUTPUT,
		c2Map[mustExtract(s.c2, portCTag)],
	)

	bMap := getMap(6)
	s.Require().Equal(2, len(bMap))
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_INPUT,
		bMap[mustExtract(s.b, portATag)],
	)
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_OUTPUT,
		bMap[mustExtract(s.b, portBTag)],
	)

	t1Map := getMap(7)
	s.Require().Equal(3, len(t1Map))
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_INPUT,
		t1Map[mustExtract(s.t1, portATag)],
	)
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_OUTPUT,
		t1Map[mustExtract(s.t1, portBTag)],
	)
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_INPUT,
		t1Map[mustExtract(s.t1, portCTag)],
	)

	t2Map := getMap(8)
	s.Require().Equal(3, len(t2Map))
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_INPUT,
		t2Map[mustExtract(s.t2, portATag)],
	)
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_OUTPUT,
		t2Map[mustExtract(s.t2, portBTag)],
	)
	s.Equal(
		pb.NodeDescription_AttachedPortDescription_INPUT,
		t2Map[mustExtract(s.t2, portCTag)],
	)
}

func (s *ContextSelectorTestSuite) TestConfigure() {
	err := s.selector.configure()
	s.Require().Nil(err)

	mustGetPorts := func(f func() ([]graph.Port, error)) []graph.Port {
		if ports, err := f(); err != nil {
			panic(err)
		} else {
			return ports
		}
	}

	inArray := func(p graph.Port, arr []graph.Port) bool {
		for _, port := range arr {
			if p == port {
				return true
			}
		}
		return false
	}

	s.Require().Equal(1, len(mustGetPorts(s.c1.GetRequirePorts)))
	s.True(inArray(mustExtract(s.c1, portATag), mustGetPorts(s.c1.GetRequirePorts)))
	s.Require().Equal(2, len(mustGetPorts(s.c1.GetUpdatePorts)))
	s.True(inArray(mustExtract(s.c1, portBTag), mustGetPorts(s.c1.GetUpdatePorts)))
	s.True(inArray(mustExtract(s.c1, portCTag), mustGetPorts(s.c1.GetUpdatePorts)))

	s.Require().Equal(1, len(mustGetPorts(s.c2.GetRequirePorts)))
	s.True(inArray(mustExtract(s.c2, portATag), mustGetPorts(s.c2.GetRequirePorts)))
	s.Require().Equal(2, len(mustGetPorts(s.c2.GetUpdatePorts)))
	s.True(inArray(mustExtract(s.c2, portBTag), mustGetPorts(s.c2.GetUpdatePorts)))
	s.True(inArray(mustExtract(s.c2, portCTag), mustGetPorts(s.c2.GetUpdatePorts)))

	s.Require().Equal(1, len(mustGetPorts(s.b.GetRequirePorts)))
	s.True(inArray(mustExtract(s.b, portATag), mustGetPorts(s.b.GetRequirePorts)))
	s.Require().Equal(1, len(mustGetPorts(s.b.GetUpdatePorts)))
	s.True(inArray(mustExtract(s.b, portBTag), mustGetPorts(s.b.GetUpdatePorts)))

	s.Require().Equal(2, len(mustGetPorts(s.t1.GetRequirePorts)))
	s.True(inArray(mustExtract(s.t1, portATag), mustGetPorts(s.t1.GetRequirePorts)))
	s.True(inArray(mustExtract(s.t1, portCTag), mustGetPorts(s.t1.GetRequirePorts)))
	s.Require().Equal(1, len(mustGetPorts(s.t1.GetUpdatePorts)))
	s.True(inArray(mustExtract(s.t1, portBTag), mustGetPorts(s.t1.GetUpdatePorts)))

	s.Require().Equal(2, len(mustGetPorts(s.t2.GetRequirePorts)))
	s.True(inArray(mustExtract(s.t2, portATag), mustGetPorts(s.t2.GetRequirePorts)))
	s.True(inArray(mustExtract(s.t2, portCTag), mustGetPorts(s.t2.GetRequirePorts)))
	s.Require().Equal(1, len(mustGetPorts(s.t2.GetUpdatePorts)))
	s.True(inArray(mustExtract(s.t2, portBTag), mustGetPorts(s.t2.GetUpdatePorts)))
}

func printConnConfigs(connConfigs [][]map[graph.Port]connType) {
	str := ""
	for i, config := range connConfigs {
		str += fmt.Sprintf("%d\n", i)
		for j, l := range config {
			str += fmt.Sprintf("\t%d\t%v\n", j, l)
		}
	}
	fmt.Printf("%s\n", str)
}

func TestContextSelectorTestSuite(t *testing.T) {
	suite.Run(t, new(ContextSelectorTestSuite))
}

func mustExtract(node RepresentationNode, portTag string) graph.Port {
	port, err := node.GetPortByName(portTag)
	if err != nil {
		panic(err)
	}
	return port
}
