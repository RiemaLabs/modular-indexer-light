package transfer

type TransferByInscription []VerifiableOrdTransfer

func (a TransferByInscription) Len() int {
	return len(a)
}

func (a TransferByInscription) Less(i, j int) bool {
	return a[i].ordTransfer.InscriptionID < a[j].ordTransfer.InscriptionID
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
