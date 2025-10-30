package constructor

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	history_api "github.com/runtime-radar/runtime-radar/history-api/api"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/server/handler"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func RuntimeHistoryListEventsSlice(svc service.RuntimeHistory) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dir, ok := mux.Vars(r)["direction"]
		if !ok {
			handler.StatusJSONResp(w, status.New(codes.InvalidArgument, "direction is empty"))
			return
		}

		cursorParam := mux.Vars(r)["cursor"]
		cursor, err := time.Parse(time.RFC3339Nano, cursorParam)
		if err != nil {
			handler.StatusJSONResp(w, status.Newf(codes.InvalidArgument, "invalid cursor: %s", cursorParam))
			return
		}

		var sliceSize uint32
		if param := r.URL.Query().Get("slice_size"); param != "" {
			ss, err := strconv.ParseUint(param, 10, 32)
			if err != nil {
				handler.StatusJSONResp(w, status.Newf(codes.InvalidArgument, "invalid slice_size: %s", param))
				return
			}

			sliceSize = uint32(ss)
		}

		req := &history_api.ListRuntimeEventSliceReq{
			Cursor:    timestamppb.New(cursor),
			Direction: dir,
			SliceSize: sliceSize,
		}
		resp, err := svc.ListRuntimeEventSlice(r.Context(), req)
		if err != nil {
			if st, ok := status.FromError(err); ok {
				handler.StatusJSONResp(w, st)
				return
			}

			handler.StatusJSONResp(w, status.Newf(codes.Internal, "can't list runtime event slice: %v", err))
			return
		}

		handler.SendProtoResp(w, resp)
	})
}
