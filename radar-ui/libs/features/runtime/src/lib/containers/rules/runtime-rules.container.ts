import { KbqAlertColors } from '@koobiq/components/alert';
import { KbqBadgeColors } from '@koobiq/components/badge';
import { PopUpPlacements } from '@koobiq/components/core';
import { BehaviorSubject, Observable, filter, forkJoin, map, of, switchMap, take } from 'rxjs';
import { ChangeDetectionStrategy, Component, Input, OnChanges, Optional } from '@angular/core';
import {
    KbqSidepanelConfig,
    KbqSidepanelPosition,
    KbqSidepanelRef,
    KbqSidepanelService
} from '@koobiq/components/sidepanel';

import { I18nService } from '@cs/i18n';
import { PermissionType } from '@cs/domains/role';
import {
    CreateRuleRequest,
    GetRuleResponse,
    Rule,
    RuleRequestService,
    RuleSeverity,
    RuleStoreService,
    RuleType,
    UpdateRuleRequest
} from '@cs/domains/rule';
import { LoadStatus, CoreUtilsService as utils } from '@cs/core';
import { Notification, NotificationStoreService } from '@cs/domains/notification';
import { SharedModalService, SharedRuleSidepanelFormComponent, SharedRuleSidepanelFormProps } from '@cs/shared';

import { RUNTIME_DETAILS_LIST_ITEMS_LIMIT } from '../../constants/runtime-config.constant';
import { RuntimeRuleForm } from '../../interfaces/runtime-form.interface';

