import { PermissionType } from '@cs/domains/role';

export type TokenSidepanelPermissionMap = {
    [key in string]: Map<PermissionType, boolean>;
};

export interface TokenSidepanelFormProps {
    permissions: TokenSidepanelPermissionMap;
}
