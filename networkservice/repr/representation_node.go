package repr

import (
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"github.com/Sovianum/turbonetwork/nodeservice/server/adapters"
	"regexp"
)

type connType = pb.NodeDescription_AttachedPortDescription_PortType

type portDescription struct {
	baseID int
	contextIDs []int
}

type RepresentationNode interface {
	graph.Node
	GetPortByName(portTag string) (graph.Port, error)
	GetConnectionLines() []map[graph.Port]connType
	SelectState(stateID int) error
}

func NewRepresentationNode(description adapters.NodeDescription, multiPortMap map[string]int) (RepresentationNode, error) {
	if err := checkArgs(description, multiPortMap); err != nil {
		return nil, err
	}

	multiPortMapCopy := make(map[string]int)
	for key, val := range multiPortMap {
		multiPortMapCopy[key] = val
	}

	result := &representationNode{
		description:  description,
		ports:        make([]graph.Port, len(description.BasePorts)),
		requirePorts: make([]graph.Port, 0),
		updatePorts:  make([]graph.Port, 0),
		portIndex:    make(map[string]int),
		descriptionIndex: make(map[graph.Port]portDescription),

	}

	for i, basePortDescription := range description.BasePorts {
		port := graph.NewAttachedPort(result)
		result.ports[i] = port

		var portTag string
		prefix := basePortDescription.Description.Prefix
		if !basePortDescription.Description.IsMulti {
			portTag = prefix
		} else {
			cnt := multiPortMapCopy[prefix]
			portTag = fmt.Sprintf("%s_%d", prefix, cnt)
			multiPortMapCopy[prefix]--
		}

		port.SetTag(portTag)
		result.portIndex[portTag] = i

		switch basePortDescription.Type {
		case pb.NodeDescription_AttachedPortDescription_INPUT:
			result.requirePorts = append(result.requirePorts, port)
		case pb.NodeDescription_AttachedPortDescription_OUTPUT:
			result.updatePorts = append(result.updatePorts, port)
		}
	}

	for portName, ind := range result.portIndex {
		result.descriptionIndex[result.ports[ind]] = getPortDescription(portName, description)
	}

	return result, nil
}

type representationNode struct {
	graph.BaseNode

	description adapters.NodeDescription
	ports       []graph.Port

	portIndex    map[string]int
	descriptionIndex map[graph.Port]portDescription
	requirePorts []graph.Port
	updatePorts  []graph.Port
}

func (node *representationNode) GetName() string {
	return common.EitherString(node.GetInstanceName(), "representationNode")
}

func (node *representationNode) Process() error {
	return nil
}

func (node *representationNode) GetRequirePorts() ([]graph.Port, error) {
	return node.requirePorts, nil
}

func (node *representationNode) GetUpdatePorts() ([]graph.Port, error) {
	return node.updatePorts, nil
}

func (node *representationNode) GetPorts() []graph.Port {
	return node.ports
}

func (node *representationNode) GetPortByName(portTag string) (graph.Port, error) {
	return node.getPortByName(portTag)
}

func (node *representationNode) ContextDefined(key int) bool {
	return true
}

func (node *representationNode) GetConnectionLines() []map[graph.Port]connType {
	var resultLength int
	if cl := len(node.description.ContextStates); cl == 0 {
		resultLength = 1
	} else {
		resultLength = cl
	}

	result := make([]map[graph.Port]connType, resultLength)
	for i := range result{
		// it is safe to pass i to getConnectionLine even for context independent nodes
		// cos it won't be used
		result[i] = node.getConnectionLine(i)
	}
	return result
}

func (node *representationNode) SelectState(stateID int) error {
	if l := len(node.description.ContextStates); l == 0 {
		return nil
	} else if l <= stateID {
		return fmt.Errorf("stateID out of range")
	}

	state := node.description.ContextStates[stateID]
	for _, pd := range state.Ports {
		// it is safe to extract port by prefix cos
		// context dependent multiport is prohibited
		port := node.ports[node.portIndex[pd.Description.Prefix]]
		switch pd.Type {
		case pb.NodeDescription_AttachedPortDescription_INPUT:
			node.requirePorts = append(node.requirePorts, port)
		case pb.NodeDescription_AttachedPortDescription_OUTPUT:
			node.updatePorts = append(node.updatePorts, port)
		}
	}
	return nil
}

func (node *representationNode) getConnectionLine(contextID int) map[graph.Port]connType {
	result := make(map[graph.Port]connType)
	for _, ind := range node.portIndex {
		port := node.ports[ind]

		d := node.descriptionIndex[port]
		connType := node.description.BasePorts[d.baseID].Type
		if len(d.contextIDs) > 0 {
			connType = node.description.ContextStates[contextID].Ports[d.contextIDs[contextID]].Type
		}

		result[port] = connType
	}
	return result
}

