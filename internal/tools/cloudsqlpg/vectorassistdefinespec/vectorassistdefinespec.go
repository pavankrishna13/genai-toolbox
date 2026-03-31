// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vectorassistdefinespec

import (
	"context"
	"fmt"
	"net/http"

	yaml "github.com/goccy/go-yaml"
	"github.com/googleapis/genai-toolbox/internal/embeddingmodels"
	"github.com/googleapis/genai-toolbox/internal/sources"
	"github.com/googleapis/genai-toolbox/internal/tools"
	"github.com/googleapis/genai-toolbox/internal/util"
	"github.com/googleapis/genai-toolbox/internal/util/parameters"
	"github.com/jackc/pgx/v5/pgxpool"
)

const resourceType string = "vector-assist-define-spec"

const defineSpecQuery = `
		SELECT recommendation_id, vector_spec_id, table_name, schema_name, query, recommendation, applied, modified, created_at 
		FROM vector_assist.define_spec(table_name => $1::TEXT, schema_name => $2::TEXT, spec_id => $3::TEXT, 
			vector_column_name => $4::TEXT, text_column_name => $5::TEXT, 
			vector_index_type => $6::TEXT, embeddings_available => $7::BOOLEAN, 
			num_vectors => $8::INTEGER, dimensionality => $9::INTEGER, 
			embedding_model => $10::TEXT, prefilter_column_names => $11, 
			distance_func => $12::TEXT, quantization => $13::TEXT, 
			memory_budget_kb => $14::INTEGER, target_recall => $15::FLOAT, 
			target_top_k =>$16::INTEGER, tune_vector_index =>$17::BOOLEAN);
`

func init() {
	if !tools.Register(resourceType, newConfig) {
		panic(fmt.Sprintf("tool type %q already registered", resourceType))
	}
}

func newConfig(ctx context.Context, name string, decoder *yaml.Decoder) (tools.ToolConfig, error) {
	actual := Config{Name: name}
	if err := decoder.DecodeContext(ctx, &actual); err != nil {
		return nil, err
	}
	return actual, nil
}

type compatibleSource interface {
	PostgresPool() *pgxpool.Pool
	RunSQL(context.Context, string, []any) (any, error)
}

type Config struct {
	Name         string   `yaml:"name" validate:"required"`
	Type         string   `yaml:"type" validate:"required"`
	Source       string   `yaml:"source" validate:"required"`
	Description  string   `yaml:"description"`
	AuthRequired []string `yaml:"authRequired"`
}

var _ tools.ToolConfig = Config{}

func (cfg Config) ToolConfigType() string {
	return resourceType
}

