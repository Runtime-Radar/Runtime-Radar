import { BehaviorSubject, Subject } from 'rxjs';
import { Inject, Injectable } from '@angular/core';

import { API_PATH, API_SINGLE_TENANT_PATHS } from '@cs/api';
import { CoreWindowService, FORM_VALIDATION_REG_EXP } from '@cs/core';

const API_CLUSTER_PATH_SESSION_KEY = 'pclstrpth';

const API_CLUSTER_URL_QUERY_PARAM_KEY = 'clusterUrl';

@Injectable({
    providedIn: 'root'
})
export class ApiPathService {
    private host = '';

    readonly host$ = new BehaviorSubject(this.host);

    readonly error$ = new Subject<string>();

    constructor(
        private readonly coreWindowService: CoreWindowService,
        @Inject(API_PATH) private readonly apiPath: string,
        @Inject(API_SINGLE_TENANT_PATHS) private readonly apiSingleTenantPaths: string[]
    ) {}

    initialize() {
        // there are links which should switch cluster based on query params into email templates
        const href = this.coreWindowService.location.href.substring(this.coreWindowService.location.href.indexOf('?'));
        let value = this.coreWindowService.sessionStorage.getItem(API_CLUSTER_PATH_SESSION_KEY) || '';
        if (href.includes(API_CLUSTER_URL_QUERY_PARAM_KEY)) {
            value = decodeURIComponent(href.substring(API_CLUSTER_URL_QUERY_PARAM_KEY.length + 2));
        }

        if (FORM_VALIDATION_REG_EXP.IP_DOMAIN_SCHEME.test(value)) {
            this.setHost(value);
        }
    }

    get(path: string): string {
        const segment = path.substring(0, path.indexOf('/')) || path;
        const host = this.apiSingleTenantPaths.includes(segment) ? '' : this.host;

        return `${host}${this.apiPath}${path}`;
    }

    setHost(value: string) {
        this.host = value;
        this.host$.next(value);
        this.coreWindowService.sessionStorage.setItem(API_CLUSTER_PATH_SESSION_KEY, value);
    }

    setError(value: string) {
        this.error$.next(value);
    }
}
