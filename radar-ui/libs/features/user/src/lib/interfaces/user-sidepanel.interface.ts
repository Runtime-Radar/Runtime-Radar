import { Observable } from 'rxjs';

import { AuthCredentials } from '@cs/domains/auth';
import { Role } from '@cs/domains/role';
import { User } from '@cs/domains/user';

export interface UserSidepanelFormProps {
    credentials$: Observable<AuthCredentials>;
    roles$: Observable<Role[]>;
    user: Partial<User>;
    isEdit: boolean;
}
