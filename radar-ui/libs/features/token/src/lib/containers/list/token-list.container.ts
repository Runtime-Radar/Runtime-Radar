import { ActivatedRoute } from '@angular/router';
import { DateTime } from 'luxon';
import { ChangeDetectionStrategy, Component } from '@angular/core';
import { KbqSidepanelConfig, KbqSidepanelPosition, KbqSidepanelService } from '@koobiq/components/sidepanel';
import { Observable, filter, map, switchMap, take } from 'rxjs';

import { ApiPathService } from '@cs/api';
import { AuthStoreService } from '@cs/domains/auth';
import { I18nService } from '@cs/i18n';
import { SharedModalService } from '@cs/shared';
import { ClusterStoreService, RegisteredCluster } from '@cs/domains/cluster';
import { CreateTokenRequest, Token, TokenPermissionName, TokenStoreService } from '@cs/domains/token';
import { LoadStatus, CoreUtilsService as utils } from '@cs/core';
import { PermissionName, PermissionType, RolePermissionMap } from '@cs/domains/role';

import { TokenFeatureSidepanelFormComponent } from '../../components/sidepanel-form/token-sidepanel-form.component';
import { TokenSidepanelFormProps } from '../../interfaces/token-sidepanel.interface';
import { TokenForm, TokenPermissionForm, TokenPermissionType } from '../../interfaces/token-form.interface';

@Component({
    templateUrl: './token-list.container.html',
    styleUrl: './token-list.container.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class TokenFeatureListContainer {
    readonly tokens$: Observable<Token[]> = this.tokenStoreService.tokens$.pipe(
        map((tokens) => tokens.sort((a, b) => (a.access_token === b.access_token ? 0 : a.access_token ? -1 : 1)))
    );

    readonly loadStatus$: Observable<LoadStatus> = this.tokenStoreService.loadStatus$;

    readonly clusters$: Observable<RegisteredCluster[]> = this.clusterStoreService.registeredClusters$;

    readonly activeClusterHost$ = this.apiPathService.host$;

    /* eslint @typescript-eslint/dot-notation: "off" */
    readonly permissions: RolePermissionMap = this.route.snapshot.data['permissions'];

    readonly permissionType = PermissionType;

    readonly permissionName = PermissionName;

    readonly loadStatus = LoadStatus;

    readonly tokenPermissionName = TokenPermissionName;

    readonly dateTimeFullFormat = DateTime.DATETIME_FULL;

    constructor(
        private readonly authStoreService: AuthStoreService,
        private readonly apiPathService: ApiPathService,
        private readonly clusterStoreService: ClusterStoreService,
        private readonly i18nService: I18nService,
        private readonly route: ActivatedRoute,
        private readonly sharedModalService: SharedModalService,
        private readonly sidepanelService: KbqSidepanelService,
        private readonly tokenStoreService: TokenStoreService
    ) {}

    openCreateSidepanel() {
        const config: KbqSidepanelConfig<TokenSidepanelFormProps> = {
            position: KbqSidepanelPosition.Right,
            hasBackdrop: true,
            data: {
                permissions: this.permissions
            }
        };

        this.sidepanelService
            .open(TokenFeatureSidepanelFormComponent, config)
            .afterClosed()
            .pipe(
                take(1),
                filter(utils.isDefined),
                switchMap((form: TokenForm) =>
                    this.authStoreService.credentials$.pipe(
                        take(1),
                        map((credentials) => ({
                            userId: credentials.userId,
                            form
                        }))
                    )
                )
            )
            .subscribe(({ form, userId }) => {
                const request: CreateTokenRequest = {
                    name: form.name,
                    expires_at: form.date ? form.date.toJSDate().toISOString() : null, // RFC3339
                    user_id: userId,
                    permissions: {
                        [TokenPermissionName.RULES]: {
                            actions: this.getPermissionActions(form.permissions[TokenPermissionName.RULES]),
                            description: ''
                        },
                        [TokenPermissionName.EVENTS]: {
                            actions: this.getPermissionActions(form.permissions[TokenPermissionName.EVENTS]),
                            description: ''
                        }
                    }
                };

                this.tokenStoreService.createToken(request);
            });
    }

    openDeleteModal(id: string) {
        this.sharedModalService.delete({
            title: this.i18nService.translate('Token.DeleteModal.Content.Title'),
            content: this.i18nService.translate('Token.DeleteModal.Content.Text'),
            confirmText: this.i18nService.translate('Token.DeleteModal.Button.Confirm'),
            cancelText: this.i18nService.translate('Token.DeleteModal.Button.Cancel'),
            confirmHandler: () => {
                this.tokenStoreService.deleteToken(id);
            }
        });
    }

    openRevokeModal() {
        this.sharedModalService.delete({
            title: this.i18nService.translate('Token.RevokeModal.Content.Title'),
            content: this.i18nService.translate('Token.RevokeModal.Content.Text'),
            confirmText: this.i18nService.translate('Token.RevokeModal.Button.Confirm'),
            cancelText: this.i18nService.translate('Token.RevokeModal.Button.Cancel'),
            confirmHandler: () => {
                this.tokenStoreService.revokeTokens();
            }
        });
    }

    switchCluster(id: string) {
        this.clusterStoreService.switchCluster(id);
    }

    private getPermissionActions(form: TokenPermissionForm): PermissionType[] {
        return Object.entries(form).reduce((acc, [key, value]) => {
            if (!value) {
                return acc;
            }

            switch (key) {
                case TokenPermissionType.CREATE:
                    acc.push(PermissionType.CREATE);
                    break;
                case TokenPermissionType.READ:
                    acc.push(PermissionType.READ);
                    break;
                case TokenPermissionType.UPDATE:
                    acc.push(PermissionType.UPDATE);
                    break;
                case TokenPermissionType.DELETE:
                    acc.push(PermissionType.DELETE);
                    break;
            }

            return acc;
        }, [] as PermissionType[]);
    }
}
