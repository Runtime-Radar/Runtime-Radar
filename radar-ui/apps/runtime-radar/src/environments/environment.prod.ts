import { isChildCluster } from './argument/child-cluster';

export const environment = {
    api: '/api/v1/',
    singleTenant: ['signin', 'tokens', 'user', 'role', 'cluster'],
    childCluster: isChildCluster,
    pollingInterval: 60000,
    refreshInterval: 300000,
    production: true
};
