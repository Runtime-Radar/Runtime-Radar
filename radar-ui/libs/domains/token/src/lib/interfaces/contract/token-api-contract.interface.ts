import { Token, TokenPermissions } from './token-contract.interface';

export interface GetTokensResponse {
    access_tokens: Token[];
    total: number;
}

export interface CreateTokenRequest {
    name: string;
    user_id: string;
    permissions: TokenPermissions;
    expires_at: string | null; // RFC3339
}

export interface CreateTokenResponse {
    id: string;
    access_token: string;
}

export type EmptyTokenResponse = Record<string, unknown>;
