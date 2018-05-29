package p2p

type Version struct {
	Version int
	Height  int
	Blocks  []string
}

func NewVersion(version int, height int, blocks []string) Version {
	return Version{
		Version: version,
		Height:  height,
		Blocks:  blocks,
	}
}

func (v Version) GetVersion() int {
	return v.Version
}

func (v Version) GetHeight() int {
	return v.Height
}

func (v Version) GetBlocks() []string {
	return v.Blocks
}
