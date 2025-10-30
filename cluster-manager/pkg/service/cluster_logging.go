package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/cluster-manager/api"
	"github.com/runtime-radar/runtime-radar/lib/server/interceptor"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClusterLogging struct {
	api.ClusterControllerServer
}

func (cl *ClusterLogging) Create(ctx context.Context, req *api.Cluster) (resp *api.CreateClusterResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called ClusterControllerServer.Create")
	}(time.Now())

	resp, err = cl.ClusterControllerServer.Create(ctx, req)
	return
}

func (cl *ClusterLogging) Read(ctx context.Context, req *api.ReadClusterReq) (resp *api.ReadClusterResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called ClusterControllerServer.Read")
	}(time.Now())

	resp, err = cl.ClusterControllerServer.Read(ctx, req)
	return
}

func (cl *ClusterLogging) Update(ctx context.Context, req *api.Cluster) (resp *emptypb.Empty, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called ClusterControllerServer.Update")
	}(time.Now())

	resp, err = cl.ClusterControllerServer.Update(ctx, req)
	return
}

func (cl *ClusterLogging) Delete(ctx context.Context, req *api.DeleteClusterReq) (resp *emptypb.Empty, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called ClusterControllerServer.Delete")
	}(time.Now())

	resp, err = cl.ClusterControllerServer.Delete(ctx, req)
	return
}

func (cl *ClusterLogging) ListPage(ctx context.Context, req *api.ListClusterPageReq) (resp *api.ListClusterPageResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Int("result", len(resp.GetClusters())).
			Stringer("correlation_id", corrID).
			Msg("Called ClusterControllerServer.ListPage")
	}(time.Now())

	resp, err = cl.ClusterControllerServer.ListPage(ctx, req)
	return
}

func (cl *ClusterLogging) Register(ctx context.Context, req *api.RegisterClusterReq) (resp *emptypb.Empty, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called ClusterControllerServer.Register")
	}(time.Now())

	resp, err = cl.ClusterControllerServer.Register(ctx, req)
	return
}

func (cl *ClusterLogging) Unregister(ctx context.Context, req *api.UnregisterClusterReq) (resp *emptypb.Empty, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Interface("result", resp).
			Stringer("correlation_id", corrID).
			Msg("Called ClusterControllerServer.Unregister")
	}(time.Now())

	resp, err = cl.ClusterControllerServer.Unregister(ctx, req)
	return
}

func (cl *ClusterLogging) GenerateUninstallCmd(ctx context.Context, req *api.GenerateUninstallCmdReq) (resp *api.GenerateUninstallCmdResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Int("result", len(resp.GetCmd())).
			Stringer("correlation_id", corrID).
			Msg("Called ClusterControllerServer.GenerateUninstallCmd")
	}(time.Now())

	resp, err = cl.ClusterControllerServer.GenerateUninstallCmd(ctx, req)
	return
}

func (cl *ClusterLogging) GenerateInstallCmd(ctx context.Context, req *api.GenerateInstallCmdReq) (resp *api.GenerateInstallCmdResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Int("result", len(resp.GetCmd())).
			Stringer("correlation_id", corrID).
			Msg("Called ClusterControllerServer.GenerateInstallCmd")
	}(time.Now())

	resp, err = cl.ClusterControllerServer.GenerateInstallCmd(ctx, req)
	return
}

func (cl *ClusterLogging) GenerateValuesYAML(ctx context.Context, req *api.GenerateValuesReq) (resp *api.GenerateValuesResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Stringer("correlation_id", corrID).
			Msg("Called ClusterControllerServer.GenerateValuesYAML")
	}(time.Now())

	resp, err = cl.ClusterControllerServer.GenerateValuesYAML(ctx, req)
	return
}

func (cl *ClusterLogging) ListRegistered(ctx context.Context, req *emptypb.Empty) (resp *api.ListRegisteredResp, err error) {
	defer func(t0 time.Time) {
		corrID, _ := interceptor.CorrelationIDFromContext(ctx)

		log.Err(err).Str("delay", time.Since(t0).String()).
			Interface("args", req).
			Int("result", len(resp.GetClusters())).
			Stringer("correlation_id", corrID).
			Msg("Called ClusterControllerServer.ListRegistered")
	}(time.Now())

	resp, err = cl.ClusterControllerServer.ListRegistered(ctx, req)
	return
}
