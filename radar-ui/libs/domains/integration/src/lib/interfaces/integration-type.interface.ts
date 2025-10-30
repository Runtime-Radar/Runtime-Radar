import { IntegrationEmailAuthType, IntegrationType } from './contract/integration-contract.interface';

export interface IntegrationTypeOption {
    id: IntegrationType;
    localizationKey: string;
}

export interface IntegrationEmailAuthTypeOption {
    id: IntegrationEmailAuthType;
    localizationKey: string;
}
