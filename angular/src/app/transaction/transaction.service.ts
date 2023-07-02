import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Transaction, TransactionList } from './transaction';
import * as moment from 'moment';

@Injectable({
  providedIn: 'root'
})
export class TransactionService {

  constructor(private http: HttpClient) {}

  list(from: Date, to: Date): Observable<TransactionList> {
    const fromString = moment(from).format('YYYY-MM-DD');
    const toString = moment(to).format('YYYY-MM-DD');

    return this.http.get<TransactionList>(
      `/api/v1/transactions?from=${fromString}&to=${toString}`
    );
  }

  listForCategoryID(from: Date, to: Date, categoryID: number): Observable<TransactionList> {
    const fromString = moment(from).format('YYYY-MM-DD');
    const toString = moment(to).format('YYYY-MM-DD');

    return this.http.get<TransactionList>(
      `/api/v1/transactions?from=${fromString}&to=${toString}&category=${categoryID}`
    );
  }

  listUnclassified(from: Date, to: Date): Observable<TransactionList> {
    const fromString = moment(from).format('YYYY-MM-DD');
    const toString = moment(to).format('YYYY-MM-DD');

    return this.http.get<TransactionList>(
      `/api/v1/transactions/unclassified?from=${fromString}&to=${toString}`
    );
  }

  uploadCSV(file: File): Observable<any> {
    const formData = new FormData();
    formData.append('file', file, file.name);
    return this.http.post('/api/v1/transactions/csv', formData);
  }

  setCategory(transactionID: number, categoryID: number | null): Observable<any> {
    return this.http.patch(`/api/v1/transaction/${transactionID}`, {
      categoryID: categoryID
    });
  }

  hide(transactionID: number, hide: boolean): Observable<any> {
    return this.http.patch(`/api/v1/transaction/${transactionID}`, {
      hidden: hide
    });
  }

}
