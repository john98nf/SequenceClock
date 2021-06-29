module john98nf/SequenceClock/deployer

go 1.15

replace john98nf/SequenceClock/deployer/internal/sequence => ./internal/sequence

replace john98nf/SequenceClock/deployer/internal/templateHandler => ./internal/templateHandler

require (
	github.com/apache/openwhisk-client-go v0.0.0-20210313152306-ea317ea2794c
	github.com/kataras/iris/v12 v12.1.8
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	john98nf/SequenceClock/deployer/internal/sequence v0.0.0-00010101000000-000000000000
	john98nf/SequenceClock/deployer/internal/templateHandler v0.0.0-00010101000000-000000000000
)
