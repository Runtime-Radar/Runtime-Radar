import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { DetectorDomainModule } from '@cs/domains/detector';
import { I18nModule } from '@cs/i18n';
import { NotificationDomainModule } from '@cs/domains/notification';
import { RuleDomainModule } from '@cs/domains/rule';
import { SharedModule } from '@cs/shared';

import { RuleFeatureFilterPipe } from './pipes/rule-filter.pipe';
import { RuleFeatureListContainer } from './containers/list/rule-list.container';
import { RuleFeaturePanelFilterComponent } from './components/panel-filter/rule-panel-filter.component';
import { RuleFeatureRoutingModule } from './rule-routing.module';
import { RuleFeatureVerdictRadioComponent } from './components/verdict-radio/rule-verdict-radio.component';

@NgModule({
    imports: [
        CommonModule,
        FormsModule,
        DetectorDomainModule,
        I18nModule,
        NotificationDomainModule,
        ReactiveFormsModule,
        RuleDomainModule,
        RuleFeatureRoutingModule,
        SharedModule
    ],
    declarations: [
        RuleFeaturePanelFilterComponent,
        RuleFeatureFilterPipe,
        RuleFeatureListContainer,
        RuleFeatureVerdictRadioComponent
    ]
})
export class RuleFeatureModule {}
