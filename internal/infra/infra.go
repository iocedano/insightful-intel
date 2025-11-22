package infra

import "context"

type contextKey string

const ExecutionIDKey contextKey = "execution_id"
const PipelineIDKey contextKey = "pipeline_id"
const StepIDKey contextKey = "step_id"

func GetExecutionID(ctx context.Context) (string, bool) {
	executionID, ok := ctx.Value(ExecutionIDKey).(string)
	if !ok {
		return "", false
	}
	return executionID, true
}

func SetExecutionID(ctx context.Context, executionID string) context.Context {
	return context.WithValue(ctx, ExecutionIDKey, executionID)
}

func GetPipelineID(ctx context.Context) (string, bool) {
	pipelineID, ok := ctx.Value(PipelineIDKey).(string)
	return pipelineID, ok
}

func SetPipelineID(ctx context.Context, pipelineID string) context.Context {
	return context.WithValue(ctx, PipelineIDKey, pipelineID)
}
