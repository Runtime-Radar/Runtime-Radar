import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { ClusterDomainModule } from '@cs/domains/cluster';
import { I18nModule } from '@cs/i18n';
import { SharedModule } from '@cs/shared';

import { SwitchFeaturePageContainer } from './containers/page/switch-page.container';
import { SwitchFeatureRoutingModule } from './switch-routing.module';

@NgModule({
    imports: [CommonModule, ClusterDomainModule, SwitchFeatureRoutingModule, I18nModule, SharedModule],
    declarations: [SwitchFeaturePageContainer]
})
export class SwitchFeatureModule {}
