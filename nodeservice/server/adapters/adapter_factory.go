package adapters

type NodeAdapterFactory interface {
	GetAdapter(nodeType string) (NodeAdapter, error)
}
