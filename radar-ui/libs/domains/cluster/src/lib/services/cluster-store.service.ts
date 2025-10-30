import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Store } from '@ngrx/store';

import { LoadStatus } from '@cs/core';

import {
    CREATE_CLUSTER_ENTITY_TODO_ACTION,
    DELETE_CLUSTER_ENTITY_TODO_ACTION,
    SWITCH_CLUSTER_ENTITY_TODO_ACTION,
    UPDATE_CLUSTER_ENTITY_TODO_ACTION
} from '../stores/cluster-action.store';
import { Cluster, ClusterState, CreateClusterRequest, RegisteredCluster, UpdateClusterRequest } from '../interfaces';
import { getCluster, getClusterLoadStatus, getClusters, getRegisteredClusters } from '../stores/cluster-selector.store';

@Injectable({
    providedIn: 'root'
})
export class ClusterStoreService {
    readonly clusters$: Observable<Cluster[]> = this.store.select(getClusters);

    readonly registeredClusters$: Observable<RegisteredCluster[]> = this.store.select(getRegisteredClusters);

    readonly cluster$ = (id: string): Observable<Cluster | undefined> => this.store.select(getCluster(id));

    readonly loadStatus$: Observable<LoadStatus> = this.store.select(getClusterLoadStatus);

    constructor(private readonly store: Store<ClusterState>) {}

    createCluster(item: CreateClusterRequest) {
        this.store.dispatch(CREATE_CLUSTER_ENTITY_TODO_ACTION({ item }));
    }

    updateCluster(id: string, item: UpdateClusterRequest) {
        this.store.dispatch(UPDATE_CLUSTER_ENTITY_TODO_ACTION({ id, item }));
    }

    deleteCluster(id: string) {
        this.store.dispatch(DELETE_CLUSTER_ENTITY_TODO_ACTION({ id }));
    }

    switchCluster(url: string) {
        this.store.dispatch(SWITCH_CLUSTER_ENTITY_TODO_ACTION({ url }));
    }
}
