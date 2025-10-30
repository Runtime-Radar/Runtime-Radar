import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { I18nModule } from '@cs/i18n';
import { SharedModule } from '@cs/shared';

import { ErrorFeatureNotFoundContainer } from './containers/not-found/error-not-found.container';
import { ErrorFeatureRoutingModule } from './error-routing.module';

@NgModule({
    imports: [CommonModule, ErrorFeatureRoutingModule, I18nModule, SharedModule],
    declarations: [ErrorFeatureNotFoundContainer]
})
export class ErrorFeatureModule {}
