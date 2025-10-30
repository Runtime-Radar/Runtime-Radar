import { ActionReducerMap, createFeatureSelector, createSelector } from '@ngrx/store';

import { ClusterEntityState, ClusterState, RegisteredClusterEntityState } from '../interfaces';
import { clusterEntitySelector, clusterReducer, clusterRegisteredEntitySelector } from './cluster-reducer.store';

export const CLUSTER_DOMAIN_KEY = 'cluster';

export interface ClusterDomainState {
    readonly domain: ClusterState;
}

const selectClusterDomainState = createFeatureSelector<ClusterDomainState>(CLUSTER_DOMAIN_KEY);
const selectClusterState = createSelector(selectClusterDomainState, (state: ClusterDomainState) => state.domain);
const selectClusterEntityState = createSelector(selectClusterState, (state: ClusterState) => state.list);
const selectRegisteredClusterEntityState = createSelector(
    selectClusterState,
    (state: ClusterState) => state.registeredClusters
);

export const getClusterLoadStatus = createSelector(selectClusterState, (state: ClusterState) => state.loadStatus);

export const getClusterLastUpdate = createSelector(selectClusterState, (state: ClusterState) => state.lastUpdate);

export const getClusters = createSelector(selectClusterEntityState, (state: ClusterEntityState) =>
    clusterEntitySelector.selectAll(state)
);

export const getRegisteredClusters = createSelector(
    selectRegisteredClusterEntityState,
    (state: RegisteredClusterEntityState) => clusterRegisteredEntitySelector.selectAll(state)
);

export const getCluster = (id: string) =>
    createSelector(selectClusterEntityState, (state: ClusterEntityState) => state.entities[id]);

export const clusterDomainReducer: ActionReducerMap<ClusterDomainState> = {
    domain: clusterReducer
};
