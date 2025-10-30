import { Observable } from 'rxjs';

import { DetectorExtended } from '@cs/domains/detector';
import { Notification } from '@cs/domains/notification';
import { PermissionType } from '@cs/domains/role';
import { Rule } from '@cs/domains/rule';

export interface SharedRuleSidepanelProps {
    isDeleted: boolean;
    permissions?: Map<PermissionType, boolean>;
    rule$: Observable<Rule | undefined>;
    detectors$: Observable<DetectorExtended[]>;
    notifications$: Observable<Notification[]>;
    updateHandler?: (rule: Rule) => void;
    deleteHandler?: (id: string) => void;
}