func (node *representationNode) getPortByName(portTag string) (graph.Port, error) {
	if index, ok := node.portIndex[portTag]; !ok {
		return nil, fmt.Errorf("port %s not found", portTag)
	} else {
		return node.ports[index], nil
	}
}

func getPortDescription(portName string, nodeDescription adapters.NodeDescription) portDescription {
	result := portDescription{}
	prefix := getPrefix(portName)
	var portType pb.NodeDescription_AttachedPortDescription_PortType

	for i, bd := range nodeDescription.BasePorts {
		if bd.Description.Prefix == prefix {
			result.baseID = i
			portType = bd.Type
		}
	}

	if portType == pb.NodeDescription_AttachedPortDescription_CONTEXT_DEPENDENT {
		result.contextIDs = make([]int, len(nodeDescription.ContextStates))
		for i, cd := range nodeDescription.ContextStates {
			for j, pd := range cd.Ports {
				if pd.Description.Prefix == prefix {
					result.contextIDs[i] = j
				}
			}
		}
	}
	return result
}

func checkArgs(description adapters.NodeDescription, multiPortMap map[string]int) error {
	basePorts := description.BasePorts
	seen := make(map[string]bool)
	contextDependentTags := make(map[string]bool)

	var errList []error

	for _, base := range basePorts {
		prefix := base.Description.Prefix
		portType := base.Type
		isMulti := base.Description.IsMulti
		contextDependent := portType == pb.NodeDescription_AttachedPortDescription_CONTEXT_DEPENDENT

		// check that all port tags are unique
		if _, ok := seen[prefix]; ok {
			errList = append(errList, fmt.Errorf(
				"duplicate port tag %s",
				seen[prefix],
			))
			continue
		}

		//check that there is no multi context defined ports
		if contextDependent && isMulti {
			errList = append(errList, fmt.Errorf(
				"multi context dependent ports are not allowed (%s)", prefix,
			))
			continue
		}

		// check that the amount of all valid multi ports is specified
		if _, ok := multiPortMap[prefix]; isMulti && !ok {
			errList = append(errList, fmt.Errorf(
				"you have not specified port number of multi port \"%s\" of node \"%s\"",
				prefix, description.NodeType,
			))
			continue
		}

		if contextDependent {
			contextDependentTags[prefix] = true
		}
	}

	if len(contextDependentTags) != 0 {
		errors := checkContextStates(description.ContextStates, contextDependentTags)
		if errors != nil {
			errList = append(errList, fmt.Errorf(
				"context states fail: %s", joinErrors(errors),
			))
		}
	}

	if errList == nil {
		return nil
	}
	return fmt.Errorf("failed at node %s: %s", description.NodeType, joinErrors(errList))
}

func checkContextStates(contextStates []*pb.NodeDescription_ContextState, contextTags map[string]bool) []error {
	var result []error

	tagExtractor := func(state *pb.NodeDescription_ContextState) (map[string]bool, []error) {
		result := make(map[string]bool)
		var errors []error

		for _, item := range state.Ports {
			if item.Type == pb.NodeDescription_AttachedPortDescription_CONTEXT_DEPENDENT {
				errors = append(errors, fmt.Errorf(
					"context depndent ports are not allowed in context states (%s)",
					item.Description.Prefix,
				))
				continue
			}
			result[item.Description.Prefix] = true
		}
		if errors != nil {
			return nil, errors
		}
		return result, nil
	}

	for i, state := range contextStates {
		tags, errors := tagExtractor(state)
		if errors != nil {
			result = append(result, fmt.Errorf("fail at context state %d: %s", i, joinErrors(errors)))
			continue
		}

		if len(tags) != len(contextTags) {
			result = append(result, fmt.Errorf(
				"fail at context state %d: %s", i,
				fmt.Sprintf(
					"number of specified context ports in state (%d) does not match number of context ports (%d)",
					len(tags), len(contextTags),
				),
			))
			continue
		}

		for tag := range tags {
			if _, ok := contextTags[tag]; !ok {
				result = append(result, fmt.Errorf(
					"fail at context state %d: %s", i,
					fmt.Sprintf(
						"port %s of context state not found in ports of the node",
						tag,
					),
				))
			}
		}
	}
	return result
}

func joinErrors(errors []error) string {
	result := "["
	for _, err := range errors {
		result += err.Error() + "; "
	}
	result += "]"
	return result
}

func getPrefix(s string) string {
	i := prefixMatcher.Find([]byte(s))
	if i == nil {
		return s
	}
	return s[:i[0]]
}

var prefixMatcher = regexp.MustCompile("_[0-9]+$")
