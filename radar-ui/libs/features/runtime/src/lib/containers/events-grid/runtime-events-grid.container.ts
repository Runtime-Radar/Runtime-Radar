import { PopUpPlacements } from '@koobiq/components/core';
import { Router } from '@angular/router';
import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output, booleanAttribute } from '@angular/core';
import { KbqSidepanelConfig, KbqSidepanelPosition, KbqSidepanelService } from '@koobiq/components/sidepanel';
import { filter, take } from 'rxjs';

import { PermissionType } from '@cs/domains/role';
import {
    CreateRuleRequest,
    RULE_SEVERITIES,
    RULE_SEVERITY_ORDER,
    RuleSeverity,
    RuleSeverityOption,
    RuleStoreService,
    RuleType
} from '@cs/domains/rule';
import {
    RUNTIME_CONTEXT,
    RuntimeCapabilityType,
    RuntimeContext,
    RuntimeDetectError,
    RuntimeEventEntity,
    RuntimeEventProcess,
    RuntimeEventThreat,
    RuntimeEventType
} from '@cs/domains/runtime';
import { RouterName, CoreUtilsService as utils } from '@cs/core';
import { SharedRuleSidepanelFormComponent, SharedRuleSidepanelFormProps } from '@cs/shared';

import { RuntimeFeatureSidepanelCodeComponent } from '../../components/sidepanel-code/runtime-sidepanel-code.component';
import { RuntimeFeatureSidepanelIncidentComponent } from '../../components/sidepanel-incident/runtime-sidepanel-incident.component';
import { RuntimeFeatureSidepanelThreatsComponent } from '../../components/sidepanel-threats/runtime-sidepanel-threats.component';
import { RuntimeRouterName } from '../../interfaces/runtime-navigation.interface';
import { RuntimeRuleForm } from '../../interfaces/runtime-form.interface';
import { RuntimeEventExtended, RuntimeEventsGridContextId } from '../../interfaces/runtime-events.interface';
import {
    RuntimeSidepanelCodeProps,
    RuntimeSidepanelIncidentProps,
    RuntimeSidepanelThreatsProps
} from '../../interfaces/runtime-sidepanel.interface';

@Component({
    selector: 'cs-runtime-feature-events-grid-container',
    templateUrl: './runtime-events-grid.container.html',
    styleUrl: './runtime-events-grid.container.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class RuntimeFeatureEventsGridContainer {
    localEvents: RuntimeEventExtended[] = [];
    @Input({ required: true }) set events(values: RuntimeEventExtended[] | null) {
        if (values) {
            this.localEvents = this.extendEventsBySeverity(values);
        }
    }

    @Input() activeContextId?: string | null;

    @Input() rulePermissions = new Map<PermissionType, boolean>();

    @Input({ transform: booleanAttribute }) isContextChangeAvailable = false;

    @Output() contextIdChange = new EventEmitter<RuntimeEventsGridContextId>();

    readonly ruleSeverityOptions: RuleSeverityOption[] = [...RULE_SEVERITIES].sort(
        (a, b) => RULE_SEVERITY_ORDER[a.id] - RULE_SEVERITY_ORDER[b.id]
    );

    readonly runtimeContextOptions = RUNTIME_CONTEXT;

    readonly runtimeContext = RuntimeContext;

    readonly runtimeEventType = RuntimeEventType;

    readonly ruleSeverity = RuleSeverity;

    readonly permissionType = PermissionType;

    readonly tooltipPlacements = PopUpPlacements;

    readonly runtimeCapabilityType = RuntimeCapabilityType;

    constructor(
        private readonly router: Router,
        private readonly ruleStoreService: RuleStoreService,
        private readonly sidepanelService: KbqSidepanelService
    ) {}

    goToDetailsPage(id: string) {
        this.router.navigate([RouterName.DEFAULT, RouterName.RUNTIME, RuntimeRouterName.EVENTS, id]);
    }

    setContextId(id: string, context: RuntimeContext, execId: string, parentExecId: string) {
        if (this.isContextChangeAvailable) {
            this.contextIdChange.emit({
                id,
                context,
                execId,
                parentExecId
            });
        }
    }

    openThreatsSidepanel(
        threats: RuntimeEventThreat[],
        effectives: RuntimeCapabilityType[],
        errors: RuntimeDetectError[]
    ) {
        const config: KbqSidepanelConfig<RuntimeSidepanelThreatsProps> = {
            position: KbqSidepanelPosition.Right,
            hasBackdrop: true,
            data: {
                threats: threats.sort((a, b) => RULE_SEVERITY_ORDER[a.severity] - RULE_SEVERITY_ORDER[b.severity]),
                effectives,
                errors
            }
        };

        this.sidepanelService
            .open(RuntimeFeatureSidepanelThreatsComponent, config)
            .afterClosed()
            .pipe(take(1))
            .subscribe();
    }

    openCreateRuleSidepanel(event: RuntimeEventExtended, process: RuntimeEventProcess) {
        const config: KbqSidepanelConfig<Partial<SharedRuleSidepanelFormProps>> = {
            position: KbqSidepanelPosition.Right,
            hasBackdrop: true,
            data: {
                isEdit: true,
                rule: {
                    rule: {
                        version: '1', // @todo: create environment constant
                        notify: {
                            targets: [],
                            severity: event.threatSeverity || RuleSeverity.NONE,
                            verdict: null
                        },
                        whitelist: {
                            threats: [],
                            binaries: []
                        }
                    },
                    scope: {
                        version: '1', // @todo: create environment constant
                        namespaces: process.process.pod?.namespace ? [process.process.pod.namespace] : [],
                        image_names: process.process.pod?.container.image.name
                            ? [process.process.pod?.container.image.name]
                            : [],
                        registries: [],
                        clusters: [],
                        pods: process.process.pod?.name ? [process.process.pod.name] : [],
                        containers: process.process.pod?.container ? [process.process.pod.container.name] : [],
                        nodes: event.event.node_name ? [event.event.node_name] : []
                    }
                }
            }
        };

        this.sidepanelService
            .open(SharedRuleSidepanelFormComponent, config)
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

    openIncidentSidepanel(type: RuntimeEventType, process: RuntimeEventProcess, event: RuntimeEventExtended) {
        const config: KbqSidepanelConfig<RuntimeSidepanelIncidentProps> = {
            position: KbqSidepanelPosition.Right,
            hasBackdrop: true,
            data: {
                type,
                process,
                time: event.event.time,
                severity: event.incident_severity,
                threats: event.threats,
                ruleIds: [...event.block_by, ...event.notify_by],
                permissions: this.rulePermissions
            }
        };

        this.sidepanelService
            .open(RuntimeFeatureSidepanelIncidentComponent, config)
            .afterClosed()
            .pipe(take(1))
            .subscribe();
    }

    openViewCodeSidepanel(entity: RuntimeEventEntity) {
        const config: KbqSidepanelConfig<RuntimeSidepanelCodeProps> = {
            position: KbqSidepanelPosition.Right,
            hasBackdrop: true,
            data: {
                content: JSON.stringify(entity)
            }
        };

        this.sidepanelService
            .open(RuntimeFeatureSidepanelCodeComponent, config)
            .afterClosed()
            .pipe(take(1))
            .subscribe();
    }

    private extendEventsBySeverity(events: RuntimeEventExtended[]): RuntimeEventExtended[] {
        return events.map((event) => ({
            ...event,
            threatSeverity: event.threats
                .sort((a, b) => RULE_SEVERITY_ORDER[a.severity] - RULE_SEVERITY_ORDER[b.severity])
                .at(0)?.severity
        }));
    }
}
