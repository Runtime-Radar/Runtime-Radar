import { Action } from '@ngrx/store';
import { HttpErrorResponse } from '@angular/common/http';
import { I18nService } from '@cs/i18n';
import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { KbqToastService, KbqToastStyle } from '@koobiq/components/toast';
import { Observable, of } from 'rxjs';
import { catchError, switchMap, take } from 'rxjs/operators';

import { ALLOW_AUTH_EVENT_ACTION } from '@cs/domains/auth';
import { LoadStatus } from '@cs/core';
import { ApiErrorCode, ApiUtilsService as apiUtils } from '@cs/api';

import { Role } from '../interfaces';
import { RoleRequestService } from '../services/role-request.service';
import {
    ROLE_LOAD_DONE_EVENT_ACTION,
    SET_ALL_ROLE_ENTITIES_DOC_ACTION,
    UPDATE_ROLE_LOAD_STATUS_DOC_ACTION
} from './role-action.store';

@Injectable({
    providedIn: 'root'
})
export class RoleEffectStore {
    readonly loadRoles$: Observable<Action> = createEffect(() =>
        this.actions$.pipe(
            ofType(ALLOW_AUTH_EVENT_ACTION),
            switchMap(() =>
                this.roleRequestService.getRoles().pipe(
                    take(1),
                    catchError((error: HttpErrorResponse) => {
                        if (apiUtils.getReasonCode(error) === ApiErrorCode.PERMISSION_DENIED) {
                            this.toastService.show({
                                style: KbqToastStyle.Warning,
                                title: this.i18nService.translate('Common.Pseudo.Notification.Forbidden')
                            });
                        }

                        return of([] as Role[]);
                    })
                )
            ),
            switchMap((list) => {
                const loadStatus = list.length ? LoadStatus.LOADED : LoadStatus.ERROR;

                return [
                    SET_ALL_ROLE_ENTITIES_DOC_ACTION({ list }),
                    UPDATE_ROLE_LOAD_STATUS_DOC_ACTION({ loadStatus }),
                    ROLE_LOAD_DONE_EVENT_ACTION({ loadStatus })
                ];
            })
        )
    );

    constructor(
        private readonly actions$: Actions,
        private readonly i18nService: I18nService,
        private readonly toastService: KbqToastService,
        private readonly roleRequestService: RoleRequestService
    ) {}
}
