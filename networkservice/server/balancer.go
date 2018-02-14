package server

// Balancer handles creation and destruction of nodes on the remote servers
// register server should not make such orders manually
type Balancer interface {
	// Create creates nodes on the remote servers and sets data about created nodes to the corresponding
	// RepresentationNodes
	Create(data *GraphData) error
	// Delete destroys nodes on the remote servers and removes data about them from corresponding
	// RepresentationNodes
	Delete(data *GraphData) error
}
