module github.com/john98nf/SequenceClock/deployer/internal/templateHandler

go 1.15

replace github.com/john98nf/SequenceClock/deployer/pkg/sequence => ../../pkg/sequence

require (
	github.com/apache/openwhisk-client-go v0.0.0-20210313152306-ea317ea2794c
	github.com/john98nf/SequenceClock/deployer/pkg/sequence v0.0.0-00010101000000-000000000000
)
