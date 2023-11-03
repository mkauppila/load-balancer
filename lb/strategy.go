package lb

type Strategy interface {
	getNextServer() (*Server, error)
}
