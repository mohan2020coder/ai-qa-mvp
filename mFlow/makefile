PROJECT_NAME=mFlow

.PHONY: bootstrap up down logs backend frontend

bootstrap:
	@mkdir -p backend/internal/app \
		backend/internal/auth \
		backend/internal/workflow/model \
		backend/internal/workflow/repo \
		backend/internal/workflow/service \
		backend/internal/execution/engine \
		backend/internal/execution/repo \
		backend/internal/node/http \
		backend/internal/node/db \
		backend/internal/node/mq \
		backend/internal/node/custom \
		backend/internal/user \
		backend/internal/graphql \
		backend/internal/infra \
		backend/internal/utils
	@mkdir -p backend/cmd/server
	# Backend init
	@cd backend && go mod init github.com/$(PROJECT_NAME)/backend && go get github.com/99designs/gqlgen
	@echo "package main\n\nimport \"fmt\"\n\nfunc main(){ fmt.Println(\"Backend running...\") }" > backend/cmd/server/main.go
	@cd backend && go run github.com/99designs/gqlgen init
	# Backend Dockerfile
	@echo 'FROM golang:1.22-alpine\nWORKDIR /app\nCOPY . .\nRUN go mod tidy\nRUN go build -o server ./cmd/server\nCMD ["./server"]' > backend/Dockerfile
	# Frontend init (requires Node 22)
	@npm create vite@latest frontend -- --template react-ts
	@cd frontend && npm install react-flow-renderer @apollo/client graphql
	# Frontend Dockerfile
	@echo 'FROM node:22-alpine\nWORKDIR /app\nCOPY package*.json ./\nRUN npm install\nCOPY . .\nCMD ["npm","run","dev","--","--host","0.0.0.0","--port","5173"]' > frontend/Dockerfile
	@echo "✅ Project bootstrap complete."

up:
	docker-compose up --build -d

down:
	docker-compose down

logs:
	docker-compose logs -f

backend:
	cd backend && go run cmd/server/main.go

frontend:
	cd frontend && npm run dev


gen-graphql:
	cd backend && go run github.com/99designs/gqlgen generate

gen-graphql-init:
	cd backend && go run github.com/99designs/gqlgen init