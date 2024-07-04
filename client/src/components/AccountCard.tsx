import type { Account } from "../types/Accounts";

export const AccountCard = ({ account }: { account: Account }) => {
  return (
    <div className="py-6 px-12 bg-gray-200 rounded-md">
      <h2 className="font-bold text-2xl text-gray-700">{account.Name}</h2>
      <span className="font-black text-3xl text-green-700">
        ${account.Balance.toString()}
      </span>
      <div className="flex flex-row w-1/2">
        <ul className="flex flex-row justify-around">
          <li>
            <li>
              <button
                type="button"
                className="font-medium text-blue-600 dark:text-blue-500 hover:underline mr-3"
              >
                Edit
              </button>
            </li>
            <button
              type="button"
              className="font-medium text-red-600 dark:text-red-500 hover:underline"
            >
              Remove
            </button>
          </li>
        </ul>
      </div>
    </div>
  );
};
