import { Directive, ElementRef, Input, OnInit, Renderer2 } from '@angular/core';

import { RuntimeEventType } from '@cs/domains/runtime';

const RUNTIME_EVENT_TYPE_ICONS = new Map<RuntimeEventType, string>([
    [RuntimeEventType.EXEC, 'exec'],
    [RuntimeEventType.EXIT, 'exit'],
    [RuntimeEventType.KPROBE, 'kprobe'],
    [RuntimeEventType.LOADER, 'loader'],
    [RuntimeEventType.TRACEPOINT, 'tracepoint'],
    [RuntimeEventType.UPROBE, 'uprobe']
]);

@Directive({
    selector: '[runtimeEventTypeIcon]'
})
export class RuntimeFeatureEventTypeIconDirective implements OnInit {
    @Input({ required: true }) type!: RuntimeEventType;

    constructor(
        private readonly el: ElementRef,
        private readonly renderer: Renderer2
    ) {}

    ngOnInit() {
        if (RUNTIME_EVENT_TYPE_ICONS.get(this.type)) {
            const element = this.el.nativeElement;
            const icon = RUNTIME_EVENT_TYPE_ICONS.get(this.type);

            this.renderer.setStyle(element, 'display', 'inline-block');
            this.renderer.setStyle(element, 'width', '18px');
            this.renderer.setStyle(element, 'height', '19px');
            this.renderer.setStyle(
                element,
                'background',
                `center / cover no-repeat url('/assets/images/runtime/icon-${icon}.svg')`
            );
        }
    }
}
