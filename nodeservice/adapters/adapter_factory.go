package adapters

// NodeAdapterFactory returns NodeAdapter by type name of the node it can handle
type NodeAdapterFactory interface {
	GetAdapter(nodeType string) (NodeAdapter, error)
}
