import { EffectsModule } from '@ngrx/effects';
import { NgModule } from '@angular/core';
import { StoreModule } from '@ngrx/store';

import { ApiModule } from '@cs/api';

import { RoleEffectStore } from './stores/role-effect.store';
import { ROLE_DOMAIN_KEY, roleDomainReducer } from './stores/role-selector.store';

@NgModule({
    imports: [
        ApiModule,
        StoreModule.forFeature(ROLE_DOMAIN_KEY, roleDomainReducer),
        EffectsModule.forFeature([RoleEffectStore])
    ]
})
export class RoleDomainModule {}
