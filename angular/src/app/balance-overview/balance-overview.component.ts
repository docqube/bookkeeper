import { Component, Input } from '@angular/core';
import { CategoryTransactions } from '../transaction/transaction';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-balance-overview',
  templateUrl: './balance-overview.component.html',
  styleUrls: ['./balance-overview.component.scss']
})
export class BalanceOverviewComponent {

  @Input('categoryTransactions') $categoryTransactions: Observable<CategoryTransactions[]> | undefined;
  categoryTransactions: CategoryTransactions[] = [];

  constructor() {}

  ngOnInit(): void {
    this.$categoryTransactions?.subscribe((data) => {
      this.categoryTransactions = data;
    });
  }

  getExpenseCategoryPercentage(category: CategoryTransactions): number {
    return Math.abs(category.transactionsSum) / this.getIncomeSum() * 100;
  }

  getExpenseCategories(): CategoryTransactions[] {
    return this.categoryTransactions.filter((data) => {
      return data.transactionsSum < 0;
    });
  }

  getIncomeSum(): number {
    return this.categoryTransactions.filter((data) => data.transactionsSum > 0).reduce((acc, data) => acc + data.transactionsSum, 0);
  }

  getExpenseSum(): number {
    return Math.abs(this.getExpenseCategories().reduce((acc, data) => acc + data.transactionsSum, 0));
  }
}
