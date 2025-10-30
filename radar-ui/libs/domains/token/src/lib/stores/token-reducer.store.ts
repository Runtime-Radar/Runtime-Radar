import { createEntityAdapter } from '@ngrx/entity';
import { Action, ActionReducer, createReducer, on } from '@ngrx/store';

import { LoadStatus } from '@cs/core';

import {
    DELETE_ALL_TOKEN_ENTITIES_DOC_ACTION,
    DELETE_TOKEN_ENTITY_DOC_ACTION,
    SET_ALL_TOKEN_ENTITIES_DOC_ACTION,
    SET_TOKEN_ENTITY_DOC_ACTION,
    UPDATE_TOKEN_STATE_DOC_ACTION,
    UPSERT_TOKEN_ENTITIES_DOC_ACTION
} from './token-action.store';
import { Token, TokenState } from '../interfaces';

const adapter = createEntityAdapter<Token>();

const INITIAL_STATE: TokenState = {
    loadStatus: LoadStatus.INIT,
    lastUpdate: 0,
    list: adapter.getInitialState()
};

const reducer: ActionReducer<TokenState, Action> = createReducer(
    INITIAL_STATE,
    on(UPDATE_TOKEN_STATE_DOC_ACTION, (state, values) => ({ ...state, ...values })),
    on(SET_ALL_TOKEN_ENTITIES_DOC_ACTION, (state, { list }) => ({
        ...state,
        list: adapter.setAll(list, state.list)
    })),
    on(DELETE_ALL_TOKEN_ENTITIES_DOC_ACTION, (state) => ({
        ...state,
        list: adapter.removeAll(state.list)
    })),
    on(SET_TOKEN_ENTITY_DOC_ACTION, (state, { item }) => ({
        ...state,
        list: adapter.setOne(item, state.list)
    })),
    on(UPSERT_TOKEN_ENTITIES_DOC_ACTION, (state, { list }) => ({
        ...state,
        list: adapter.upsertMany(list, state.list)
    })),
    on(DELETE_TOKEN_ENTITY_DOC_ACTION, (state, { id }) => ({
        ...state,
        list: adapter.removeOne(id, state.list)
    }))
);

export const tokenEntitySelector = adapter.getSelectors();

export function tokenReducer(state: TokenState | undefined, action: Action): TokenState {
    return reducer(state, action);
}
