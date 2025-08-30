# This script will create a repository scaffold for the automation platform
# and also write a standalone Python script that reproduces the same files.
# It writes files under /mnt/data/automation_scaffold and creates a zip at /mnt/data/automation_scaffold.zip
#
# After this cell runs you will find:
#   - /mnt/data/automation_scaffold/  (the full scaffold)
#   - /mnt/data/create_automation_scaffold.py  (the standalone script you can run locally)
#   - /mnt/data/automation_scaffold.zip  (zipped scaffold)
#

import os, stat, zipfile, textwrap, json

ROOT = "automation_scaffold"

files = {
    "README.md": textwrap.dedent("""\
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
    """),

    ".env.example": textwrap.dedent("""\
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
    """),

    "docker-compose.yml": textwrap.dedent("""\
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
        command: bash -c \"npm install && npm run dev -- --host 0.0.0.0\"
        ports:
          - '5173:5173'

      migrations:
        image: migrate/migrate
        depends_on:
          - postgres
        entrypoint: [\"/bin/sh\", \"-c\", \"sleep 5; migrate -path /migrations -database ${DATABASE_URL} up || true\"]
        volumes:
          - ./migrations:/migrations

    volumes:
      pgdata:
      rabbitmqdata:
    """),

    "docker-compose.ci.yml": textwrap.dedent("""\
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
    """),

    "devops/up.sh": textwrap.dedent("""\
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
    """),

    "devops/down.sh": textwrap.dedent("""\
    #!/usr/bin/env bash
    set -e
    ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
    cd "$ROOT_DIR"

    docker compose down
    """),

    "devops/reset.sh": textwrap.dedent("""\
    #!/usr/bin/env bash
    set -e
    ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
    cd "$ROOT_DIR"

    docker compose down -v --remove-orphans
    rm -rf ./.tmp

    # bring up fresh
    ./devops/up.sh
    """),

    "devops/wait_for_services.sh": textwrap.dedent("""\
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
    """),

    "devops/provision_and_start_ec2.sh": textwrap.dedent("""\
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
    """),

    "devops/provision_and_deploy_prod.sh": textwrap.dedent("""\
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
    """),

    "devops/ec2_user_data.sh": textwrap.dedent("""\
    #!/bin/bash
    yum update -y
    amazon-linux-extras install docker -y
    service docker start
    usermod -a -G docker ec2-user
    mkdir -p /opt/automation
    # The repo should be copied or pulled here by external provisioning
    """),

    "infra/prod/main.tf": textwrap.dedent("""\
    provider "aws" {
      region = var.aws_region
    }

    variable "aws_region" { default = "us-east-1" }
    variable "ssh_key_name" { default = "" }

    data "aws_ami" "amazon_linux" {
      most_recent = true
      owners      = ["amazon"]
      filter {
        name   = "name"
        values = ["amzn2-ami-hvm-*-x86_64-gp2"]
      }
    }

    variable "instance_type" { default = "t3.small" }

    resource "aws_instance" "app" {
      ami           = data.aws_ami.amazon_linux.id
      instance_type = var.instance_type
      key_name      = var.ssh_key_name
      user_data     = file("${path.module}/../../devops/ec2_user_data.sh")
      tags = { Name = "automation-app" }
    }

    output "public_ip" {
      value = aws_instance.app.public_ip
    }
    """),

    "infra/qa/main.tf": textwrap.dedent("""\
    # QA environment terraform skeleton. Mirror prod but smaller instance sizes and separate state.
    """),

    "migrations/000001_initial.sql": textwrap.dedent("""\
    -- Create schema (from user-provided SQL). Replace or extend as needed.
    -- Enable UUID
   -- Enable UUID support
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    -- ============= USERS ====================
    CREATE TABLE users (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        email VARCHAR(255) UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        name VARCHAR(255),
        role VARCHAR(50) DEFAULT 'user', -- 'user', 'admin', 'owner'
        is_active BOOLEAN DEFAULT true,
        created_at TIMESTAMPTZ DEFAULT now(),
        updated_at TIMESTAMPTZ DEFAULT now()
    );

    -- ============= WORKFLOWS ====================
    CREATE TABLE workflows (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR(255) NOT NULL,
        description TEXT,
        created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        is_active BOOLEAN DEFAULT true,
        settings JSONB,
        version INT DEFAULT 1,
        is_template BOOLEAN DEFAULT false,
        forked_from UUID REFERENCES workflows(id) ON DELETE SET NULL,
        created_at TIMESTAMPTZ DEFAULT now(),
        updated_at TIMESTAMPTZ DEFAULT now()
    );

    CREATE INDEX idx_workflows_created_by ON workflows(created_by);

    -- ============= WORKFLOW SHARING ====================
    CREATE TABLE workflow_shares (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        access_level VARCHAR(50) NOT NULL CHECK (access_level IN ('view', 'edit', 'owner')),
        created_at TIMESTAMPTZ DEFAULT now(),
        UNIQUE (workflow_id, user_id)
    );

    -- ============= NODES ====================
    CREATE TABLE nodes (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
        type VARCHAR(100) NOT NULL,
        name VARCHAR(255),
        position_x FLOAT NOT NULL,
        position_y FLOAT NOT NULL,
        parameters JSONB NOT NULL,
        credentials_id UUID REFERENCES credentials(id),
        created_at TIMESTAMPTZ DEFAULT now(),
        updated_at TIMESTAMPTZ DEFAULT now()
    );

    CREATE INDEX idx_nodes_workflow ON nodes(workflow_id);

    -- ============= CONNECTIONS ====================
    CREATE TABLE connections (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
        from_node_id UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
        to_node_id UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
        output_index INT DEFAULT 0,
        input_index INT DEFAULT 0
    );

    -- ============= CREDENTIALS ====================
    CREATE TABLE credentials (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR(255) NOT NULL,
        type VARCHAR(100) NOT NULL,
        data BYTEA NOT NULL, -- Encrypted binary data
        iv BYTEA NOT NULL,   -- Initialization vector for encryption
        tag BYTEA,           -- Auth tag (for AES-GCM)
        owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        created_at TIMESTAMPTZ DEFAULT now()
    );

    CREATE INDEX idx_credentials_owner ON credentials(owner_id);

    -- ============= EXECUTIONS ====================
    CREATE TABLE executions (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
        status VARCHAR(20) NOT NULL CHECK (status IN ('success', 'error', 'running', 'cancelled')),
        started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        ended_at TIMESTAMPTZ,
        error_message TEXT,
        execution_data JSONB NOT NULL,
        run_by_user_id UUID REFERENCES users(id)
    );

    CREATE INDEX idx_executions_workflow ON executions(workflow_id);

    -- ============= TRIGGERS ====================
    CREATE TABLE triggers (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
        type VARCHAR(50) NOT NULL CHECK (type IN ('cron', 'webhook')),
        config JSONB NOT NULL,
        is_enabled BOOLEAN DEFAULT true
    );

    -- ============= WEBHOOKS ====================
    CREATE TABLE webhooks (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        trigger_id UUID NOT NULL REFERENCES triggers(id) ON DELETE CASCADE,
        path VARCHAR(255) UNIQUE NOT NULL,
        method VARCHAR(10) CHECK (method IN ('GET', 'POST', 'PUT', 'DELETE')),
        secret VARCHAR(255)
    );

    -- ============= VARIABLES (WORKFLOW-SCOPED) ====================
    CREATE TABLE workflow_variables (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
        key VARCHAR(100) NOT NULL,
        value TEXT NOT NULL,
        is_secret BOOLEAN DEFAULT false,
        UNIQUE (workflow_id, key)
    );

    -- ============= VERSIONS (Workflow Snapshots) ====================
    CREATE TABLE workflow_versions (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
        version INT NOT NULL,
        definition JSONB NOT NULL, -- Full snapshot of nodes/connections/settings
        created_at TIMESTAMPTZ DEFAULT now(),
        UNIQUE (workflow_id, version)
    );

    -- ============= TAGS ====================
    CREATE TABLE tags (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR(100) UNIQUE NOT NULL
    );

    CREATE TABLE workflow_tags (
        workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
        tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
        PRIMARY KEY (workflow_id, tag_id)
    );

    -- ============= AUDIT LOGS ====================
    CREATE TABLE audit_logs (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        user_id UUID REFERENCES users(id) ON DELETE SET NULL,
        action VARCHAR(100) NOT NULL,
        entity_type VARCHAR(50) NOT NULL,
        entity_id UUID NOT NULL,
        metadata JSONB,
        timestamp TIMESTAMPTZ DEFAULT now()
    );

    -- ============= INDEXING STRATEGIES ====================
    -- Helpful for filtering and performance
    CREATE INDEX idx_executions_status ON executions(status);
    CREATE INDEX idx_workflow_variables_key ON workflow_variables(key);
    CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);

    -- ============= SECURITY PLACEHOLDERS ====================
    -- NOTE: Actual encryption/decryption should be done in application layer.
    -- Optionally use PostgreSQL's pgcrypto for basic encryption.

    -- ============= OPTIONAL RLS (Row Level Security) ====================
    -- To enforce per-user access on workflows:
    -- ALTER TABLE workflows ENABLE ROW LEVEL SECURITY;
    -- CREATE POLICY user_own_workflows ON workflows
    --   USING (created_by = current_setting('app.current_user')::uuid);

    -- Then in app layer, set:
    -- SET app.current_user = 'user-id-here';


    """),

    "backend/api/main.go": textwrap.dedent("""\
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
    """),

    "backend/api/go.mod": textwrap.dedent("""\
    module automation-api

    go 1.21
    """),

    "backend/worker/main.go": textwrap.dedent("""\
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
    """),

    "backend/worker/go.mod": textwrap.dedent("""\
    module automation-worker

    go 1.21
    """),

    ".github/workflows/ci.yml": textwrap.dedent("""\
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
    """),

    ".github/workflows/deploy.yml": textwrap.dedent("""\
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
    """),

    "Makefile": textwrap.dedent("""\
    .PHONY: dev up down reset build test

    up:
    \t./devops/up.sh

    down:
    \t./devops/down.sh

    reset:
    \t./devops/reset.sh

    build:
    \tcd backend/api && go build -o bin/api ./...
    \tcd backend/worker && go build -o bin/worker ./...

    test:
    \tcd backend/api && go test ./...
    \tcd backend/worker && go test ./...
    """),

    # Minimal frontend placeholder
    "frontend/package.json": textwrap.dedent("""\
    {
      "name": "automation-frontend",
      "version": "0.0.1",
      "private": true,
      "scripts": {
        "dev": "vite"
      },
      "dependencies": {}
    }
    """),

    "frontend/index.html": "<!-- Frontend placeholder -->\n<html><body><h1>Automation Frontend</h1></body></html>\n",
}

