scalar JSON
scalar DateTime
scalar UUID

# ========== TYPES ==========
type User {
  id: UUID!
  email: String!
  name: String
  role: String!
  isActive: Boolean!
  createdAt: DateTime!
  updatedAt: DateTime!
}

type Workflow {
  id: UUID!
  name: String!
  description: String
  createdBy: UUID!
  isActive: Boolean!
  settings: JSON
  version: Int!
  isTemplate: Boolean!
  forkedFrom: UUID
  createdAt: DateTime!
  updatedAt: DateTime!

  nodes: [Node!]!
  connections: [Connection!]!
  triggers: [Trigger!]!
  executions(limit: Int = 10, offset: Int = 0): [Execution!]!
  variables: [WorkflowVariable!]!
  tags: [Tag!]!
  owner: User!
}

type Node {
  id: UUID!
  workflowId: UUID!
  type: String!
  name: String
  positionX: Float!
  positionY: Float!
  parameters: JSON!
  credentialsId: UUID
  createdAt: DateTime!
  updatedAt: DateTime!
}

type Connection {
  id: UUID!
  workflowId: UUID!
  fromNodeId: UUID!
  toNodeId: UUID!
  outputIndex: Int!
  inputIndex: Int!
}

type Credential {
  id: UUID!
  name: String!
  type: String!
  ownerId: UUID!
  createdAt: DateTime!
  # DO NOT expose data in queries; use a mutation to reveal masked or decrypted info if authorized
}

type Execution {
  id: UUID!
  workflowId: UUID!
  status: String! # success, error, running, cancelled
  startedAt: DateTime!
  endedAt: DateTime
  errorMessage: String
  executionData: JSON!    # snapshot of node-level statuses, metadata
  runByUserId: UUID
}

type Trigger {
  id: UUID!
  workflowId: UUID!
  type: String!  # 'cron' | 'webhook'
  config: JSON!
  isEnabled: Boolean!
}

type Webhook {
  id: UUID!
  triggerId: UUID!
  path: String!
  method: String
  secret: String
}

type WorkflowVariable {
  id: UUID!
  workflowId: UUID!
  key: String!
  value: String!
  isSecret: Boolean!
}

type WorkflowVersion {
  id: UUID!
  workflowId: UUID!
  version: Int!
  definition: JSON!
  createdAt: DateTime!
}

type Tag { id: UUID!, name: String! }

type AuditLog {
  id: UUID!
  userId: UUID
  action: String!
  entityType: String!
  entityId: UUID!
  metadata: JSON
  timestamp: DateTime!
}

# ========== INPUTS ==========
input WorkflowInput {
  name: String!
  description: String
  isActive: Boolean
  settings: JSON
  isTemplate: Boolean
  forkedFrom: UUID
}

input NodeInput {
  id: UUID
  type: String!
  name: String
  positionX: Float!
  positionY: Float!
  parameters: JSON!
  credentialsId: UUID
}

input ConnectionInput {
  fromNodeId: UUID!
  toNodeId: UUID!
  outputIndex: Int
  inputIndex: Int
}

input StartExecutionInput {
  workflowId: UUID!
  runByUserId: UUID
  # optional: initial payload for the run
  payload: JSON
}

# ========== QUERIES ==========
type Query {
  me: User
  users(limit: Int = 100, offset: Int = 0): [User!]!

  workflow(id: UUID!): Workflow
  workflows(limit: Int = 100, offset: Int = 0): [Workflow!]!

  nodes(workflowId: UUID!): [Node!]!
  node(id: UUID!): Node

  connections(workflowId: UUID!): [Connection!]!

  executions(workflowId: UUID!, limit: Int = 20, offset: Int = 0): [Execution!]!
  execution(id: UUID!): Execution

  triggers(workflowId: UUID!): [Trigger!]!
  webhooks(triggerId: UUID!): [Webhook!]!

  credentials(ownerId: UUID!): [Credential!]!
  credential(id: UUID!): Credential
}

# ========== MUTATIONS ==========
type Mutation {
  # Users
  createUser(email: String!, password: String!, name: String, role: String = "user"): User!
  updateUser(id: UUID!, name: String, role: String, isActive: Boolean): User!
  deleteUser(id: UUID!): Boolean!

  # Workflows & composition
  createWorkflow(input: WorkflowInput!): Workflow!
  updateWorkflow(id: UUID!, input: WorkflowInput!): Workflow!
  deleteWorkflow(id: UUID!): Boolean!

  addNode(workflowId: UUID!, input: NodeInput!): Node!
  updateNode(id: UUID!, input: NodeInput!): Node!
  deleteNode(id: UUID!): Boolean!

  addConnection(workflowId: UUID!, input: ConnectionInput!): Connection!
  deleteConnection(id: UUID!): Boolean!

  createOrUpdateCredential(name: String!, type: String!, data: JSON!, ownerId: UUID!): Credential!
  deleteCredential(id: UUID!): Boolean!

  createTrigger(workflowId: UUID!, type: String!, config: JSON!): Trigger!
  updateTrigger(id: UUID!, config: JSON!, isEnabled: Boolean): Trigger!
  deleteTrigger(id: UUID!): Boolean!

  createWebhook(triggerId: UUID!, path: String!, method: String, secret: String): Webhook!
  rotateWebhookSecret(webhookId: UUID!): Webhook!

  setWorkflowVariable(workflowId: UUID!, key: String!, value: String!, isSecret: Boolean = false): WorkflowVariable!
  deleteWorkflowVariable(id: UUID!): Boolean!

  # Execution control
  startExecution(input: StartExecutionInput!): Execution!    # enqueues and returns execution row
  cancelExecution(id: UUID!): Boolean!

  # Versioning
  snapshotWorkflowVersion(workflowId: UUID!): WorkflowVersion!

  # Audit helpers (internal)
  appendAuditLog(userId: UUID, action: String!, entityType: String!, entityId: UUID!, metadata: JSON): AuditLog!
}

# ========== SUBSCRIPTIONS (optional Stage-0 extra) ==========
type Subscription {
  executionUpdated(executionId: UUID!): Execution
  nodeLogAdded(executionId: UUID!, nodeId: UUID!): LogEntry
}


go install github.com/99designs/gqlgen@latest
go run github.com/99designs/gqlgen generate

docker-compose up --build


Backend: http://localhost:8080

Frontend: http://localhost:5173

RabbitMQ UI: http://localhost:15672

Postgres: localhost:5432



------------------------


/project
  /db
    schema.sql          <-- your PostgreSQL schema (users, workflows, ai_agents, etc.)
    queries/            <-- sqlc SQL files
      users.sql
      workflows.sql
      executions.sql
      ai_tasks.sql
  /graph
    schema.graphql      <-- GraphQL schema (from previous message)
    resolver.go         <-- gqlgen resolvers
    models.go           <-- gqlgen types (auto-generated)
  /cmd/server
    main.go             <-- start HTTP server with gqlgen
  go.mod
