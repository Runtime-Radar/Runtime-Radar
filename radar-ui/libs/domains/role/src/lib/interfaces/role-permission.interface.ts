import { PermissionName, PermissionType } from './contract/role-contract.interface';

export type RolePermissionMap = {
    [key in PermissionName]: Map<PermissionType, boolean>;
};
