package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"gx/ipfs/QmcTzQXRcU2vf8yX5EEboz1BSvWC7wWmeYAKVQmhp8WZYU/sha256-simd" 
)

var totalMiners int
const bigOlNum 100000

type Block struct {
	Nonce int64
	Parents []*Block
	Owner int
	Height int
	Null bool
	Weight int
	Seed int64
}

// Hash returns the hash of this block
func (b *Block) Hash() [32]byte {
	if b.hash == [32]byte{} {
		d, _ := json.Marshal(b)
		h := sha256.Sum256(d)
		b.hash = h
	}

	return b.hash
}

func (b *Block) ShortName() string {
	h := b.Hash()
	return fmt.Sprintf("%x", h[:8])
}

type RationalMiner struct {
	Power float64
	PrivateForks map[[32]byte][]*Block
	ID int
}

func NewRationalMiner(id int, power float64) *RationalMiner {
	return &RationalMiner{
		Power: power,
		PrivateForks: make(map[[32]byte][]*Block, 0),
		ID: id,
	}
}

func parentHeight(parents []*Block) {
	if len(parents) == 0 {
		panic("Don't call height on no parents")
	}
	return parents[0].Height
}

func parentWeight(parents []*Block) {
	if len(parents) == 0 {
		panic("Don't call weight on no parents")
	}
	return len(parents) + parents[0].Weight - 1
}

// maybeGenerateBlock maybe makes a new block with the given parents
func (m *RationalMiner) generateBlock(parents []*Block) *Block {
	// Given parents and id we have a unique source for new ticket
	t := m.generateTicket(minTicket)
	nextBlock := Block{
		Nonce: getUniqueID()
		Parents: parents,
		Owner: m.ID,
		Height: parentHeight(parents) + 1,
		Weight: parentWeight(parents),
		Seed: t,
	}
	
	if isWinningTicket(t) {
		nextBlock.Null = false
		nextBlock.Weight += 1
	} else {
		nextBlock.Null = true
		// TODO: update private forks		
	}

	return nextBlock
}

// generateTicket
func (m *RationalMiner) generateTicket(minTicket int) int64 {
	seed := minTicket + m.ID
	r := rand.New(rand.NewSource(seed))
	ticket := rand.Int63n(int64(bigOlNum * totalMiners))
	return ticket
}

// MaybeTrimForks purges the private fork slice of 

// Mine outputs the block that a miner mines in a round where the leaves of
// the block tree are given by liveHeads.  A miner will only ever mine one
// block in a round because if it mines two or more it gets slashed.  #Incentives #Blockchain
func (m *RationalMiner) Mine(liveHeads []*Block) *Block {
	m.SourceAllForks(liveHeads)
	
	maxWeight := 0
	var bestBlock *Block
	for i := 0; i <= len(m.PrivateForks); i++ {
		blk := generateBlock(parents)
		if !blk.Null && blk.Weight > maxWeight {
			bestBlock = blk
			maxWeight = weight
		}
	}
	// Get rid of private forks that we could not release without slashing.
	// Only occurs if a block is found.
	m.MaybeTrimForks(bestBlock)
	return bestBlock
}

func main() {
	rand.Seed(time.Now().UnixNano())
	gen := &Block{
		Nonce: getUniqueID(),
		Parents: nil,
		Owner: -1,
		Height: 0,
		Null: false,
		Weight: 0,
	}
	liveHeads := [][]*Block{[]*Block{gen}}
	roundNum := 1000
	totalMiners = 30
	miners := make([]*RationalMiner, totalMiners)
	for m := 0; m < totalMiners; m++ {
		miners[m] = NewRationalMiner(m, 1.0/totalMiners)
	}
	for round := 0; round < roundNum; round++ {
		var newBlocks []*Block
		for m := 0; m < totalMiners; m++ {
			// Each miner mines
			blk := Mine(liveHeads)
			if blk != nil {
				newBlocks = append(newBlocks, blk)
			}
		}

		// Network updates forks
		liveHeads = mergeNewBlocks(liveHeads, newBlocks)
	}
}
