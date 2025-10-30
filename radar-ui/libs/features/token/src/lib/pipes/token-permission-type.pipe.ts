import { Pipe, PipeTransform } from '@angular/core';

import { I18nService } from '@cs/i18n';
import { PermissionType } from '@cs/domains/role';
import { TokenPermissionName } from '@cs/domains/token';

@Pipe({
    name: 'tokenPermissionType',
    pure: false
})
export class TokenFeaturePermissionTypePipe implements PipeTransform {
    constructor(private readonly i18nService: I18nService) {}

    transform(type?: PermissionType, permissionName?: TokenPermissionName): string {
        switch (type) {
            case PermissionType.CREATE:
                return this.i18nService.translate('Token.CreateForm.RulePermissions.Label.CanCreate');
            case PermissionType.READ:
                if (permissionName === TokenPermissionName.EVENTS) {
                    return this.i18nService.translate('Token.CreateForm.EventPermissions.Label.CanRead');
                }

                return this.i18nService.translate('Token.CreateForm.RulePermissions.Label.CanRead');
            case PermissionType.UPDATE:
                return this.i18nService.translate('Token.CreateForm.RulePermissions.Label.CanUpdate');
            case PermissionType.DELETE:
                return this.i18nService.translate('Token.CreateForm.RulePermissions.Label.CanDelete');
            default:
                return '—';
        }
    }
}
