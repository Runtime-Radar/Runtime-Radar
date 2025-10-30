import {
    IntegrationEmailAuthType,
    IntegrationEmailAuthTypeOption,
    IntegrationType,
    IntegrationTypeOption
} from '../interfaces';

export const INTEGRATION_TYPE: IntegrationTypeOption[] = [
    {
        id: IntegrationType.EMAIL,
        localizationKey: 'Integration.Pseudo.Type.Email'
    },
    {
        id: IntegrationType.SYSLOG,
        localizationKey: 'Integration.Pseudo.Type.Syslog'
    },
    {
        id: IntegrationType.WEBHOOK,
        localizationKey: 'Integration.Pseudo.Type.Webhook'
    }
];

export const INTEGRATION_EMAIL_AUTH_TYPE: IntegrationEmailAuthTypeOption[] = [
    {
        id: IntegrationEmailAuthType.NONE,
        localizationKey: 'Integration.Pseudo.AuthType.None'
    },
    {
        id: IntegrationEmailAuthType.CRAM_MD5,
        localizationKey: 'Integration.Pseudo.AuthType.Crammd5'
    },
    {
        id: IntegrationEmailAuthType.LOGIN,
        localizationKey: 'Integration.Pseudo.AuthType.Login'
    },
    {
        id: IntegrationEmailAuthType.PLAIN,
        localizationKey: 'Integration.Pseudo.AuthType.Plain'
    }
];
