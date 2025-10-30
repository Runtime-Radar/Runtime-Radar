import { KbqAlertColors } from '@koobiq/components/alert';
import { KbqBadgeColors } from '@koobiq/components/badge';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ActivatedRoute, Router } from '@angular/router';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, DestroyRef, OnDestroy, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormRecord } from '@angular/forms';
import { KbqSidepanelConfig, KbqSidepanelPosition, KbqSidepanelService } from '@koobiq/components/sidepanel';
import {
    Observable,
    combineLatest,
    debounceTime,
    distinctUntilChanged,
    distinctUntilKeyChanged,
    filter,
    map,
    startWith,
    take,
    tap
} from 'rxjs';

import { ApiPathService } from '@cs/api';
import { AuthStoreService } from '@cs/domains/auth';
import { I18nService } from '@cs/i18n';
import { SharedModalService } from '@cs/shared';
import { ClusterStoreService, RegisteredCluster } from '@cs/domains/cluster';
import { DEFAULT_ADMIN_ROLE_ID, PermissionName, PermissionType, RolePermissionMap } from '@cs/domains/role';
import { FormScheme, RouterName, CoreUtilsService as utils } from '@cs/core';
import {
    RUNTIME_EVENT_PROCESSOR,
    RuntimeEventProcessorHistoryControl,
    RuntimeMonitorConfig,
    RuntimeMonitorPermission,
    RuntimeMonitorTracingPolicies,
    RuntimeMonitorTracingPolicy,
    RuntimeStoreService
} from '@cs/domains/runtime';

import { RUNTIME_NAVIGATION_TABS } from '../../constants/runtime-navigation.constant';
import { RUNTIME_SETTINGS_TRACING_POLICIES_PROCESSES_KEY } from '../../constants/runtime-config.constant';
import { RuntimeFeaturePolicyNameService } from '../../services/runtime-policy-name.service';
import { RuntimeFeatureSidepanelPermissionFormComponent } from '../../components/sidepanel-permission-form/runtime-sidepanel-permission-form.component';
import { RuntimeFeatureSidepanelPolicyComponent } from '../../components/sidepanel-policy/runtime-sidepanel-policy.component';
import { RuntimeFeatureSidepanelPolicyFormComponent } from '../../components/sidepanel-policy-form/runtime-sidepanel-policy-form.component';
import { RuntimeFeatureConfigUtilsService as runtimeConfigUtils } from '../../services/runtime-utils.service';
import {
    RuntimeExpertModeForm,
    RuntimeSettingForm,
    RuntimeSettingPermissionCounter,
    RuntimeSettingPermissionForm,
    RuntimeSettingPermissionRecord,
    RuntimeSettingPolicyForm,
    RuntimeSettingPolicyRecord
} from '../../interfaces/runtime-form.interface';
import {
    RuntimeSidepanelPermissionFormProps,
    RuntimeSidepanelPolicyFormProps,
    RuntimeSidepanelPolicyProps
} from '../../interfaces/runtime-sidepanel.interface';

