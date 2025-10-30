import { EntityState } from '@ngrx/entity';

import { LoadStatus } from '@cs/core';

import { Cluster, RegisteredCluster } from '../contract/cluster-contract.interface';

export type ClusterEntityState = EntityState<Cluster>;
export type RegisteredClusterEntityState = EntityState<RegisteredCluster>;

export interface ClusterState {
    loadStatus: LoadStatus;
    lastUpdate: number;
    list: ClusterEntityState;
    registeredClusters: RegisteredClusterEntityState;
}
