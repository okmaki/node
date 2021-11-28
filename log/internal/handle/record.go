package handle

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/okmaki/node/log/internal/core"
	"github.com/okmaki/node/log/internal/port"
)

func isValidId(field string, value string) error {
	_, err := uuid.Parse(value)
	if err != nil {
		return errors.New("invalid " + field + " " + value)
	}

	return nil
}

func Record(recorder port.LogRecorder, notifier port.HubNotifier) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var log core.Log

		if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
			core.Warning("GetRecordHandler", "failed to decode request - %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := log.Validate(isValidId); err != nil {
			core.Warning("GetRecordHandler", "invalid log - %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := recorder.Record(log); err != nil {
			core.Error("GetRecordHandler", "failed to record log - %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		notifier.Notify(log)
		w.WriteHeader(http.StatusOK)
	}
}
