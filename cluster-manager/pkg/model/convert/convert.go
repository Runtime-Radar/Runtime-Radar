package convert

import (
	"fmt"

	"github.com/runtime-radar/runtime-radar/cluster-manager/api"
	"github.com/runtime-radar/runtime-radar/cluster-manager/pkg/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ClusterToProto(c *model.Cluster, maskSensitive bool) *api.Cluster {
	if maskSensitive {
		c.Config.MaskSensitive()
	}

	return &api.Cluster{
		Id:        c.ID.String(),
		Name:      c.Name,
		Config:    (*api.Cluster_Config)(c.Config),
		Status:    ClusterStatusToProto(c.Status),
		CreatedAt: timestamppb.New(c.CreatedAt),
	}
}

func ClustersToProto(cs []*model.Cluster, maskSensitive bool) []*api.Cluster {
	pbs := make([]*api.Cluster, 0, len(cs))
	for _, c := range cs {
		pbs = append(pbs, ClusterToProto(c, maskSensitive))
	}
	return pbs
}

func ClusterStatusFromProto(proto api.Cluster_Status) model.ClusterStatus {
	switch proto {
	case api.Cluster_STATUS_UNREGISTERED:
		return model.ClusterStatusUnregistered
	case api.Cluster_STATUS_REGISTERED:
		return model.ClusterStatusRegistered
	default: // normally should not happen
		panic(fmt.Sprintf("invalid cluster status given: %s", proto))
	}
}

func ClusterStatusToProto(cs model.ClusterStatus) api.Cluster_Status {
	switch cs {
	case model.ClusterStatusUnregistered:
		return api.Cluster_STATUS_UNREGISTERED
	case model.ClusterStatusRegistered:
		return api.Cluster_STATUS_REGISTERED
	default: // normally should not happen
		panic(fmt.Sprintf("invalid cluster status given: %s", cs))
	}
}
