import { Directive, ElementRef, Input, OnChanges, Renderer2 } from '@angular/core';

import { RuleSeverity, RuleVerdict } from '@cs/domains/rule';

const SEVERITY_COLORS = new Map<RuleSeverity, string>([
    [RuleSeverity.CRITICAL, '#ff0000'],
    [RuleSeverity.HIGH, '#e8612d'],
    [RuleSeverity.MEDIUM, '#e19f12'],
    [RuleSeverity.LOW, '#4187ff'],
    [RuleSeverity.NONE, 'inherit']
]);

const VERDICT_COLORS = new Map<RuleVerdict, string>([
    [RuleVerdict.DANGEROUS, '#e8612d'],
    [RuleVerdict.UNWANTED, '#e19f12'],
    [RuleVerdict.CLEAN, '#4187ff'],
    [RuleVerdict.NONE, 'inherit']
]);

@Directive({
    selector: '[severityColor]'
})
export class SharedSeverityColorDirective implements OnChanges {
    @Input() severity?: RuleSeverity | null;

    @Input() verdict?: RuleVerdict | null;

    constructor(
        private readonly el: ElementRef,
        private readonly renderer: Renderer2
    ) {}

    ngOnChanges() {
        if (this.severity) {
            this.renderer.setStyle(this.el.nativeElement, 'color', SEVERITY_COLORS.get(this.severity));
        }

        if (this.verdict) {
            this.renderer.setStyle(this.el.nativeElement, 'color', VERDICT_COLORS.get(this.verdict));
        }
    }
}
