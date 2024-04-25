package transfer

import "github.com/RiemaLabs/modular-indexer-committee/ord/getter"

type ByNewSatpoint []getter.OrdTransfer

func (a ByNewSatpoint) Len() int {
	return len(a)
}

func (a ByNewSatpoint) Less(i, j int) bool {
	return a[i].NewSatpoint < a[j].NewSatpoint
}

func (a ByNewSatpoint) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
