package transfer

type TransferByInscription []OrdTransfer

func (a TransferByInscription) Len() int {
	return len(a)
}

func (a TransferByInscription) Less(i, j int) bool {
	return a[i].NewSatpoint < a[j].NewSatpoint
}

func (a TransferByInscription) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type ArrayFloatsam []Flotsam

func (a ArrayFloatsam) Len() int {
	return len(a)
}

func (a ArrayFloatsam) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ArrayFloatsam) Less(i, j int) bool {
	return a[i].Offset < a[j].Offset
}

func FindIndex(data []string, in string) (bool, int) {
	for i := 0; i < len(data); i++ {
		if data[i] == in {
			return true, i
		}
	}
	return false, 0
}