@Component({
    templateUrl: './runtime-settings.container.html',
    styleUrl: './runtime-settings.container.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class RuntimeFeatureSettingsContainer implements OnInit, OnDestroy {
    readonly expertModeForm: FormGroup<FormScheme<RuntimeExpertModeForm>> = this.formBuilder.group({
        isExpert: [false]
    });

    readonly form: FormGroup<FormScheme<RuntimeSettingForm, 'policies' | 'permissions'>> = this.formBuilder.group({
        policies: this.formBuilder.record<RuntimeSettingPolicyRecord>({}),
        permissions: this.formBuilder.record<RuntimeSettingPermissionRecord>({}),
        historyControl: [RuntimeEventProcessorHistoryControl.NONE]
    });

    readonly runtimeHasChanges$: Observable<boolean> = this.runtimeStoreService.runtimeHasChanges$;

    readonly runtimeHasPoliciesChanges$: Observable<boolean> = this.runtimeStoreService.runtimeHasPoliciesChanges$;

    readonly runtimeIsExpertMode$: Observable<boolean> = this.runtimeStoreService.runtimeIsExpertMode$.pipe(
        tap((isExpertMode) => {
            if (isExpertMode) {
                this.expertModeForm.get('isExpert')?.setValue(true, { onlySelf: true });
            }
        })
    );

    readonly runtimeIsOverlayed$: Observable<boolean> = this.runtimeStoreService.runtimeIsOverlayed$;

    readonly clusters$: Observable<RegisteredCluster[]> = this.clusterStoreService.registeredClusters$;

    readonly activeClusterHost$ = this.apiPathService.host$;

    // @todo: replace observable to separate permission
    readonly isAdminRole$: Observable<boolean> = this.authStoreService.credentials$.pipe(
        take(1),
        map((credentials) => credentials.roleId === DEFAULT_ADMIN_ROLE_ID)
    );

    private configSnapshot?: RuntimeMonitorConfig;
    private historyControlSnapshot?: RuntimeEventProcessorHistoryControl;
    private readonly config$: Observable<RuntimeMonitorConfig> = combineLatest([
        this.runtimeStoreService.runtimeMonitorConfig$.pipe(
            map((config) => this.addPseudoTracingPoliciesToConfig(config))
        ),
        this.runtimeStoreService.eventProcessorHistoryControl$
    ]).pipe(
        tap(([config, historyControl]) => {
            if (historyControl) {
                this.historyControlSnapshot = historyControl;
            }
            this.configSnapshot = config;
            this.setForms(config, historyControl);
        }),
        map(([config]) => config)
    );

    private readonly hasChanges$ = this.form.valueChanges.pipe(
        startWith(this.form.value),
        /* eslint @typescript-eslint/no-magic-numbers: "off" */
        debounceTime(250),
        distinctUntilChanged(),
        tap(() => {
            const formValues = utils.getFormValues<RuntimeSettingForm>(this.form.controls);
            this.runtimeStoreService.checkChanges({
                ...runtimeConfigUtils.convertSettingFormToMonitorConfig(formValues),
                historyControl: formValues.historyControl
            });
        })
    );

    readonly permissionCounter$: Observable<RuntimeSettingPermissionCounter> = this.form.valueChanges.pipe(
        startWith(this.form.value),
        distinctUntilKeyChanged('permissions'),
        map(() => utils.getFormValues<RuntimeSettingForm>(this.form.controls).permissions),
        map((permissions) => {
            const values = Object.values(permissions);
            const allowLength = values.filter((item) => !!item.isAllowedType).length;

            return {
                allow: allowLength,
                deny: values.length - allowLength
            };
        })
    );

    private readonly toggleExpertMode$ = this.expertModeForm.valueChanges.pipe(
        map((values) => values.isExpert),
        tap((isExpert) => {
            if (isExpert) {
                this.openSwitchExpertModeModal();
            } else {
                this.runtimeStoreService.switchExpertMode();
            }
        })
    );

    /* eslint @typescript-eslint/dot-notation: "off" */
    readonly permissions: RolePermissionMap = this.route.snapshot.data['permissions'];

    readonly permissionType = PermissionType;

    readonly permissionName = PermissionName;

    readonly runtimeEventProcessorOptions = RUNTIME_EVENT_PROCESSOR;

    readonly runtimeNavigationTabs = RUNTIME_NAVIGATION_TABS;

    readonly alertColors = KbqAlertColors;

    readonly badgeColors = KbqBadgeColors;

    readonly tracingPoliciesProcessesKey = RUNTIME_SETTINGS_TRACING_POLICIES_PROCESSES_KEY;

    private openPolicySidepanelKey: string | null = null;

    get policiesFormGroup(): FormRecord {
        return this.form?.get('policies') as FormRecord;
    }

    get permissionsFormGroup(): FormRecord {
        return this.form?.get('permissions') as FormRecord;
    }

    constructor(
        private readonly authStoreService: AuthStoreService,
        private readonly cdr: ChangeDetectorRef,
        private readonly destroyRef: DestroyRef,
        private readonly formBuilder: FormBuilder,
        private readonly i18nService: I18nService,
        private readonly clusterStoreService: ClusterStoreService,
        private readonly apiPathService: ApiPathService,
        private readonly sharedModalService: SharedModalService,
        private readonly route: ActivatedRoute,
        private readonly router: Router,
        private readonly runtimeStoreService: RuntimeStoreService,
        private readonly runtimeFeaturePolicyNameService: RuntimeFeaturePolicyNameService,
        private readonly sidepanelService: KbqSidepanelService
    ) {}

    ngOnInit() {
        this.config$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe();
        this.hasChanges$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe();
        this.toggleExpertMode$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe();
    }

    ngOnDestroy() {
        if (this.openPolicySidepanelKey) {
            this.sidepanelService.closeAll();
        }
    }

    tabChange(path?: string) {
        this.router.navigate([RouterName.RUNTIME, path]);
    }

    openViewPolicySidepanel(key: string, policy: RuntimeSettingPolicyForm) {
        const config: KbqSidepanelConfig<Partial<RuntimeSidepanelPolicyProps>> = {
            position: KbqSidepanelPosition.Right,
            hasBackdrop: false,
            data: {
                name: policy.name,
                description: policy.description,
                yaml: policy.yaml
            }
        };

        if (this.openPolicySidepanelKey) {
            this.sidepanelService.closeAll();
        }

        if (this.openPolicySidepanelKey !== key) {
            const sidepanel = this.sidepanelService.open(RuntimeFeatureSidepanelPolicyComponent, config);
            sidepanel
                .afterOpened()
                .pipe(take(1))
                .subscribe(() => {
                    this.openPolicySidepanelKey = key;
                });
            sidepanel
                .afterClosed()
                .pipe(take(1))
                .subscribe(() => {
                    this.openPolicySidepanelKey = null;
                });
        }
    }

    openCreatePolicySidepanel() {
        const config: KbqSidepanelConfig<Partial<RuntimeSidepanelPolicyFormProps>> = {
            position: KbqSidepanelPosition.Right,
            hasBackdrop: true,
            data: {}
        };

        if (this.openPolicySidepanelKey) {
            this.sidepanelService.closeAll();
        }

        this.sidepanelService
            .open(RuntimeFeatureSidepanelPolicyFormComponent, config)
            .afterClosed()
            .pipe(take(1), filter(utils.isDefined))
            .subscribe((form: RuntimeSettingPolicyForm) => {
                const control: FormGroup<FormScheme<RuntimeSettingPolicyForm>> = this.formBuilder.group({
                    isEnabled: [false],
                    name: [form.name],
                    description: [form.description],
                    yaml: [form.yaml]
                });

                this.policiesFormGroup.addControl(utils.generateUuid('zz'), control);
                this.runtimeFeaturePolicyNameService.set(form.name);
            });
    }

    openEditPolicySidepanel(policy: RuntimeSettingPolicyForm, key: string) {
        const config: KbqSidepanelConfig<Partial<RuntimeSidepanelPolicyFormProps>> = {
            position: KbqSidepanelPosition.Right,
            hasBackdrop: true,
            data: {
                isEdit: true,
                name: policy.name,
                description: policy.description,
                yaml: policy.yaml
            }
        };

        if (this.openPolicySidepanelKey) {
            this.sidepanelService.closeAll();
        }

        this.sidepanelService
            .open(RuntimeFeatureSidepanelPolicyFormComponent, config)
            .afterClosed()
            .pipe(take(1), filter(utils.isDefined))
            .subscribe((form: RuntimeSettingPolicyForm) => {
                const control: FormGroup<FormScheme<RuntimeSettingPolicyForm>> = this.formBuilder.group({
                    isEnabled: [policy.isEnabled],
                    name: [form.name],
                    description: [form.description],
                    yaml: [form.yaml]
                });

                this.policiesFormGroup.setControl(key, control);
                this.runtimeFeaturePolicyNameService.replace(form.name, policy.name);
                // @todo: refactor this part to remove cdr service
                this.cdr.markForCheck();
            });
    }

    openDeletePolicyModal(key: string, name: string) {
        this.sharedModalService.delete({
            title: this.i18nService.translate('Runtime.PolicyDeleteModal.Content.Title'),
            content: this.i18nService.translate('Runtime.PolicyDeleteModal.Content.Text'),
            confirmText: this.i18nService.translate('Runtime.PolicyDeleteModal.Button.Confirm'),
            cancelText: this.i18nService.translate('Runtime.PolicyDeleteModal.Button.Cancel'),
            confirmHandler: () => {
                this.policiesFormGroup.removeControl(key);
                this.runtimeFeaturePolicyNameService.remove(name);
                // @todo: refactor this part to remove cdr service
                this.cdr.markForCheck();
            }
        });
    }

    openCreatePermissionSidepanel() {
        const config: KbqSidepanelConfig<Partial<RuntimeSidepanelPermissionFormProps>> = {
            position: KbqSidepanelPosition.Right,
            hasBackdrop: true,
            data: {}
        };

        this.sidepanelService
            .open(RuntimeFeatureSidepanelPermissionFormComponent, config)
            .afterClosed()
            .pipe(take(1), filter(utils.isDefined))
            .subscribe((form: RuntimeSettingPermissionForm) => {
                this.permissionsFormGroup.addControl(utils.generateUuid(), this.getPrefilledPermissionsFormGroup(form));
            });
    }

    openEditPermissionSidepanel(key: string, permission: RuntimeSettingPermissionForm) {
        const config: KbqSidepanelConfig<Partial<RuntimeSidepanelPermissionFormProps>> = {
            position: KbqSidepanelPosition.Right,
            hasBackdrop: true,
            data: {
                isEdit: true,
                isAllowedType: permission.isAllowedType,
                namespaces: permission.namespaces,
                pods: permission.pods,
                labels: permission.labels
            }
        };

        this.sidepanelService
            .open(RuntimeFeatureSidepanelPermissionFormComponent, config)
            .afterClosed()
            .pipe(take(1), filter(utils.isDefined))
            .subscribe((form: RuntimeSettingPermissionForm) => {
                this.permissionsFormGroup.setControl(key, this.getPrefilledPermissionsFormGroup(form));
            });
    }

    removePermissionEntity(key: string) {
        this.sharedModalService.delete({
            content: this.i18nService.translate('Runtime.PermissionDeleteModal.Content.Text'),
            confirmText: this.i18nService.translate('Runtime.PermissionDeleteModal.Button.Confirm'),
            cancelText: this.i18nService.translate('Runtime.PermissionDeleteModal.Button.Cancel'),
            confirmHandler: () => {
                this.permissionsFormGroup.removeControl(key);
            }
        });
    }

    create() {
        const formValues = utils.getFormValues<RuntimeSettingForm>(this.form.controls);
        this.runtimeStoreService.createConfig(
            runtimeConfigUtils.convertSettingFormToMonitorConfig(formValues),
            formValues.historyControl
        );

        this.form.updateValueAndValidity();
        this.runtimeFeaturePolicyNameService.clear();
    }

    cancel() {
        this.runtimeStoreService.hideOverlay();
    }

    reset() {
        if (this.configSnapshot && this.historyControlSnapshot) {
            this.setForms(this.configSnapshot, this.historyControlSnapshot);
            this.runtimeFeaturePolicyNameService.clear();
        }
    }

    switchCluster(id: string) {
        this.clusterStoreService.switchCluster(id);
    }

    private openSwitchExpertModeModal() {
        this.sharedModalService.delete({
            title: this.i18nService.translate('Runtime.SwitchModeModal.Content.Title'),
            content: this.i18nService.translate('Runtime.SwitchModeModal.Content.Text'),
            confirmText: this.i18nService.translate('Runtime.SwitchModeModal.Button.Confirm'),
            cancelText: this.i18nService.translate('Runtime.SwitchModeModal.Button.Cancel'),
            confirmHandler: () => {
                this.runtimeStoreService.switchExpertMode();
            },
            cancelHandler: () => {
                this.expertModeForm.get('isExpert')?.setValue(false, { onlySelf: true });
            }
        });
    }

    private getPrefilledPermissionsFormGroup(
        form: RuntimeSettingPermissionForm
    ): FormGroup<FormScheme<RuntimeSettingPermissionForm, never, 'namespaces' | 'pods' | 'labels'>> {
        return this.formBuilder.group({
            isAllowedType: [form.isAllowedType],
            namespaces: this.formBuilder.array<string>(form.namespaces),
            pods: this.formBuilder.array<string>(form.pods),
            labels: this.formBuilder.array<string>(form.labels)
        });
    }

    private setForms(config: RuntimeMonitorConfig, historyControl?: RuntimeEventProcessorHistoryControl) {
        if (historyControl !== undefined) {
            this.form.patchValue({
                historyControl
            });
        }

        this.setPoliciesForm(config.tracing_policies);
        this.setPermissionsForm(config.allow_list, config.deny_list);
        if (!this.permissions[PermissionName.SYSTEM].has(PermissionType.UPDATE)) {
            this.form.get('historyControl')?.disable({ onlySelf: true });
        }
    }

    private setPoliciesForm(tracingPolicies: RuntimeMonitorTracingPolicies) {
        Object.keys(this.policiesFormGroup.controls).forEach((key) => {
            this.policiesFormGroup.removeControl(key);
        });

        Object.keys(tracingPolicies).forEach((key) => {
            const policyValue: RuntimeMonitorTracingPolicy | null = tracingPolicies[key];
            const policyForm: FormGroup<FormScheme<RuntimeSettingPolicyForm>> = this.formBuilder.group({
                isEnabled: [policyValue.enabled || false],
                name: [policyValue.name || ''],
                description: [policyValue.description || ''],
                yaml: [policyValue.yaml || '']
            });

            this.policiesFormGroup.addControl(key, policyForm);
        });
    }

    private setPermissionsForm(
        allowValues: Partial<RuntimeMonitorPermission>[],
        denyValues: Partial<RuntimeMonitorPermission>[]
    ) {
        Object.keys(this.permissionsFormGroup.controls).forEach((key) => {
            this.permissionsFormGroup.removeControl(key);
        });

        allowValues.forEach((item, index) => {
            this.permissionsFormGroup.addControl(
                utils.generateUuid(index.toString()),
                this.getPrefilledPermissionsFormGroup({
                    isAllowedType: true,
                    namespaces: item.namespace || [],
                    pods: item.pod_regex || [],
                    labels: item.labels || []
                })
            );
        });

        denyValues.forEach((item, index) => {
            this.permissionsFormGroup.addControl(
                utils.generateUuid(index.toString()),
                this.getPrefilledPermissionsFormGroup({
                    isAllowedType: false,
                    namespaces: item.namespace || [],
                    pods: item.pod_regex || [],
                    labels: item.labels || []
                })
            );
        });
    }

    private addPseudoTracingPoliciesToConfig(config: RuntimeMonitorConfig): RuntimeMonitorConfig {
        const pseudoTracingPoliciesConfig: RuntimeMonitorTracingPolicies = {
            [RUNTIME_SETTINGS_TRACING_POLICIES_PROCESSES_KEY]: {
                name: this.i18nService.translate('Runtime.SettingsPage.Policies.Checkbox.RunStopProcesses'),
                description: this.i18nService.translate('Runtime.SettingsPage.Policies.Popover.RunStopProcesses'),
                enabled: true
            }
        };

        return {
            ...config,
            tracing_policies: { ...pseudoTracingPoliciesConfig, ...config.tracing_policies }
        };
    }
}
