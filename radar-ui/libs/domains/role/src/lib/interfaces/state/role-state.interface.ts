import { EntityState } from '@ngrx/entity';

import { LoadStatus } from '@cs/core';

import { Role } from '../contract/role-contract.interface';

export type RoleEntityState = EntityState<Role>;

export interface RoleState {
    loadStatus: LoadStatus;
    list: RoleEntityState;
}
