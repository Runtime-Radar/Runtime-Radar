package convert

import (
	"github.com/runtime-radar/runtime-radar/public-api/pkg/model"
)

func AccessTokensToResponse(accessTokens []*model.AccessToken) []*model.AccessTokenResp {
	resps := make([]*model.AccessTokenResp, 0, len(accessTokens))
	for _, at := range accessTokens {
		resps = append(resps, &model.AccessTokenResp{
			at.ID,
			at.Name,
			at.UserID,
			at.Permissions,
			at.ExpiresAt,
			at.CreatedAt,
			at.InvalidatedAt,
		})
	}

	return resps
}
