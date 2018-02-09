package networkservice

import (
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"github.com/Sovianum/turbocycle/core/graph"
	"fmt"
)

const (
	inputTag  = "inputTag"
	outputTag = "outputTag"
	portATag  = "portATag"
	portBTag  = "portBTag"
	portCTag  = "portCTag"
)

func TestCheckContextStates(t *testing.T) {
	tc := []struct {
		states      []*pb.NodeDescription_ContextState
		contextTags map[string]bool
		containsErr bool
		errFragment string
	}{
		{
			containsErr: false,
		},
		{
			states: []*pb.NodeDescription_ContextState{
				{
					Ports: []*pb.NodeDescription_AttachedPortDescription{
						{
							Type:        pb.NodeDescription_AttachedPortDescription_CONTEXT_DEPENDENT,
							Description: &pb.PortDescription{},
						},
					},
				},
			},
			containsErr: true,
			errFragment: "context depndent ports are not allowed in",
		},
		{
			states: []*pb.NodeDescription_ContextState{
				{
					Ports: []*pb.NodeDescription_AttachedPortDescription{
						{
							Type:        pb.NodeDescription_AttachedPortDescription_INPUT,
							Description: &pb.PortDescription{},
						},
					},
				},
			},
			containsErr: true,
			errFragment: "number of specified context ports in state",
		},
		{
			states: []*pb.NodeDescription_ContextState{
				{
					Ports: []*pb.NodeDescription_AttachedPortDescription{
						{
							Type:        pb.NodeDescription_AttachedPortDescription_INPUT,
							Description: &pb.PortDescription{},
						},
					},
				},
			},
			containsErr: true,
			contextTags: map[string]bool{
				"wrong": false,
			},
			errFragment: "of context state not found in ports ",
		},
	}

	for i, c := range tc {
		errs := checkContextStates(c.states, c.contextTags)
		if c.containsErr {
			assert.NotNil(t, errs)
			text := joinErrors(errs)
			assert.True(t, strings.Contains(text, c.errFragment), "%d", i)
		} else {
			assert.Nil(t, errs)
		}
	}

	errs := checkContextStates(nil, nil)
	assert.Nil(t, errs)
}

func TestRepresentationNode_ContextDefined_Single_OK(t *testing.T) {
	node, err := NewRepresentationNode(getSinkDescription(), nil)
	assert.Nil(t, err)
	assert.True(t, node.ContextDefined(0))
}

func TestRepresentationNode_ContextDefined_Single_Fail(t *testing.T) {
	node, err := NewRepresentationNode(getBipoleDescription(), nil)
	assert.Nil(t, err)
	assert.False(t, node.ContextDefined(0))
}

func TestRepresentationNode_ContextDefined_Chain_OK(t *testing.T) {
	source, _ := NewRepresentationNode(getSourceDescription(), nil)
	sink, _ := NewRepresentationNode(getSinkDescription(), nil)
	bipole, _ := NewRepresentationNode(getBipoleDescription(), nil)

	var err error

	out, err := source.GetPortByName(outputTag)
	assert.Nil(t, err)

	in, err := sink.GetPortByName(inputTag)
	assert.Nil(t, err)

	a, err := bipole.GetPortByName(portATag)
	assert.Nil(t, err)

	b, err := bipole.GetPortByName(portBTag)
	assert.Nil(t, err)

	graph.Link(out, a)
	graph.Link(b, in)

	assert.True(t, bipole.ContextDefined(0))
}

func TestRepresentationNode_ContextDefined_Chain_Fail(t *testing.T) {
	bipole1, _ := NewRepresentationNode(getBipoleDescription(), nil)
	bipole2, _ := NewRepresentationNode(getBipoleDescription(), nil)

	a1, _ := bipole1.GetPortByName(portATag)
	b1, _ := bipole1.GetPortByName(portBTag)

	a2, _ := bipole2.GetPortByName(portATag)
	b2, _ := bipole2.GetPortByName(portBTag)

	graph.Link(b1, a2)
	graph.Link(b2, a1)

	assert.False(t, bipole1.ContextDefined(0))
	assert.False(t, bipole2.ContextDefined(0))
}

