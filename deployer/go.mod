module github.com/john98nf/SequenceClock/deployer

go 1.15

replace github.com/john98nf/SequenceClock/deployer/internal/sequence => ./internal/sequence

replace github.com/john98nf/SequenceClock/deployer/internal/templateHandler => ./internal/templateHandler

require (
	github.com/apache/openwhisk-client-go v0.0.0-20210313152306-ea317ea2794c
	github.com/gin-gonic/gin v1.7.4
	github.com/john98nf/SequenceClock/deployer/internal/sequence v0.0.0-00010101000000-000000000000
	github.com/john98nf/SequenceClock/deployer/internal/templateHandler v0.0.0-00010101000000-000000000000
)
