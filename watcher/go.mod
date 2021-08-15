module github.com/john98nf/SequenceClock/watcher

go 1.15

replace github.com/john98nf/SequenceClock/watcher/pkg/request => ./pkg/request

require (
	github.com/containerd/containerd v1.5.5 // indirect
	github.com/docker/docker v20.10.8+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/gin-gonic/gin v1.7.4
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/john98nf/SequenceClock/watcher/pkg/request v0.0.0-00010101000000-000000000000
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	google.golang.org/grpc v1.40.0 // indirect
)
