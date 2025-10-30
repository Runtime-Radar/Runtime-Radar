package config

import (
	"github.com/google/go-cmp/cmp"
	"github.com/runtime-radar/runtime-radar/event-processor/pkg/model"
	"google.golang.org/protobuf/testing/protocmp"
)

type Selector struct {
	HistoryControl bool
}

func Diff(oldCfg, newCfg *model.Config) (sel Selector, changed bool) {
	// TODO: disabled for debugging purposes, anyways it's a small optimization
	// if !newCfg.CreatedAt.After(oldCfg.CreatedAt) {
	// 	return
	// }

	// This logic can be much simpler, but it's supposed to become more complex over time,
	// so it's kept in line with runtime-monitor approach
	if !cmp.Equal(oldCfg.Config.HistoryControl, newCfg.Config.HistoryControl, protocmp.Transform()) {
		sel.HistoryControl, changed = true, true
	}

	return
}