@Component({
    selector: 'cs-runtime-feature-rules-container',
    templateUrl: './runtime-rules.container.html',
    styleUrl: './runtime-rules.container.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class RuntimeFeatureRulesContainer implements OnChanges {
    @Input() ruleIds?: string[];

    @Input() permissions = new Map<PermissionType, boolean>();

    readonly loadStatus$: Observable<LoadStatus> = this.ruleStoreService.loadStatus$;

    readonly notifications$ = (ids?: string[]): Observable<Notification[]> =>
        this.notificationStoreService.notifications$.pipe(
            map((list) => (ids ? list.filter((item) => ids.includes(item.id)) : []))
        );

    // @todo: refactor getRule to get an ability to provide ids
    private readonly ruleIds$ = new BehaviorSubject<string[] | undefined>(undefined);
    readonly rulesResponse$: Observable<GetRuleResponse[]> = this.ruleStoreService.rules$.pipe(
        map((rules) => rules.filter((item) => item.type === RuleType.TYPE_RUNTIME)),
        switchMap((runtimeRules) =>
            this.ruleIds$.pipe(
                map((ruleIds) => {
                    if (!ruleIds) {
                        return {
                            rules: runtimeRules.map<GetRuleResponse>((item) => ({ rule: item, deleted: false })),
                            deletedRuleIds: []
                        };
                    }

                    const rules = ruleIds.length ? runtimeRules.filter((item) => ruleIds?.includes(item.id)) : [];
                    const actualRuleIds = rules.map((item) => item.id);

                    return {
                        rules: rules.map<GetRuleResponse>((item) => ({ rule: item, deleted: false })),
                        deletedRuleIds: ruleIds.filter((id) => !actualRuleIds.includes(id))
                    };
                })
            )
        ),
        switchMap(({ rules, deletedRuleIds }) => {
            if (!deletedRuleIds.length) {
                return of(rules);
            }

            return forkJoin(deletedRuleIds.map((id) => this.ruleRequestService.getRule(id))).pipe(
                map((response: GetRuleResponse[]) => [...response, ...rules])
            );
        })
    );

    readonly ruleNotificationsLimit = RUNTIME_DETAILS_LIST_ITEMS_LIMIT;

    readonly tooltipPlacements = PopUpPlacements;

    readonly permissionType = PermissionType;

    readonly loadStatus = LoadStatus;

    readonly alertColors = KbqAlertColors;

    readonly badgeColors = KbqBadgeColors;

    expandLimit: number | undefined = RUNTIME_DETAILS_LIST_ITEMS_LIMIT;

    constructor(
        private readonly sidepanelService: KbqSidepanelService,
        private readonly notificationStoreService: NotificationStoreService,
        private readonly ruleStoreService: RuleStoreService,
        private readonly i18nService: I18nService,
        private readonly ruleRequestService: RuleRequestService,
        private readonly sharedModalService: SharedModalService,
        @Optional() readonly sidepanelRef: KbqSidepanelRef
    ) {}

    ngOnChanges() {
        this.ruleIds$.next(this.ruleIds);
    }

    openEditRuleSidepanel(rule: Rule) {
        const config: KbqSidepanelConfig<Partial<SharedRuleSidepanelFormProps>> = {
            position: KbqSidepanelPosition.Right,
            hasBackdrop: true,
            data: {
                isEdit: true,
                rule
            }
        };

        this.sidepanelService
            .open(SharedRuleSidepanelFormComponent, config)
            .afterClosed()
            .pipe(take(1), filter(utils.isDefined))
            .subscribe((form: RuntimeRuleForm) => {
                const request: UpdateRuleRequest = {
                    name: form.name,
                    type: RuleType.TYPE_RUNTIME,
                    rule: {
                        version: '1', // @todo: create environment constant
                        notify:
                            form.notifySeverity === RuleSeverity.NONE
                                ? null
                                : {
                                      severity: form.notifySeverity,
                                      verdict: null,
                                      targets: form.mailIds
                                  },
                        whitelist: {
                            threats: form.detectors,
                            binaries: form.binaries
                        }
                    },
                    scope: {
                        version: '1', // @todo: create environment constant
                        image_names: form.imageNames,
                        registries: form.registries,
                        clusters: [],
                        namespaces: form.namespaces,
                        pods: form.pods,
                        containers: form.containers,
                        nodes: form.nodes
                    }
                };

                this.ruleStoreService.updateRule(rule.id, request);
            });
    }

    openDeleteRuleModal(id: string) {
        if (this.sidepanelRef) {
            this.sidepanelRef.close(undefined);
        }

        this.sharedModalService.delete({
            content: this.i18nService.translate('Rule.DeleteModal.Title.Text'),
            confirmText: this.i18nService.translate('Rule.DeleteModal.Button.Confirm'),
            cancelText: this.i18nService.translate('Rule.DeleteModal.Button.Cancel'),
            confirmHandler: () => {
                this.ruleStoreService.deleteRule(id);
            }
        });
    }

    openCreateRuleSidepanel() {
        this.sidepanelService
            .open(SharedRuleSidepanelFormComponent, {
                position: KbqSidepanelPosition.Right,
                hasBackdrop: true,
                data: {}
            })
            .afterClosed()
            .pipe(take(1), filter(utils.isDefined))
            .subscribe((form: RuntimeRuleForm) => {
                const request: CreateRuleRequest = {
                    name: form.name,
                    type: RuleType.TYPE_RUNTIME,
                    rule: {
                        version: '1', // @todo: create environment constant
                        notify:
                            form.notifySeverity === RuleSeverity.NONE
                                ? null
                                : {
                                      severity: form.notifySeverity,
                                      verdict: null,
                                      targets: form.mailIds
                                  },
                        whitelist: {
                            threats: form.detectors,
                            binaries: form.binaries
                        }
                    },
                    scope: {
                        version: '1', // @todo: create environment constant
                        image_names: form.imageNames,
                        registries: form.registries,
                        clusters: [],
                        namespaces: form.namespaces,
                        pods: form.pods,
                        containers: form.containers,
                        nodes: form.nodes
                    }
                };

                this.ruleStoreService.createRule(request);
            });
    }

    deleteRuleThreat(rule: Rule, threatKey: string) {
        const threats: string[] = rule.rule.whitelist.threats.filter((item) => item !== threatKey);
        const request: UpdateRuleRequest = {
            name: rule.name,
            type: rule.type,
            scope: rule.scope,
            rule: {
                version: rule.rule.version,
                notify: rule.rule.notify,
                whitelist: {
                    threats,
                    binaries: rule.rule.whitelist.binaries
                }
            }
        };

        this.ruleStoreService.updateRule(rule.id, request);
    }

    deleteRuleBinary(rule: Rule, binaryKey: string) {
        const binaries: string[] = rule.rule.whitelist.binaries.filter((item) => item !== binaryKey);
        const request: UpdateRuleRequest = {
            name: rule.name,
            type: rule.type,
            scope: rule.scope,
            rule: {
                version: rule.rule.version,
                notify: rule.rule.notify,
                whitelist: {
                    threats: rule.rule.whitelist.threats,
                    binaries
                }
            }
        };

        this.ruleStoreService.updateRule(rule.id, request);
    }

    expand() {
        this.expandLimit = undefined;
    }
}
