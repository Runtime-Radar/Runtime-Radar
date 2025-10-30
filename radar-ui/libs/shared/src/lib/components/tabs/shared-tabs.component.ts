import { KbqTabChangeEvent } from '@koobiq/components/tabs';
import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { Observable, map } from 'rxjs';

import { CoreNavigationStoreService } from '@cs/core';

interface AbstractTabOption {
    path: string;
    localizationKey: string;
}

@Component({
    selector: 'cs-tabs-component',
    templateUrl: './shared-tabs.component.html',
    styleUrl: './shared-tabs.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class SharedTabsComponent<T extends AbstractTabOption> {
    @Input({ required: true }) tabs!: T[];

    @Input({ required: true }) localizationFn!: any;

    @Input() testId?: string;

    @Output() tabChange = new EventEmitter<string>();

    readonly activeTabIndex$: Observable<number> = this.coreNavigationStoreService.routeSlug$.pipe(
        map((slug) => this.tabs.findIndex((tab) => tab.path === slug))
    );

    constructor(private readonly coreNavigationStoreService: CoreNavigationStoreService) {}

    onSelectedTabChange(tab: KbqTabChangeEvent) {
        const item = this.tabs[tab.index];
        this.tabChange.emit(item ? item.path : undefined);
    }
}
