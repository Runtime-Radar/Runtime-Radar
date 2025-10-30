export enum RuntimeRouterName {
    SETTINGS = 'settings',
    EVENTS = 'events',
    DETECTORS = 'detectors'
}

export interface RuntimeNavigationTab {
    path: RuntimeRouterName;
    localizationKey: string;
}
