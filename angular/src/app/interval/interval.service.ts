import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { FiscalMonth } from './interval';

@Injectable({
  providedIn: 'root'
})
export class IntervalService {

  constructor(private http: HttpClient) {}

  getFiscalMonth(month: number, year: number, incomeCategoryID: number): Observable<FiscalMonth> {
    return this.http.get<FiscalMonth>(
      `/api/v1/interval/fiscal-month?month=${month}&year=${year}&income_category_id=${incomeCategoryID}`
    );
  }

}
