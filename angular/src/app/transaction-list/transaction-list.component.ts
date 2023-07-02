import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Transaction } from '../transaction/transaction';
import { Observable, Subject, takeUntil } from 'rxjs';
import { Category } from '../category/category';
import * as moment from 'moment';

@Component({
  selector: 'app-transaction-list',
  templateUrl: './transaction-list.component.html'
})
export class TransactionListComponent {

  @Input('transactions') $transactions: Observable<Transaction[]>;
  @Input('categories') $categories: Observable<Category[]>;

  @Output() transactionChanged: EventEmitter<Transaction> = new EventEmitter();

  transactions: Transaction[] = [];
  categories: Category[];

  datedTransactions: { [key: string]: Transaction[] } = {};

  private ngUnsubscribe: Subject<any> = new Subject();

  constructor() {}

  ngOnInit(): void {
    this.$transactions?.subscribe((data) => {
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

    this.$categories?.
      pipe(takeUntil(this.ngUnsubscribe)).
      subscribe((data) => {
        this.categories = data;
      });
  }

  ngOnDestroy(): void {
    this.ngUnsubscribe.next(null);
    this.ngUnsubscribe.complete();
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

  handleTransactionChange(transaction: Transaction): void {
    this.transactionChanged.emit(transaction);
  }
}
