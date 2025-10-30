import { AbstractControl, AsyncValidatorFn, ValidationErrors } from '@angular/forms';
import { Observable, map, of, switchMap, take } from 'rxjs';

import { Cluster, ClusterStoreService } from '@cs/domains/cluster';

export class ClusterFeatureUrlValidator {
    static isUrlUnique(clusterStoreService: ClusterStoreService): AsyncValidatorFn {
        return (control: AbstractControl): Observable<ValidationErrors | null> => {
            if (!control.value) {
                return of(null);
            }

            return of(control.value).pipe(
                switchMap((value) =>
                    clusterStoreService.clusters$.pipe(
                        take(1),
                        map((clusters: Cluster[]) =>
                            clusters.some((item) => item.config.own_cs_url === value) ? { urlExists: true } : null
                        )
                    )
                )
            );
        };
    }
}
