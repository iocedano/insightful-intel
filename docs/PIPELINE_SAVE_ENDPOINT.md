# Pipeline Save Endpoint Documentation

## Overview

The pipeline save endpoint allows you to save pipeline responses to the database. This endpoint supports both dynamic pipeline results and domain search results.

## Endpoints

### 1. Dedicated Save Endpoint

**POST** `/api/pipeline/save`

This endpoint is specifically designed for saving pipeline responses to the database.

### 2. Enhanced Dynamic Pipeline Endpoint

**GET** `/dynamic?save=true`

The existing dynamic pipeline endpoint now supports an optional `save` parameter that automatically saves the results to the database.

## Usage Examples

### Saving Dynamic Pipeline Results

```bash
curl -X POST http://localhost:8080/api/pipeline/save \
  -H "Content-Type: application/json" \
  -d '{
    "total_steps": 3,
    "successful_steps": 2,
    "failed_steps": 1,
    "max_depth_reached": 2,
    "config": {
      "max_depth": 2,
      "max_concurrent_steps": 3,
      "delay_between_steps": 5,
      "skip_duplicates": true
    },
    "steps": [
      {
        "domain_type": "ONAPI",
        "search_parameter": "Novasco",
        "category": "COMPANY_NAME",
        "keywords": ["Novasco", "company"],
        "success": true,
        "error": null,
        "output": {
          "entities_found": 5,
          "search_time": "1.2s"
        },
        "keywords_per_category": {
          "COMPANY_NAME": ["Novasco"],
          "INDUSTRY": ["technology"]
        },
        "depth": 1
      }
    ]
  }'
```

### Saving Domain Search Results

```bash
curl -X POST http://localhost:8080/api/pipeline/save \
  -H "Content-Type: application/json" \
  -d '{
    "success": true,
    "error": null,
    "domain_type": "DGII",
    "search_parameter": "123456789",
    "keywords_per_category": {
      "COMPANY_NAME": ["Test Company"],
      "RNC": ["123456789"]
    },
    "output": {
      "company_name": "Test Company SRL",
      "status": "active",
      "category": "commercial"
    }
  }'
```

### Auto-save with Dynamic Pipeline

```bash
curl "http://localhost:8080/dynamic?q=Novasco&save=true&depth=3"
```

## Response Format

### Success Response

```json
{
  "success": true,
  "message": "Dynamic pipeline result saved successfully",
  "type": "dynamic_pipeline"
}
```

### Error Response

```json
{
  "error": "Unsupported pipeline result format"
}
```

## Database Schema

The endpoint saves data to the following tables:

### Dynamic Pipeline Results
- `dynamic_pipeline_results` - Main pipeline metadata
- `dynamic_pipeline_steps` - Individual step details

### Domain Search Results
- `domain_search_results` - Single domain search results

## Error Handling

The endpoint handles various error scenarios:

1. **Invalid JSON**: Returns 400 Bad Request
2. **Unsupported Format**: Returns 400 Bad Request for unrecognized data structures
3. **Database Errors**: Returns 500 Internal Server Error
4. **Method Not Allowed**: Returns 405 for non-POST requests

## Testing

Use the provided test file to verify the endpoint functionality:

```bash
go run test_pipeline_save.go
```

## Integration with Existing Pipeline

The dynamic pipeline endpoint (`/dynamic`) now supports an optional `save=true` parameter that automatically saves results to the database without affecting the response format. This allows for seamless integration with existing workflows.

## Notes

- The endpoint automatically detects the data type (DynamicPipelineResult vs DomainSearchResult)
- Database errors are logged but don't fail the dynamic pipeline execution when using the `save` parameter
- All timestamps are automatically managed by the database
- The endpoint supports CORS for web applications

