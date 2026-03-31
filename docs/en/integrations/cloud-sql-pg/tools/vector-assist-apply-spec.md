---
title: "vector-assist-apply-spec"
type: docs
weight: 1
description: >
  The "vector-assist-apply-spec" tool automatically executes all SQL recommendations
  associated with a specific vector specification or table to finalize the
  vector search setup.
---

## About

The `vector-assist-apply-spec` tool automatically executes all the SQL recommendations associated with a specific vector specification (spec_id) or table. It runs the necessary commands in the correct sequence to provision the workload, marking each step as applied once successful. 

Use this tool when the user has reviewed the generated recommendations from a defined (or modified) spec and is ready to apply the changes directly to their database instance to finalize the vector search setup. Under the hood, this tool connects to the target database and executes the `vector_assist.apply_spec` function.

## Compatible Sources

{{< compatible-sources >}}

## Requirements

{{< notice tip >}} 
Ensure that your target PostgreSQL database has the required `vector_assist` extension installed, in order for this tool to execute successfully.
{{< /notice >}}

## Example

```yaml
kind: tool
name: apply_spec
type: vector-assist-apply-spec
source: my-database-source
description: "This tool automatically executes all the SQL recommendations associated with a specific vector specification (spec_id) or table. It runs the necessary commands in the correct sequence to provision the workload, marking each step as applied once successful. Use this tool when the user has reviewed the generated recommendations from a defined (or modified) spec and is ready to apply the changes directly to their database instance to finalize the vector search setup."
```

## Reference

| **field**   | **type** | **required** | **description**                                      |
|-------------|:--------:|:------------:|------------------------------------------------------|
| type        |  string  |     true     | Must be "vector-assist-apply-spec".                 |
| source      |  string  |     true     | Name of the source the SQL should execute on.        |
| description |  string  |    false     | Description of the tool that is passed to the agent. |