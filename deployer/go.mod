module john98nf/SequenceClock/deployer

go 1.15

replace internal/sequence => ./internal/sequence

replace internal/templateHandler => ./internal/templateHandler

require (
	github.com/kataras/iris/v12 v12.1.8
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	internal/sequence v0.0.0-00010101000000-000000000000
	internal/templateHandler v0.0.0-00010101000000-000000000000
)
