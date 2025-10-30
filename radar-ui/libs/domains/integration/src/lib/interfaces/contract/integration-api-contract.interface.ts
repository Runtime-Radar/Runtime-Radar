import { Integration } from './integration-contract.interface';

export interface GetIntegrationsResponse<T extends Integration> {
    integrations: T[];
}

export type CreateIntegrationRequest<T extends Integration> = Omit<T, 'id'>;

export interface CreateIntegrationResponse {
    id: string;
}

export type UpdateIntegrationRequest<T extends Integration> = CreateIntegrationRequest<T>;

export type EmptyIntegrationResponse = Record<string, unknown>;
