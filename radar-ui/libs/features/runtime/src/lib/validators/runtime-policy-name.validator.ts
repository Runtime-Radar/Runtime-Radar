import { AbstractControl, AsyncValidatorFn, ValidationErrors } from '@angular/forms';
import { Observable, map, of, switchMap, take } from 'rxjs';

import { RuntimeStoreService } from '@cs/domains/runtime';

import { RuntimeFeaturePolicyNameService } from '../services/runtime-policy-name.service';

export class RuntimeFeaturePolicyNameValidator {
    static isNameUnique(
        runtimeStoreService: RuntimeStoreService,
        runtimeFeaturePolicyNameService: RuntimeFeaturePolicyNameService,
        originName: string | undefined
    ): AsyncValidatorFn {
        return (control: AbstractControl): Observable<ValidationErrors | null> => {
            if (!control.value) {
                return of(null);
            }

            return of(control.value).pipe(
                switchMap((value) =>
                    runtimeStoreService.runtimeMonitorConfig$.pipe(
                        take(1),
                        map((config) => {
                            const names = Object.values(config.tracing_policies)
                                .map((item) => item.name)
                                .concat(runtimeFeaturePolicyNameService.get())
                                .filter((name) => name !== originName);

                            return names.includes(value) ? { nameExists: true } : null;
                        })
                    )
                )
            );
        };
    }
}
