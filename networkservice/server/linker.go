package server

// Linker groups nodes in the graph according to the call order into such groups
// that data can be passed between nodes without network usage
type Linker interface {
	// Link links all possible nodes on the remote server and fills GraphData.domainCallOrder member
	Link(data *GraphData) error
}
