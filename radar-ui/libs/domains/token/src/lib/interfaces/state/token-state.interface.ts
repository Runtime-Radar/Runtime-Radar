import { EntityState } from '@ngrx/entity';

import { LoadStatus } from '@cs/core';

import { Token } from '../contract/token-contract.interface';

export type TokenEntityState = EntityState<Token>;

export interface TokenState {
    loadStatus: LoadStatus;
    lastUpdate: number;
    list: TokenEntityState;
}
