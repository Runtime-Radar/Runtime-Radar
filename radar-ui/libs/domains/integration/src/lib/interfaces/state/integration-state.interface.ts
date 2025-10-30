import { EntityState } from '@ngrx/entity';

import { LoadStatus } from '@cs/core';

import {
    IntegrationEmail,
    IntegrationSyslog,
    IntegrationType,
    IntegrationWebhook
} from '../contract/integration-contract.interface';

export type IntegrationEmailEntityState = EntityState<IntegrationEmail>;

export type IntegrationSyslogEntityState = EntityState<IntegrationSyslog>;

export type IntegrationWebhookEntityState = EntityState<IntegrationWebhook>;

export interface IntegrationState {
    loadStatus: LoadStatus;
    loadedTypes: IntegrationType[];
    lastUpdate: number;
    email: IntegrationEmailEntityState;
    syslog: IntegrationSyslogEntityState;
    webhook: IntegrationWebhookEntityState;
}
