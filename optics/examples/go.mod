module github.com/fogfish/golem/optics/examples

go 1.21

require (
	github.com/fogfish/golem/optics v0.11.0
	github.com/fogfish/golem/pure v0.0.0-00010101000000-000000000000
	github.com/fogfish/guid/v2 v2.0.4
)

replace github.com/fogfish/golem/optics => ../

require github.com/fogfish/golem/hseq v1.1.1 // indirect

replace github.com/fogfish/golem/hseq => ../../hseq

replace github.com/fogfish/golem/pure => ../../pure
