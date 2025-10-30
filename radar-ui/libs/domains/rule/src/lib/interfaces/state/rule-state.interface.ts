import { EntityState } from '@ngrx/entity';

import { LoadStatus } from '@cs/core';

import { Rule } from '../contract/rule-contract.interface';

export type RuleEntityState = EntityState<Rule>;

export interface RuleState {
    loadStatus: LoadStatus;
    lastUpdate: number;
    list: RuleEntityState;
}
