import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Store } from '@ngrx/store';

import { RuleType } from '@cs/domains/rule';

import { NotificationState } from '../interfaces/state/notification-state.interface';
import {
    CREATE_NOTIFICATION_ENTITY_TODO_ACTION,
    DELETE_NOTIFICATION_ENTITY_TODO_ACTION,
    UPDATE_NOTIFICATION_ENTITY_TODO_ACTION
} from '../stores/notification-action.store';
import { CreateNotificationRequest, Notification, UpdateNotificationRequest } from '../interfaces';
import {
    getNotifications,
    getNotificationsByEventType,
    getNotificationsByIntegrationId
} from '../stores/notification-selector.store';

@Injectable({
    providedIn: 'root'
})
export class NotificationStoreService {
    readonly notifications$: Observable<Notification[]> = this.store.select(getNotifications);

    readonly notificationsByEventType$ = (type: RuleType): Observable<Notification[]> =>
        this.store.select(getNotificationsByEventType(type));

    readonly notificationsByIntegrationId$ = (id: string): Observable<Notification[]> =>
        this.store.select(getNotificationsByIntegrationId(id));

    constructor(private readonly store: Store<NotificationState>) {}

    createNotification(item: CreateNotificationRequest) {
        this.store.dispatch(CREATE_NOTIFICATION_ENTITY_TODO_ACTION({ item }));
    }

    updateNotification(id: string, item: UpdateNotificationRequest) {
        this.store.dispatch(UPDATE_NOTIFICATION_ENTITY_TODO_ACTION({ id, item }));
    }

    deleteNotification(id: string) {
        this.store.dispatch(DELETE_NOTIFICATION_ENTITY_TODO_ACTION({ id }));
    }
}
