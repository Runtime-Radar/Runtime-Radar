package service

import history_api "github.com/runtime-radar/runtime-radar/history-api/api"

type RuntimeHistory interface {
	history_api.RuntimeHistoryClient
}
