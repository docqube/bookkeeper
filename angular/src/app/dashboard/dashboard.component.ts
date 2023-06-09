import { Component } from '@angular/core';
import * as moment from 'moment';
import { Subject, takeUntil } from 'rxjs';
import { CategoryTransactions, Transaction } from '../transaction/transaction';
import { CategoryService } from '../category/category.service';
import { TransactionService } from '../transaction/transaction.service';
import { FormControl } from '@angular/forms';
import { Category, CategoryColors } from '../category/category';

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

  categoryTransactionsSubject: Subject<CategoryTransactions[]> = new Subject();
  transactionsSubject: Subject<Transaction[]> = new Subject();
  categoriesSubject: Subject<Category[]> = new Subject();

  startDateFormControl: FormControl = new FormControl<string>(moment().startOf('month').format('YYYY-MM-DD'));
  endDateFormControl: FormControl = new FormControl<string>(moment().endOf('month').format('YYYY-MM-DD'));

  private categoryColorPalette: string[] = CategoryColors;

  private ngUnsubscribe: Subject<any> = new Subject();

  constructor(private categoryService: CategoryService,
              private transactionService: TransactionService) {}

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

    this.loadCategories();
    this.loadTransactions();
  }

  loadCategories(): void {
    this.categoryTransactions = [];
    this.categoryService.list()
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe({
        next: (categories) => {
          this.categories = categories;
          for (const category of categories) {
            this.categoryTransactions.push({
              category: category,
              transactions: [],
              transactionsSum: 0,
              transactionsLoaded: true
            });
            this.loadTransactionsForCategory(category.id);
          }
        },
        error: (err) => {
          this.loaded = true;
          this.error = true;
        }
      });
  }

  loadTransactionsForCategory(categoryID: number): void {
    this.transactionService.listForCategoryID(this.startDateFormControl.value, this.endDateFormControl.value, categoryID)
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe({
        next: (transactions) => {
          if (!transactions) {
            return;
          }
          const categoryRowData = this.categoryTransactions.find(row => row.category?.id === categoryID);
          if (!categoryRowData) {
            return;
          }
          categoryRowData.transactionsSum = transactions.reduce((sum, transaction) => sum + transaction.amount, 0);
          categoryRowData.transactionsLoaded = true;

          if (this.categoryTransactions.filter(row => row.transactionsLoaded).length == this.categoryTransactions.length) {
            this.loadUnclassifiedTransactions();
          }
        }
      });
  }

  loadUnclassifiedTransactions(): void {
    this.transactionService.listUnclassified(this.startDateFormControl.value, this.endDateFormControl.value)
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe({
        next: (transactions) => {
          if (!transactions) {
            return;
          }
          let categoryRowData = this.categoryTransactions.find(row => !row.category);
          if (!categoryRowData) {
            categoryRowData = {
              category: undefined,
              transactions: [],
              transactionsSum: 0,
              transactionsLoaded: true
            };
            this.categoryTransactions.push(categoryRowData);
          }
          categoryRowData.transactionsSum = transactions.reduce((sum, transaction) => sum + transaction.amount, 0);
          categoryRowData.transactionsLoaded = true;

          this.emitCategoryTransactions();
        }
      });
  }

  emitCategoryTransactions(): void {
    this.categoryTransactions.sort((a, b) => {
      return Math.abs(b.transactionsSum) - Math.abs(a.transactionsSum);
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

    this.categoryTransactionsSubject.next(this.categoryTransactions);
    this.categoriesSubject.next(this.categories);
  }

  loadTransactions(): void {
    this.transactionService.list(this.startDateFormControl.value, this.endDateFormControl.value)
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe({
        next: (transactions) => {
          this.transactions = transactions;
          this.transactionsSubject.next(this.transactions);
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
    this.startDateFormControl.setValue(moment(this.startDateFormControl.value).add(count, 'month').startOf('month').format('YYYY-MM-DD'));
    this.endDateFormControl.setValue(moment(this.endDateFormControl.value).add(count, 'month').endOf('month').format('YYYY-MM-DD'));
    this.loadData();
  }
}
