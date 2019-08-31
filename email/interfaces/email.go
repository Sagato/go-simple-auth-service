package interfaces

type Sender interface {
	Send(to []string, body []byte) error
}
