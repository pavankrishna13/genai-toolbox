---
title: "vector-assist-generate-query Tool"
type: docs
weight: 1
description: >
  The "vector-assist-generate-query" tool produces optimized SQL queries for
  vector search, leveraging metadata and specifications to enable semantic
  and similarity searches.
---

## About

The `vector-assist-generate-query` tool generates optimized SQL queries for vector search by leveraging the metadata and vector specifications defined in a specific spec_id. It serves as the primary actionable tool for generating the executable SQL required to retrieve relevant results based on vector similarity.

The tool contextually understands requirements such as distance functions, quantization, and filtering to ensure the resulting query is compatible with the corresponding vector index. Additionally, it can automatically handle iterative index scans for filtered queries and calculate the necessary search parameters (like ef_search) to meet a target recall.
## Compatible Sources

{{< compatible-sources >}}

## Requirements

{{< notice tip >}} 
Ensure that your target PostgreSQL database has the required `vector_assist` extension installed, in order for this tool to execute successfully.
{{< /notice >}}

## Example

```yaml
kind: tool
name: generate_query
type: vector-assist-generate-query
source: my-database-source
description: "This tool generates optimized SQL queries for vector search by leveraging the metadata and vector specifications defined in a specific spec_id. It may return a single query or a sequence of multiple SQL queries that can be executed sequentially. Use this tool when a user wants to perform semantic or similarity searches on their data. It serves as the primary actionable tool to invoke for generating the executable SQL required to retrieve relevant results based on vector similarity."
```

## Reference

| **field**   | **type** | **required** | **description**                                      |
|-------------|:--------:|:------------:|------------------------------------------------------|
| type        |  string  |     true     | Must be "vector-assist-generate-query".                 |
| source      |  string  |     true     | Name of the source the SQL should execute on.        |
| description |  string  |    false     | Description of the tool that is passed to the agent. |