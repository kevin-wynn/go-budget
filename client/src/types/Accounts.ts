export type Account = {
  CreatedAt: Date;
  UpdatedAt: Date;
  DeletedAt: Date | null;
  ID: Number;
  Name: String;
  Type: "savings" | "checking" | "credit card";
  Balance: Number;
};
