import { Injectable } from '@angular/core';
import { KbqModalService } from '@koobiq/components/modal';
import { take } from 'rxjs';

import { SharedModalComponent } from './shared-modal.component';
import { SharedModalParams } from './shared-modal.interface';

@Injectable({
    providedIn: 'root'
})
export class SharedModalService {
    constructor(private readonly modalService: KbqModalService) {}

    delete(params: SharedModalParams) {
        this.modalService
            .open<SharedModalComponent, boolean | undefined>({
                kbqContent: SharedModalComponent,
                // @todo: prop is deprecated, replace after koobiq@18 migration
                kbqComponentParams: params,
                kbqClosable: false
            })
            .afterClose.pipe(take(1))
            .subscribe((isSuccessful?: boolean) => {
                if (isSuccessful) {
                    params.confirmHandler();
                } else if (params.cancelHandler) {
                    params.cancelHandler();
                }
            });
    }
}
