import { DOCUMENT } from '@angular/common';
import { Inject, Injectable } from '@angular/core';

@Injectable({
    providedIn: 'root'
})
export class CoreWindowService {
    readonly document: Document;

    private readonly window: Window;

    get location(): Location {
        return this.window.location;
    }

    get navigator(): Navigator {
        return this.window.navigator;
    }

    get localStorage(): Storage {
        return this.window.localStorage;
    }

    get sessionStorage(): Storage {
        return this.window.sessionStorage;
    }

    get isSecureContext(): boolean {
        return this.window.isSecureContext;
    }

    constructor(@Inject(DOCUMENT) private readonly nativeDocument: Document) {
        // it's a workaround to have document property properly typed
        // see: https://github.com/angular/angular/issues/15640
        if (!this.nativeDocument.defaultView) {
            throw new Error('window is not available');
        }

        this.window = this.nativeDocument.defaultView;
        this.document = this.nativeDocument;
    }

    atob(data: string): string {
        return this.window.atob(data);
    }

    btoa(data: string): string {
        return this.window.btoa(data);
    }
}
