package view

type View uint64

const (
	ViewMain View = iota
	ViewInvoices
	ViewReceivers
	ViewReceiverEdit
	ViewConfig
)
