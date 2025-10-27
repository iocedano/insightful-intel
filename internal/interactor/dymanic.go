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
	pipelineResultRepo *repositories.PipelineRepository
	scjRepo            *repositories.ScjRepository
	dgiiRegisterRepo   *repositories.DgiiRepository
	pgrNewsRepo        *repositories.PgrRepository
	googleDockingRepo  *repositories.DockingRepository
	onapiRepo          *repositories.OnapiRepository
}

func NewDynamicPipelineInteractor(
	pipelineResultRepo *repositories.PipelineRepository,
	scjRepo *repositories.ScjRepository,
	dgiiRegisterRepo *repositories.DgiiRepository,
	pgrNewsRepo *repositories.PgrRepository,
	googleDockingRepo *repositories.DockingRepository,
	onapiRepo *repositories.OnapiRepository,
) *DynamicPipelineInteractor {
	return &DynamicPipelineInteractor{
		pipelineResultRepo: pipelineResultRepo,
		scjRepo:            scjRepo,
		dgiiRegisterRepo:   dgiiRegisterRepo,
		pgrNewsRepo:        pgrNewsRepo,
		googleDockingRepo:  googleDockingRepo,
		onapiRepo:          onapiRepo,
	}
}

func (d *DynamicPipelineInteractor) ExecuteDynamicPipeline(ctx context.Context, query string, maxDepth int, skipDuplicates bool) {
	// Create a channel to receive pipeline steps
	stepChan := make(chan module.DynamicPipelineStep, 100)
	done := make(chan bool)

	// Configure the dynamic pipeline
	config := module.DynamicPipelineConfig{
		MaxDepth:           maxDepth,
		MaxConcurrentSteps: 10,
		DelayBetweenSteps:  2,
		SkipDuplicates:     skipDuplicates,
	}

	// Available domains
	availableDomains := []domain.DomainType{
		domain.DomainTypeONAPI,
		domain.DomainTypeSCJ,
		domain.DomainTypeDGII,
		domain.DomainTypePGR,
		domain.DomainTypeGoogleDocking,
	}

	// Start pipeline execution in a goroutine
	go func() {
		defer close(stepChan)
		defer close(done)

		// Execute the dynamic pipeline with step callback
		dynamicResult, err := d.executeDynamicPipelineWithCallback(ctx, query, availableDomains, config, stepChan)
		if err != nil {
			// Send error as a step
			errorStep := module.DynamicPipelineStep{
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

		if err := d.pipelineResultRepo.Create(ctx, dynamicResult); err != nil {
			log.Println("Error creating  dynamicResult pipeline result", err)
			errorStep := module.DynamicPipelineStep{
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
		summaryStep := module.DynamicPipelineStep{
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
	// Flush the response to ensure immediate delivery
	// flusher, ok := w.(http.Flusher)
	// if ok {
	// 	flusher.Flush()
	// }
	// if !ok {
	// 	http.Error(w, "Streaming not supported", http.StatusInternalServerError)
	// 	return
	// }

	// Stream the steps as they come
	stepCount := 0
	for {
		select {
		case step, ok := <-stepChan:
			if !ok {
				// Channel closed, send completion event
				// s.writeSSEEvent(w, "complete", map[string]interface{}{
				// 	"message":     "Pipeline execution completed",
				// 	"total_steps": stepCount,
				// }, flusher)
				return
			}

			stepCount++

			// Convert step to ConnectorPipeline format
			pipelineStep := domain.DomainSearchResult{
				Success:             step.Success,
				Error:               step.Error,
				DomainType:          step.DomainType,
				SearchParameter:     step.SearchParameter,
				Output:              step.Output,
				KeywordsPerCategory: step.KeywordsPerCategory,
			}

			spew.Dump(pipelineStep)

			// // Send step as SSE event
			// eventData := map[string]interface{}{
			// 	"step_number": stepCount,
			// 	"step":        pipelineStep,
			// 	"depth":       step.Depth,
			// 	"category":    string(step.Category),
			// 	"keywords":    step.Keywords,
			// }

			// eventType := "step"
			// switch step.DomainType {
			// case "error":
			// 	eventType = "error"
			// case "SUMMARY":
			// 	eventType = "sumary"
			// }

			// s.writeSSEEvent(w, eventType, eventData, flusher)

		case <-ctx.Done():
			// Client disconnected
			return
		}
	}
}

// executeDynamicPipelineWithCallback executes the dynamic pipeline and sends steps to a channel
func (s *DynamicPipelineInteractor) executeDynamicPipelineWithCallback(ctx context.Context, query string, availableDomains []domain.DomainType, config module.DynamicPipelineConfig, stepChan chan<- module.DynamicPipelineStep) (*module.DynamicPipelineResult, error) {
	// Create a custom pipeline executor that streams steps
	return s.executeStreamingPipeline(ctx, query, availableDomains, config, stepChan)
}

// executeStreamingPipeline executes the pipeline with real-time streaming
func (d *DynamicPipelineInteractor) executeStreamingPipeline(ctx context.Context, query string, availableDomains []domain.DomainType, config module.DynamicPipelineConfig, stepChan chan<- module.DynamicPipelineStep) (*module.DynamicPipelineResult, error) {
	// Create the initial pipeline steps
	initialResult, err := module.CreateDynamicPipeline(query, availableDomains, config)
	if err != nil {
		return nil, err
	}

	// Get initial steps from the result
	initialSteps := initialResult.Steps

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
	processedSteps := make([]module.DynamicPipelineStep, 0)

	// Create a queue for steps to process
	stepQueue := make([]module.DynamicPipelineStep, len(initialSteps))
	copy(stepQueue, initialSteps)

	for len(stepQueue) > 0 {
		// Get next step from queue
		step := stepQueue[0]
		stepQueue = stepQueue[1:]

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

		// ---
		log.Println("Creating pipeline result")
		if err := d.pipelineResultRepo.Create(ctx, result); err != nil {
			log.Println("Error creating pipeline ---result", err)
			return nil, err
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
	dynamicResult := &module.DynamicPipelineResult{
		Steps:           processedSteps,
		TotalSteps:      totalSteps,
		SuccessfulSteps: successfulSteps,
		FailedSteps:     failedSteps,
		MaxDepthReached: maxDepthReached,
		Config:          config,
	}

	d.pipelineResultRepo.Create(ctx, dynamicResult)

	return dynamicResult, nil
}

func (*DynamicPipelineInteractor) generateNextSteps(completedStep module.DynamicPipelineStep, availableDomains []domain.DomainType, searchedKeywordsPerDomain map[domain.DomainType]map[string]bool, config module.DynamicPipelineConfig) []module.DynamicPipelineStep {
	var newSteps []module.DynamicPipelineStep

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
				newStep := module.DynamicPipelineStep{
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
