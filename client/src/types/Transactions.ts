import type { Account } from "./Accounts";
import type { Category } from "./Category";
import type { Payee } from "./Payee";

export type Transaction = {
  ID: Number;
  CreatedAt: Date;
  UpdatedAt: Date;
  DeletedAt: Date | null;
  Date: Date;
  AccountID: Number;
  Account: Account;
  CategoryID: Number;
  Category: Category;
  Amount: Number;
  PayeeID: Number;
  Payee: Payee;
};
