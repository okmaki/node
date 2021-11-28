package port

import "github.com/okmaki/node/log/internal/core"

type HubClient interface {
	GetFilter() core.LogFilter
	Send(log core.Log)
	Disconnect()
}

type HubManager interface {
	Join(client HubClient)
	Leave(client HubClient)
}

type HubNotifier interface {
	Notify(log core.Log)
}
