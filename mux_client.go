package alpacamux

type MuxClient struct {
	Servers []string // upstream IP:Port pairs
	peer    MuxPeer
}

func (f *MuxClient) Listen() error {
	return f.peer.Init(0, f.Servers)
}

func (f *MuxClient) Read(buf []byte) (int, error) {
	return f.peer.Read(buf)
}

func (f *MuxClient) Write(buf []byte) error {
	return f.peer.Write(buf)
}
