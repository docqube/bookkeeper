import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Transaction } from './transaction';
import * as moment from 'moment';

@Injectable({
  providedIn: 'root'
})
export class TransactionService {

  constructor(private http: HttpClient) {}

  list(from: Date, to: Date): Observable<Transaction[]> {
    const fromString = moment(from).format('YYYY-MM-DD');
    const toString = moment(to).format('YYYY-MM-DD');

    return this.http.get<Transaction[]>(
      `/api/v1/transactions?from=${fromString}&to=${toString}`
    );
  }

  listForCategoryID(from: Date, to: Date, categoryID: number): Observable<Transaction[]> {
    const fromString = moment(from).format('YYYY-MM-DD');
    const toString = moment(to).format('YYYY-MM-DD');

    return this.http.get<Transaction[]>(
      `/api/v1/transactions?from=${fromString}&to=${toString}&category=${categoryID}`
    );
  }

  listUnclassified(from: Date, to: Date): Observable<Transaction[]> {
    const fromString = moment(from).format('YYYY-MM-DD');
    const toString = moment(to).format('YYYY-MM-DD');

    return this.http.get<Transaction[]>(
      `/api/v1/transactions/unclassified?from=${fromString}&to=${toString}`
    );
  }

  uploadCSV(file: File): Observable<any> {
    const formData = new FormData();
    formData.append('file', file, file.name);
    return this.http.post('/api/v1/transactions/csv', formData);
  }

}
