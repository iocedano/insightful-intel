package infra

import "context"

type contextKey string

const executionIDKey contextKey = "executionID"
const pipelineIDKey contextKey = "pipelineID"
const stepIDKey contextKey = "stepID"

func GetStepID(ctx context.Context) (string, bool) {
	stepID, ok := ctx.Value(stepIDKey).(string)
	return stepID, ok
}

func SetStepID(ctx context.Context, stepID string) context.Context {
	return context.WithValue(ctx, stepIDKey, stepID)
}

func GetExecutionID(ctx context.Context) (string, bool) {
	executionID, ok := ctx.Value(executionIDKey).(string)
	return executionID, ok
}

func SetExecutionID(ctx context.Context, executionID string) context.Context {
	return context.WithValue(ctx, executionIDKey, executionID)
}

func GetPipelineID(ctx context.Context) (string, bool) {
	pipelineID, ok := ctx.Value(pipelineIDKey).(string)
	return pipelineID, ok
}

func SetPipelineID(ctx context.Context, pipelineID string) context.Context {
	return context.WithValue(ctx, pipelineIDKey, pipelineID)
}
