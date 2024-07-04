# go-budget

A self host-able go powered budgeting tool

## Local Development

### Backend

Create a `budget.yaml` file inside of the `data` directory that utilizes the following example format:

```yaml
categories:
  - name: Mortgage
    due: 1
  - name: Internet
    due: 1
  - name: Utilities
    due: 1
  - name: Amex
    due: 18
accounts:
  - name: Capital One Checking
    type: checking
  - name: Capital One Savings
    type: savings
```

Run `go run server/main.go` to start the server

### Frontend

`cd client` and `yarn install` to install dependencies
`yarn dev` to start local web server

### Notes

In order to work with both client and server running you will need to run each in their own terminal windows. Then set the ENV variable `PUBLIC_SERVER_PORT=8000` in order for the front end to proxy to the correct localhost port
