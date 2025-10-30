import { CommonModule } from '@angular/common';
import { KbqButtonModule } from '@koobiq/components/button';
import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { KbqModalModule, KbqModalRef } from '@koobiq/components/modal';

@Component({
    templateUrl: './shared-modal.component.html',
    styleUrl: './shared-modal.component.scss',
    imports: [KbqModalModule, KbqButtonModule, CommonModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: true
})
export class SharedModalComponent {
    @Input() title?: string;

    @Input({ required: true }) content!: string;

    @Input({ required: true }) confirmText!: string;

    @Input({ required: true }) cancelText!: string;

    constructor(private readonly modal: KbqModalRef) {}

    close(isSuccessful: boolean) {
        this.modal.destroy(isSuccessful || undefined);
    }
}
