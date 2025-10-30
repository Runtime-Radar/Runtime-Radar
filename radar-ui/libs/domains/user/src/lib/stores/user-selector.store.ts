import { ActionReducerMap, createFeatureSelector, createSelector } from '@ngrx/store';

import { UserEntityState, UserState } from '../interfaces/state/user-state.interface';
import { userEntitySelector, userReducer } from './user-reducer.store';

export const USER_DOMAIN_KEY = 'user';

export interface UserDomainState {
    readonly domain: UserState;
}

const selectUserDomainState = createFeatureSelector<UserDomainState>(USER_DOMAIN_KEY);
const selectUserState = createSelector(selectUserDomainState, (state: UserDomainState) => state.domain);
const selectUserEntityState = createSelector(selectUserState, (state: UserState) => state.list);

export const getUserLoadStatus = createSelector(selectUserState, (state: UserState) => state.loadStatus);

export const getUserLastUpdate = createSelector(selectUserState, (state: UserState) => state.lastUpdate);

export const getUsers = createSelector(selectUserEntityState, (state: UserEntityState) =>
    userEntitySelector.selectAll(state)
);

export const userDomainReducer: ActionReducerMap<UserDomainState> = {
    domain: userReducer
};
