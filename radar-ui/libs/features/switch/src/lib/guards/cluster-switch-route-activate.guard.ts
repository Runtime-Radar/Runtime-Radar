import { Router } from '@angular/router';
import { inject } from '@angular/core';
import { Observable, of, tap } from 'rxjs';

import { IS_CHILD_CLUSTER, RouterName } from '@cs/core';

const clusterFeatureSwitchRouteActivate = (): Observable<boolean> => {
    const router = inject(Router);
    const isChildCluster = inject(IS_CHILD_CLUSTER);

    return of(isChildCluster).pipe(
        tap((isChild) => {
            if (!isChild) {
                router.navigate([RouterName.DEFAULT]);
            }
        })
    );
};

export const clusterFeatureSwitchRouteActivateGuard = () => clusterFeatureSwitchRouteActivate();

export const clusterFeatureSwitchRouteActivateChildGuard = () => clusterFeatureSwitchRouteActivate();
