import { Component } from '@angular/core';
import * as moment from 'moment';
import { Observable, Subject, combineLatest, forkJoin, merge, takeUntil } from 'rxjs';
import { CategoryTransactions, Transaction, TransactionList } from '../transaction/transaction';
import { CategoryService } from '../category/category.service';
import { TransactionService } from '../transaction/transaction.service';
import { FormControl } from '@angular/forms';
import { Category, CategoryColors } from '../category/category';
import { IntervalService } from '../interval/interval.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent {

  loaded = false;
  error = false;

  categoryTransactions: CategoryTransactions[] = [];
  transactions: Transaction[] = [];
  categories: Category[] = [];

  $categoryTransactions: Subject<CategoryTransactions[]> = new Subject();
  $transactions: Subject<Transaction[]> = new Subject();
  $categories: Subject<Category[]> = new Subject();

  monthFormControl: FormControl = new FormControl<number>(moment().month());
  yearFormControl: FormControl = new FormControl<number>(moment().year());
  startDate: Date = moment().startOf('month').toDate();
  endDate: Date = this.startDate;

  private categoryColorPalette: string[] = CategoryColors;

  private ngUnsubscribe: Subject<any> = new Subject();

  constructor(private categoryService: CategoryService,
              private transactionService: TransactionService,
              private intervalService: IntervalService) {}

  ngOnInit(): void {
    this.loadData();
  }

  ngOnDestroy(): void {
    this.ngUnsubscribe.next(null);
    this.ngUnsubscribe.complete();
  }

  loadData(): void {
    this.loaded = false;
    this.error = false;

    this.intervalService.getFiscalMonth(this.monthFormControl.value + 1, this.yearFormControl.value, 1)
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe({
        next: (fiscalMonth) => {
          this.startDate = fiscalMonth.start;
          this.endDate = fiscalMonth.end;

          this.loadCategories();
          this.loadTransactions();
        },
        error: (err) => {
          this.loaded = true;
          this.error = true;
        }
      });
  }

  loadCategories(): void {
    this.categoryTransactions = [];
    this.categoryService.list()
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe({
        next: (categories) => {
          this.categories = categories;

          const transactionRequests: Observable<TransactionList>[] = [];
          for (const category of categories) {
            transactionRequests.push(this.transactionService.listForCategoryID(this.startDate, this.endDate, category.id));
          }
          transactionRequests.push(this.transactionService.listUnclassified(this.startDate, this.endDate));

          combineLatest(transactionRequests)
            .pipe(takeUntil(this.ngUnsubscribe))
            .subscribe({
              next: (transactionLists) => {
                for (let i = 0; i < transactionLists.length; i++) {
                  let category: Category | undefined = undefined;
                  if (i < categories.length) {
                    category = categories[i];
                  }
                  this.categoryTransactions.push({
                    category: category,
                    transactionList: transactionLists[i]
                  });
                }

                this.emitCategoryTransactions();
              }
            });
        },
        error: (err) => {
          this.loaded = true;
          this.error = true;
        }
      });
  }

  emitCategoryTransactions(): void {
    this.categoryTransactions.sort((a, b) => {
      return Math.abs(b.transactionList.sum) - Math.abs(a.transactionList.sum);
    });

    let colorIndex = 0;
    for (let i = 0; i < this.categoryTransactions.length; i++) {
      const row = this.categoryTransactions[i];
      if (row.category) {
        if (row.category.color) {
          continue;
        }

        row.category.color = this.categoryColorPalette[colorIndex % this.categoryColorPalette.length];
        const categoryIndex = this.categories.findIndex(category => category.id === row.category?.id);
        if (categoryIndex >= 0) {
          this.categories[categoryIndex].color = row.category.color;
        }
        colorIndex++;
      }
    }

    this.$categoryTransactions.next(this.categoryTransactions);
    this.$categories.next(this.categories);
  }

  loadTransactions(): void {
    this.transactionService.list(this.startDate, this.endDate)
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe({
        next: (list) => {
          this.transactions = list.items;
          this.$transactions.next(this.transactions);
        }
      });
  }

  importFile(event: any): void {
    const file = event.target?.files[0];
    if (!file) {
      return;
    }
    this.transactionService.uploadCSV(file)
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe({
        next: () => {
          this.loadData();
        }
      });
  }

  moveMonth(count: number): void {
    const newMonth = moment().set({
      'month': this.monthFormControl.value,
      'year': this.yearFormControl.value
    }).add(count, 'month');
    this.monthFormControl.setValue(newMonth.month());
    this.yearFormControl.setValue(newMonth.year());
    this.loadData();
  }

  handleTransactionChange(transaction: Transaction): void {
    this.loadCategories();
  }
}
