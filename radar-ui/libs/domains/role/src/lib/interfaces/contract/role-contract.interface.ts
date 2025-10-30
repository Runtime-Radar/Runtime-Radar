export enum PermissionName {
    CLUSTERS = 'clusters',
    EVENTS = 'events',
    INTEGRATIONS = 'integrations',
    NOTIFICATIONS = 'notifications',
    REGISTRIES = 'registries',
    ROLES = 'roles',
    RULES = 'rules',
    SCANNING = 'scanning',
    SYSTEM = 'system_settings',
    TOKENS = 'public_access_tokens',
    INVALIDATE_TOKENS = 'invalidate_public_access_tokens',
    USERS = 'users'
}

export enum PermissionType {
    CREATE = 'create',
    READ = 'read',
    UPDATE = 'update',
    DELETE = 'delete',
    EXECUTE = 'execute'
}

export interface RolePermission {
    actions: PermissionType[];
    description: string;
}

export type RolePermissions = {
    [key in PermissionName]: RolePermission;
};

export interface Role {
    id: string;
    description: string;
    global_flag: boolean;
    role_name: string;
    role_permissions: RolePermissions;
}
