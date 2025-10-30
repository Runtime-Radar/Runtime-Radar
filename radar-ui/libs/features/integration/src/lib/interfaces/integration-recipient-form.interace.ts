import { NotificationEventType } from '@cs/domains/notification';

export interface IntegrationRecipientTemplateHeaderForm {
    key: string;
    value: string;
}

export type IntegrationRecipientTemplateRecord = {
    [key: string]: IntegrationRecipientTemplateHeaderForm;
};

export interface IntegrationRecipientForm {
    name: string;
    recipients: string[];
    eventType: NotificationEventType;
    clusterId: string;
    clusterUrl: string;
    clusterName: string;
    centralUrl: string;
    template: string;
    isTemplateDefault: boolean;
    subjectTemplate: string; // email
    path: string; // webhook
    header: IntegrationRecipientTemplateRecord; // webhook
}
