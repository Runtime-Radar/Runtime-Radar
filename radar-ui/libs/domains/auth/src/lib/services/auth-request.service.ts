import { HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, map } from 'rxjs';

import { ApiEmptyRequest, ApiService } from '@cs/api';

import {
    GetAppVersionResponse,
    GetCentralUrlResponse,
    GetLoginRequest,
    GetLoginResponse,
    GetTokenResponse
} from '../interfaces/contract/auth-api-contract.interface';

@Injectable({
    providedIn: 'root'
})
export class AuthRequestService {
    constructor(private readonly apiService: ApiService) {}

    getLogin(request: GetLoginRequest): Observable<GetLoginResponse> {
        return this.apiService.post<GetLoginRequest, GetLoginResponse>('signin', request);
    }

    getTokens(headers: HttpHeaders): Observable<GetTokenResponse> {
        return this.apiService.get<ApiEmptyRequest, GetTokenResponse>('tokens', {}, headers);
    }

    getAppVersion(): Observable<string> {
        return this.apiService
            .get<ApiEmptyRequest, GetAppVersionResponse>('info/version')
            .pipe(map((response) => response.version));
    }

    /** @external */
    getCentralUrl(): Observable<string> {
        return this.apiService
            .get<ApiEmptyRequest, GetCentralUrlResponse>('info/central-cs-url')
            .pipe(map((response) => response.url));
    }
}
