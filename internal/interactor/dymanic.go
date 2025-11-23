package interactor

import (
	"context"
	"insightful-intel/internal/domain"
	"insightful-intel/internal/module"
	"insightful-intel/internal/repositories"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type DynamicPipelineInteractor struct {
	repositories *repositories.RepositoryFactory
}

func NewDynamicPipelineInteractor(
	repositoryFactory *repositories.RepositoryFactory,
) *DynamicPipelineInteractor {
	return &DynamicPipelineInteractor{
		repositories: repositoryFactory,
	}
}

func (d *DynamicPipelineInteractor) ExecuteDynamicPipeline(ctx context.Context, query string, maxDepth int, skipDuplicates bool) (*domain.DynamicPipelineResult, error) {
	spew.Dump("ExecuteDynamicPipeline")
	spew.Dump(ctx.Err() == context.Canceled)
	spew.Dump(ctx.Err() == context.DeadlineExceeded)
	spew.Dump(ctx.Err())

	// Create a channel to receive pipeline steps
	stepChan := make(chan domain.DynamicPipelineStep, 100)
	done := make(chan bool)

	// Configure the dynamic pipeline
	config := domain.DynamicPipelineConfig{
		MaxDepth:           maxDepth,
		MaxConcurrentSteps: 10,
		DelayBetweenSteps:  2,
		SkipDuplicates:     skipDuplicates,
	}

	// Available domains
	availableDomains := domain.AllDomainTypes()

	// Start pipeline execution in a goroutine
	go func() {
		defer close(stepChan)
		defer close(done)

		// Execute the dynamic pipeline with step callback
		dynamicResult, err := d.executeDynamicPipelineWithCallback(ctx, query, availableDomains, config, stepChan)
		if err != nil {
			// Send error as a step
			errorStep := domain.DynamicPipelineStep{
				DomainType:      "ERROR",
				SearchParameter: query,
				Success:         false,
				Error:           err,
				Output:          nil,
				Depth:           0,
			}
			stepChan <- errorStep
			return
		}

		// Send final summary
		summaryStep := domain.DynamicPipelineStep{
			DomainType:      "SUMMARY",
			SearchParameter: query,
			Success:         true,
			Error:           nil,
			Output: map[string]interface{}{
				"total_steps":       dynamicResult.TotalSteps,
				"successful_steps":  dynamicResult.SuccessfulSteps,
				"failed_steps":      dynamicResult.FailedSteps,
				"max_depth_reached": dynamicResult.MaxDepthReached,
			},
			Depth: dynamicResult.MaxDepthReached,
		}
		stepChan <- summaryStep
	}()

	// Stream the steps as they come
	stepCount := 0
	for {
		select {
		case step, ok := <-stepChan:
			if !ok {
				return nil, nil
			}

			stepCount++

			spew.Dump("step", step)

		case <-ctx.Done():
			// Client disconnected
			return nil, nil
		}
	}

	return nil, nil
}

// executeDynamicPipelineWithCallback executes the dynamic pipeline and sends steps to a channel
func (s *DynamicPipelineInteractor) executeDynamicPipelineWithCallback(ctx context.Context, query string, availableDomains []domain.DomainType, config domain.DynamicPipelineConfig, stepChan chan<- domain.DynamicPipelineStep) (*domain.DynamicPipelineResult, error) {
	// Create a custom pipeline executor that streams steps
	return s.executeStreamingPipeline(ctx, query, availableDomains, config, stepChan)
}

// executeStreamingPipeline executes the pipeline with real-time streaming
func (d *DynamicPipelineInteractor) executeStreamingPipeline(ctx context.Context, query string, availableDomains []domain.DomainType, config domain.DynamicPipelineConfig, stepChan chan<- domain.DynamicPipelineStep) (*domain.DynamicPipelineResult, error) {
	// Create the initial pipeline steps
	createdPipelineResult, err := module.CreateDynamicPipeline(ctx, query, availableDomains, config)
	if err != nil {
		return nil, err
	}

	_, err = d.repositories.GetPipelineRepository().CreateDynamicPipelineResult(ctx, createdPipelineResult)
	if err != nil {
		return nil, err
	}

	// Get initial steps from the result
	initialSteps := createdPipelineResult.Steps

	totalSteps := 0
	successfulSteps := 0
	failedSteps := 0
	maxDepthReached := 0

	// Track searched keywords per domain to avoid duplicates
	searchedKeywordsPerDomain := make(map[domain.DomainType]map[string]bool)
	for _, domainType := range availableDomains {
		searchedKeywordsPerDomain[domainType] = make(map[string]bool)
	}

	// Process steps with streaming
	processedSteps := make([]domain.DynamicPipelineStep, 0)

	// Create a queue for steps to process
	stepQueue := make([]domain.DynamicPipelineStep, len(initialSteps))
	copy(stepQueue, initialSteps)

	for len(stepQueue) > 0 {
		// Get next step from queue
		step := stepQueue[0]
		stepQueue = stepQueue[1:]

		step.PipelineID = createdPipelineResult.ID

		// Send step start event
		startStep := step
		startStep.Success = false
		startStep.Output = nil
		stepChan <- startStep

		// Execute the step
		result, err := module.SearchDomain(step.DomainType, domain.DomainSearchParams{Query: step.SearchParameter})

		// Update step with results
		step.Success = err == nil
		step.Error = err
		if result != nil {
			step.Output = result.Output
			step.KeywordsPerCategory = result.KeywordsPerCategory
		}

		err = d.repositories.GetPipelineRepository().CreateDynamicPipelineStep(ctx, &step)
		if err != nil {
			log.Println("Error creating pipeline step ----> ", err)
			return nil, err
		}

		result.PipelineStepsID = step.ID

		created, err := d.repositories.GetPipelineRepository().CreateDomainSearchResult(ctx, result)
		if err != nil {
			log.Println("Error creating pipeline CreateDomainSearchResult ----> ", err)
			return nil, err
		}

		switch step.DomainType {
		case domain.DomainTypeONAPI:
			entities, ok := created.Output.([]domain.Entity)
			if !ok {
				log.Println("Error casting result output to []domain.Entity ----> ", err)
				return nil, err
			}
			for _, entity := range entities {
				entity.DomainSearchResultID = created.ID
				if err := d.repositories.GetOnapiRepository().Create(ctx, entity); err != nil {
					log.Println("Error creating onapi repository ----> ", err)
					return nil, err
				}
			}
		case domain.DomainTypeSCJ:
			cases, ok := created.Output.([]domain.ScjCase)
			if !ok {
				log.Println("Error casting result output to []domain.ScjCase ", err)
				return nil, err
			}
			for _, c := range cases {
				c.DomainSearchResultID = created.ID
				if err := d.repositories.GetScjRepository().Create(ctx, c); err != nil {
					log.Println("Error creating scj repository ", err, c)
					return nil, err
				}
			}
		case domain.DomainTypeDGII:
			results, ok := created.Output.([]domain.Register)
			if !ok {
				log.Println("Error casting result output to []domain.Register ", err)
				return nil, err
			}
			for _, result := range results {
				result.DomainSearchResultID = created.ID
				if err := d.repositories.GetDgiiRepository().Create(ctx, result); err != nil {
					log.Println("Error creating dgii repository ", err)
					return nil, err
				}
			}
		case domain.DomainTypePGR:
			results, ok := created.Output.([]domain.PGRNews)
			if !ok {
				log.Println("Error casting result output to []domain.PGRNews ", err)
				return nil, err
			}
			for _, result := range results {
				result.DomainSearchResultID = created.ID
				if err := d.repositories.GetPgrRepository().Create(ctx, result); err != nil {
					log.Println("Error creating pgr repository ", err)
					return nil, err
				}
			}
		case domain.DomainTypeGoogleDocking, domain.DomainTypeSocialMedia, domain.DomainTypeFileType, domain.DomainTypeXSocialMedia:
			results, ok := created.Output.([]domain.GoogleDockingResult)
			if !ok {
				log.Println("Error casting result output to []domain.GoogleDockingResult ", err)
				return nil, err
			}
			for _, result := range results {
				result.DomainSearchResultID = created.ID
				if err := d.repositories.GetDockingRepository().Create(ctx, result); err != nil {
					log.Println("Error creating docking repository ", err)
					return nil, err
				}
			}
		}

		// Update counters
		totalSteps++
		if step.Success {
			successfulSteps++
		} else {
			failedSteps++
		}

		if step.Depth > maxDepthReached {
			maxDepthReached = step.Depth
		}

		// Send completed step
		stepChan <- step
		processedSteps = append(processedSteps, step)

		// Add delay between steps for better streaming experience
		time.Sleep(time.Duration(config.DelayBetweenSteps) * time.Second)

		// Generate new steps from keywords if not at max depth
		if step.Depth < config.MaxDepth && step.Success && step.Output != nil {
			newSteps := d.generateNextSteps(step, availableDomains, searchedKeywordsPerDomain, config)
			stepQueue = append(stepQueue, newSteps...)
		}
	}

	// Create final result
	createdPipelineResult.TotalSteps = totalSteps
	createdPipelineResult.SuccessfulSteps = successfulSteps
	createdPipelineResult.FailedSteps = failedSteps
	createdPipelineResult.MaxDepthReached = maxDepthReached
	createdPipelineResult.Config = config

	err = d.repositories.GetPipelineRepository().UpdateDynamicPipelineResult(ctx, createdPipelineResult)
	if err != nil {
		log.Println("Error creating pipeline result ----> ", err)
		return nil, err
	}

	return createdPipelineResult, nil
}

func (*DynamicPipelineInteractor) generateNextSteps(completedStep domain.DynamicPipelineStep, availableDomains []domain.DomainType, searchedKeywordsPerDomain map[domain.DomainType]map[string]bool, config domain.DynamicPipelineConfig) []domain.DynamicPipelineStep {
	var newSteps []domain.DynamicPipelineStep

	// Extract keywords from the completed step
	keywordsPerCategory := completedStep.KeywordsPerCategory
	if keywordsPerCategory == nil {
		return newSteps
	}

	// Generate new steps for each keyword category
	for category, keywords := range keywordsPerCategory {
		for _, keyword := range keywords {
			// Skip if already searched or if keyword is too short
			if len(keyword) < 3 {
				continue
			}

			// Generate steps for each available domain
			for _, domainType := range availableDomains {
				// Skip if already searched this keyword for this domain
				if searchedKeywordsPerDomain[domainType][keyword] {
					continue
				}

				// Skip if same domain as current step
				if domainType == completedStep.DomainType {
					continue
				}

				// Mark as searched
				searchedKeywordsPerDomain[domainType][keyword] = true

				// Create new step
				newStep := domain.DynamicPipelineStep{
					DomainType:          domainType,
					SearchParameter:     keyword,
					Category:            category,
					Keywords:            []string{keyword},
					Success:             false,
					Error:               nil,
					Output:              nil,
					KeywordsPerCategory: nil,
					Depth:               completedStep.Depth + 1,
				}

				newSteps = append(newSteps, newStep)
			}
		}
	}

	return newSteps
}
