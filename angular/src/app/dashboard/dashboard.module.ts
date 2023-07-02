import { NgModule } from '@angular/core';
import { DashboardRoutingModule } from './dashboard-routing.module';
import { DashboardComponent } from './dashboard.component';
import { SharedModule } from '../shared/shared.module';
import { CategoryListComponent } from '../category-list/category-list.component';
import { TransactionListComponent } from '../transaction-list/transaction-list.component';
import { BalanceOverviewComponent } from '../balance-overview/balance-overview.component';
import { TransactionComponent } from '../transaction/transaction.component';

@NgModule({
  imports: [
    SharedModule.forRoot(),
    DashboardRoutingModule,
  ],
  declarations: [
    DashboardComponent,
    CategoryListComponent,
    TransactionListComponent,
    TransactionComponent,
    BalanceOverviewComponent,
  ],
  exports: [
    DashboardComponent,
  ],
})
export class DashboardModule {}
