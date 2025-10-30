import { ActivatedRoute, Router } from '@angular/router';
import { ChangeDetectionStrategy, Component } from '@angular/core';
import { IModalOptionsForService, KbqModalService, ModalSize } from '@koobiq/components/modal';
import { Observable, map, take } from 'rxjs';

import { ApiPathService } from '@cs/api';
import { I18nService } from '@cs/i18n';
import { RouterName } from '@cs/core';
import { SharedModalService } from '@cs/shared';
import { ClusterStoreService, RegisteredCluster } from '@cs/domains/cluster';
import { DetectorExtended, DetectorStoreService, DetectorType } from '@cs/domains/detector';
import { PermissionName, PermissionType, RolePermissionMap } from '@cs/domains/role';

import { RUNTIME_NAVIGATION_TABS } from '../../constants/runtime-navigation.constant';
import { RuntimeFeatureUploadDetectorModalComponent } from '../../components/upload-detector-modal/runtime-upload-detector-modal.component';

@Component({
    templateUrl: './runtime-detectors.container.html',
    styleUrl: './runtime-detectors.container.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class RuntimeFeatureDetectorsContainer {
    readonly detectors$: Observable<DetectorExtended[]> = this.detectorStoreService
        .detectors$([DetectorType.RUNTIME])
        .pipe(map((detectors) => detectors.sort((a, b) => a.id.localeCompare(b.id))));

    readonly clusters$: Observable<RegisteredCluster[]> = this.clusterStoreService.registeredClusters$;

    readonly activeClusterHost$ = this.apiPathService.host$;

    /* eslint @typescript-eslint/dot-notation: "off" */
    readonly permissions: RolePermissionMap = this.route.snapshot.data['permissions'];

    readonly permissionType = PermissionType;

    readonly permissionName = PermissionName;

    readonly runtimeNavigationTabs = RUNTIME_NAVIGATION_TABS;

    constructor(
        private readonly detectorStoreService: DetectorStoreService,
        private readonly i18nService: I18nService,
        private readonly clusterStoreService: ClusterStoreService,
        private readonly apiPathService: ApiPathService,
        private readonly modalService: KbqModalService,
        private readonly route: ActivatedRoute,
        private readonly router: Router,
        private readonly sharedModalService: SharedModalService
    ) {}

    tabChange(path?: string) {
        this.router.navigate([RouterName.RUNTIME, path]);
    }

    openCreateModal() {
        const config: IModalOptionsForService = {
            kbqComponent: RuntimeFeatureUploadDetectorModalComponent,
            kbqSize: ModalSize.Medium,
            kbqClosable: false
        };

        this.modalService
            .open<RuntimeFeatureUploadDetectorModalComponent, string[] | undefined>(config)
            .afterClose.pipe(take(1))
            .subscribe((base64list?: string[]) => {
                if (base64list && base64list.length) {
                    this.detectorStoreService.createRuntimeDetectors(base64list);
                }
            });
    }

    openDeleteModal(key: string, version: number) {
        this.sharedModalService.delete({
            title: this.i18nService.translate('Runtime.DeleteModal.Content.Title'),
            content: this.i18nService.translate('Runtime.DeleteModal.Content.Text'),
            confirmText: this.i18nService.translate('Runtime.DeleteModal.Button.Confirm'),
            cancelText: this.i18nService.translate('Runtime.DeleteModal.Button.Cancel'),
            confirmHandler: () => {
                this.detectorStoreService.deleteRuntimeDetector(key, version);
            }
        });
    }

    switchCluster(id: string) {
        this.clusterStoreService.switchCluster(id);
    }
}
