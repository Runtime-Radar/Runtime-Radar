import { ChangeDetectionStrategy, Component, Input, booleanAttribute } from '@angular/core';

@Component({
    selector: 'cs-event-action-component',
    templateUrl: './shared-event-action.component.html',
    styleUrl: './shared-event-action.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class SharedEventActionComponent {
    @Input() blockIds?: string[];

    @Input() notifyIds?: string[];

    @Input() testLocator?: string;

    @Input({ transform: booleanAttribute }) isIconVisible = false;

    @Input({ transform: booleanAttribute }) isEmptyStateShown = false;
}
