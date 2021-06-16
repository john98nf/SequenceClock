module github.com/john98nf/SequenceClock

go 1.15

replace SequenceClock/cmd => ./cmd

replace SequenceClock/internal/controllerTemplates => ./internal/controllerTemplates

require (
	SequenceClock/cmd v0.0.0-00010101000000-000000000000
	SequenceClock/internal/controllerTemplates v0.0.0-00010101000000-000000000000 // indirect
	github.com/spf13/cobra v1.1.3 // indirect
)
