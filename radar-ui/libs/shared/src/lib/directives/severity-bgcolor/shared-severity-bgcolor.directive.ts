import { Directive, ElementRef, Input, OnChanges, Renderer2 } from '@angular/core';

import { RuleSeverity } from '@cs/domains/rule';

const SEVERITY_BACKGROUND_COLORS = new Map<RuleSeverity, string>([
    [RuleSeverity.CRITICAL, '#fff3f3'],
    [RuleSeverity.HIGH, '#fef2ef'],
    [RuleSeverity.MEDIUM, '#fff4dd'],
    [RuleSeverity.LOW, '#eff6ff'],
    [RuleSeverity.NONE, 'inherit']
]);

@Directive({
    selector: '[severityBgColor]'
})
export class SharedSeverityBgColorDirective implements OnChanges {
    @Input() severity?: RuleSeverity | null;

    constructor(
        private readonly el: ElementRef,
        private readonly renderer: Renderer2
    ) {}

    ngOnChanges() {
        if (this.severity) {
            this.renderer.setStyle(
                this.el.nativeElement,
                'background-color',
                SEVERITY_BACKGROUND_COLORS.get(this.severity)
            );
        }
    }
}
