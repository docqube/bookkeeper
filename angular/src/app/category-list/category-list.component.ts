import { Component, Input } from '@angular/core';
import { CategoryTransactions } from '../transaction/transaction';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-category-list',
  templateUrl: './category-list.component.html',
  styleUrls: ['./category-list.component.scss']
})
export class CategoryListComponent {

  @Input() categoryTransactions: Observable<CategoryTransactions[]> | undefined;
  dataSource: CategoryTransactions[] = [];

  displayedColumns: string[] = ['name', 'amount'];

  constructor() {}

  ngOnInit(): void {
    this.categoryTransactions?.subscribe((data) => {
      this.dataSource = data;
    });
  }
}
