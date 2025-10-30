import { createEntityAdapter } from '@ngrx/entity';
import { Action, ActionReducer, createReducer, on } from '@ngrx/store';

import {
    DELETE_DETECTOR_ENTITY_DOC_ACTION,
    SET_ALL_DETECTOR_ENTITIES_DOC_ACTION,
    SET_MANY_DETECTOR_CONFIG_ENTITIES_DOC_ACTION,
    SET_MANY_DETECTOR_ENTITIES_DOC_ACTION,
    UPSERT_DETECTOR_CONFIG_ENTITY_DOC_ACTION
} from './detector-action.store';
import { DetectorConfig, DetectorExtended, DetectorState } from '../interfaces/state/detector-state.interface';

const configAdapter = createEntityAdapter<DetectorConfig>();
const listAdapter = createEntityAdapter<DetectorExtended>();

const INITIAL_STATE: DetectorState = {
    config: configAdapter.getInitialState(),
    list: listAdapter.getInitialState()
};

const reducer: ActionReducer<DetectorState, Action> = createReducer(
    INITIAL_STATE,
    on(SET_MANY_DETECTOR_CONFIG_ENTITIES_DOC_ACTION, (state, { config }) => ({
        ...state,
        config: configAdapter.setMany(config, state.config)
    })),
    on(UPSERT_DETECTOR_CONFIG_ENTITY_DOC_ACTION, (state, { config }) => ({
        ...state,
        config: configAdapter.upsertOne(config, state.config)
    })),
    on(SET_ALL_DETECTOR_ENTITIES_DOC_ACTION, (state, { list }) => ({
        ...state,
        list: listAdapter.setAll(list, state.list)
    })),
    on(SET_MANY_DETECTOR_ENTITIES_DOC_ACTION, (state, { list }) => ({
        ...state,
        list: listAdapter.setMany(list, state.list)
    })),
    on(DELETE_DETECTOR_ENTITY_DOC_ACTION, (state, { id }) => ({
        ...state,
        list: listAdapter.removeOne(id, state.list)
    }))
);

export const detectorListEntitySelector = listAdapter.getSelectors();

export function detectorReducer(state: DetectorState | undefined, action: Action): DetectorState {
    return reducer(state, action);
}