func TestRepresentationNode_ContextDefined_Network_OK(t *testing.T) {
	source, _ := NewRepresentationNode(getSourceDescription(), nil)
	sink1, _ := NewRepresentationNode(getSinkDescription(), nil)
	sink2, _ := NewRepresentationNode(getSinkDescription(), nil)
	sink3, _ := NewRepresentationNode(getSinkDescription(), nil)

	c1, _ := NewRepresentationNode(get1In2Out(), nil)
	c2, _ := NewRepresentationNode(get1In2Out(), nil)
	b, _ := NewRepresentationNode(getBipoleDescription(), nil)
	t1, _ := NewRepresentationNode(get2In1Out(), nil)
	t2, _ := NewRepresentationNode(get2In1Out(), nil)

	mustLink(source, c1, outputTag, portATag)
	mustLink(c1, c2, portBTag, portATag)
	mustLink(c1, t1, portCTag, portCTag)
	mustLink(c2, b, portBTag, portATag)
	mustLink(c2, sink1, portCTag, inputTag)
	mustLink(b, t2, portBTag, portATag)
	mustLink(t2, t1, portBTag, portATag)
	mustLink(sink2, t2, inputTag, portCTag)
	mustLink(t1, sink3, portBTag, inputTag)


	//assert.True(t, source.ContextDefined(0))
	//assert.True(t, c1.ContextDefined(0))
	//assert.True(t, c2.ContextDefined(0))
	//assert.True(t, b.ContextDefined(0))
	//assert.True(t, t2.ContextDefined(0))
	assert.True(t, t1.ContextDefined(0))
	//assert.True(t, sum.ContextDefined(0))
	//assert.True(t, pSource.ContextDefined(0))
	//assert.True(t, sink1.ContextDefined(0))

	fmt.Println(describeReprNode(c1, "c1"))
	fmt.Println(describeReprNode(c2, "c2"))
	fmt.Println(describeReprNode(b, "b"))
	fmt.Println(describeReprNode(t1, "t1"))
	fmt.Println(describeReprNode(t2, "t2"))
}

//func TestRepresentationNode_ContextDefined_Network_Fail(t *testing.T) {
//	gSource, _ := NewRepresentationNode(getSourceDescription(), nil)
//	c1, _ := NewRepresentationNode(get1In2Out(), nil)
//	c2, _ := NewRepresentationNode(get1In2Out(), nil)
//	b, _ := NewRepresentationNode(getBipoleDescription(), nil)
//	t1, _ := NewRepresentationNode(get2In1Out(), nil)
//	t2, _ := NewRepresentationNode(get2In1Out(), nil)
//	sum, _ := NewRepresentationNode(get2In1Out(), nil)
//
//	mustLink(gSource, c1, outputTag, portATag)
//	mustLink(c1, c2, portBTag, portATag)
//	mustLink(c1, t1, portCTag, portCTag)
//	mustLink(c2, b, portBTag, portATag)
//	mustLink(c2, sum, portCTag, portATag)
//	mustLink(b, t2, portBTag, portATag)
//	mustLink(t2, t1, portBTag, portATag)
//	mustLink(sum, t2, portCTag, portCTag)
//	mustLink(t1, sum, portBTag, portBTag)
//
//
//	assert.True(t, gSource.ContextDefined(0))
//	assert.True(t, c1.ContextDefined(0))
//	assert.True(t, c2.ContextDefined(0))
//	assert.True(t, b.ContextDefined(0))
//	assert.False(t, t2.ContextDefined(0))
//	assert.False(t, t1.ContextDefined(0))
//	assert.False(t, sum.ContextDefined(0))
//
//	fmt.Println(describeReprNode(c1, "c1"))
//	fmt.Println(describeReprNode(c2, "c2"))
//	fmt.Println(describeReprNode(b, "b"))
//	fmt.Println(describeReprNode(t1, "t1"))
//	fmt.Println(describeReprNode(t2, "t2"))
//	fmt.Println(describeReprNode(sum, "sum"))
//}

func mustLink(node1, node2 RepresentationNode, tag1, tag2 string) {
	port1, _ := node1.GetPortByName(tag1)
	port2, _ := node2.GetPortByName(tag2)
	graph.Link(port1, port2)
}

func describeReprNode(node RepresentationNode, name string) string {
	c := node.(*representationNode)
	result := fmt.Sprintf("node: %s\n", name)
	for i, state := range c.contextMatchingStates {
		result += fmt.Sprintf("\tstate_%d:\n", i)
		for _, port := range state.Ports {
			result += fmt.Sprintf("\t\t%s: %s\n", port.Description.Prefix, getStringStatus(port.Type))
		}
	}
	return result
}

