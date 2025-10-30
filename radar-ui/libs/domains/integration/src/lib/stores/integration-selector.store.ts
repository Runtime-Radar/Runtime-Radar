import { ActionReducerMap, createFeatureSelector, createSelector } from '@ngrx/store';

import {
    IntegrationEmailEntityState,
    IntegrationState,
    IntegrationSyslogEntityState,
    IntegrationWebhookEntityState
} from '../interfaces/state/integration-state.interface';
import {
    integrationEmailEntitySelector,
    integrationReducer,
    integrationSyslogEntitySelector,
    integrationWebhookEntitySelector
} from './integration-reducer.store';

export const INTEGRATION_DOMAIN_KEY = 'integration';

export interface IntegrationDomainState {
    readonly domain: IntegrationState;
}

const selectIntegrationDomainState = createFeatureSelector<IntegrationDomainState>(INTEGRATION_DOMAIN_KEY);
const selectIntegrationState = createSelector(
    selectIntegrationDomainState,
    (state: IntegrationDomainState) => state.domain
);
const selectIntegrationEmailEntityState = createSelector(
    selectIntegrationState,
    (state: IntegrationState) => state.email
);
const selectIntegrationSyslogEntityState = createSelector(
    selectIntegrationState,
    (state: IntegrationState) => state.syslog
);
const selectIntegrationWebhookEntityState = createSelector(
    selectIntegrationState,
    (state: IntegrationState) => state.webhook
);

export const getIntegrationLoadStatus = createSelector(
    selectIntegrationState,
    (state: IntegrationState) => state.loadStatus
);

export const getIntegrationLoadedTypes = createSelector(
    selectIntegrationState,
    (state: IntegrationState) => state.loadedTypes
);

export const getIntegrationLastUpdate = createSelector(
    selectIntegrationState,
    (state: IntegrationState) => state.lastUpdate
);

export const getEmailIntegrations = createSelector(
    selectIntegrationEmailEntityState,
    (state: IntegrationEmailEntityState) => integrationEmailEntitySelector.selectAll(state)
);

export const getSyslogIntegrations = createSelector(
    selectIntegrationSyslogEntityState,
    (state: IntegrationSyslogEntityState) => integrationSyslogEntitySelector.selectAll(state)
);

export const getWebhookIntegrations = createSelector(
    selectIntegrationWebhookEntityState,
    (state: IntegrationWebhookEntityState) => integrationWebhookEntitySelector.selectAll(state)
);

export const integrationDomainReducer: ActionReducerMap<IntegrationDomainState> = {
    domain: integrationReducer
};
