import { createEntityAdapter } from '@ngrx/entity';
import { Action, ActionReducer, createReducer, on } from '@ngrx/store';

import { LoadStatus } from '@cs/core';

import { IntegrationState } from '../interfaces/state/integration-state.interface';
import {
    DELETE_INTEGRATION_ENTITY_DOC_ACTION,
    SET_ALL_INTEGRATION_ENTITIES_DOC_ACTION,
    SET_INTEGRATION_ENTITY_DOC_ACTION,
    SET_INTEGRATION_LOADED_TYPE_DOC_ACTION,
    UPDATE_INTEGRATION_ENTITY_DOC_ACTION,
    UPDATE_INTEGRATION_STATE_DOC_ACTION
} from './integration-action.store';
import { IntegrationEmail, IntegrationSyslog, IntegrationWebhook } from '../interfaces';

const emailAdapter = createEntityAdapter<IntegrationEmail>();
const syslogAdapter = createEntityAdapter<IntegrationSyslog>();
const webhookAdapter = createEntityAdapter<IntegrationWebhook>();

const INITIAL_STATE: IntegrationState = {
    loadStatus: LoadStatus.INIT,
    loadedTypes: [],
    lastUpdate: 0,
    email: emailAdapter.getInitialState(),
    syslog: syslogAdapter.getInitialState(),
    webhook: webhookAdapter.getInitialState()
};

const reducer: ActionReducer<IntegrationState, Action> = createReducer(
    INITIAL_STATE,
    on(UPDATE_INTEGRATION_STATE_DOC_ACTION, (state, values) => ({ ...state, ...values })),
    on(SET_INTEGRATION_LOADED_TYPE_DOC_ACTION, (state, { integrationType }) => ({
        ...state,
        loadedTypes: [...state.loadedTypes, integrationType]
    })),
    on(SET_ALL_INTEGRATION_ENTITIES_DOC_ACTION, (state, { email, syslog, webhook }) => ({
        ...state,
        email: email ? emailAdapter.setAll(email, state.email) : state.email,
        syslog: syslog ? syslogAdapter.setAll(syslog, state.syslog) : state.syslog,
        webhook: webhook ? webhookAdapter.setAll(webhook, state.webhook) : state.webhook
    })),
    on(SET_INTEGRATION_ENTITY_DOC_ACTION, (state, { email, syslog, webhook }) => ({
        ...state,
        email: email ? emailAdapter.setOne(email, state.email) : state.email,
        syslog: syslog ? syslogAdapter.setOne(syslog, state.syslog) : state.syslog,
        webhook: webhook ? webhookAdapter.setOne(webhook, state.webhook) : state.webhook
    })),
    on(UPDATE_INTEGRATION_ENTITY_DOC_ACTION, (state, { email, syslog, webhook }) => ({
        ...state,
        email: email ? emailAdapter.updateOne(email, state.email) : state.email,
        syslog: syslog ? syslogAdapter.updateOne(syslog, state.syslog) : state.syslog,
        webhook: webhook ? webhookAdapter.updateOne(webhook, state.webhook) : state.webhook
    })),
    on(DELETE_INTEGRATION_ENTITY_DOC_ACTION, (state, { emailId, syslogId, webhookId }) => ({
        ...state,
        email: emailId ? emailAdapter.removeOne(emailId, state.email) : state.email,
        syslog: syslogId ? syslogAdapter.removeOne(syslogId, state.syslog) : state.syslog,
        webhook: webhookId ? webhookAdapter.removeOne(webhookId, state.webhook) : state.webhook
    }))
);

export const integrationEmailEntitySelector = emailAdapter.getSelectors();
export const integrationSyslogEntitySelector = syslogAdapter.getSelectors();
export const integrationWebhookEntitySelector = webhookAdapter.getSelectors();

export function integrationReducer(state: IntegrationState | undefined, action: Action): IntegrationState {
    return reducer(state, action);
}
