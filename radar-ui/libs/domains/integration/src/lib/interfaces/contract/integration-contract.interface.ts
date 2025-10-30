interface AbstractIntegration {
    id: string;
    type: IntegrationType;
    name: string;
    skip_check: boolean;
}

export enum IntegrationType {
    EMAIL = 'email',
    SYSLOG = 'syslog',
    WEBHOOK = 'webhook'
}

export enum IntegrationEmailAuthType {
    NONE = 'AUTH_TYPE_NONE',
    CRAM_MD5 = 'AUTH_TYPE_CRAM_MD5',
    LOGIN = 'AUTH_TYPE_LOGIN',
    PLAIN = 'AUTH_TYPE_PLAIN'
}

export interface IntegrationEmailEntity {
    auth_type: IntegrationEmailAuthType;
    from: string;
    server: string;
    username: string;
    password: string;
    ca: string;
    use_tls: boolean;
    use_start_tls: boolean;
    insecure: boolean;
}

export interface IntegrationEmail extends AbstractIntegration {
    type: IntegrationType.EMAIL;
    email: IntegrationEmailEntity;
}

export interface IntegrationSyslogEntity {
    address: string;
}

export interface IntegrationSyslog extends AbstractIntegration {
    type: IntegrationType.SYSLOG;
    syslog: IntegrationSyslogEntity;
}

export interface IntegrationWebhookEntity {
    url: string;
    login: string;
    password: string;
    ca: string;
    insecure: boolean;
}

export interface IntegrationWebhook extends AbstractIntegration {
    type: IntegrationType.WEBHOOK;
    webhook: IntegrationWebhookEntity;
}

export type Integration = IntegrationEmail | IntegrationSyslog | IntegrationWebhook;
