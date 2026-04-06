---
title: "vector-assist-modify-spec Tool"
type: docs
weight: 1
description: >
  The "vector-assist-modify-spec" tool modifies an existing vector specification
  with new parameters or overrides, recalculating the generated SQL
  recommendations to match the updated requirements.
---

## About

The `vector-assist-modify-spec` tool modifies an existing vector specification (identified by a required `spec_id`) with new parameters or overrides. Upon modification, it automatically recalculates and refreshes the list of generated recommendations by `vector_assist.define-spec` to match the updated spec requirements. 

Use this tool when a user or agent wants to adjust or fine-tune the configuration of an already defined vector spec (such as changing the target recall, embedding model, or quantization) before actually executing the setup commands. Under the hood, this tool connects to the target database and executes the `vector_assist.modify_spec` function to generate the updated specifications.

## Compatible Sources

{{< compatible-sources >}}

## Requirements

{{< notice tip >}} 
Ensure that your target PostgreSQL database has the required `vector_assist` extension installed, and that the `vector_assist.modify_spec` function is available in order for this tool to execute successfully.
{{< /notice >}}

## Example

```yaml
kind: tool
name: modify_spec
type: vector-assist-modify-spec
source: my-database-source
description: "This tool modifies an existing vector specification (identified by a required spec_id) with new parameters or overrides. Upon modification, it automatically recalculates and refreshes the list of generated SQL recommendations to match the updated requirements. Use this tool when a user wants to adjust or fine-tune the configuration of an already defined vector spec (such as changing the target recall, embedding model, or quantization) before actually executing the setup commands."
```

## Reference

| **field**   | **type** | **required** | **description**                                      |
|-------------|:--------:|:------------:|------------------------------------------------------|
| type        |  string  |     true     | Must be "vector-assist-modify-spec".                 |
| source      |  string  |     true     | Name of the source the SQL should execute on.        |
| description |  string  |    false     | Description of the tool that is passed to the agent. |