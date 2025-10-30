import { Detector } from './detector-contract.interface';

export interface GetDetectorsRequest {
    page_size: number;
}

export interface GetDetectorsResponse {
    detectors: Detector[];
    total: number;
}

export interface CreateDetectorRequest {
    wasmBase64: string;
}

export interface CreateDetectorResponse {
    detector: Detector;
}

export interface DeleteDetectorRequest {
    id: string;
    version: number;
}

export type EmptyDetectorResponse = Record<string, unknown>;
