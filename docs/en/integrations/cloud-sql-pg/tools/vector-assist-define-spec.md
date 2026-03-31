---
title: "vector-assist-define-spec"
type: docs
weight: 1
description: >
  The "vector-assist-define-spec" tool defines a new vector specification by
  capturing the user's intent and requirements for a vector search workload,
  generating SQL recommendations for setting up database, embeddings, and
  vector indexes.
---

## About

The `vector-assist-define-spec` tool defines a new vector specification by capturing the user's intent and requirements for a vector search workload. It generates a complete, ordered set of SQL recommendations required to set up the database, embeddings, and vector indexes. 

Use this tool at the very beginning of the vector setup process when an agent or user first wants to configure a table for vector search, generate embeddings, or create a new vector index. Under the hood, this tool connects to the target database and executes the `vector_assist.define_spec` function to generate the necessary specifications.

## Compatible Sources

{{< compatible-sources >}}

## Requirements

{{< notice tip >}} 
Ensure that your target PostgreSQL database has the required `vector_assist` extension installed, in order for this tool to execute successfully.
{{< /notice >}}

## Example

```yaml
kind: tool
name: define_spec
type: vector-assist-define-spec
source: my-database-source
description: "This tool defines a new vector specification by capturing the user's intent and requirements for a vector search workload. This generates a complete, ordered set of SQL recommendations required to set up the database, embeddings, and vector indexes. Use this tool at the very beginning of the vector setup process when a user first wants to configure a table for vector search, generate embeddings, or create a new vector index."
```

## Reference

| **field**   | **type** | **required** | **description**                                      |
|-------------|:--------:|:------------:|------------------------------------------------------|
| type        |  string  |     true     | Must be "vector-assist-define-spec".                 |
| source      |  string  |     true     | Name of the source the SQL should execute on.        |
| description |  string  |    false     | Description of the tool that is passed to the agent. |