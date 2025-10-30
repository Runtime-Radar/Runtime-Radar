export enum DetectorType {
    RUNTIME = 'RUNTIME'
}

export interface Detector {
    id: string;
    name: string;
    description: string;
    version: number;
    author?: string;
    contact?: string;
    license?: string;
}
