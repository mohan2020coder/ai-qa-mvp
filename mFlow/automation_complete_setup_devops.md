# Complete Setup — DevOps-first Repo Skeleton

This document contains a complete, ready-to-copy repository scaffold for the n8n-like automation platform. It includes **development**, **CI/QA**, and **production (Terraform)** setups, plus scripts that let you bring each environment up with a single command.

---

## File tree (what's included)

```
repo-root/
├── README.md
├── .env.example
├── docker-compose.yml
├── docker-compose.ci.yml
├── devops/
│   ├── up.sh
│   ├── down.sh
│   ├── reset.sh
│   ├── wait_for_services.sh
│   ├── provision_and_start_ec2.sh
│   └── provision_and_deploy_prod.sh
├── infra/
│   ├── qa/
│   │   └── main.tf
│   └── prod/
│       └── main.tf
├── backend/
│   ├── api/            # Go gqlgen service (skeleton)
│   │   ├── go.mod
│   │   ├── main.go
│   │   ├── gqlgen.yml
│   │   └── graph/       # gqlgen generated/resolvers
│   └── worker/         # Go worker & orchestrator skeleton
│       ├── go.mod
│       └── main.go
├── migrations/
│   └── 000001_initial.sql
├── .github/
│   └── workflows/
│       ├── ci.yml
│       └── deploy.yml
└── Makefile
```

---

> **Important**: The repository files below use placeholders for secrets and AWS IDs. Replace `${...}` placeholders before running in production.

---

## README.md

```markdown
# Automation Platform — DevOps-first Scaffold

This repository provides a devops-first scaffold for the automation platform (React TS frontend, Go gqlgen API, Go workers/orchestrator, Postgres, RabbitMQ).

## Quick start (development)

1. Copy `.env.example` to `.env` and edit if you wish.
2. Start the local stack:

```bash
./devops/up.sh
```

3. API GraphQL endpoint: http://localhost:8080/graphql
4. RabbitMQ management: http://localhost:15672 (guest/guest)
5. Frontend dev: http://localhost:5173

## One-line EC2 bootstrap (dev/stage)

```bash
./devops/provision_and_start_ec2.sh
```

## Production provisioning (Terraform)

Edit `infra/prod/main.tf` with your AWS account values, then:

```bash
cd infra/prod
terraform init
terraform apply -auto-approve
```

See `devops/provision_and_deploy_prod.sh` to deploy docker images.

```
```

---

## .env.example

```env
# Database
DATABASE_URL=postgres://automation:automation@postgres:5432/automation?sslmode=disable

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/

# App secrets
JWT_SECRET=replace-with-secure-secret
KMS_KEY_ID=dev-kms

# App config
PORT=8080

# AWS (used by infra scripts)
AWS_REGION=us-east-1
AWS_SSH_KEY_NAME=your-keypair-name
EC2_INSTANCE_TYPE=t3.small
```

---

## docker-compose.yml (development)

```yaml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: automation
      POSTGRES_USER: automation
      POSTGRES_PASSWORD: automation
    ports:
      - '5432:5432'
    volumes:
      - pgdata:/var/lib/postgresql/data

  rabbitmq:
    image: rabbitmq:3.10-management
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
    ports:
      - '5672:5672'
      - '15672:15672'
    volumes:
      - rabbitmqdata:/var/lib/rabbitmq

  api:
    build: ./backend/api
    depends_on:
      - postgres
      - rabbitmq
    environment:
      DATABASE_URL: ${DATABASE_URL}
      RABBITMQ_URL: ${RABBITMQ_URL}
      JWT_SECRET: ${JWT_SECRET}
    ports:
      - '8080:8080'

  worker:
    build: ./backend/worker
    depends_on:
      - postgres
      - rabbitmq
    environment:
      DATABASE_URL: ${DATABASE_URL}
      RABBITMQ_URL: ${RABBITMQ_URL}

  frontend:
    image: node:20
    working_dir: /app
    volumes:
      - ./frontend:/app
    command: bash -c "npm install && npm run dev -- --host 0.0.0.0"
    ports:
      - '5173:5173'

  migrations:
    image: migrate/migrate
    depends_on:
      - postgres
    entrypoint: ["/bin/sh", "-c", "sleep 5; migrate -path /migrations -database ${DATABASE_URL} up || true"]
    volumes:
      - ./migrations:/migrations

