import { NgModule, ModuleWithProviders, SecurityContext } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { SvgIconComponent } from './svg-icon/svg-icon.component';


@NgModule({
  imports: [
    CommonModule,
  ],
  exports: [
    CommonModule,
    ReactiveFormsModule,
    SvgIconComponent,
  ],
  declarations: [
    SvgIconComponent,
  ],
})
export class SharedModule {
  static forRoot(): ModuleWithProviders<SharedModule> {
    return {
      ngModule: SharedModule,
      providers: [
      ]
    };
  }
}
