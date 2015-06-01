package qlog

type TopAccessResource struct {
	Url   string
	Host  string
	Path  string
	Count int
}

type AccessCntOfSupplier struct {
	Supplier string
	Count    int
}
