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
  @Input('categories') categoriesObservable: Observable<Category[]> | undefined;

  transactions: Transaction[] = [];
  categories: Category[] | undefined;

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

    this.categoriesObservable?.subscribe((data) => {
      this.categories = data;
    });
  }

  getCategory(id: number): Category | undefined {
    return this.categories?.find((category) => category.id === id);
  }

  getCategoryColor(id?: number): string {
    let color = '#64748b';
    if (!id) {
      return color;
    }
    const category = this.getCategory(id);
    if (!category) {
      return color;
    }

    if (category.color) {
      return `${category.color}`;
    }
    return color;
  }
}
