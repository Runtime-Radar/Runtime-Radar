import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { I18nModule } from '@cs/i18n';
import { SharedModule } from '@cs/shared';

import { ForbiddenFeaturePermissionContainer } from './containers/permission/forbidden-permission.container';
import { ForbiddenFeatureRoutingModule } from './forbidden-routing.module';

@NgModule({
    imports: [CommonModule, ForbiddenFeatureRoutingModule, I18nModule, SharedModule],
    declarations: [ForbiddenFeaturePermissionContainer]
})
export class ForbiddenFeatureModule {}
