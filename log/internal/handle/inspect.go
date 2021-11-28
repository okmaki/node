package handle

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/okmaki/node/log/internal/port"
)

func Inspect(upgrader websocket.Upgrader, hubManager port.HubManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
