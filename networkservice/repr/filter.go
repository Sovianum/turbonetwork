package repr

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
)

func NewUnitFilter(val bool) Filter {
	return &complexFilter{
		v: func(state NodePortState) bool {
			return val
		},
	}
}

func NewFilterFromState(nps NodePortState) Filter {
	return &complexFilter{
		v: func(npsInner NodePortState) bool {
			for innerPort, innerState := range npsInner {
				// absence of port is interpreted as no requirements
				if state, stateOk := nps[innerPort]; stateOk {
					if state == pb.NodeDescription_AttachedPortDescription_NEUTRAL {
						continue
					}
					if innerState != state {
						return false
					}
				}
			}
			return true
		},
	}
}

type NodePortState map[graph.Port]pb.NodeDescription_AttachedPortDescription_PortType

type Filter interface {
	Validate(state NodePortState) bool
}

func Any(filters ...Filter) Filter {
	return &complexFilter{
		v: func(state NodePortState) bool {
			for _, filter := range filters {
				if filter.Validate(state) {
					return true
				}
			}
			return false
		},
	}
}

func All(filters ...Filter) Filter {
	return &complexFilter{
		v: func(state NodePortState) bool {
			for _, filter := range filters {
				if !filter.Validate(state) {
					return false
				}
			}
			return true
		},
	}
}

type complexFilter struct {
	v func(state NodePortState) bool
}

func (c *complexFilter) Validate(state NodePortState) bool {
	return c.v(state)
}
