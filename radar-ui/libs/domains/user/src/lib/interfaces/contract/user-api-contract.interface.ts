import { User } from './user-contract.interface';

export interface GetUsersResponse {
    users: User[];
}

export interface CreateUserRequest {
    username: string;
    email: string;
    password: string;
    role_id: string;
}

export interface UpdateUserRequest {
    email: string;
    role_id: string;
}

export interface UpdatePasswordRequest {
    password: string;
}

export interface DeleteUserResponse {
    id: string;
}
