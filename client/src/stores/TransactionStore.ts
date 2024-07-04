import { map } from "nanostores";
import type { Transaction } from "../types/Transactions";

export const getPaginatedTransactions = async (page = 1) => {
  const transactionsReq = await fetch(
    `http://localhost:${
      import.meta.env.PUBLIC_SERVER_PORT
    }/transactions/${page}`
  );
  return (await transactionsReq.json()) as Transaction[];
};

export const transactions = map([] as Transaction[]);
