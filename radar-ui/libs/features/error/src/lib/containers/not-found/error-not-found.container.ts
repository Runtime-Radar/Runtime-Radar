import { ChangeDetectionStrategy, Component } from '@angular/core';

import { RouterName } from '@cs/core';

@Component({
    templateUrl: './error-not-found.container.html',
    styleUrl: './error-not-found.container.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class ErrorFeatureNotFoundContainer {
    readonly routerName = RouterName;
}
