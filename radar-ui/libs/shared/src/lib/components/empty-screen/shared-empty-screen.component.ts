import { ChangeDetectionStrategy, Component, Input } from '@angular/core';

@Component({
    selector: 'cs-empty-screen-component',
    templateUrl: './shared-empty-screen.component.html',
    styleUrl: './shared-empty-screen.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class SharedEmptyScreenComponent {
    @Input({ required: true }) title!: string;

    @Input({ required: true }) description!: string;

    @Input() imageUrl?: string;
}
