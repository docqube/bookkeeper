import { Component, Input } from '@angular/core';
import { Transaction } from '../transaction/transaction';
import { Observable } from 'rxjs';
import { Category } from '../category/category';
import * as moment from 'moment';

@Component({
  selector: 'app-transaction-list',
  templateUrl: './transaction-list.component.html'
})
export class TransactionListComponent {

  @Input('transactions') transactionsObservable: Observable<Transaction[]> | undefined;
  @Input() categories: Category[] | undefined;
  transactions: Transaction[] = [];

  datedTransactions: { [key: string]: Transaction[] } = {};

  constructor() {}

  ngOnInit(): void {
    this.transactionsObservable?.subscribe((data) => {
      this.transactions = data;

      if (!this.transactions) {
        return;
      }

      this.datedTransactions = {};
      for (const transaction of this.transactions) {
        const dateString = moment(transaction.bookingDate).format('YYYY-MM-DD');
        if (!this.datedTransactions[dateString]) {
          this.datedTransactions[dateString] = [];
        }
        this.datedTransactions[dateString].push(transaction);
      }
    });
  }

  getCategory(id: number): Category | undefined {
    return this.categories?.find((category) => category.id === id);
  }
}
