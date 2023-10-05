module invoice-maker

go 1.19

require (
	github.com/gomarkdown/markdown v0.0.0-20221013030248-663e2500819c
	github.com/google/uuid v1.3.0
	github.com/phpdave11/gofpdf v1.4.2
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gdamore/tcell/v2 v2.6.0
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/rivo/tview v0.0.0-20230928053139-9bc1d28d88a9
	github.com/rivo/uniseg v0.4.3 // indirect
	github.com/shopspring/decimal v1.3.1
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/term v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
)

replace golang.org/x/text => golang.org/x/text v0.7.0

replace github.com/gdamore/tcell/v2 => github.com/gdamore/tcell/v2 v2.6.0
