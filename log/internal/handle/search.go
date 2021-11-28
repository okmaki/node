package handle

import (
	"encoding/json"
	"net/http"

	"github.com/okmaki/node/log/internal/core"
	"github.com/okmaki/node/log/internal/port"
)

const defaultLimit = 5

func Search(searcher port.LogSearcher) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var filter core.LogFilter

		if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
			core.Warning("GetSearchHandler", "failed to decode request - %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if filter.SourceType == "" {
			core.Warning("GetSearchHandler", "invalid filter - source type required")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var logs []core.Log
		var err error

		if filter.TransactionId != "" {
			logs, err = searcher.SearchByTransaction(filter, defaultLimit)
		} else {
			logs, err = searcher.SearchBySource(filter, defaultLimit)
		}

		if err != nil {
			// TODO: save logs that failed to record, to local file
			core.Error("GetSearchHandler", "failed to search logs using filter %s - %v", core.GetJson(filter), err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		core.Info("GetSearchHandler", "found %d results", len(logs))

		data, err := json.Marshal(logs)
		if err != nil {
			core.Error("GetSearchHandler", "failed to encode response - %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(data); err != nil {
			core.Error("GetSearchHandler", "failed to write response - %v", err)
		}
	}
}
