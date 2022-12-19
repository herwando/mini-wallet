# Mini Wallet

#### Prerequisite

- Git
- Golang 1.18
- [PostgreSQL](https://www.postgresql.org/download/)

#### Setup

- Install Git
  See [Git Installation](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

- Install Golang

- Install Docker
  See [docker-compose](https://docs.docker.com/compose/install/

- Setup PostgreSQL with docker-compose. Adjust the configuration as needed, by default PostgreSQL will be exposed via port 15432. Also adjust PostgreSQL env as needed in [database.env](https://github.com/herwando/mini-wallet/blob/main/dev/database.env). After that, run docker-compose up.
  ```sh
  cd dev && docker-compose up
  ```

- Clone this repository
  ```sh
  git clone git@github.com:herwando/mini-wallet.git
  cd mini-wallet
  ```

- Install dependencies
  ```sh
  make dep
  ```

- Copy env.sample and if necessary, modify the env value(s)
  ```sh
  cp env.sample .env
  ```

- Download database migration tools
  ```sh
  make tool-migrate
  ```

- Run migration for each module. 
  ```sh
  make migrate-up
  ```

- Run golang main on cmd
  ```sh
  make run
  ```

#### API Documentaion

- POST Initialize my account for wallet
```sh
curl --location --request POST 'http://localhost/api/v1/init' \
--form 'customer_xid="ea0212d3-abd6-406f-8c67-868e814a2436"'
```

- POST Enable my wallet
```sh
curl --location --request POST 'http://localhost/api/v1/wallet' \
--header 'Authorization: Token 6b3f7dc70abe8aed3e56658b86fa508b472bf238'
```

- GET View my wallet balance
```sh
curl --location --request GET 'http://localhost/api/v1/wallet' \
--header 'Authorization: Token 6b3f7dc70abe8aed3e56658b86fa508b472bf238'
```

- POST Add virtual money to my wallet
```sh
curl --location --request POST 'http://localhost/api/v1/wallet/deposits' \
--header 'Authorization: Token 6b3f7dc70abe8aed3e56658b86fa508b472bf238' \
--form 'amount="100000"' \
--form 'reference_id="50535246-dcb2-4929-8cc9-004ea06f5241"'
```

- POST Use virtual money from my wallet
```sh
curl --location --request POST 'http://localhost/api/v1/wallet/withdrawals' \
--header 'Authorization: Token 6b3f7dc70abe8aed3e56658b86fa508b472bf238' \
--form 'amount="60000"' \
--form 'reference_id="4b01c9bb-3acd-47dc-87db-d9ac483d20b2"'
```

- PATCH Disable my wallet
```sh
curl --location --request PATCH 'http://localhost/api/v1/wallet' \
--header 'Authorization: Token 6b3f7dc70abe8aed3e56658b86fa508b472bf238' \
--form 'is_disabled="true"'
```