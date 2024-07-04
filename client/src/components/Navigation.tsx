export const Navigation = ({ pathname }: { pathname: string }) => {
  let active = "transactions";
  switch (pathname) {
    case "/accounts":
      active = "accounts";
      break;
    case "/categories":
      active = "categories";
      break;
    case "/":
    case "/transactions":
    default:
      active = "transactions";
      break;
  }
  return (
    <div>
      <h1 className="text-gray-50 text-4xl font-bold">Go Budget!</h1>
      <div className="flex flex-row my-6">
        <ul className="flex flex-row w-1/2 justify-between text-gray-50">
          <li
            className={
              active === "transactions" ? "text-blue-500" : "text-gray-50"
            }
          >
            <a href="/">Transactions</a>
          </li>
          <li
            className={
              active === "categories" ? "text-blue-500" : "text-gray-50"
            }
          >
            <a href="categories">Categories</a>
          </li>
          <li
            className={active === "accounts" ? "text-blue-500" : "text-gray-50"}
          >
            <a href="accounts">Accounts</a>
          </li>
        </ul>
      </div>
    </div>
  );
};
