package p2p

//DispatcherInfo message used to create and maintain adjacency between dispatcher nodes
type DispatcherInfo struct {
	Pub   []byte
	Peers []string
}

//NewDispatcherInfo initializes dispatcher info
func NewDispatcherInfo(pub []byte, p []string) *DispatcherInfo {
	return &DispatcherInfo{Pub: pub, Peers: p}
}

//GetPub return public key
func (d DispatcherInfo) GetPub() []byte {
	return d.Pub
}

//GetPeers returns peers
func (d DispatcherInfo) GetPeers() []string {
	return d.Peers
}

//SetPub sets public key
func (d *DispatcherInfo) SetPub(pub []byte) {
	d.Pub = pub
}

//SetPeers sets peers
func (d *DispatcherInfo) SetPeers(n []string) {
	d.Peers = n
}

//AddPeer adds peers
func (d *DispatcherInfo) AddPeer(n string) {
	d.Peers = append(d.Peers, n)
}
