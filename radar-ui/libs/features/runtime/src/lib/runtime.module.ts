import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { DetectorDomainModule } from '@cs/domains/detector';
import { I18nModule } from '@cs/i18n';
import { NotificationDomainModule } from '@cs/domains/notification';
import { RuleDomainModule } from '@cs/domains/rule';
import { RuntimeDomainModule } from '@cs/domains/runtime';
import { SharedModule } from '@cs/shared';

import { RuntimeFeatureByteFormatterPipe } from './pipes/runtime-byte-formatter.pipe';
import { RuntimeFeatureContextPopoverComponent } from './components/context-popover/runtime-context-popover.component';
import { RuntimeFeatureDateTimePeriodPickerComponent } from './components/datetime-period-picker/runtime-datetime-period-picker.component';
import { RuntimeFeatureDetailsContainer } from './containers/details/runtime-details.container';
import { RuntimeFeatureDetectorsContainer } from './containers/detectors/runtime-detectors.container';
import { RuntimeFeatureEventCounterComponent } from './components/event-counter/runtime-event-counter.component';
import { RuntimeFeatureEventTypeIconDirective } from './directives/runtime-event-type-icon.directive';
import { RuntimeFeatureEventsContainer } from './containers/events/runtime-events.container';
import { RuntimeFeatureEventsGridContainer } from './containers/events-grid/runtime-events-grid.container';
import { RuntimeFeatureFilterPopoverComponent } from './components/filter-popover/runtime-filter-popover.component';
import { RuntimeFeatureHistoryDropdownComponent } from './components/history-dropdown/runtime-history-dropdown.component';
import { RuntimeFeatureHistoryLabelPipe } from './pipes/runtime-history-label.pipe';
import { RuntimeFeatureNanosecondsFormatterPipe } from './pipes/runtime-nanoseconds-formatter.pipe';
import { RuntimeFeaturePermissionsFilterPipe } from './pipes/runtime-permissions-filter.pipe';
import { RuntimeFeaturePresetDropdownComponent } from './components/preset-dropdown/runtime-preset-dropdown.component';
import { RuntimeFeatureRoutingModule } from './runtime-routing.module';
import { RuntimeFeatureRulesContainer } from './containers/rules/runtime-rules.container';
import { RuntimeFeatureSettingsContainer } from './containers/settings/runtime-settings.container';
import { RuntimeFeatureSeverityThreatsCounterPipe } from './pipes/runtime-severity-threats-counter.pipe';
import { RuntimeFeatureSidepanelCodeComponent } from './components/sidepanel-code/runtime-sidepanel-code.component';
import { RuntimeFeatureSidepanelIncidentComponent } from './components/sidepanel-incident/runtime-sidepanel-incident.component';
import { RuntimeFeatureSidepanelPermissionFormComponent } from './components/sidepanel-permission-form/runtime-sidepanel-permission-form.component';
import { RuntimeFeatureSidepanelPolicyComponent } from './components/sidepanel-policy/runtime-sidepanel-policy.component';
import { RuntimeFeatureSidepanelPolicyFormComponent } from './components/sidepanel-policy-form/runtime-sidepanel-policy-form.component';
import { RuntimeFeatureSidepanelThreatsComponent } from './components/sidepanel-threats/runtime-sidepanel-threats.component';
import { RuntimeFeatureUploadDetectorModalComponent } from './components/upload-detector-modal/runtime-upload-detector-modal.component';

@NgModule({
    imports: [
        CommonModule,
        FormsModule,
        I18nModule,
        DetectorDomainModule,
        NotificationDomainModule,
        ReactiveFormsModule,
        RuleDomainModule,
        RuntimeDomainModule,
        RuntimeFeatureRoutingModule,
        SharedModule
    ],
    declarations: [
        RuntimeFeatureByteFormatterPipe,
        RuntimeFeatureContextPopoverComponent,
        RuntimeFeatureDateTimePeriodPickerComponent,
        RuntimeFeatureEventTypeIconDirective,
        RuntimeFeatureFilterPopoverComponent,
        RuntimeFeatureHistoryDropdownComponent,
        RuntimeFeatureHistoryLabelPipe,
        RuntimeFeaturePresetDropdownComponent,
        RuntimeFeatureDetailsContainer,
        RuntimeFeatureDetectorsContainer,
        RuntimeFeatureEventCounterComponent,
        RuntimeFeatureEventsContainer,
        RuntimeFeatureEventsGridContainer,
        RuntimeFeatureRulesContainer,
        RuntimeFeatureNanosecondsFormatterPipe,
        RuntimeFeaturePermissionsFilterPipe,
        RuntimeFeatureSettingsContainer,
        RuntimeFeatureSeverityThreatsCounterPipe,
        RuntimeFeatureSidepanelCodeComponent,
        RuntimeFeatureSidepanelIncidentComponent,
        RuntimeFeatureSidepanelPermissionFormComponent,
        RuntimeFeatureSidepanelPolicyComponent,
        RuntimeFeatureSidepanelPolicyFormComponent,
        RuntimeFeatureSidepanelThreatsComponent,
        RuntimeFeatureUploadDetectorModalComponent
    ]
})
export class RuntimeFeatureModule {}
