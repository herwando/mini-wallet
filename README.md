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

