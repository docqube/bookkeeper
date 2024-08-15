import { LOCALE_ID, NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { DashboardModule } from './dashboard/dashboard.module';
import { SharedModule } from './shared/shared.module';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import * as de from '@angular/common/locales/de';
import { registerLocaleData } from '@angular/common';

@NgModule({ declarations: [
        AppComponent
    ],
    bootstrap: [AppComponent], imports: [SharedModule.forRoot(),
        BrowserModule,
        AppRoutingModule,
        BrowserAnimationsModule,
        DashboardModule], providers: [
        {
            provide: LOCALE_ID,
            useValue: 'de-DE'
        },
        provideHttpClient(withInterceptorsFromDi()),
    ] })
export class AppModule {
  constructor() {
    registerLocaleData(de.default);
  }
}
