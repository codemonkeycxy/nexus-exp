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
2025/05/03 19:12:13 INFO  No logger configured for temporal client. Created default one.
2025/05/03 19:12:13 Started workflow WorkflowID nexus_greet_caller_workflow_20250503191213 RunID 01969910-bec6-78d5-965f-2f0106ab98dd
2025/05/03 19:12:13 Workflow result: Hello, World
2025/05/03 19:12:13 Started workflow WorkflowID nexus_greet_caller_workflow_20250503191213 RunID 01969910-bed8-774f-9a7f-a7efa6c49791
2025/05/03 19:12:18 Workflow result: After 5.10463s, the sloth Snugglemuffin responded: Hello, World
After 5.10463s, the sloth Snugglemuffin responded: Hello, World
After 5.10463s, the sloth Snugglemuffin responded: Hello, World
2025/05/03 19:12:18 Started workflow WorkflowID nexus_greet_caller_workflow_20250503191218 RunID 01969910-d2d0-7785-a676-4a2586eb2cea
2025/05/03 19:12:23 Workflow result: After 5.109416s, the sloth Snugglemuffin responded: Hello, World
After 5.113452s, the sloth Lazeberry responded: Hello, World
After 5.11534s, the sloth Mochapaws responded: Hello, World
```

### Regenerate proto code

```
buf generate
```
The proto gen library is a little broken right now so you will see some compile errors in the generated file. I'm sure
you will figure it out though :P
