import { ClusterStatus, ClusterStatusOption } from '../interfaces';

export const CLUSTER_CREATE_FRAGMENT = 'create';

export const CLUSTER_STATUS: ClusterStatusOption[] = [
    {
        id: ClusterStatus.REGISTERED,
        localizationKey: 'Common.Pseudo.ClusterStatus.Registered'
    },
    {
        id: ClusterStatus.UNREGISTERED,
        localizationKey: 'Common.Pseudo.ClusterStatus.Unregistered'
    }
];
