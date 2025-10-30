package service

import (
	"context"

	history_api "github.com/runtime-radar/runtime-radar/history-api/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const runtimeHistoryMaxSliceSize = 1000

type RuntimeHistoryGeneric struct {
	history_api.RuntimeHistoryClient
}

func (rhg *RuntimeHistoryGeneric) ListRuntimeEventSlice(ctx context.Context, req *history_api.ListRuntimeEventSliceReq, opts ...grpc.CallOption) (*history_api.ListRuntimeEventSliceResp, error) {
	if req.GetSliceSize() > runtimeHistoryMaxSliceSize {
		return nil, status.Newf(codes.InvalidArgument, "slice_size is more than %d", runtimeHistoryMaxSliceSize).Err()
	}

	return rhg.RuntimeHistoryClient.ListRuntimeEventSlice(ctx, req, opts...)
}
