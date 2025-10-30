//go:build !tinygo.wasm

package processor

import (
	"context"
	"fmt"

	"github.com/runtime-radar/runtime-radar/event-processor/api"
	"github.com/runtime-radar/runtime-radar/event-processor/pkg/build"
	"github.com/runtime-radar/runtime-radar/lib/docker"
	enforcer_api "github.com/runtime-radar/runtime-radar/policy-enforcer/api"
)

func (wp *WorkersPool) evaluatePolicy(ctx context.Context, ed *eventData, threats []*api.Threat) (*enforcer_api.EvaluatePolicyRuntimeEventReq, error) {
	events := []*enforcer_api.EvaluatePolicyRuntimeEventReq_Result_Event{}
	for _, t := range threats {
		events = append(events, &enforcer_api.EvaluatePolicyRuntimeEventReq_Result_Event{
			DetectorId: t.GetDetector().GetId(),
			Severity:   t.GetSeverity(),
		})
	}

	image, registry := "", ""

	if ed.ContainerImage != "" {
		ref, err := docker.ParseReference(ed.ContainerImage)
		if err != nil {
			return nil, fmt.Errorf("can't parse image reference: %w", err)
		}

		image, registry = ref.Image, ref.Registry
	}

	req := &enforcer_api.EvaluatePolicyRuntimeEventReq{
		Actor: build.AppName,
		Action: &enforcer_api.EvaluatePolicyRuntimeEventReq_Action{
			Type: actionType,
			Args: &enforcer_api.EvaluatePolicyRuntimeEventReq_Action_Args{
				Namespace: ed.PodNamespace,
				Pod:       ed.PodName,
				Container: ed.ContainerName,
				Node:      ed.NodeName,
				ImageName: image,
				Registry:  registry,
				Binary:    ed.ProcessBinary,
			},
		},
		Result: &enforcer_api.EvaluatePolicyRuntimeEventReq_Result{
			Events: events,
		},
	}

	return wp.enforcer.EvaluatePolicyRuntimeEvent(ctx, req)
}