volumes:
  pgdata:
  rabbitmqdata:
```

---

## docker-compose.ci.yml (CI / ephemeral)

```yaml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: automation
      POSTGRES_USER: automation
      POSTGRES_PASSWORD: automation
  rabbitmq:
    image: rabbitmq:3.10
  api:
    build: ./backend/api
    depends_on: [postgres, rabbitmq]
  worker:
    build: ./backend/worker
    depends_on: [postgres, rabbitmq]
```

---

## devops/up.sh

```bash
#!/usr/bin/env bash
set -e
ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
cd "$ROOT_DIR"

echo "Starting local dev stack..."

docker compose up --build -d

# wait for postgres
./devops/wait_for_services.sh

# run migrations (uses migrate tool container)
docker compose run --rm migrations || true

echo "Dev stack started. GraphQL: http://localhost:8080/graphql"
```

Make executable: `chmod +x devops/up.sh`.

---

## devops/down.sh

```bash
#!/usr/bin/env bash
set -e
ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
cd "$ROOT_DIR"

docker compose down
```

---

## devops/reset.sh

```bash
#!/usr/bin/env bash
set -e
ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
cd "$ROOT_DIR"

docker compose down -v --remove-orphans
rm -rf ./.tmp

# bring up fresh
./devops/up.sh
```

---

## devops/wait_for_services.sh

```bash
#!/usr/bin/env bash
set -e

function wait_for_port() {
  host=$1
  port=$2
  retries=30
  while ! nc -z $host $port; do
    retries=$((retries-1))
    if [ $retries -le 0 ]; then
      echo "Timeout waiting for $host:$port"
      exit 1
    fi
    sleep 1
  done
}

wait_for_port postgres 5432
wait_for_port rabbitmq 5672
echo "Services up"
```

---

## devops/provision_and_start_ec2.sh (single-EC2 bootstrap, dev/stage)

```bash
#!/usr/bin/env bash
set -e
ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
cd "$ROOT_DIR/infra/prod"

if ! command -v terraform >/dev/null 2>&1; then
  echo "terraform not installed"
  exit 1
fi

terraform init
terraform apply -auto-approve

# Terraform outputs should include public_ip and ssh_user
PUB_IP=$(terraform output -raw public_ip)
SSH_USER=ec2-user

echo "Waiting for instance to be reachable..."
while ! ssh -o StrictHostKeyChecking=no ${SSH_USER}@${PUB_IP} true 2>/dev/null; do
  sleep 3
done

# Copy repo and run start script remotely (or use S3 to stash artifacts)
scp -o StrictHostKeyChecking=no -r $ROOT_DIR ${SSH_USER}@${PUB_IP}:/opt/automation
ssh -o StrictHostKeyChecking=no ${SSH_USER}@${PUB_IP} 'cd /opt/automation && ./devops/up.sh'

echo "Provisioned and started stack on EC2: ${PUB_IP}"
```

> NOTE: This script expects `infra/prod/main.tf` to provision an EC2 with a public IP and `ec2-user` available via your SSH key in the default location. Customize to your org's setup.

---

## devops/provision_and_deploy_prod.sh (deploy new images to prod EC2 fleet - simple example)

```bash
#!/usr/bin/env bash
set -e
# Assumes docker images have been pushed to ECR and production EC2s are accessible via SSH
# Example usage: ./devops/provision_and_deploy_prod.sh prod-repo:tag

IMAGE_TAG=$1
if [ -z "$IMAGE_TAG" ]; then
  echo "Usage: $0 <image-tag>"
  exit 1
fi

# Replace with your list of instance IPs or use Terraform output
INSTANCES=("${PROD_HOSTS}")
for host in "${INSTANCES[@]}"; do
  ssh -o StrictHostKeyChecking=no ec2-user@${host} "docker pull ${IMAGE_TAG} && docker compose pull && docker compose up -d"
done
```

---

## infra/prod/main.tf (skeleton)

```hcl
provider "aws" {
  region = var.aws_region
}

variable "aws_region" { default = "us-east-1" }
variable "ssh_key_name" { default = "${AWS_SSH_KEY_NAME}" }

resource "aws_instance" "app" {
  ami           = data.aws_ami.amazon_linux.id
  instance_type = var.instance_type
  key_name      = var.ssh_key_name
  user_data     = file("../../devops/ec2_user_data.sh")
  tags = { Name = "automation-app" }
}

