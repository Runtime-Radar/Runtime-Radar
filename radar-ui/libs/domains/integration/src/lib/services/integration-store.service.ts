import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Store } from '@ngrx/store';

import { LoadStatus } from '@cs/core';

import { IntegrationState } from '../interfaces/state/integration-state.interface';
import {
    CREATE_EMAIL_INTEGRATION_ENTITY_TODO_ACTION,
    CREATE_SYSLOG_INTEGRATION_ENTITY_TODO_ACTION,
    CREATE_WEBHOOK_INTEGRATION_ENTITY_TODO_ACTION,
    DELETE_EMAIL_INTEGRATION_ENTITY_TODO_ACTION,
    DELETE_SYSLOG_INTEGRATION_ENTITY_TODO_ACTION,
    DELETE_WEBHOOK_INTEGRATION_ENTITY_TODO_ACTION,
    UPDATE_EMAIL_INTEGRATION_ENTITY_TODO_ACTION,
    UPDATE_SYSLOG_INTEGRATION_ENTITY_TODO_ACTION,
    UPDATE_WEBHOOK_INTEGRATION_ENTITY_TODO_ACTION
} from '../stores/integration-action.store';
import {
    CreateIntegrationRequest,
    IntegrationEmail,
    IntegrationSyslog,
    IntegrationWebhook,
    UpdateIntegrationRequest
} from '../interfaces';
import {
    getEmailIntegrations,
    getIntegrationLoadStatus,
    getSyslogIntegrations,
    getWebhookIntegrations
} from '../stores/integration-selector.store';

@Injectable({
    providedIn: 'root'
})
export class IntegrationStoreService {
    readonly emailIntegrations$: Observable<IntegrationEmail[]> = this.store.select(getEmailIntegrations);

    readonly syslogIntegrations$: Observable<IntegrationSyslog[]> = this.store.select(getSyslogIntegrations);

    readonly webhookIntegrations$: Observable<IntegrationWebhook[]> = this.store.select(getWebhookIntegrations);

    readonly loadStatus$: Observable<LoadStatus> = this.store.select(getIntegrationLoadStatus);

    constructor(private readonly store: Store<IntegrationState>) {}

    createEmailIntegration(item: CreateIntegrationRequest<IntegrationEmail>) {
        this.store.dispatch(CREATE_EMAIL_INTEGRATION_ENTITY_TODO_ACTION({ item }));
    }

    createSyslogIntegration(item: CreateIntegrationRequest<IntegrationSyslog>) {
        this.store.dispatch(CREATE_SYSLOG_INTEGRATION_ENTITY_TODO_ACTION({ item }));
    }

    createWebhookIntegration(item: CreateIntegrationRequest<IntegrationWebhook>) {
        this.store.dispatch(CREATE_WEBHOOK_INTEGRATION_ENTITY_TODO_ACTION({ item }));
    }

    updateEmailIntegration(id: string, item: UpdateIntegrationRequest<IntegrationEmail>) {
        this.store.dispatch(UPDATE_EMAIL_INTEGRATION_ENTITY_TODO_ACTION({ id, item }));
    }

    updateSyslogIntegration(id: string, item: UpdateIntegrationRequest<IntegrationSyslog>) {
        this.store.dispatch(UPDATE_SYSLOG_INTEGRATION_ENTITY_TODO_ACTION({ id, item }));
    }

    updateWebhookIntegration(id: string, item: UpdateIntegrationRequest<IntegrationWebhook>) {
        this.store.dispatch(UPDATE_WEBHOOK_INTEGRATION_ENTITY_TODO_ACTION({ id, item }));
    }

    deleteEmailIntegration(id: string) {
        this.store.dispatch(DELETE_EMAIL_INTEGRATION_ENTITY_TODO_ACTION({ id }));
    }

    deleteSyslogIntegration(id: string) {
        this.store.dispatch(DELETE_SYSLOG_INTEGRATION_ENTITY_TODO_ACTION({ id }));
    }

    deleteWebhookIntegration(id: string) {
        this.store.dispatch(DELETE_WEBHOOK_INTEGRATION_ENTITY_TODO_ACTION({ id }));
    }
}
