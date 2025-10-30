import { Role } from '@cs/domains/role';

export interface AuthJwtData {
    user_id: string;
    username: string;
    email: string;
    token_type: string;
    auth_type: string;
    role: Role; // @todo: replace to roleId
    exp: number;
    iat: number;
    last_password_changed_at: number;
}