func getStringStatus(t pb.NodeDescription_AttachedPortDescription_PortType) string {
	switch t {
	case pb.NodeDescription_AttachedPortDescription_INPUT:
		return "input"
	case pb.NodeDescription_AttachedPortDescription_OUTPUT:
		return "output"
	default:
		return "neutral"
	}
}

func getSourceDescription() *pb.NodeDescription {
	return &pb.NodeDescription{
		NodeType: "source",
		BasePorts: []*pb.NodeDescription_AttachedPortDescription{
			{
				Type: pb.NodeDescription_AttachedPortDescription_OUTPUT,
				Description: &pb.PortDescription{
					Prefix:  outputTag,
					IsMulti: false,
				},
			},
		},
	}
}

func getSinkDescription() *pb.NodeDescription {
	return &pb.NodeDescription{
		NodeType: "sink",
		BasePorts: []*pb.NodeDescription_AttachedPortDescription{
			{
				Type: pb.NodeDescription_AttachedPortDescription_INPUT,
				Description: &pb.PortDescription{
					Prefix:  inputTag,
					IsMulti: false,
				},
			},
		},
	}
}

func getBipoleDescription() *pb.NodeDescription {
	return &pb.NodeDescription{
		NodeType: "bipole",
		BasePorts: []*pb.NodeDescription_AttachedPortDescription{
			{
				Type: pb.NodeDescription_AttachedPortDescription_CONTEXT_DEPENDENT,
				Description: &pb.PortDescription{
					Prefix:  portATag,
					IsMulti: false,
				},
			},
			{
				Type: pb.NodeDescription_AttachedPortDescription_CONTEXT_DEPENDENT,
				Description: &pb.PortDescription{
					Prefix:  portBTag,
					IsMulti: false,
				},
			},
		},
		ContextStates: []*pb.NodeDescription_ContextState{
			{
				Ports: []*pb.NodeDescription_AttachedPortDescription{
					{
						Type: pb.NodeDescription_AttachedPortDescription_INPUT,
						Description: &pb.PortDescription{
							Prefix:  portATag,
							IsMulti: false,
						},
					},
					{
						Type: pb.NodeDescription_AttachedPortDescription_OUTPUT,
						Description: &pb.PortDescription{
							Prefix:  portBTag,
							IsMulti: false,
						},
					},
				},
			},
			{
				Ports: []*pb.NodeDescription_AttachedPortDescription{
					{
						Type: pb.NodeDescription_AttachedPortDescription_OUTPUT,
						Description: &pb.PortDescription{
							Prefix:  portATag,
							IsMulti: false,
						},
					},
					{
						Type: pb.NodeDescription_AttachedPortDescription_INPUT,
						Description: &pb.PortDescription{
							Prefix:  portBTag,
							IsMulti: false,
						},
					},
				},
			},
		},
	}
}

func get2In1Out() *pb.NodeDescription {
	return &pb.NodeDescription{
		NodeType: "tripole",
		BasePorts: []*pb.NodeDescription_AttachedPortDescription{
			{
				Type: pb.NodeDescription_AttachedPortDescription_CONTEXT_DEPENDENT,
				Description: &pb.PortDescription{
					Prefix:  portATag,
					IsMulti: false,
				},
			},
			{
				Type: pb.NodeDescription_AttachedPortDescription_CONTEXT_DEPENDENT,
				Description: &pb.PortDescription{
					Prefix:  portBTag,
					IsMulti: false,
				},
			},
			{
				Type: pb.NodeDescription_AttachedPortDescription_CONTEXT_DEPENDENT,
				Description: &pb.PortDescription{
					Prefix:  portCTag,
					IsMulti: false,
				},
			},
		},
		ContextStates: []*pb.NodeDescription_ContextState{
			{
				Ports: []*pb.NodeDescription_AttachedPortDescription{
					{
						Type: pb.NodeDescription_AttachedPortDescription_OUTPUT,
						Description: &pb.PortDescription{
							Prefix:  portATag,
							IsMulti: false,
						},
					},
					{
						Type: pb.NodeDescription_AttachedPortDescription_INPUT,
						Description: &pb.PortDescription{
							Prefix:  portBTag,
							IsMulti: false,
						},
					},
					{
						Type: pb.NodeDescription_AttachedPortDescription_INPUT,
						Description: &pb.PortDescription{
							Prefix:  portCTag,
							IsMulti: false,
						},
					},
				},
			},
			{
				Ports: []*pb.NodeDescription_AttachedPortDescription{
					{
						Type: pb.NodeDescription_AttachedPortDescription_INPUT,
						Description: &pb.PortDescription{
							Prefix:  portATag,
							IsMulti: false,
						},
					},
					{
						Type: pb.NodeDescription_AttachedPortDescription_OUTPUT,
						Description: &pb.PortDescription{
							Prefix:  portBTag,
							IsMulti: false,
						},
					},
					{
						Type: pb.NodeDescription_AttachedPortDescription_INPUT,
						Description: &pb.PortDescription{
							Prefix:  portCTag,
							IsMulti: false,
						},
					},
				},
			},
			{
				Ports: []*pb.NodeDescription_AttachedPortDescription{
					{
						Type: pb.NodeDescription_AttachedPortDescription_INPUT,
						Description: &pb.PortDescription{
							Prefix:  portATag,
							IsMulti: false,
						},
					},
					{
						Type: pb.NodeDescription_AttachedPortDescription_INPUT,
						Description: &pb.PortDescription{
							Prefix:  portBTag,
							IsMulti: false,
						},
					},
					{
						Type: pb.NodeDescription_AttachedPortDescription_OUTPUT,
						Description: &pb.PortDescription{
							Prefix:  portCTag,
							IsMulti: false,
						},
					},
				},
			},
		},
	}
}

