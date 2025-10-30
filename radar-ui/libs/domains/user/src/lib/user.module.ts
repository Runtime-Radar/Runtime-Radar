import { EffectsModule } from '@ngrx/effects';
import { NgModule } from '@angular/core';
import { StoreModule } from '@ngrx/store';

import { ApiModule } from '@cs/api';

import { UserEffectStore } from './stores/user-effect.store';
import { USER_DOMAIN_KEY, userDomainReducer } from './stores/user-selector.store';

@NgModule({
    imports: [
        ApiModule,
        StoreModule.forFeature(USER_DOMAIN_KEY, userDomainReducer),
        EffectsModule.forFeature([UserEffectStore])
    ]
})
export class UserDomainModule {}
