package alpacamux

type Forwarder interface {
	Listen() error
	Read(buf []byte) (int, error)
	Write(buf []byte) error
}