func get1In2Out() *pb.NodeDescription {
	return &pb.NodeDescription{
		NodeType: "tripole",
		BasePorts: []*pb.NodeDescription_AttachedPortDescription{
			{
				Type: pb.NodeDescription_AttachedPortDescription_CONTEXT_DEPENDENT,
				Description: &pb.PortDescription{
					Prefix:  portATag,
					IsMulti: false,
				},
			},
			{
				Type: pb.NodeDescription_AttachedPortDescription_CONTEXT_DEPENDENT,
				Description: &pb.PortDescription{
					Prefix:  portBTag,
					IsMulti: false,
				},
			},
			{
				Type: pb.NodeDescription_AttachedPortDescription_CONTEXT_DEPENDENT,
				Description: &pb.PortDescription{
					Prefix:  portCTag,
					IsMulti: false,
				},
			},
		},
		ContextStates: []*pb.NodeDescription_ContextState{
			{
				Ports: []*pb.NodeDescription_AttachedPortDescription{
					{
						Type: pb.NodeDescription_AttachedPortDescription_INPUT,
						Description: &pb.PortDescription{
							Prefix:  portATag,
							IsMulti: false,
						},
					},
					{
						Type: pb.NodeDescription_AttachedPortDescription_OUTPUT,
						Description: &pb.PortDescription{
							Prefix:  portBTag,
							IsMulti: false,
						},
					},
					{
						Type: pb.NodeDescription_AttachedPortDescription_OUTPUT,
						Description: &pb.PortDescription{
							Prefix:  portCTag,
							IsMulti: false,
						},
					},
				},
			},
			{
				Ports: []*pb.NodeDescription_AttachedPortDescription{
					{
						Type: pb.NodeDescription_AttachedPortDescription_OUTPUT,
						Description: &pb.PortDescription{
							Prefix:  portATag,
							IsMulti: false,
						},
					},
					{
						Type: pb.NodeDescription_AttachedPortDescription_INPUT,
						Description: &pb.PortDescription{
							Prefix:  portBTag,
							IsMulti: false,
						},
					},
					{
						Type: pb.NodeDescription_AttachedPortDescription_OUTPUT,
						Description: &pb.PortDescription{
							Prefix:  portCTag,
							IsMulti: false,
						},
					},
				},
			},
			{
				Ports: []*pb.NodeDescription_AttachedPortDescription{
					{
						Type: pb.NodeDescription_AttachedPortDescription_OUTPUT,
						Description: &pb.PortDescription{
							Prefix:  portATag,
							IsMulti: false,
						},
					},
					{
						Type: pb.NodeDescription_AttachedPortDescription_OUTPUT,
						Description: &pb.PortDescription{
							Prefix:  portBTag,
							IsMulti: false,
						},
					},
					{
						Type: pb.NodeDescription_AttachedPortDescription_INPUT,
						Description: &pb.PortDescription{
							Prefix:  portCTag,
							IsMulti: false,
						},
					},
				},
			},
		},
	}
}
