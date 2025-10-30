import { ActionReducerMap, createFeatureSelector, createSelector } from '@ngrx/store';

import { TokenEntityState, TokenState } from '../interfaces';
import { tokenEntitySelector, tokenReducer } from './token-reducer.store';

export const TOKEN_DOMAIN_KEY = 'token';

export interface TokenDomainState {
    readonly domain: TokenState;
}

const selectTokenDomainState = createFeatureSelector<TokenDomainState>(TOKEN_DOMAIN_KEY);
const selectTokenState = createSelector(selectTokenDomainState, (state: TokenDomainState) => state.domain);
const selectTokenEntityState = createSelector(selectTokenState, (state: TokenState) => state.list);

export const getTokenLoadStatus = createSelector(selectTokenState, (state: TokenState) => state.loadStatus);

export const getTokenLastUpdate = createSelector(selectTokenState, (state: TokenState) => state.lastUpdate);

export const getTokens = createSelector(selectTokenEntityState, (state: TokenEntityState) =>
    tokenEntitySelector.selectAll(state)
);

export const tokenDomainReducer: ActionReducerMap<TokenDomainState> = {
    domain: tokenReducer
};
