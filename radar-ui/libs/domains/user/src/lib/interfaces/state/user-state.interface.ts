import { EntityState } from '@ngrx/entity';

import { LoadStatus } from '@cs/core';

import { User } from '../contract/user-contract.interface';

export type UserEntityState = EntityState<User>;

export interface UserState {
    loadStatus: LoadStatus;
    lastUpdate: number;
    list: UserEntityState;
}
