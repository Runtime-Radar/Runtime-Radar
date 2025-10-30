import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { BrowserModule } from '@angular/platform-browser';
import { StoreDevtoolsModule } from '@ngrx/store-devtools';
import { APP_INITIALIZER, NgModule } from '@angular/core';
import { NavigationActionTiming, RouterState, StoreRouterConnectingModule } from '@ngrx/router-store';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';

import { I18nModule } from '@cs/i18n';
import { SharedModule } from '@cs/shared';
import { API_PATH, API_SINGLE_TENANT_PATHS } from '@cs/api';
import { CoreInitService, CoreModule, IS_CHILD_CLUSTER, POLLING_INTERVAL, REFRESH_INTERVAL } from '@cs/core';

import { AppContainer } from './app.container';
import { AppRoutingModule } from './app-routing.module';
import { NavbarComponent } from './components/navbar/navbar.component';
import { environment } from '../environments/environment';

function initializeFactory(initService: CoreInitService): () => Promise<void> {
    return (): Promise<void> => initService.initialize();
}

@NgModule({
    imports: [
        AppRoutingModule,
        BrowserModule,
        BrowserAnimationsModule,
        CoreModule,
        SharedModule,
        I18nModule.forRoot({
            prodMode: environment.production
        }),
        StoreRouterConnectingModule.forRoot({
            navigationActionTiming: NavigationActionTiming.PostActivation,
            routerState: RouterState.Full
        }),
        StoreDevtoolsModule.instrument({
            logOnly: environment.production
        })
    ],
    providers: [
        {
            provide: APP_INITIALIZER,
            useFactory: initializeFactory,
            deps: [CoreInitService],
            multi: true
        },
        {
            provide: API_PATH,
            useValue: environment.api
        },
        {
            provide: API_SINGLE_TENANT_PATHS,
            useValue: environment.singleTenant
        },
        {
            provide: POLLING_INTERVAL,
            useValue: environment.pollingInterval
        },
        {
            provide: REFRESH_INTERVAL,
            useValue: environment.refreshInterval
        },
        {
            provide: IS_CHILD_CLUSTER,
            useValue: environment.childCluster
        },
        provideHttpClient(withInterceptorsFromDi())
    ],
    declarations: [AppContainer, NavbarComponent],
    bootstrap: [AppContainer]
})
export class AppModule {}
