import { EffectsModule } from '@ngrx/effects';
import { NgModule } from '@angular/core';
import { StoreModule } from '@ngrx/store';

import { ApiModule } from '@cs/api';

import { NotificationEffectStore } from './stores/notification-effect.store';
import { NOTIFICATION_DOMAIN_KEY, notificationDomainReducer } from './stores/notification-selector.store';

@NgModule({
    imports: [
        ApiModule,
        StoreModule.forFeature(NOTIFICATION_DOMAIN_KEY, notificationDomainReducer),
        EffectsModule.forFeature([NotificationEffectStore])
    ]
})
export class NotificationDomainModule {}
