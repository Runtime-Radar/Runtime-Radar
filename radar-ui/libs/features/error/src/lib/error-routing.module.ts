import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ErrorFeatureNotFoundContainer } from './containers/not-found/error-not-found.container';

const routes: Routes = [
    {
        path: '',
        component: ErrorFeatureNotFoundContainer,
        data: {
            localizationTitleKey: 'Common.ErrorPage.Header.Title'
        }
    }
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule]
})
export class ErrorFeatureRoutingModule {}