output "public_ip" {
  value = aws_instance.app.public_ip
}
```

> This is intentionally minimal — add VPC, SG, ALB, ASG, IAM roles, RDS, and Amazon MQ resources in your full production TF.

---

## infra/qa/main.tf (skeleton) — similar to prod but smaller

```hcl
# QA environment terraform skeleton. Mirror prod but smaller instance sizes and separate state.
```

---

## migrations/000001_initial.sql

```sql
-- Create schema (from user-provided SQL). Replace or extend as needed.
-- Enable UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name VARCHAR(255),
    role VARCHAR(50) DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Include additional tables (workflows, nodes, etc.) — paste your full schema here.
```

> **Note**: Replace the `-- Include additional tables` line with the full schema you provided earlier (users, workflows, nodes...)

---

## backend/api/main.go (skeleton)

```go
package main

import (
    "log"
    "net/http"
)

func main() {
    // Load config & env
    // Connect to Postgres (pgx)
    // Connect to RabbitMQ
    // Setup gqlgen GraphQL server

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })

    log.Println("starting api on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Add `gqlgen.yml` and schema at `backend/api/graph/schema.graphqls` — use `go run github.com/99designs/gqlgen@latest` to generate.

---

## backend/worker/main.go (skeleton)

```go
package main

import (
    "log"
)

func main() {
    // Connect to RabbitMQ
    // Consume node_tasks queue
    // Execute node logic (HTTP, script, etc.)
    log.Println("worker started")
    select {}
}
```

---

## .github/workflows/ci.yml

```yaml
name: CI
on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_DB: automation
          POSTGRES_USER: automation
          POSTGRES_PASSWORD: automation
        ports: [5432]
      rabbitmq:
        image: rabbitmq:3.10
        ports: [5672]
    steps:
      - uses: actions/checkout@v4
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Build
        run: |
          cd backend/api && go build ./...
          cd backend/worker && go build ./...
      - name: Run unit tests
        run: |
          cd backend/api && go test ./...
          cd backend/worker && go test ./...
      - name: Run integration (compose)
        run: |
          docker compose -f docker-compose.ci.yml up --build -d
          ./devops/wait_for_services.sh
          # Run basic integration checks here
          docker compose -f docker-compose.ci.yml down -v
```

---

## .github/workflows/deploy.yml (example)

```yaml
name: Deploy
on:
  push:
    branches: [ main ]

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build Docker images
        run: |
          docker build -t ${{ secrets.ECR_REGISTRY }}/automation-api:${{ github.sha }} ./backend/api
          docker build -t ${{ secrets.ECR_REGISTRY }}/automation-worker:${{ github.sha }} ./backend/worker
      - name: Login to ECR
        uses: aws-actions/amazon-ecr-login@v2
      - name: Push images to ECR
        run: |
          docker push ${{ secrets.ECR_REGISTRY }}/automation-api:${{ github.sha }}
          docker push ${{ secrets.ECR_REGISTRY }}/automation-worker:${{ github.sha }}
      - name: Trigger deploy
        run: |
          ssh -o StrictHostKeyChecking=no ec2-user@${{ secrets.PROD_HOST }} 'cd /opt/automation && ./devops/provision_and_deploy_prod.sh ${{ secrets.ECR_REGISTRY }}/automation-api:${{ github.sha }}'
```

---

## Makefile

```makefile
.PHONY: dev up down reset build test

up:
	./devops/up.sh

down:
	./devops/down.sh

reset:
	./devops/reset.sh

build:
	cd backend/api && go build -o bin/api ./...
	cd backend/worker && go build -o bin/worker ./...

test:
	cd backend/api && go test ./...
	cd backend/worker && go test ./...
```

---

## Next steps
1. Copy this document's file contents into files in your repo (or ask me to emit individual files for copy/paste).  
2. Populate `backend/api/graph/schema.graphqls` with the full GraphQL schema you provided and run `gqlgen generate`.  
3. Implement DB migrations fully (paste the full SQL schema into `migrations/000001_initial.sql`).  
4. Wire up actual service code (DB connections, gql resolvers, RabbitMQ publishers/consumers).  

---

If you'd like, I can now:
- write out each file separately into the canvas (as individual code files) so you can copy them directly, **or**
- generate a zip bundle and provide a download link (if you want the repo as a zip).

Tell me which format you prefer and I will produce it immediately.

