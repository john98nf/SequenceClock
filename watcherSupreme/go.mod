module john98nf/SequenceClock/watcherSupreme

go 1.15

replace john98nf/SequenceClock/watcher/pkg/request => ../watcher/pkg/request

replace john98nf/SequenceClock/watcherSupreme/pkg/watcherClient => ./pkg/watcherClient

require github.com/gin-gonic/gin v1.7.2
