package alpacamux

const MAX_MTU = 1500

type Forwarder interface {
	Listen() error
	Read(buf []byte) (int, error)
	Write(buf []byte) error
}
