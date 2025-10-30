import { createAction, props } from '@ngrx/store';

import { DetectorType } from '@cs/domains/detector';

import { DetectorConfig, DetectorExtended } from '../interfaces/state/detector-state.interface';

export const LOAD_DETECTOR_ENTITIES_TODO_ACTION = createAction(
    '[Detector] Load',
    props<{ detectorType: DetectorType }>()
);

export const POLLING_LOAD_DETECTOR_ENTITIES_TODO_ACTION = createAction(
    '[Detector] Polling Load',
    props<{ detectorType: DetectorType }>()
);

export const CREATE_RUNTIME_DETECTOR_ENTITIES_TODO_ACTION = createAction(
    '[Detector] Create Runtime',
    props<{ base64list: string[] }>()
);

export const DELETE_RUNTIME_DETECTOR_ENTITY_TODO_ACTION = createAction(
    '[Detector] Delete Runtime',
    props<{ key: string; version: number }>()
);

export const SET_MANY_DETECTOR_CONFIG_ENTITIES_DOC_ACTION = createAction(
    '[Detector] (Doc) Set Many Config',
    props<{ config: DetectorConfig[] }>()
);

export const UPSERT_DETECTOR_CONFIG_ENTITY_DOC_ACTION = createAction(
    '[Detector] (Doc) Upsert',
    props<{ config: DetectorConfig }>()
);

export const SET_ALL_DETECTOR_ENTITIES_DOC_ACTION = createAction(
    '[Detector] (Doc) Set All',
    props<{ list: DetectorExtended[] }>()
);

export const SET_MANY_DETECTOR_ENTITIES_DOC_ACTION = createAction(
    '[Detector] (Doc) Set Many',
    props<{ list: DetectorExtended[] }>()
);

export const DELETE_DETECTOR_ENTITY_DOC_ACTION = createAction('[Detector] (Doc) Delete', props<{ id: string }>());
