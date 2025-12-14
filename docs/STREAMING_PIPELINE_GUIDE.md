# Dynamic Pipeline Streaming Guide

This guide explains how to use the enhanced `dynamicPipelineHandler` that can return chunks of every pipeline step result in real-time.

## Overview

The dynamic pipeline handler now supports two modes:
1. **Standard Mode**: Returns complete results after all steps are finished
2. **Streaming Mode**: Returns results as each step completes using Server-Sent Events (SSE)

## API Endpoints

### Standard Dynamic Pipeline
```
GET /dynamic?q=query&depth=3&skip_duplicates=true
```

**Parameters:**
- `q` (optional): Search query (default: "Novasco")
- `depth` (optional): Maximum depth (1-10, default: 53)
- `skip_duplicates` (optional): Skip duplicate searches (default: true)

**Response:** Complete JSON response with all pipeline results

### Streaming Dynamic Pipeline
```
GET /dynamic?q=query&depth=3&skip_duplicates=true&stream=true
```

**Parameters:**
- `q` (optional): Search query (default: "Novasco")
- `depth` (optional): Maximum depth (1-10, default: 53)
- `skip_duplicates` (optional): Skip duplicate searches (default: true)
- `stream` (required): Set to "true" to enable streaming

**Response:** Server-Sent Events stream with real-time step results

## Server-Sent Events Format

The streaming endpoint returns data in Server-Sent Events format:

```
event: step
data: {"step_number": 1, "step": {...}, "depth": 0, "category": "company_name", "keywords": [...]}

event: summary
data: {"step_number": 5, "step": {"Output": {"total_steps": 5, "successful_steps": 4, "failed_steps": 1, "max_depth_reached": 2}}}

event: complete
data: {"message": "Pipeline execution completed", "total_steps": 5}
```

### Event Types

1. **`step`**: Individual pipeline step result
2. **`summary`**: Final pipeline summary with statistics
3. **`error`**: Error occurred during execution
4. **`complete`**: Pipeline execution finished

### Step Event Data Structure

```json
{
  "step_number": 1,
  "step": {
    "Success": true,
    "Error": null,
    "Name": "ONAPI",
    "SearchParameter": "Novasco",
    "Output": {...},
    "keywordsPerCategory": {...}
  },
  "depth": 0,
  "category": "company_name",
  "keywords": ["Novasco", "Company"]
}
```

## Usage Examples

### JavaScript Client

```javascript
// Start streaming pipeline
function startStreamingPipeline(query, depth = 3) {
    const url = new URL('/dynamic', window.location.origin);
    url.searchParams.set('q', query);
    url.searchParams.set('depth', depth);
    url.searchParams.set('stream', 'true');
    
    const eventSource = new EventSource(url.toString());
    
    eventSource.addEventListener('step', function(event) {
        const data = JSON.parse(event.data);
        console.log('Step completed:', data);
        // Handle step result
        displayStep(data);
    });
    
    eventSource.addEventListener('summary', function(event) {
        const data = JSON.parse(event.data);
        console.log('Pipeline summary:', data);
        // Handle final summary
        displaySummary(data);
    });
    
    eventSource.addEventListener('complete', function(event) {
        const data = JSON.parse(event.data);
        console.log('Pipeline completed:', data);
        eventSource.close();
    });
    
    eventSource.onerror = function(event) {
        console.error('Stream error:', event);
        eventSource.close();
    };
}

// Display step result
function displayStep(stepData) {
    const stepDiv = document.createElement('div');
    stepDiv.className = 'step';
    stepDiv.innerHTML = `
        <h3>Step ${stepData.step_number}: ${stepData.step.Name}</h3>
        <p>Query: ${stepData.step.SearchParameter}</p>
        <p>Status: ${stepData.step.Success ? 'Success' : 'Failed'}</p>
        <p>Depth: ${stepData.depth}</p>
        ${stepData.step.Error ? `<p>Error: ${stepData.step.Error}</p>` : ''}
    `;
    document.getElementById('steps').appendChild(stepDiv);
}
```

### cURL Example

```bash
# Standard mode
curl "http://localhost:8080/dynamic?q=Novasco&depth=3"

# Streaming mode
curl -N "http://localhost:8080/dynamic?q=Novasco&depth=3&stream=true"
```

