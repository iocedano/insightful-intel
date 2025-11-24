package domain

// DynamicPipelineConfig holds configuration for the dynamic pipeline
type DynamicPipelineConfig struct {
	Query              string       `json:"query"`
	MaxDepth           int          `json:"max_depth"`
	MaxConcurrentSteps int          `json:"max_concurrent_steps"`
	DelayBetweenSteps  int          `json:"delay_between_steps"`
	SkipDuplicates     bool         `json:"skip_duplicates"`
	AvailableDomains   []DomainType `json:"available_domains"`
}

// DynamicPipelineStep represents a single step in the pipeline
type DynamicPipelineStep struct {
	ID                  ID                           `json:"id"`
	PipelineID          ID                           `json:"pipeline_id"`
	DomainType          DomainType                   `json:"domain_type"`
	SearchParameter     string                       `json:"search_parameter"`
	Category            KeywordCategory              `json:"category"`
	Keywords            []string                     `json:"keywords"`
	Success             bool                         `json:"success"`
	Error               error                        `json:"error"`
	Output              any                          `json:"output"`
	KeywordsPerCategory map[KeywordCategory][]string `json:"keywords_per_category"`
	Depth               int                          `json:"depth"`
}

// DynamicPipelineResult represents the complete pipeline result
type DynamicPipelineResult struct {
	ID              ID                    `json:"id"`
	Steps           []DynamicPipelineStep `json:"steps"`
	TotalSteps      int                   `json:"total_steps"`
	SuccessfulSteps int                   `json:"successful_steps"`
	FailedSteps     int                   `json:"failed_steps"`
	MaxDepthReached int                   `json:"max_depth_reached"`
	Config          DynamicPipelineConfig `json:"config"`
}
