# nexus-exp
Experiment with Nexus + Temporal Workflow + Proto Code Gen

## Getting started locally

### Get `temporal` CLI to enable local development

1. Follow the instructions on the [docs
   site](https://learn.temporal.io/getting_started/go/dev_environment/#set-up-a-local-temporal-service-for-development-with-temporal-cli)
   to install Temporal CLI.

> NOTE: The recommended version is at least v1.3.0.

### Spin up environment

#### Start temporal server

```
temporal server start-dev
```

### Initialize environment

In a separate terminal window

#### Create caller and target namespaces

```
temporal operator namespace create --namespace my-target-namespace
temporal operator namespace create --namespace my-caller-namespace
```

#### Create Nexus endpoint

```
temporal operator nexus endpoint create \
  --name my-nexus-endpoint-name \
  --target-namespace my-target-namespace \
  --target-task-queue my-handler-task-queue
```

### Make Nexus calls across namespace boundaries

In separate terminal windows:

### Nexus server

```
go run ./server/starter \
    -target-host localhost:7233 \
    -namespace my-target-namespace
```

### Nexus client worker

```
go run ./client/worker \
    -target-host localhost:7233 \
    -namespace my-caller-namespace
```

### Start client workflow

```
go run ./client/starter \
    -target-host localhost:7233 \
    -namespace my-caller-namespace
```

### Output

which should result in:
```
2025/04/18 17:05:15 INFO  No logger configured for temporal client. Created default one.
2025/04/18 17:05:15 Started workflow WorkflowID nexus_greet_caller_workflow_20250418170515 RunID 9f6f97bd-ffd5-4244-9c25-2b13225974f8
2025/04/18 17:05:15 Workflow result: Hello, World
2025/04/18 17:05:15 Started workflow WorkflowID nexus_greet_caller_workflow_20250418170515 RunID dd5f4b70-46a0-4050-b876-9e56b89fda7c
2025/04/18 17:05:20 Workflow result: After 5.021175s, the sloth responded: Hello, World
```
