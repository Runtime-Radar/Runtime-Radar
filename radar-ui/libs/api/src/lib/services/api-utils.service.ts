import { HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';

import { ApiErrorCode, ApiErrorResponse } from '../interfaces/contract/api-error-contract.interface';

@Injectable({
    providedIn: 'root'
})
export class ApiUtilsService {
    static getReasonCode(error: HttpErrorResponse): ApiErrorCode {
        const detail = (error.error as ApiErrorResponse)?.details?.at(0);

        return detail ? detail.reason : ApiErrorCode.UNKNOWN_ISSUE;
    }
}
