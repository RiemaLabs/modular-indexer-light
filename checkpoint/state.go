package checkpoint

import (
	"log"
	"sync"

	"github.com/RiemaLabs/indexer-committee/ord"
	"github.com/RiemaLabs/indexer-committee/ord/getter"
)

// Maintain a queue of states to prepare for the re-org.
// TODO: Use the first state and stateDiffs to represent states.
type StateQueue struct {
	States []ord.State
	sync.RWMutex
}

// Build the queue from the start height.
func NewQueues(getter getter.OrdGetter, initState ord.State, queryHash bool, startHeight uint, BitcoinConfirmations uint) (*StateQueue, error) {
	var states []ord.State
	states = make([]ord.State, BitcoinConfirmations)
	state := initState
	for i := startHeight; i <= startHeight+BitcoinConfirmations-1; i++ {
		ordTransfer, err := getter.GetOrdTransfers(i)
		if err != nil {
			return nil, err
		}
		state = ord.Exec(state.Copy(), ordTransfer)
		var hash string
		if queryHash {
			hash, err = getter.GetBlockHash(i)
			if err != nil {
				return nil, err
			}
		} else {
			hash = ""
		}
		state.Height = i
		state.Hash = hash
		states[i-startHeight] = state
	}
	queue := StateQueue{
		States: states,
	}
	return &queue, nil
}

func (queue *StateQueue) StartHeight() uint {
	return queue.States[0].Height
}

func (queue *StateQueue) LastestHeight() uint {
	return queue.StartHeight() + uint(len(queue.States))
}

func (queue *StateQueue) LastestState() ord.State {
	return queue.States[len(queue.States)-1]
}

func (queue *StateQueue) State(blockHeight uint) ord.State {
	return queue.States[blockHeight-queue.StartHeight()]
}

// Offer the latest state and pop the oldest state.
func (queue *StateQueue) Offer(element ord.State) {
	for i := 0; i <= len(queue.States)-2; i++ {
		queue.States[i] = queue.States[i+1]
	}
	queue.States[len(queue.States)-1] = element
}

func (queue *StateQueue) Println() {
	log.Println("====", len(queue.States), "====", queue.StartHeight(), "====")
	for _, node := range queue.States {
		log.Print(node.Height, "*")
	}
}

func (queue *StateQueue) Update(getter getter.OrdGetter, initState ord.State, latestHeight uint) error {
	state := initState
	curHeight := state.Height
	for i := curHeight + 1; i <= latestHeight; i++ {
		ordTransfer, err := getter.GetOrdTransfers(i)
		if err != nil {
			return err
		}
		state = ord.Exec(state.Copy(), ordTransfer)
		hash, err := getter.GetBlockHash(i)
		if err != nil {
			return err
		}
		state.Height = i
		state.Hash = hash
		queue.Offer(state)
	}
	return nil
}

// Check if the reorganization happened.
// If so, return the height where the reorganization happened, else, return 0.
func (queue *StateQueue) CheckForReorg(getter getter.OrdGetter) (uint, error) {
	for i := 0; i <= len(queue.States)-1; i++ {
		state := queue.States[i]
		height := state.Height
		hash := state.Hash
		newHash, err := getter.GetBlockHash(height)
		if err != nil {
			return 0, err
		}
		if hash == newHash {
			continue
		} else {
			return height, nil
		}
	}
	return 0, nil
}