func (cfg Config) Initialize(srcs map[string]sources.Source) (tools.Tool, error) {
	allParameters := parameters.Parameters{
		parameters.NewStringParameterWithRequired("table_name", "Table name on which vector workload needs to be set up.", true),
		parameters.NewStringParameterWithRequired("schema_name", "Schema containing the given table.", false),
		parameters.NewStringParameterWithRequired("spec_id", "Unique ID for the vector spec. Auto-generated, if not specified.", false),
		parameters.NewStringParameterWithRequired("vector_column_name", "Column name for the column with vector embeddings.", false),
		parameters.NewStringParameterWithRequired("text_column_name", "Column name for the column with text on which vector search needs to be set up.", false),
		parameters.NewStringParameterWithRequired("vector_index_type", "Type of the vector index to be created (Allowed inputs: 'hnsw', 'ivfflat', 'scann').", false),
		parameters.NewBooleanParameterWithRequired("embeddings_available", "Boolean parameter to know if vector embeddings are already available in the table.", false),
		parameters.NewIntParameterWithRequired("num_vectors", "Number of vectors expected in the dataset.", false),
		parameters.NewIntParameterWithRequired("dimensionality", "If vectors are already generated, set to dimension of vectors. If not, set to dimensionality of the embedding_model.", false),
		parameters.NewStringParameterWithRequired("embedding_model", "Optional parameter: Model to be used for generating embeddings.", false),
		parameters.NewArrayParameterWithRequired("prefilter_column_names", "Columns based on which prefiltering will happen in vector search queries.", false, parameters.NewStringParameter("prefilter_column_name", "Pre filter column name")),
		parameters.NewStringParameterWithRequired("distance_func", "Distance function to be used for comparing vectors (Allowed inputs: 'cosine', 'ip', 'l2', 'l1').", false),
		parameters.NewStringParameterWithRequired("quantization", "Quantization to be used for creating the vector indexes (Allowed inputs: 'none', 'halfvec', 'bit').", false),
		parameters.NewIntParameterWithRequired("memory_budget_kb", "Maximum size in KB that the index can consume in memory while building.", false),
		parameters.NewFloatParameterWithRequired("target_recall", "The recall that the user would like to target with the given index for standard vector queries.", false),
		parameters.NewIntParameterWithRequired("target_top_k", "The top-K values that need to be retrieved for the given query.", false),
		parameters.NewBooleanParameterWithRequired("tune_vector_index", "Boolean parameter to specify if the auto tuning is required for the index.", false),
	}
	paramManifest := allParameters.Manifest()

	if cfg.Description == "" {
		cfg.Description = "This tool defines a new vector specification by capturing the user's intent and requirements for a vector search workload. This generates a complete, ordered set of SQL recommendations required to set up the database, embeddings, and vector indexes. Use this tool at the very beginning of the vector setup process when a user first wants to configure a table for vector search, generate embeddings, or create a new vector index."
	}

	mcpManifest := tools.GetMcpManifest(cfg.Name, cfg.Description, cfg.AuthRequired, allParameters, nil)

	return Tool{
		Config:    cfg,
		allParams: allParameters,
		manifest: tools.Manifest{
			Description:  cfg.Description,
			Parameters:   paramManifest,
			AuthRequired: cfg.AuthRequired,
		},
		mcpManifest: mcpManifest,
	}, nil
}

var _ tools.Tool = Tool{}

type Tool struct {
	Config
	allParams   parameters.Parameters `yaml:"allParams"`
	manifest    tools.Manifest
	mcpManifest tools.McpManifest
}

func (t Tool) ToConfig() tools.ToolConfig {
	return t.Config
}

func (t Tool) Invoke(ctx context.Context, resourceMgr tools.SourceProvider, params parameters.ParamValues, accessToken tools.AccessToken) (any, util.ToolboxError) {
	source, err := tools.GetCompatibleSource[compatibleSource](resourceMgr, t.Source, t.Name, t.Type)
	if err != nil {
		return nil, util.NewClientServerError("source used is not compatible with the tool", http.StatusInternalServerError, err)
	}
	paramsMap := params.AsMap()

	newParams, err := parameters.GetParams(t.allParams, paramsMap)
	if err != nil {
		return nil, util.NewAgentError("unable to extract standard params", err)
	}
	sliceParams := newParams.AsSlice()
	resp, err := source.RunSQL(ctx, defineSpecQuery, sliceParams)
	if err != nil {
		return nil, util.ProcessGeneralError(err)
	}
	return resp, nil
}

func (t Tool) EmbedParams(ctx context.Context, paramValues parameters.ParamValues, embeddingModelsMap map[string]embeddingmodels.EmbeddingModel) (parameters.ParamValues, error) {
	return parameters.EmbedParams(ctx, t.allParams, paramValues, embeddingModelsMap, nil)
}

func (t Tool) Manifest() tools.Manifest {
	return t.manifest
}

func (t Tool) McpManifest() tools.McpManifest {
	return t.mcpManifest
}

func (t Tool) Authorized(verifiedAuthServices []string) bool {
	return tools.IsAuthorized(t.AuthRequired, verifiedAuthServices)
}

func (t Tool) RequiresClientAuthorization(resourceMgr tools.SourceProvider) (bool, error) {
	return false, nil
}

func (t Tool) GetAuthTokenHeaderName(resourceMgr tools.SourceProvider) (string, error) {
	return "Authorization", nil
}

func (t Tool) GetParameters() parameters.Parameters {
	return t.allParams
}
