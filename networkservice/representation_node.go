package networkservice

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/nodeservice/server/adapters"
	"fmt"
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"regexp"
	"math/rand"
)

type RepresentationNode interface {
	graph.Node
	GetPortByName(portTag string) (graph.Port, error)
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
		description:description,
		ports:make([]graph.Port, len(description.BasePorts)),
		requirePorts:make([]graph.Port, 0),
		updatePorts:make([]graph.Port, 0),
		portIndex:make(map[string]int),
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

	return result, nil
}

type representationNode struct {
	description adapters.NodeDescription
	ports []graph.Port

	portIndex map[string]int
	requirePorts []graph.Port
	updatePorts []graph.Port

	contextDefined bool
	contextCallKey int

	contextMatchingStates  []*pb.NodeDescription_ContextState
	contextState *pb.NodeDescription_ContextState
}

func (node *representationNode) SetName(name string) {}

func (node *representationNode) GetName() string {
	return "representationNode"
}

func (node *representationNode) Process() error {
	var matchStates []*pb.NodeDescription_ContextState
	var indices []int

	for i, state := range node.contextMatchingStates {
		if node.matchDataFlowDirection(state) {
			matchStates = append(matchStates, state)
			indices = append(indices, i)
		}
	}

	switch len(matchStates) {
	case 0:
		return fmt.Errorf("matching state not found")
	case 1:
		node.contextState = matchStates[0]
		return nil
	default:
		return fmt.Errorf("found multiple matching states %v", indices)
	}
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
	if len(node.description.ContextStates) == 0 {
		return true
	}
	if node.contextDefined {
		return true
	}

	if !rootKey(key) {
		if node.contextCallKey == key {
			return node.contextDefined
		}
		node.contextCallKey = key
	}

	node.contextMatchingStates = make([]*pb.NodeDescription_ContextState, 0)
	foundMatching := false
	for _, contextState := range node.description.ContextStates {
		if rootKey(key) {
			node.contextCallKey = newKey()
		}

		if node.matchContextDefinition(contextState) {
			foundMatching = true
			node.contextMatchingStates = append(node.contextMatchingStates, contextState)
		}
	}
	//return len(node.contextMatchingStates) != 0
	return foundMatching
}

func (node *representationNode) matchContextDefinition(contextState *pb.NodeDescription_ContextState) bool {
	r := regexp.MustCompile("_[0-9]+$")

	getPrefix := func(s string) string {
		i := r.Find([]byte(s))
		if i == nil {
			return s
		}
		return s[:i[0]]
	}

	stateMap := make(map[string]bool)
	for _, state := range contextState.Ports {
		if state.Type == pb.NodeDescription_AttachedPortDescription_INPUT {
			stateMap[state.Description.Prefix] = true
		}
	}

	for portName, index := range node.portIndex {
		prefix := getPrefix(portName)
		if _, ok := stateMap[prefix]; ok {
			outerNode := node.ports[index].GetOuterNode()
			if outerNode == nil {
				return false
			}
			if !outerNode.ContextDefined(node.contextCallKey) {
				return false
			}
		}
	}
	return true
}

func (node *representationNode) matchDataFlowDirection(contextState *pb.NodeDescription_ContextState) bool {
	var requiredPorts []*pb.NodeDescription_AttachedPortDescription
	var updatedPorts []*pb.NodeDescription_AttachedPortDescription
	for _, port := range contextState.Ports {
		switch port.Type {
		case pb.NodeDescription_AttachedPortDescription_INPUT:
			requiredPorts = append(requiredPorts, port)
		case pb.NodeDescription_AttachedPortDescription_OUTPUT:
			updatedPorts = append(updatedPorts, port)
		}
	}

	if len(requiredPorts) == 0 && len(updatedPorts) == 0 {
		return true
	}

	if len(requiredPorts) > 0 {
		for _, required := range requiredPorts {
			port, portErr := node.getPortByName(required.Description.Prefix)
			if portErr != nil {
				return false
			}

			isSource, sourceErr := nodes.IsDataSource(port)
			if sourceErr != nil || !isSource {
				return false
			}
		}
	}

	if len(updatedPorts) > 0 {
		for _, updated := range updatedPorts {
			port, portErr := node.getPortByName(updated.Description.Prefix)
			if portErr != nil {
				return false
			}

			isSink, sourceErr := nodes.IsDataSink(port)
			if sourceErr != nil || !isSink {
				return false
			}
		}
	}

	return true
}

func (node *representationNode) getPortByName(portTag string) (graph.Port, error) {
	if index, ok := node.portIndex[portTag]; !ok {
		return nil, fmt.Errorf("port %s not found", portTag)
	} else {
		return node.ports[index], nil
	}
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

func rootKey(key int) bool {
	return key <= 0
}

func newKey() int {
	return rand.Int() / 2 + 10
}