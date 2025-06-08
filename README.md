# README

## Overview
- A Go-based backend providing a TaskService over gRPC, backed by MySQL.
- `Server`: exposes TodoService RPCs (AddTask, ListTasks, CompleteTask) on port 50051.
- `Client`: a Cobra-powered CLI that can start the server and invoke those RPCs.
- `Auth`: simple Bearer-token interceptor; set AUTH_TOKEN in .env and pass --token on the client.


## Prerequisites
- **Go 1.23+**  
- **Docker** & **Docker Compose**  
  - [Install Docker](https://docs.docker.com/get-docker/)  
  - [Install Docker Compose](https://docs.docker.com/compose/install/)  
- **Ginkgo** & **mockgen** (for unit tests) in your `$PATH`  
  ```bash
   go install github.com/onsi/ginkgo/ginkgo@latest
   go install github.com/golang/mock/mockgen@latest
   export PATH="$(go env GOPATH)/bin:$PATH"
   ```

## Architecture
   ### Backend
   - Language: Go 1.23
   - Dependency injection: Uber FX
   - Logging: Zap

   ### Repository layer:
   - raw database/sql + github.com/go-sql-driver/mysql

   ### Service layer: 
   - TaskService interface → business logic

   ### Transport: 
   - gRPC server with a unary interceptor for token auth

   ### Database
   - MySQL containerized via Docker Compose

   ### Schema 
   - files in ./schema (e.g. 01_create_tasks_table.sql)

   - Connection configured via `.env`

## Configuration
   - Create a `.env` configuration file.
   ```bash
      touch .env
   ```
   - An example of a configuration file is
   ```bash
      # .env
      MYSQL_ROOT_PASSWORD=rootpassword
      MYSQL_HOST=mysql
      MYSQL_PORT=3306
      MYSQL_USER=user
      MYSQL_PASSWORD=userpassword
      MYSQL_DATABASE=project_db
      GRPC_PORT=50051
      HTTP_PORT=8000
      AUTH_TOKEN=test_auth
   ```

## Running the application
   
   ### Set the envionment variables
   ```bash
      set -o allexport
      source .env
      set +o allexport
   ```

   ### Bring up everything
   ```bash
      docker compose up
   ```

   ### Starting the Server
   - Once MySQL is healthy, launch the gRPC server:
   ```bash
      docker compose exec todo todo server \
      --grpc-port 50051 \
      --mysql-host mysql \
      --mysql-port 3306 \
      --mysql-user "$MYSQL_USER" \
      --mysql-pass "$MYSQL_PASSWORD" \
      --mysql-db   "$MYSQL_DATABASE"
   ```
   - This will block and run the TaskService until you Ctrl+C.

   ### Using the CLI Client
   - In a separate shell (after the server is running), you can manage tasks:
   ```bash
      # Add a task
      docker compose exec todo todo client add \
      --host localhost --port 50051 \
      --token "$AUTH_TOKEN" \
      --title "Buy eggs" --desc "A dozen"

      # List all tasks
      docker compose exec todo todo client get \
      --host localhost --port 50051 \
      --token "$AUTH_TOKEN"

      # Mark task #1 complete
      docker compose exec todo todo client complete \
      --host localhost --port 50051 \
      --token "$AUTH_TOKEN" \
      --id 1
   ```

   ### Running Unit Tests
   - Mocks live under `pkg/repository/mock_repository` and `pkg/service/mock_service`. Regenerate them if you change interfaces.
   ```bash
      go generate ./pkg/repository
      go generate ./pkg/service
   ```
   - Run the ginkgo tests
   ```bash
      ginkgo -r pkg/service
      ginkgo -r pkg/transport/grpc
   ```

   ### Inspecting MySQL
   ```bash
      docker exec -it hearx-mysql-1 bash
   ```

   ```bash
      mysql -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE"
   ```

   ### Cleanup
   ```bash
      docker compose down
   ```
   - `down` stops containers (data persists)
   - `down -v` also removes volumes (data wiped)