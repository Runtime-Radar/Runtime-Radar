import { createEntityAdapter } from '@ngrx/entity';
import { Action, ActionReducer, createReducer, on } from '@ngrx/store';

import { LoadStatus } from '@cs/core';

import { Role } from '../interfaces';
import { RoleState } from '../interfaces/state/role-state.interface';
import { SET_ALL_ROLE_ENTITIES_DOC_ACTION, UPDATE_ROLE_LOAD_STATUS_DOC_ACTION } from './role-action.store';

const adapter = createEntityAdapter<Role>();

const INITIAL_STATE: RoleState = {
    loadStatus: LoadStatus.INIT,
    list: adapter.getInitialState()
};

const reducer: ActionReducer<RoleState, Action> = createReducer(
    INITIAL_STATE,
    on(UPDATE_ROLE_LOAD_STATUS_DOC_ACTION, (state, { loadStatus }) => ({ ...state, loadStatus })),
    on(SET_ALL_ROLE_ENTITIES_DOC_ACTION, (state, { list }) => ({
        ...state,
        list: adapter.setAll(list, state.list)
    }))
);

export const roleEntitySelector = adapter.getSelectors();

export function roleReducer(state: RoleState | undefined, action: Action): RoleState {
    return reducer(state, action);
}
