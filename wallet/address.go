package wallet

type Address struct {
	Address          string
	Change           bool
	DeriviationIndex uint32
}

func NewAddress(
	address string,
	change bool,
	deriviationIndex uint32,
) *Address {
	return &Address{
		Address:          address,
		Change:           change,
		DeriviationIndex: deriviationIndex,
	}
}
