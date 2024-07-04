import type { Account } from "../types/Accounts";
import { AccountCard } from "./AccountCard";

export const Accounts = ({ accounts }: { accounts: Account[] }) => {
  return (
    <div className="grid gap-4 grid-cols-3">
      {accounts.map((account) => (
        <AccountCard account={account} key={account.ID.toString()} />
      ))}
    </div>
  );
};
