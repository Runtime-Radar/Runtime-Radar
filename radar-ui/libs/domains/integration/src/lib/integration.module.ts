import { EffectsModule } from '@ngrx/effects';
import { NgModule } from '@angular/core';
import { StoreModule } from '@ngrx/store';

import { ApiModule } from '@cs/api';

import { IntegrationEffectStore } from './stores/integration-effect.store';
import { IntegrationTypeLocalizationPipe } from './pipes/integration-type.pipe';
import { INTEGRATION_DOMAIN_KEY, integrationDomainReducer } from './stores/integration-selector.store';

@NgModule({
    imports: [
        ApiModule,
        StoreModule.forFeature(INTEGRATION_DOMAIN_KEY, integrationDomainReducer),
        EffectsModule.forFeature([IntegrationEffectStore])
    ],
    declarations: [IntegrationTypeLocalizationPipe],
    exports: [IntegrationTypeLocalizationPipe]
})
export class IntegrationDomainModule {}
