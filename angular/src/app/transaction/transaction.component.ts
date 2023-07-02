import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Transaction } from '../transaction/transaction';
import { Observable, Subject, takeUntil } from 'rxjs';
import { Category } from '../category/category';
import { TransactionService } from './transaction.service';
import { FormControl } from '@angular/forms';

@Component({
  selector: 'app-transaction',
  templateUrl: './transaction.component.html'
})
export class TransactionComponent {

  @Input() transaction: Transaction;
  @Input('categories') $categories: Observable<Category[]>;

  @Output() transactionChanged: EventEmitter<Transaction> = new EventEmitter();

  categories: Category[];
  editMode = false;
  categoryFormControl = new FormControl(0);

  private ngUnsubscribe: Subject<any> = new Subject();

  constructor(private transactionService: TransactionService) {}

  ngOnInit(): void {
    this.$categories?.
      pipe(takeUntil(this.ngUnsubscribe)).
      subscribe((data) => {
        this.categories = data;
      });

    this.categoryFormControl.setValue(this.transaction.category?.id ?? 0);
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

  setCategory(): void {
    const category = this.categories?.find((category) => category.id == this.categoryFormControl.value);
    this.transaction.category = category;
    this.transactionService.setCategory(this.transaction.id, category?.id ?? null)
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe({
        next: () => {
          this.editMode = false;
          this.transactionChanged.emit(this.transaction);
        }
      });
  }

  hide(hide: boolean): void {
    this.transactionService.hide(this.transaction.id, hide)
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe({
        next: () => {
          this.transaction.hidden = hide;
          this.transactionChanged.emit(this.transaction);
        }
      });
  }
}
