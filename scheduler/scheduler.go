package main

func newServer() *Server {
	s := &Server{}
	return s
}

func main() {
	s := newServer()
	s.InitServer()

	s.ServerForever()
}
