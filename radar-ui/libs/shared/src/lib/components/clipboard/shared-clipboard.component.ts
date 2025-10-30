import { ChangeDetectionStrategy, Component, Input, booleanAttribute } from '@angular/core';

import { SharedClipboardService } from '@cs/shared';

@Component({
    selector: 'cs-clipboard-component',
    templateUrl: './shared-clipboard.component.html',
    styleUrl: './shared-clipboard.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class SharedClipboardComponent {
    @Input({ required: true }) value!: string | null;

    @Input({ transform: booleanAttribute }) isButtonTextVisible = false;

    @Input() isDisabled? = false;

    constructor(private readonly clipboardService: SharedClipboardService) {}

    copyToClipboard() {
        if (this.value) {
            this.clipboardService.copyToClipboard(this.value);
        }
    }
}
