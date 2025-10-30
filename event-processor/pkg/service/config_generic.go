//go:build !tinygo.wasm

package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/event-processor/api"
	"github.com/runtime-radar/runtime-radar/event-processor/pkg/database"
	"github.com/runtime-radar/runtime-radar/event-processor/pkg/model"
	"github.com/runtime-radar/runtime-radar/event-processor/pkg/processor"
	"github.com/runtime-radar/runtime-radar/event-processor/pkg/processor/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type ConfigGeneric struct {
	api.UnimplementedConfigControllerServer

	ConfigRepository database.ConfigRepository
	Processor        processor.Processor
}

func (cg *ConfigGeneric) Add(ctx context.Context, req *api.Config) (*emptypb.Empty, error) {
	if reason, ok := cg.validateConfig(req); !ok {
		return nil, status.Error(codes.InvalidArgument, reason)
	}

	idStr := req.GetId()
	var id uuid.UUID
	var err error

	if idStr != "" {
		id, err = uuid.Parse(idStr)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "can't parse ID: %v", err)
		}
	}

	cfg := &model.Config{
		model.Base{ID: id},
		(*model.ConfigJSON)(req.GetConfig()),
	}
	if err := cg.ConfigRepository.Add(ctx, cfg); err != nil {
		return nil, status.Errorf(codes.Internal, "can't add config: %v", err)
	}

	// Changes applied instantly when requested. However, there can be multiple instances, each of which can update config.
	// To handle this there will be background worker doing same check periodically and calling update procedure when needed.
	oldCfg := cg.Processor.Config()

	log.Debug().Interface("old_config", oldCfg).Msgf("Old processor config")
	log.Debug().Interface("new_config", cfg).Msgf("New processor config")

	sel, changed := config.Diff(oldCfg, cfg)
	if changed {
		log.Info().
			Interface("config", cfg).
			Interface("selector", sel).
			Msgf("Processor config changed")

		cg.Processor.SetConfig(cfg)
	} else {
		log.Debug().Msgf("Processor config didn't change")
	}

	resp := &emptypb.Empty{}

	return resp, nil
}

func (cg *ConfigGeneric) Read(ctx context.Context, _ *emptypb.Empty) (*api.Config, error) {
	cfg, err := cg.ConfigRepository.GetLast(ctx, false)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "config not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "can't read config: %v", err)
	}

	resp := &api.Config{
		Id:     cfg.ID.String(),
		Config: (*api.Config_ConfigJSON)(cfg.Config),
	}

	return resp, nil
}

func (cg *ConfigGeneric) validateConfig(req *api.Config) (string, bool) {
	if req.Config == nil {
		return "no config", false
	}

	if req.Config.GetVersion() == "" {
		return "empty or missing config version", false
	} else if ver := req.Config.GetVersion(); ver != string(model.ConfigVersion) {
		return fmt.Sprintf("config version mismatch: expected %s, got %s", model.ConfigVersion, ver), false
	}

	return "", true
}
