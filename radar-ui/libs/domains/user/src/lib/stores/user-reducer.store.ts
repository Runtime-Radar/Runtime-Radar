import { createEntityAdapter } from '@ngrx/entity';
import { Action, ActionReducer, createReducer, on } from '@ngrx/store';

import { LoadStatus } from '@cs/core';

import { User } from '../interfaces';
import { UserState } from '../interfaces/state/user-state.interface';
import {
    DELETE_USER_ENTITIES_DOC_ACTION,
    DELETE_USER_ENTITY_DOC_ACTION,
    SET_ALL_USER_ENTITIES_DOC_ACTION,
    SET_USER_ENTITY_DOC_ACTION,
    UPDATE_USER_ENTITY_DOC_ACTION,
    UPDATE_USER_STATE_DOC_ACTION
} from './user-action.store';

const adapter = createEntityAdapter<User>();

const INITIAL_STATE: UserState = {
    loadStatus: LoadStatus.INIT,
    lastUpdate: 0,
    list: adapter.getInitialState()
};

const reducer: ActionReducer<UserState, Action> = createReducer(
    INITIAL_STATE,
    on(UPDATE_USER_STATE_DOC_ACTION, (state, values) => ({ ...state, ...values })),
    on(SET_ALL_USER_ENTITIES_DOC_ACTION, (state, { list }) => ({
        ...state,
        list: adapter.setAll(list, state.list)
    })),
    on(SET_USER_ENTITY_DOC_ACTION, (state, { item }) => ({
        ...state,
        list: adapter.setOne(item, state.list)
    })),
    on(UPDATE_USER_ENTITY_DOC_ACTION, (state, { item }) => ({
        ...state,
        list: adapter.updateOne(item, state.list)
    })),
    on(DELETE_USER_ENTITY_DOC_ACTION, (state, { id }) => ({
        ...state,
        list: adapter.removeOne(id, state.list)
    })),
    on(DELETE_USER_ENTITIES_DOC_ACTION, (state, { ids }) => ({
        ...state,
        list: adapter.removeMany(ids, state.list)
    }))
);

export const userEntitySelector = adapter.getSelectors();

export function userReducer(state: UserState | undefined, action: Action): UserState {
    return reducer(state, action);
}
