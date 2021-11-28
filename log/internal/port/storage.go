package port

import "github.com/okmaki/node/log/internal/core"

type LogRecorder interface {
	Record(log core.Log) error
}

type LogSearcher interface {
	SearchBySource(filter core.LogFilter, limit int) ([]core.Log, error)
	SearchByTransaction(filter core.LogFilter, limit int) ([]core.Log, error)
}
