import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Store } from '@ngrx/store';

import { getDetectors } from '../stores/detector-selector.store';
import {
    CREATE_RUNTIME_DETECTOR_ENTITIES_TODO_ACTION,
    DELETE_RUNTIME_DETECTOR_ENTITY_TODO_ACTION
} from '../stores/detector-action.store';
import { DetectorExtended, DetectorState, DetectorType } from '../interfaces';

@Injectable({
    providedIn: 'root'
})
export class DetectorStoreService {
    readonly detectors$ = (types?: DetectorType[]): Observable<DetectorExtended[]> =>
        this.store.select(getDetectors(types));

    constructor(private readonly store: Store<DetectorState>) {}

    createRuntimeDetectors(base64list: string[]) {
        this.store.dispatch(CREATE_RUNTIME_DETECTOR_ENTITIES_TODO_ACTION({ base64list }));
    }

    deleteRuntimeDetector(key: string, version: number) {
        this.store.dispatch(DELETE_RUNTIME_DETECTOR_ENTITY_TODO_ACTION({ key, version }));
    }
}
