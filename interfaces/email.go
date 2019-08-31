package interfaces

type Email interface {
	Send(addr string, from string, to []string, msg []byte) error
}
