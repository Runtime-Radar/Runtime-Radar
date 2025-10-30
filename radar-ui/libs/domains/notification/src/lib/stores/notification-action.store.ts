import { Update } from '@ngrx/entity';
import { createAction, props } from '@ngrx/store';

import { Notification } from '../interfaces/contract/notification-contract.interface';
import { NotificationState } from '../interfaces/state/notification-state.interface';
import { CreateNotificationRequest, UpdateNotificationRequest } from '../interfaces';

export const LOAD_NOTIFICATION_ENTITIES_TODO_ACTION = createAction('[Notification] Load');

export const POLLING_LOAD_NOTIFICATION_ENTITIES_TODO_ACTION = createAction('[Notification] Polling Load');

export const CREATE_NOTIFICATION_ENTITY_TODO_ACTION = createAction(
    '[Notification] Create',
    props<{ item: CreateNotificationRequest }>()
);

export const UPDATE_NOTIFICATION_ENTITY_TODO_ACTION = createAction(
    '[Notification] Update',
    props<{ id: string; item: UpdateNotificationRequest }>()
);

export const DELETE_NOTIFICATION_ENTITY_TODO_ACTION = createAction('[Notification] Delete', props<{ id: string }>());

export const UPDATE_NOTIFICATION_STATE_DOC_ACTION = createAction(
    '[Notification] (Doc) Update State',
    props<Partial<Omit<NotificationState, 'list'>>>()
);

export const SET_ALL_NOTIFICATION_ENTITIES_DOC_ACTION = createAction(
    '[Notification] (Doc) Set All',
    props<{ list: Notification[] }>()
);

export const SET_NOTIFICATION_ENTITY_DOC_ACTION = createAction(
    '[Notification] (Doc) Set One',
    props<{ item: Notification }>()
);

export const UPDATE_NOTIFICATION_ENTITY_DOC_ACTION = createAction(
    '[Notification] (Doc) Update',
    props<{ item: Update<Notification> }>()
);

export const DELETE_NOTIFICATION_ENTITIES_DOC_ACTION = createAction(
    '[Notification] (Doc) Delete Many',
    props<{ ids: string[] }>()
);

export const DELETE_NOTIFICATION_ENTITY_DOC_ACTION = createAction(
    '[Notification] (Doc) Delete',
    props<{ id: string }>()
);