### Python Client

```python
import requests
import json

def stream_pipeline(query, depth=3):
    url = "http://localhost:8080/dynamic"
    params = {
        'q': query,
        'depth': depth,
        'stream': 'true'
    }
    
    response = requests.get(url, params=params, stream=True)
    
    for line in response.iter_lines():
        if line:
            line = line.decode('utf-8')
            if line.startswith('data: '):
                data = json.loads(line[6:])  # Remove 'data: ' prefix
                print(f"Received: {data}")

# Usage
stream_pipeline("Novasco", 3)
```

## HTML Demo

A complete HTML demo is available at `streaming_pipeline_demo.html` that demonstrates:

- Real-time step visualization
- Progress tracking
- Error handling
- Summary statistics
- Interactive controls

To use the demo:
1. Start your Go server
2. Open `streaming_pipeline_demo.html` in a web browser
3. Enter a search query and click "Start Streaming"

## Pipeline Flow

The streaming pipeline follows this flow:

1. **Initial Steps**: Create initial search steps for each domain
2. **Step Execution**: Execute each step and stream results
3. **Keyword Extraction**: Extract keywords from successful results
4. **New Steps Generation**: Create new steps from extracted keywords
5. **Depth Control**: Continue until max depth is reached
6. **Summary**: Send final statistics

### Step Types

- **ONAPI**: Company name searches
- **SCJ**: Legal case searches  
- **DGII**: Tax registry searches
- **PGR**: News searches
- **GOOGLE_DOCKING**: String search with fraud detection

## Configuration Options

### DynamicPipelineConfig

```go
type DynamicPipelineConfig struct {
    MaxDepth           int  // Maximum search depth
    MaxConcurrentSteps int  // Maximum concurrent steps
    DelayBetweenSteps  int  // Delay between steps (seconds)
    SkipDuplicates     bool // Skip duplicate keyword searches
}
```

### Default Configuration

```go
config := domain.DynamicPipelineConfig{
    MaxDepth:           53,
    MaxConcurrentSteps: 10,
    DelayBetweenSteps:  2,
    SkipDuplicates:     true,
}
```

## Error Handling

The streaming pipeline handles errors gracefully:

1. **Step Errors**: Individual step failures are streamed as error events
2. **Connection Errors**: Client disconnection stops the pipeline
3. **Timeout Errors**: Long-running steps can be cancelled
4. **Validation Errors**: Invalid parameters return HTTP 400

## Performance Considerations

- **Memory Usage**: Streaming reduces memory usage for large pipelines
- **Network**: SSE provides efficient real-time communication
- **Concurrency**: Steps are executed sequentially for better streaming experience
- **Caching**: Consider caching results for frequently searched terms

## Browser Compatibility

Server-Sent Events are supported in:
- Chrome 6+
- Firefox 6+
- Safari 5+
- Edge 12+
- Internet Explorer 10+

## Troubleshooting

### Common Issues

1. **Connection Drops**: Check network stability and server logs
2. **Missing Events**: Ensure proper event listener setup
3. **Memory Issues**: Reduce max depth or concurrent steps
4. **Timeout Errors**: Increase delay between steps

### Debug Mode

Enable debug logging by setting log level to debug in your Go application:

```go
log.SetLevel(log.DebugLevel)
```

## Advanced Usage

### Custom Step Processing

```javascript
eventSource.addEventListener('step', function(event) {
    const data = JSON.parse(event.data);
    
    // Custom processing based on domain type
    switch(data.step.Name) {
        case 'ONAPI':
            processCompanyData(data.step.Output);
            break;
        case 'GOOGLE_DOCKING':
            processSearchResults(data.step.Output);
            break;
        // ... other domains
    }
});
```

### Progress Tracking

```javascript
let totalSteps = 0;
let completedSteps = 0;

eventSource.addEventListener('step', function(event) {
    completedSteps++;
    updateProgressBar(completedSteps, totalSteps);
});

eventSource.addEventListener('summary', function(event) {
    totalSteps = event.data.step.Output.total_steps;
    updateProgressBar(completedSteps, totalSteps);
});
```

This streaming functionality provides a powerful way to monitor and interact with the dynamic pipeline execution in real-time, making it ideal for interactive applications and monitoring dashboards.

