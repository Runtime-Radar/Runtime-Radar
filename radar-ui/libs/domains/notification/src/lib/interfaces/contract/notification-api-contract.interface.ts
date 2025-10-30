import {
    Notification,
    NotificationEmail,
    NotificationSyslog,
    NotificationWebhook
} from './notification-contract.interface';

export interface GetNotificationsRequest {
    integration_id: string;
    event_type: string;
}

export interface GetNotificationsResponse {
    notifications: Notification[];
}

export interface GetNotificationResponse {
    notification: Notification;
    deleted: boolean;
}

export interface GetNotificationTemplateResponse {
    template: string;
}

export type CreateNotificationRequest =
    | Omit<NotificationEmail, 'id'>
    | Omit<NotificationSyslog, 'id'>
    | Omit<NotificationWebhook, 'id'>;

export interface CreateNotificationResponse {
    id: string;
}

export type UpdateNotificationRequest = CreateNotificationRequest;

export type EmptyNotificationResponse = Record<string, unknown>;