# Create directories and files
created = []
for path, content in files.items():
    full_path = os.path.join(ROOT, path)
    dirpath = os.path.dirname(full_path)
    if not os.path.exists(dirpath):
        os.makedirs(dirpath, exist_ok=True)
    with open(full_path, "w", encoding="utf-8") as f:
        f.write(content)
    created.append(full_path)

# Make scripts executable where appropriate
exec_paths = [
    os.path.join(ROOT, "devops/up.sh"),
    os.path.join(ROOT, "devops/down.sh"),
    os.path.join(ROOT, "devops/reset.sh"),
    os.path.join(ROOT, "devops/wait_for_services.sh"),
    os.path.join(ROOT, "devops/provision_and_start_ec2.sh"),
    os.path.join(ROOT, "devops/provision_and_deploy_prod.sh"),
    os.path.join(ROOT, "devops/ec2_user_data.sh"),
]
for p in exec_paths:
    if os.path.exists(p):
        st = os.stat(p)
        os.chmod(p, st.st_mode | stat.S_IEXEC)

# Also create a standalone python script that reproduces the scaffold when run.
standalone_script = "create_automation_scaffold.py"
script_content = textwrap.dedent(f'''\
#!/usr/bin/env python3
# Standalone script to create the automation scaffold in the current directory.
# Run: python3 create_automation_scaffold.py
import os, stat, textwrap

files = {json.dumps(files)}

root = os.path.abspath("automation_scaffold")

for path, content in files.items():
    full = os.path.join(root, path)
    d = os.path.dirname(full)
    if not os.path.exists(d):
        os.makedirs(d, exist_ok=True)
    with open(full, "w", encoding="utf-8") as f:
        f.write(content)

# Make scripts executable
execs = [
    "devops/up.sh",
    "devops/down.sh",
    "devops/reset.sh",
    "devops/wait_for_services.sh",
    "devops/provision_and_start_ec2.sh",
    "devops/provision_and_deploy_prod.sh",
    "devops/ec2_user_data.sh",
]
for e in execs:
    p = os.path.join(root, e)
    if os.path.exists(p):
        st = os.stat(p)
        os.chmod(p, st.st_mode | stat.S_IEXEC)

print("Scaffold created at", root)
''')

with open(standalone_script, "w", encoding="utf-8") as f:
    f.write(script_content)

os.chmod(standalone_script, os.stat(standalone_script).st_mode | stat.S_IEXEC)
created.append(standalone_script)

# Create a zip of the scaffold
zip_path = "automation_scaffold.zip"
with zipfile.ZipFile(zip_path, "w", zipfile.ZIP_DEFLATED) as z:
    for folder, _, filenames in os.walk(ROOT):
        for filename in filenames:
            filepath = os.path.join(folder, filename)
            arcname = os.path.relpath(filepath, ROOT)
            z.write(filepath, arcname)

print("Created files:")
for p in created:
    print(" -", p)
print("\nZipped scaffold at:", zip_path)
print("Standalone script at:", standalone_script)
