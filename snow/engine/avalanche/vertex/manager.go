// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package vertex

import (
	"github.com/ava-labs/gecko/ids"
	"github.com/ava-labs/gecko/snow/consensus/avalanche"
	"github.com/ava-labs/gecko/snow/consensus/snowstorm"
)

// Manager defines the persistant storage that is required by the consensus
// engine
type Manager interface {
	// Create a new vertex from the contents of a vertex
	BuildVertex(parentIDs ids.Set, txs []snowstorm.Tx) (avalanche.Vertex, error)

	// Attempt to convert a stream of bytes into a vertex
	ParseVertex(vertex []byte) (avalanche.Vertex, error)

	// GetVertex attempts to load a vertex by hash from storage
	GetVertex(vtxID ids.ID) (avalanche.Vertex, error)

	// Edge returns a list of accepted vertex IDs with no accepted children
	Edge() (vtxIDs []ids.ID)
}
