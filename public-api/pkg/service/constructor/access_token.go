package constructor

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/runtime-radar/runtime-radar/lib/errcommon"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/model"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/server/handler"
	"github.com/runtime-radar/runtime-radar/public-api/pkg/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func AccessTokenCreate(svc service.AccessToken) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &model.CreateAccessTokenReq{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			status := errcommon.StatusWithReason(codes.InvalidArgument, errcommon.CodeBadRequest, err.Error())
			handler.StatusJSONResp(w, status)
			return
		}

		id, token, err := svc.Create(r.Context(), req)
		if err != nil {
			if st, ok := status.FromError(err); ok {
				handler.StatusJSONResp(w, st)
				return
			}
			st := status.Newf(codes.Internal, "can't create rule: %v", err)
			handler.StatusJSONResp(w, st)

			return
		}

		handler.SendJSONResp(w, &model.CreateAccessTokenResp{id, token})
	})
}

func AccessTokenListPage(svc service.AccessToken) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageSize := service.DefaultPageSize
		order := service.DefaultOrder

		var err error
		query := r.URL.Query()

		pageNumber, err := strconv.Atoi(mux.Vars(r)["page_num"])
		if err != nil {
			status := errcommon.StatusWithReason(codes.InvalidArgument, errcommon.CodeBadRequest, fmt.Sprintf("invalid page_num: %v", err))
			handler.StatusJSONResp(w, status)
			return
		}
		if pageNumber <= 0 {
			status := errcommon.StatusWithReason(codes.InvalidArgument, errcommon.CodeBadRequest, "page_num must be positive")
			handler.StatusJSONResp(w, status)
			return
		}

		if pageSizeParam := query.Get("page_size"); pageSizeParam != "" {
			pageSize, err = strconv.Atoi(pageSizeParam)
			if err != nil {
				status := errcommon.StatusWithReason(codes.InvalidArgument, errcommon.CodeBadRequest, fmt.Sprintf("invalid page_size: %v", err))
				handler.StatusJSONResp(w, status)
				return
			}

			if pageSize <= 0 {
				status := errcommon.StatusWithReason(codes.InvalidArgument, errcommon.CodeBadRequest, "page_size must be positive")
				handler.StatusJSONResp(w, status)
				return
			}
		}

		if orderParam := query.Get("order"); orderParam != "" {
			order = orderParam
		}

		tokens, total, err := svc.ListPage(r.Context(), pageNumber, pageSize, order)
		if err != nil {
			if st, ok := status.FromError(err); ok {
				handler.StatusJSONResp(w, st)
				return
			}
			st := status.Newf(codes.Internal, "can't get list rule: %v", err)
			handler.StatusJSONResp(w, st)

			return
		}

		handler.SendJSONResp(w, &model.ListAccessTokenResp{total, tokens})
	})
}

func AccessTokenGetByID(svc service.AccessToken) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(mux.Vars(r)["id"])
		if err != nil {
			handler.StatusJSONResp(w, errcommon.StatusWithReason(codes.InvalidArgument, errcommon.CodeBadRequest, "id must be valid UUID"))
			return
		}

		token, err := svc.GetByID(r.Context(), id)
		if err != nil {
			if st, ok := status.FromError(err); ok {
				handler.StatusJSONResp(w, st)
				return
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				status := errcommon.StatusWithReason(codes.NotFound, errcommon.CodeNotFound, "record not found")
				handler.StatusJSONResp(w, status)
				return
			}
			st := status.Newf(codes.Internal, "can't get access token: %v", err)
			handler.StatusJSONResp(w, st)

			return
		}

		handler.SendJSONResp(w, token)
	})
}

func AccessTokenDelete(svc service.AccessToken) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(mux.Vars(r)["id"])
		if err != nil {
			status := errcommon.StatusWithReason(codes.InvalidArgument, errcommon.CodeBadRequest, "id must be valid UUID")
			handler.StatusJSONResp(w, status)
			return
		}

		err = svc.Delete(r.Context(), id)
		if err != nil {
			if st, ok := status.FromError(err); ok {
				handler.StatusJSONResp(w, st)
				return
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				status := errcommon.StatusWithReason(codes.NotFound, errcommon.CodeNotFound, "record not found")
				handler.StatusJSONResp(w, status)
				return
			}
			st := status.Newf(codes.Internal, "can't delete rule: %v", err)
			handler.StatusJSONResp(w, st)

			return
		}

		handler.SendJSONResp(w, struct{}{})
	})
}

func AccessTokenInvalidateAll(svc service.AccessToken) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := svc.InvalidateAll(r.Context())
		if err != nil {
			if st, ok := status.FromError(err); ok {
				handler.StatusJSONResp(w, st)
				return
			}
			st := status.Newf(codes.Internal, "can't invalidate tokens: %v", err)
			handler.StatusJSONResp(w, st)
			return
		}

		handler.SendJSONResp(w, struct{}{})
	})
}
