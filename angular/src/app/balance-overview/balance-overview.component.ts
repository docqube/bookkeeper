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
    return Math.abs(category.transactionList.sum) / this.getIncomeSum() * 100;
  }

  getExpenseCategories(): CategoryTransactions[] {
    return this.categoryTransactions.filter((data) => {
      return data.transactionList.sum < 0;
    });
  }

  getIncomeSum(): number {
    return this.categoryTransactions.filter((data) => data.transactionList.sum > 0).reduce((acc, data) => acc + data.transactionList.sum, 0);
  }

  getExpenseSum(): number {
    return Math.abs(this.getExpenseCategories().reduce((acc, data) => acc + data.transactionList.sum, 0));
  }
}
