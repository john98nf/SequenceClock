module github.com/john98nf/SequenceClock/cli

go 1.15

replace SequenceClock/cli/cmd => ./cmd

replace SequenceClock/cli/internal/controllerTemplates => ./internal/controllerTemplates

require (
	SequenceClock/cli/cmd v0.0.0-00010101000000-000000000000
	SequenceClock/cli/internal/controllerTemplates v0.0.0-00010101000000-000000000000 // indirect
	github.com/apache/openwhisk-client-go v0.0.0-20210313152306-ea317ea2794c // indirect
	github.com/spf13/cobra v1.1.3 // indirect
	github.com/spf13/viper v1.8.0 // indirect
)
