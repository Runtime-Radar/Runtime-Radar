import { KbqAlertColors } from '@koobiq/components/alert';
import { KbqBadgeColors } from '@koobiq/components/badge';
import { ChangeDetectionStrategy, Component, Inject } from '@angular/core';
import { KBQ_SIDEPANEL_DATA, KbqSidepanelRef } from '@koobiq/components/sidepanel';
import { Observable, map } from 'rxjs';

import { PermissionType } from '@cs/domains/role';
import { DetectorExtended, DetectorType } from '@cs/domains/detector';
import { Rule, RuleWhiteList } from '@cs/domains/rule';

import { SharedRuleSidepanelProps } from './shared-rule-sidepanel.interface';

@Component({
    templateUrl: './shared-rule-sidepanel.component.html',
    styleUrl: './shared-rule-sidepanel.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class SharedRuleSidepanelComponent {
    readonly alertColors = KbqAlertColors;

    readonly badgeColors = KbqBadgeColors;

    readonly permissionType = PermissionType;

    readonly detectors$ = (whiteList: RuleWhiteList): Observable<DetectorExtended[]> =>
        this.props.detectors$.pipe(
            map((detectors) => {
                const { binaries, ...rest } = whiteList;
                const keys = Object.values(rest).flat();

                return keys.map((key) => {
                    const detector = detectors.find((item) => item.key === key);
                    const defaultVersion = 1;
                    const emptyDetector: DetectorExtended = {
                        id: `${key}${defaultVersion}`,
                        key,
                        name: '',
                        description: '',
                        type: DetectorType.RUNTIME,
                        version: defaultVersion
                    };

                    return detector || emptyDetector;
                });
            })
        );

    constructor(
        private readonly sidepanelRef: KbqSidepanelRef,
        @Inject(KBQ_SIDEPANEL_DATA) public readonly props: SharedRuleSidepanelProps
    ) {}

    updateRule(rule: Rule) {
        if (this.props.updateHandler) {
            this.props.updateHandler(rule);
        }
    }

    deleteRule(id: string) {
        if (this.props.deleteHandler) {
            this.props.deleteHandler(id);
            this.sidepanelRef.close(undefined);
        }
    }
}
