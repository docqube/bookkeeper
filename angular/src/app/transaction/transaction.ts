import { Category } from "../category/category";


export interface Transaction {
  id: number;
  bookingDate: Date;
  valueDate: Date;
  recipient?: string;
  bookingText: string;
  purpose?: string;
  balance: number;
  amount: number;
  category?: Category;
}

export interface TransactionRequest {
  bookingDate?: Date;
  valueDate?: Date;
  recipient?: string;
  bookingText?: string;
  purpose?: string;
  balance?: number;
  amount?: number;
  categoryID?: number;
}

export interface CategoryTransactions {
  category?: Category;
  transactions: Transaction[];
  transactionsSum: number;
  transactionsLoaded: boolean;
}
