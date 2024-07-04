import { useStore } from "@nanostores/react";
import {
  transactions as ts,
  getPaginatedTransactions,
} from "../stores/TransactionStore";
import { useCallback, useEffect } from "react";
export const TransactionsTable = () => {
  const $transactions = useStore(ts);

  const getTransactions = useCallback(async () => {
    const res = await getPaginatedTransactions();
    ts.set(res);
  }, []);

  useEffect(() => {
    getTransactions();
  });

  const deleteTransaction = async (id: Number) => {
    await fetch(
      `http://localhost:${
        import.meta.env.PUBLIC_SERVER_PORT
      }/transactions/delete/${id}`,
      {
        method: "POST",
      }
    );

    const res = await getPaginatedTransactions();
    ts.set(res);
  };
  return (
    <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
      {$transactions.length > 0 && (
        <table className="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
          <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
            <tr>
              <th scope="col" className="p-4">
                <div className="flex items-center">
                  <input
                    id="checkbox-all-search"
                    type="checkbox"
                    className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:focus:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                  />
                  <label htmlFor="checkbox-all-search" className="sr-only">
                    checkbox
                  </label>
                </div>
              </th>
              <th scope="col" className="px-6 py-3">
                Date
              </th>
              <th scope="col" className="px-6 py-3">
                Payee
              </th>
              <th scope="col" className="px-6 py-3">
                Category
              </th>
              <th scope="col" className="px-6 py-3">
                Amount
              </th>
              <th scope="col" className="px-6 py-3">
                Action
              </th>
            </tr>
          </thead>
          <tbody>
            {$transactions.map((transaction) => (
              <tr
                key={transaction.ID.toString()}
                className="bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
              >
                <td className="w-4 p-4">
                  <div className="flex items-center">
                    <input
                      id="checkbox-table-search-1"
                      type="checkbox"
                      className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:focus:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                    />
                    <label
                      htmlFor="checkbox-table-search-1"
                      className="sr-only"
                    >
                      checkbox
                    </label>
                  </div>
                </td>
                <th
                  scope="row"
                  className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white"
                >
                  {new Date(transaction.Date).toLocaleDateString()}
                </th>
                <td className="px-6 py-4">{transaction.Payee.Name}</td>
                <td className="px-6 py-4">{transaction.Category.Name}</td>
                <td className="px-6 py-4">${transaction.Amount.toString()}</td>
                <td className="flex items-center px-6 py-4">
                  <button
                    type="button"
                    className="font-medium text-blue-600 dark:text-blue-500 hover:underline"
                  >
                    Edit
                  </button>
                  <button
                    type="button"
                    className="font-medium text-red-600 dark:text-red-500 hover:underline ms-3"
                    onClick={() => deleteTransaction(transaction.ID)}
                  >
                    Remove
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
};
