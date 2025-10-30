declare const window: {
    env?: {
        [key: string]: unknown;
    };
};

export const isChildCluster: boolean = ((window['env'] && window['env']['isChildCluster']) as boolean) || false;
