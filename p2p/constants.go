package p2p

import "errors"

const (
	//NodeDB file name of nodeinfo database
	NodeDB = "%s.nodeinfo.db"
	//NodeBucket bolt db bucket
	NodeBucket = "node"
	//DispatcherScheme scheme of gizo dispatcher url
	DispatcherScheme = "gizo"
	//MaxWorkers number of workers allowed per standard area
	MaxWorkers = 128
	//DefaultPort default port
	DefaultPort = 9999
	//CentrumURL url of centrum api
	// CentrumURL = "https://centrum-dev.herokuapp.com"
	CentrumURL = "https://ab92c3fb.ngrok.io"
	//GizoVersion version of gizo
	GizoVersion = 1
)

// node states
const (
	// when a node is not connected to the network
	DOWN = "DOWN"
	// worker - when a worker connects to a dispatchers standard area
	// dispatcher - when an adjacency is created and topology table, peer table and blockchain have not been synced
	INIT = "INIT"
	// worker - when a node starts receiving and crunching jobs
	LIVE = "LIVE"
	// dispatcher - when an adjacency is created and topology table, peer table and blockchain have been sync
	FULL = "FULL"
)

var (
	//ErrNoDispatchers occurs when there are no dispaters return from centrum
	ErrNoDispatchers = errors.New("Centrum: no dispatchers available")
)
