import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { SignInFeaturePageContainer } from './containers/page/sign-in-page.container';

const routes: Routes = [
    {
        path: '',
        component: SignInFeaturePageContainer,
        data: {
            localizationTitleKey: 'Auth.SignInPage.Header.Title'
        }
    }
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule]
})
export class SignInFeatureRoutingModule {}
