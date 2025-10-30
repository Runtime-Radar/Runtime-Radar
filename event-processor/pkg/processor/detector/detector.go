//go:build !tinygo.wasm

package detector

import (
	"context"
	"fmt"
	"time"

	tetragon_api "github.com/cilium/tetragon/api/v1/tetragon"
	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/event-processor/api"
	detector_api "github.com/runtime-radar/runtime-radar/event-processor/detector/api"
	detector_tetragon_api "github.com/runtime-radar/runtime-radar/event-processor/detector/api/tetragon"
	"github.com/runtime-radar/runtime-radar/event-processor/pkg/model"
	"github.com/runtime-radar/runtime-radar/lib/security"
	enf_model "github.com/runtime-radar/runtime-radar/policy-enforcer/pkg/model"
	"github.com/tetratelabs/wazero"
	"google.golang.org/protobuf/proto"
)

// Key is used to identify particular detector. We do not allow same combination of ID and version, but we do allow same ID with different versions.
// This is in line with model.Detector, where primary key is also set to ID and version.
type Key struct {
	ID      string
	Version uint
}

// Wrapper represents a detector itself with info about it alongside, so that it can be used without calling Detector.Info.
type Wrapper struct {
	detector_api.Detector
	Info *api.Detector
}

type Chain map[Key]Wrapper

// ChainResult represents a result of execution of detectors' chain.
// If at least one of detectors returned an error, it's appended to Errors
// and it doesn't prevent remaining chain from being executed.
type ChainResult struct {
	Threats []*api.Threat
	Errors  []*api.DetectError
}

func NewPlugin(ctx context.Context) (*detector_api.DetectorPlugin, error) {
	mc := wazero.NewModuleConfig()

	p, err := detector_api.NewDetectorPlugin(ctx, detector_api.WazeroModuleConfig(mc))
	if err != nil {
		return nil, fmt.Errorf("can't instantiate plugin: %w", err)
	}

	return p, nil
}

func NewChain(ctx context.Context, plugin *detector_api.DetectorPlugin, binaries [][]byte) (Chain, error) {
	c := make(Chain, len(binaries))

	for _, binary := range binaries {
		d, err := plugin.LoadBinary(ctx, binary)
		if err != nil {
			return nil, fmt.Errorf("can't load binary: %w", err)
		}

		resp, err := d.Info(ctx, &detector_api.InfoReq{})
		if err != nil {
			return nil, fmt.Errorf("can't get detector info: %w", err)
		}

		k := Key{resp.GetId(), uint(resp.GetVersion())}
		if _, ok := c[k]; ok {
			return nil, fmt.Errorf("duplicate detector key: %+v", k)
		}

		info := &api.Detector{
			Id:          resp.GetId(),
			Name:        resp.GetName(),
			Description: resp.GetDescription(),
			Version:     resp.GetVersion(),
			Author:      resp.GetAuthor(),
			Contact:     resp.GetContact(),
			License:     resp.GetLicense(),
		}

		c[k] = Wrapper{d, info}
	}

	return c, nil
}

// Detect returns a result of executing chain of detectors.
// If detector resulted in error, an error is appended to result and remaining chain continues to execute.
// If an error occurred outside a detector, it's returned.
func (c Chain) Detect(ctx context.Context, event *tetragon_api.GetEventsResponse) (*ChainResult, error) {
	res := &ChainResult{}

	dEvent, err := convertEvent(event)
	if err != nil {
		return nil, fmt.Errorf("can't convert event: %w", err)
	}

	for _, d := range c {
		t0 := time.Now()

		resp, err := d.Detect(ctx, &detector_api.DetectReq{Event: dEvent})
		if err != nil {
			res.Errors = append(res.Errors, &api.DetectError{
				Detector: d.Info,
				Error:    err.Error(),
			})

			continue
		}

		log.Debug().Str("delay", time.Since(t0).String()).Interface("detector_event", dEvent).Interface("detector_resp", resp).Msgf("Detector[%s] is done", resp.GetId())

		if resp.GetSeverity() > detector_api.DetectResp_NONE {
			s := enf_model.Severity(resp.GetSeverity())

			res.Threats = append(res.Threats, &api.Threat{
				Detector: &api.Detector{
					Id:          resp.GetId(),
					Name:        resp.GetName(),
					Description: resp.GetDescription(),
					Version:     resp.GetVersion(),
					Author:      resp.GetAuthor(),
					Contact:     resp.GetContact(),
				},
				Severity: s.String(),
			})
		}
	}

	return res, nil
}

func ModelFromBinary(ctx context.Context, plugin *detector_api.DetectorPlugin, bin []byte) (*model.Detector, error) {
	detector, err := plugin.LoadBinary(ctx, bin)
	if err != nil {
		return nil, fmt.Errorf("can't load wasm module: %w", err)
	}

	infoResp, err := detector.Info(ctx, &detector_api.InfoReq{})
	if err != nil {
		return nil, fmt.Errorf("can't get detector info: %w", err)
	}

	d := &model.Detector{
		ID:          infoResp.GetId(),
		Name:        infoResp.GetName(),
		Description: infoResp.GetDescription(),
		Version:     uint(infoResp.GetVersion()),
		Author:      infoResp.GetAuthor(),
		Contact:     infoResp.GetContact(),
		License:     infoResp.GetLicense(),
		WasmBinary:  bin,
		WasmHash:    security.HashSHA512AsHex(bin),
	}

	return d, err
}

func convertEvent(in *tetragon_api.GetEventsResponse) (*detector_tetragon_api.GetEventsResponse, error) {
	b, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	out := &detector_tetragon_api.GetEventsResponse{}
	if err := out.UnmarshalVT(b); err != nil {
		return nil, err
	}

	return out, nil
}
