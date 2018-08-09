package p2p

//Version used for node synchornization
type Version struct {
	Version int
	Height  int
	Blocks  []string
}

//NewVersion initializes a version
func NewVersion(version int, height int, blocks []string) Version {
	return Version{
		Version: version,
		Height:  height,
		Blocks:  blocks,
	}
}

//GetVersion returns version
func (v Version) GetVersion() int {
	return v.Version
}

//GetHeight return height
func (v Version) GetHeight() int {
	return v.Height
}

//GetBlocks returns blocks
func (v Version) GetBlocks() []string {
	return v.Blocks
}
