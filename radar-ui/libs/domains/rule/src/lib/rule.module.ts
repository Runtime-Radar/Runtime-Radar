import { EffectsModule } from '@ngrx/effects';
import { NgModule } from '@angular/core';
import { StoreModule } from '@ngrx/store';

import { ApiModule } from '@cs/api';

import { RuleEffectStore } from './stores/rule-effect.store';
import { RuleTypeLocalizationPipe } from './pipes/rule-type.pipe';
import { RULE_DOMAIN_KEY, ruleDomainReducer } from './stores/rule-selector.store';

@NgModule({
    imports: [
        ApiModule,
        StoreModule.forFeature(RULE_DOMAIN_KEY, ruleDomainReducer),
        EffectsModule.forFeature([RuleEffectStore])
    ],
    declarations: [RuleTypeLocalizationPipe],
    exports: [RuleTypeLocalizationPipe]
})
export class RuleDomainModule {}
