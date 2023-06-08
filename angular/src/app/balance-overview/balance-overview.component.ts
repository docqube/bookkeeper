import { Component, Input } from '@angular/core';
import { CategoryTransactions } from '../transaction/transaction';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-balance-overview',
  templateUrl: './balance-overview.component.html',
  styleUrls: ['./balance-overview.component.scss']
})
export class BalanceOverviewComponent {

  @Input() categoryTransactions: Observable<CategoryTransactions[]> | undefined;
  dataSource: CategoryTransactions[] = [];

  displayedColumns: string[] = ['name', 'amount'];

  constructor() {}

  ngOnInit(): void {
    this.categoryTransactions?.subscribe((data) => {
      this.dataSource = data;
    });
  }

  getExpenseCategoryPercentage(category: CategoryTransactions): number {
    return Math.abs(category.transactionsSum) / this.getIncomeSum() * 100;
  }

  getExpenseCategories(): CategoryTransactions[] {
    return this.dataSource.filter((data) => {
      return data.transactionsSum < 0;
    });
  }

  getIncomeSum(): number {
    return this.dataSource.filter((data) => data.transactionsSum > 0).reduce((acc, data) => acc + data.transactionsSum, 0);
  }

  getExpenseSum(): number {
    return Math.abs(this.getExpenseCategories().reduce((acc, data) => acc + data.transactionsSum, 0));
  }
}
