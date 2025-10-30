import { EntityState } from '@ngrx/entity';

import { LoadStatus } from '@cs/core';

import { Notification } from '../contract/notification-contract.interface';

export type NotificationEntityState = EntityState<Notification>;

export interface NotificationState {
    loadStatus: LoadStatus;
    lastUpdate: number;
    list: NotificationEntityState;
}
