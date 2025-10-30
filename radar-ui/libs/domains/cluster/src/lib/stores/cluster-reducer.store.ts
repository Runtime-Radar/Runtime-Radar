import { createEntityAdapter } from '@ngrx/entity';
import { Action, ActionReducer, createReducer, on } from '@ngrx/store';

import { LoadStatus } from '@cs/core';

import { Cluster, ClusterState, RegisteredCluster } from '../interfaces';
import {
    DELETE_CLUSTER_ENTITY_DOC_ACTION,
    SET_ALL_CLUSTER_ENTITIES_DOC_ACTION,
    SET_ALL_REGISTERED_CLUSTER_ENTITIES_DOC_ACTION,
    SET_CLUSTER_ENTITY_DOC_ACTION,
    UPDATE_CLUSTER_ENTITY_DOC_ACTION,
    UPDATE_CLUSTER_STATE_DOC_ACTION,
    UPDATE_REGISTERED_CLUSTER_ENTITY_DOC_ACTION
} from './cluster-action.store';

const adapter = createEntityAdapter<Cluster>();
const registeredClusterAdapter = createEntityAdapter<RegisteredCluster>();

const INITIAL_STATE: ClusterState = {
    loadStatus: LoadStatus.INIT,
    lastUpdate: 0,
    list: adapter.getInitialState(),
    registeredClusters: registeredClusterAdapter.getInitialState()
};

const reducer: ActionReducer<ClusterState, Action> = createReducer(
    INITIAL_STATE,
    on(UPDATE_CLUSTER_STATE_DOC_ACTION, (state, values) => ({ ...state, ...values })),
    on(SET_ALL_CLUSTER_ENTITIES_DOC_ACTION, (state, { list }) => ({
        ...state,
        list: adapter.setAll(list, state.list)
    })),
    on(SET_CLUSTER_ENTITY_DOC_ACTION, (state, { item }) => ({
        ...state,
        list: adapter.setOne(item, state.list)
    })),
    on(UPDATE_CLUSTER_ENTITY_DOC_ACTION, (state, { item }) => ({
        ...state,
        list: adapter.updateOne(item, state.list)
    })),
    on(DELETE_CLUSTER_ENTITY_DOC_ACTION, (state, { id }) => ({
        ...state,
        list: adapter.removeOne(id, state.list)
    })),
    on(SET_ALL_REGISTERED_CLUSTER_ENTITIES_DOC_ACTION, (state, { list }) => ({
        ...state,
        registeredClusters: registeredClusterAdapter.setAll(list, state.registeredClusters)
    })),
    on(UPDATE_REGISTERED_CLUSTER_ENTITY_DOC_ACTION, (state, { item }) => ({
        ...state,
        registeredClusters: registeredClusterAdapter.updateOne(item, state.registeredClusters)
    }))
);

export const clusterEntitySelector = adapter.getSelectors();
export const clusterRegisteredEntitySelector = registeredClusterAdapter.getSelectors();

export function clusterReducer(state: ClusterState | undefined, action: Action): ClusterState {
    return reducer(state, action);
}
