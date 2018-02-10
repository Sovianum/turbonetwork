package repr

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFilter(t *testing.T) {
	p1 := graph.NewPort()
	p2 := graph.NewPort()

	tc := []struct {
		baseState NodePortState
		testState NodePortState
		expected  bool
	}{
		{
			baseState: NodePortState{
				p1: pb.NodeDescription_AttachedPortDescription_INPUT,
				p2: pb.NodeDescription_AttachedPortDescription_OUTPUT,
			},
			testState: NodePortState{
				p1: pb.NodeDescription_AttachedPortDescription_INPUT,
				p2: pb.NodeDescription_AttachedPortDescription_OUTPUT,
			},
			expected: true,
		},
		{
			baseState: NodePortState{
				p2: pb.NodeDescription_AttachedPortDescription_OUTPUT,
			},
			testState: NodePortState{
				p1: pb.NodeDescription_AttachedPortDescription_INPUT,
				p2: pb.NodeDescription_AttachedPortDescription_OUTPUT,
			},
			expected: true,
		},
		{
			baseState: NodePortState{
				p1: pb.NodeDescription_AttachedPortDescription_INPUT,
				p2: pb.NodeDescription_AttachedPortDescription_OUTPUT,
			},
			testState: NodePortState{
				p2: pb.NodeDescription_AttachedPortDescription_OUTPUT,
			},
			expected: true,
		},
		{
			baseState: NodePortState{
				p1: pb.NodeDescription_AttachedPortDescription_INPUT,
				p2: pb.NodeDescription_AttachedPortDescription_OUTPUT,
			},
			testState: NodePortState{
				p1: pb.NodeDescription_AttachedPortDescription_INPUT,
				p2: pb.NodeDescription_AttachedPortDescription_INPUT,
			},
			expected: false,
		},
	}

	for i, c := range tc {
		f := NewFilterFromState(c.baseState)
		assert.Equal(t, f.Validate(c.testState), c.expected, "%d", i)
	}
}

func TestAll(t *testing.T) {
	tc := []struct {
		filters  []Filter
		expected bool
	}{
		{
			filters:  []Filter{NewUnitFilter(true), NewUnitFilter(true)},
			expected: true,
		},
		{
			filters:  []Filter{NewUnitFilter(false), NewUnitFilter(true)},
			expected: false,
		},
		{
			filters:  []Filter{NewUnitFilter(false), NewUnitFilter(false)},
			expected: false,
		},
		{
			filters:  []Filter{},
			expected: true,
		},
	}

	for i, c := range tc {
		assert.Equal(t, c.expected, All(c.filters...).Validate(nil), "%d", i)
	}
}

func TestAny(t *testing.T) {
	tc := []struct {
		filters  []Filter
		expected bool
	}{
		{
			filters:  []Filter{NewUnitFilter(true), NewUnitFilter(true)},
			expected: true,
		},
		{
			filters:  []Filter{NewUnitFilter(false), NewUnitFilter(true)},
			expected: true,
		},
		{
			filters:  []Filter{NewUnitFilter(false), NewUnitFilter(false)},
			expected: false,
		},
		{
			filters:  []Filter{},
			expected: false,
		},
	}

	for i, c := range tc {
		assert.Equal(t, c.expected, Any(c.filters...).Validate(nil), "%d", i)
	}
}

func TestComplex(t *testing.T) {
	tc := []struct {
		filter   Filter
		expected bool
	}{
		{
			filter: All(
				Any(NewUnitFilter(true), NewUnitFilter(false)),
				Any(NewUnitFilter(true), NewUnitFilter(false)),
			),
			expected: true,
		},
		{
			filter: All(
				Any(NewUnitFilter(true), NewUnitFilter(false)),
				Any(NewUnitFilter(false), NewUnitFilter(false)),
			),
			expected: false,
		},
		{
			filter: Any(
				All(NewUnitFilter(true), NewUnitFilter(false)),
				All(NewUnitFilter(true), NewUnitFilter(false)),
			),
			expected: false,
		},
		{
			filter: Any(
				All(NewUnitFilter(true), NewUnitFilter(false)),
				All(NewUnitFilter(true), NewUnitFilter(true)),
			),
			expected: true,
		},
	}

	for i, c := range tc {
		assert.Equal(t, c.expected, c.filter.Validate(nil), "%d", i)
	}
}
