package loadBalancer

type Strategy interface {
	getNextServer() (*Server, error)
}
