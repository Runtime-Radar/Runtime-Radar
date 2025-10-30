import { Update } from '@ngrx/entity';
import { createAction, props } from '@ngrx/store';

import { Cluster, ClusterState, CreateClusterRequest, RegisteredCluster, UpdateClusterRequest } from '../interfaces';

export const LOAD_CLUSTER_ENTITIES_TODO_ACTION = createAction('[Cluster] Load');

export const POLLING_LOAD_CLUSTER_ENTITIES_TODO_ACTION = createAction('[Cluster] Polling Load');

export const CREATE_CLUSTER_ENTITY_TODO_ACTION = createAction(
    '[Cluster] Create',
    props<{ item: CreateClusterRequest }>()
);

export const UPDATE_CLUSTER_ENTITY_TODO_ACTION = createAction(
    '[Cluster] Update',
    props<{ id: string; item: UpdateClusterRequest }>()
);

export const DELETE_CLUSTER_ENTITY_TODO_ACTION = createAction('[Cluster] Delete', props<{ id: string }>());

export const SWITCH_CLUSTER_ENTITY_TODO_ACTION = createAction('[Cluster] Switch', props<{ url: string }>());

export const VALIDATE_CLUSTER_HOST_TODO_ACTION = createAction(
    '[Cluster] Validate Host',
    props<{ list: RegisteredCluster[] }>()
);

export const UPDATE_CLUSTER_STATE_DOC_ACTION = createAction(
    '[Cluster] (Doc) Update State',
    props<Partial<Omit<ClusterState, 'list'>>>()
);

export const SET_ALL_CLUSTER_ENTITIES_DOC_ACTION = createAction(
    '[Cluster] (Doc) Set All',
    props<{ list: Cluster[] }>()
);

export const SET_ALL_REGISTERED_CLUSTER_ENTITIES_DOC_ACTION = createAction(
    '[Cluster] (Doc) Set All Registered',
    props<{ list: RegisteredCluster[] }>()
);

export const SET_CLUSTER_ENTITY_DOC_ACTION = createAction('[Cluster] (Doc) Set One', props<{ item: Cluster }>());

export const UPDATE_CLUSTER_ENTITY_DOC_ACTION = createAction(
    '[Cluster] (Doc) Update',
    props<{ item: Update<Cluster> }>()
);

export const UPDATE_REGISTERED_CLUSTER_ENTITY_DOC_ACTION = createAction(
    '[Cluster] (Doc) Update Registered',
    props<{ item: Update<RegisteredCluster> }>()
);

export const DELETE_CLUSTER_ENTITY_DOC_ACTION = createAction('[Cluster] (Doc) Delete', props<{ id: string }>());

export const SWITCH_CLUSTER_EVENT_ACTION = createAction('[Cluster] {Event} Switch');
