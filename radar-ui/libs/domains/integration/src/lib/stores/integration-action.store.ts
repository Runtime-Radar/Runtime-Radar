import { Update } from '@ngrx/entity';
import { createAction, props } from '@ngrx/store';

import { IntegrationState } from '../interfaces/state/integration-state.interface';
import {
    CreateIntegrationRequest,
    IntegrationEmail,
    IntegrationSyslog,
    IntegrationType,
    IntegrationWebhook,
    UpdateIntegrationRequest
} from '../interfaces';

type ActionPropsType<E, S, W> = Partial<{
    email: E;
    syslog: S;
    webhook: W;
}>;

export const LOAD_INTEGRATION_ENTITIES_TODO_ACTION = createAction('[Integration] Load');

export const POLLING_LOAD_INTEGRATION_ENTITIES_TODO_ACTION = createAction('[Integration] Polling Load');

export const APPLY_INTEGRATION_LOAD_STATUS_TODO_ACTION = createAction(
    '[Integration] Apply Load Status',
    props<{ isLoaded: boolean }>()
);

export const CREATE_EMAIL_INTEGRATION_ENTITY_TODO_ACTION = createAction(
    '[Integration] Create Email',
    props<{ item: CreateIntegrationRequest<IntegrationEmail> }>()
);

export const CREATE_SYSLOG_INTEGRATION_ENTITY_TODO_ACTION = createAction(
    '[Integration] Create Syslog',
    props<{ item: CreateIntegrationRequest<IntegrationSyslog> }>()
);

export const CREATE_WEBHOOK_INTEGRATION_ENTITY_TODO_ACTION = createAction(
    '[Integration] Create Webhook',
    props<{ item: CreateIntegrationRequest<IntegrationWebhook> }>()
);

export const UPDATE_EMAIL_INTEGRATION_ENTITY_TODO_ACTION = createAction(
    '[Integration] Update Email',
    props<{ id: string; item: UpdateIntegrationRequest<IntegrationEmail> }>()
);

export const UPDATE_SYSLOG_INTEGRATION_ENTITY_TODO_ACTION = createAction(
    '[Integration] Update Syslog',
    props<{ id: string; item: UpdateIntegrationRequest<IntegrationSyslog> }>()
);

export const UPDATE_WEBHOOK_INTEGRATION_ENTITY_TODO_ACTION = createAction(
    '[Integration] Update Webhook',
    props<{ id: string; item: UpdateIntegrationRequest<IntegrationWebhook> }>()
);

export const DELETE_EMAIL_INTEGRATION_ENTITY_TODO_ACTION = createAction(
    '[Integration] Delete Email',
    props<{ id: string }>()
);

export const DELETE_SYSLOG_INTEGRATION_ENTITY_TODO_ACTION = createAction(
    '[Integration] Delete Syslog',
    props<{ id: string }>()
);

export const DELETE_WEBHOOK_INTEGRATION_ENTITY_TODO_ACTION = createAction(
    '[Integration] Delete Webhook',
    props<{ id: string }>()
);

export const UPDATE_INTEGRATION_STATE_DOC_ACTION = createAction(
    '[Integration] (Doc) Update State',
    props<Partial<Pick<IntegrationState, 'loadStatus' | 'lastUpdate'>>>()
);

export const SET_INTEGRATION_LOADED_TYPE_DOC_ACTION = createAction(
    '[Integration] (Doc) Set Loaded Type',
    props<{ integrationType: IntegrationType }>()
);

export const SET_ALL_INTEGRATION_ENTITIES_DOC_ACTION = createAction(
    '[Integration] (Doc) Set All',
    props<ActionPropsType<IntegrationEmail[], IntegrationSyslog[], IntegrationWebhook[]>>()
);

export const SET_INTEGRATION_ENTITY_DOC_ACTION = createAction(
    '[Integration] (Doc) Set One',
    props<ActionPropsType<IntegrationEmail, IntegrationSyslog, IntegrationWebhook>>()
);

export const UPDATE_INTEGRATION_ENTITY_DOC_ACTION = createAction(
    '[Integration] (Doc) Update',
    props<ActionPropsType<Update<IntegrationEmail>, Update<IntegrationSyslog>, Update<IntegrationWebhook>>>()
);

export const DELETE_INTEGRATION_ENTITY_DOC_ACTION = createAction(
    '[Integration] (Doc) Delete',
    props<Partial<{ emailId: string; syslogId: string; webhookId: string }>>()
);

export const DELETE_CONNECTED_INTEGRATION_EVENT_ACTION = createAction(
    '[Integration] {Event} Delete',
    props<{ integrationId: string }>()
);
