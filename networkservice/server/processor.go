package server

// Processor processes nodes on the remote servers and transmits data between them
// according to GraphData.domainCallOrder member
type Processor interface {
	Process(data *GraphData) error
}
