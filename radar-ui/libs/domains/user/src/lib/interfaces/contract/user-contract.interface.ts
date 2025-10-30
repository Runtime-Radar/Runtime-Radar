export interface User {
    id: string;
    username: string;
    email: string;
    role_id: string;
}

export interface UserEditRequest {
    email: string;
    password: string;
    roleId: string;
}
