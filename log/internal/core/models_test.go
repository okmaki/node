package core_test

import (
	"testing"

	"github.com/okmaki/node/log/internal/core"
)

func TestLogFilterShouldFilter(t *testing.T) {
	log := core.Log{
		SourceId:      "source-id",
		SourceType:    "source-type",
		TransactionId: "transaction-id",
		Level:         core.LevelDebug,
		Timestamp:     420,
	}

	filter := core.NewLogFilter()

	// source id

	filter.SourceId = "wasd"
	if res := filter.ShouldFilter(log); !res {
		filterError(t, "source id", filter.SourceId, res)
	}

	filter.SourceId = ""
	if res := filter.ShouldFilter(log); res {
		filterError(t, "source id", filter.SourceId, res)
	}

	filter.SourceId = log.SourceId
	if res := filter.ShouldFilter(log); res {
		filterError(t, "source id", filter.SourceId, res)
	}

	// transaction id

	filter.TransactionId = "wasd"
	if res := filter.ShouldFilter(log); !res {
		filterError(t, "transaction id", filter.TransactionId, res)
	}

	filter.TransactionId = ""
	if res := filter.ShouldFilter(log); res {
		filterError(t, "transaction id", filter.TransactionId, res)
	}

	filter.TransactionId = log.TransactionId
	if res := filter.ShouldFilter(log); res {
		filterError(t, "transaction id", filter.TransactionId, res)
	}

	// source type

	filter.SourceType = "wasd"
	if res := filter.ShouldFilter(log); !res {
		filterError(t, "source type", filter.SourceType, res)
	}

	filter.SourceType = ""
	if res := filter.ShouldFilter(log); res {
		filterError(t, "source type", filter.SourceType, res)
	}

	filter.SourceType = log.SourceType
	if res := filter.ShouldFilter(log); res {
		filterError(t, "source type", filter.SourceType, res)
	}

	// levels

	filter.Levels = core.LevelWarning
	if res := filter.ShouldFilter(log); !res {
		filterError(t, "levels", filter.Levels, res)
	}

	filter.Levels = 0
	if res := filter.ShouldFilter(log); res {
		filterError(t, "levels", filter.Levels, res)
	}

	filter.Levels = log.Level
	if res := filter.ShouldFilter(log); res {
		filterError(t, "levels", filter.Levels, res)
	}

	filter.Levels = core.LevelWarning | log.Level
	if res := filter.ShouldFilter(log); res {
		filterError(t, "levels", filter.Levels, res)
	}

	// after

	filter.After = log.Timestamp + 1
	if res := filter.ShouldFilter(log); !res {
		filterError(t, "after", filter.After, res)
	}

	filter.After = 0
	if res := filter.ShouldFilter(log); res {
		filterError(t, "after", filter.After, res)
	}

	filter.After = log.Timestamp
	if res := filter.ShouldFilter(log); !res {
		filterError(t, "after", filter.After, res)
	}

	filter.After--
	if res := filter.ShouldFilter(log); res {
		filterError(t, "after", filter.After, res)
	}

	// before

	filter.Before = log.Timestamp - 1
	if res := filter.ShouldFilter(log); !res {
		filterError(t, "before", filter.Before, res)
	}

	filter.Before = 0
	if res := filter.ShouldFilter(log); res {
		filterError(t, "before", filter.Before, res)
	}

	filter.Before = log.Timestamp
	if res := filter.ShouldFilter(log); !res {
		filterError(t, "before", filter.Before, res)
	}

	filter.Before++
	if res := filter.ShouldFilter(log); res {
		filterError(t, "before", filter.Before, res)
	}
}

func filterError(t *testing.T, field string, value interface{}, res bool) {
	t.Errorf("%s %v - exp: %v | act: %v", field, value, !res, res)
}
