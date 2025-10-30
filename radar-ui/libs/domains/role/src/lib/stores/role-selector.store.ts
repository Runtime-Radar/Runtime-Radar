import { ActionReducerMap, createFeatureSelector, createSelector } from '@ngrx/store';

import { RoleEntityState, RoleState } from '../interfaces/state/role-state.interface';
import { roleEntitySelector, roleReducer } from './role-reducer.store';

export const ROLE_DOMAIN_KEY = 'role';

export interface RoleDomainState {
    readonly domain: RoleState;
}

const selectRoleDomainState = createFeatureSelector<RoleDomainState>(ROLE_DOMAIN_KEY);
const selectRoleState = createSelector(selectRoleDomainState, (state: RoleDomainState) => state.domain);
const selectRoleEntityState = createSelector(selectRoleState, (state: RoleState) => state.list);

export const getRoleLoadStatus = createSelector(selectRoleState, (state: RoleState) => state.loadStatus);

export const getRoles = createSelector(selectRoleEntityState, (state: RoleEntityState) =>
    roleEntitySelector.selectAll(state)
);

export const getRole = (id: string) =>
    createSelector(selectRoleEntityState, (state: RoleEntityState) => state.entities[id]);

export const roleDomainReducer: ActionReducerMap<RoleDomainState> = {
    domain: roleReducer
};
