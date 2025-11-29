package module

import (
	"encoding/json"
	"fmt"
	"insightful-intel/internal/custom"
	"insightful-intel/internal/domain"
	"io"
	"log"
	"math"
	"net/url"
	"os"
	"strings"
	"unicode"

	"github.com/davecgh/go-spew/spew"
)

var _ domain.DomainConnector[domain.GoogleDorkingResult] = &GoogleDorking{}

// GoogleDorking represents a Google Docking string search connector
type GoogleDorking struct {
	Stuff    custom.Client
	BasePath string
	PathMap  custom.CustomPathMap
}

// NewGoogleDorkingDomain creates a new Google Docking domain instance
func NewGoogleDorkingDomain() GoogleDorking {
	googleApiKey := os.Getenv("GOOGLE_API_KEY")
	googleSearchEngineId := os.Getenv("GOOGLE_CX_KEY")
	if googleApiKey == "" || googleSearchEngineId == "" {
		log.Fatal("GOOGLE_API_KEY and GOOGLE_CX_KEY are not set")
	}

	googleSearchUrl := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s", googleApiKey, googleSearchEngineId)

	return GoogleDorking{
		BasePath: googleSearchUrl,
		Stuff:    *custom.NewClient(),
	}
}

type GoogleDorkingSearchResponse struct {
	Items []domain.GoogleDorkingResult `json:"items"`
}

// GoogleDorkingBuilder provides a fluent interface for building Google Docking searches
type GoogleDorkingBuilder struct {
	params domain.GoogleDorkingSearchParams
	gd     *GoogleDorking
}

// NewGoogleDorkingBuilder creates a new Google Docking builder
func NewGoogleDorkingBuilder() *GoogleDorkingBuilder {
	return &GoogleDorkingBuilder{
		params: domain.GoogleDorkingSearchParams{
			MaxResults:    10,
			MinRelevance:  0.1,
			ExactMatch:    false,
			CaseSensitive: false,
		},
		gd: func() *GoogleDorking {
			gd := NewGoogleDorkingDomain()
			return &gd
		}(),
	}
}

// Query sets the search query
func (b *GoogleDorkingBuilder) Query(query string) *GoogleDorkingBuilder {
	b.params.Query = query
	return b
}

// MaxResults sets the maximum number of results
func (b *GoogleDorkingBuilder) MaxResults(max int) *GoogleDorkingBuilder {
	if max > 0 {
		b.params.MaxResults = max
	}
	return b
}

// MinRelevance sets the minimum relevance threshold
func (b *GoogleDorkingBuilder) MinRelevance(min float64) *GoogleDorkingBuilder {
	if min >= 0 && min <= 1 {
		b.params.MinRelevance = min
	}
	return b
}

// ExactMatch enables exact matching
func (b *GoogleDorkingBuilder) ExactMatch(exact bool) *GoogleDorkingBuilder {
	b.params.ExactMatch = exact
	return b
}

// CaseSensitive enables case-sensitive search
func (b *GoogleDorkingBuilder) CaseSensitive(caseSensitive bool) *GoogleDorkingBuilder {
	b.params.CaseSensitive = caseSensitive
	return b
}

// IncludeKeywords adds keywords that must be present
func (b *GoogleDorkingBuilder) IncludeKeywords(keywords ...string) *GoogleDorkingBuilder {
	b.params.IncludeKeywords = append(b.params.IncludeKeywords, keywords...)
	return b
}

// FileTypeKeywords adds file type keywords
func (b *GoogleDorkingBuilder) FileTypeKeywords(keywords ...string) *GoogleDorkingBuilder {
	b.params.FileTypeKeywords = append(b.params.FileTypeKeywords, keywords...)
	return b
}

// SitesKeywords adds sites keywords
func (b *GoogleDorkingBuilder) SitesKeywords(keywords ...string) *GoogleDorkingBuilder {
	b.params.SitesKeywords = append(b.params.SitesKeywords, keywords...)
	return b
}

// InURLKeywords adds keywords to the URL
func (b *GoogleDorkingBuilder) InURLKeywords(keywords ...string) *GoogleDorkingBuilder {
	b.params.InURLKeywords = append(b.params.InURLKeywords, keywords...)
	return b
}

// ExcludeKeywords adds keywords to exclude
func (b *GoogleDorkingBuilder) ExcludeKeywords(keywords ...string) *GoogleDorkingBuilder {
	b.params.ExcludeKeywords = append(b.params.ExcludeKeywords, keywords...)
	return b
}

// Build executes the search and returns results
func (b *GoogleDorkingBuilder) Build() ([]domain.GoogleDorkingResult, error) {
	if b.params.Query == "" {
		return nil, fmt.Errorf("query is required")
	}
	return b.gd.SearchWithParams(b.params)
}

// BuildWithStats executes the search and returns results with statistics
func (b *GoogleDorkingBuilder) BuildWithStats() ([]domain.GoogleDorkingResult, map[string]interface{}, error) {
	results, err := b.Build()
	if err != nil {
		return nil, nil, err
	}

	stats := b.gd.GetSearchStatistics(results)
	return results, stats, nil
}

// GetDomainType returns the domain type for Google Docking
func (*GoogleDorking) GetDomainType() domain.DomainType {
	return domain.DomainTypeGoogleDorking
}

// Search performs a Google Docking string search
func (gd *GoogleDorking) Search(query string) ([]domain.GoogleDorkingResult, error) {
	params := domain.GoogleDorkingSearchParams{
		Query:         query,
		MaxResults:    10,
		MinRelevance:  0.1,
		ExactMatch:    false,
		CaseSensitive: false,
	}
	return gd.SearchWithParams(params)
}

// SearchWithParams performs a Google Docking search with custom parameters
func (gd *GoogleDorking) SearchWithParams(params domain.GoogleDorkingSearchParams) ([]domain.GoogleDorkingResult, error) {
	if params.Query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	var q string

	if len(params.IncludeKeywords) > 0 {
		result := strings.Join(params.IncludeKeywords, " OR ")

		q = fmt.Sprintf("%s intext:(%s)", q, result)
	}

	if len(params.FileTypeKeywords) > 0 {
		quotedElements := make([]string, len(params.FileTypeKeywords))
		for i, s := range params.FileTypeKeywords {
			quotedElements[i] = "filetype:" + s // Add double quotes to each element
		}
		result := strings.Join(quotedElements, " OR ")

		q = fmt.Sprintf("%s (%s)", q, result)
	}
	if len(params.SitesKeywords) > 0 {
		quotedElements := make([]string, len(params.SitesKeywords))
		for i, s := range params.SitesKeywords {
			quotedElements[i] = "site:" + s // Add double quotes to each element
		}
		result := strings.Join(quotedElements, " OR ")

		q = fmt.Sprintf("%s %s", q, result)
	}

	if len(params.InURLKeywords) > 0 {
		quotedElements := make([]string, len(params.InURLKeywords))
		for i, s := range params.InURLKeywords {
			quotedElements[i] = "inurl:" + s // Add double quotes to each element
		}
		result := strings.Join(quotedElements, " OR ")

		q = fmt.Sprintf("%s (%s)", q, result)
	}

	if len(params.ExcludeKeywords) > 0 {
		quotedElements := make([]string, len(params.ExcludeKeywords))
		for i, s := range params.ExcludeKeywords {
			quotedElements[i] = "-\"" + s + "\"" // Add double quotes to each element
		}
		result := strings.Join(quotedElements, " OR ")

		q = fmt.Sprintf("%s intext:(%s)", q, result)
	}

	q = fmt.Sprintf("%s %s", params.Query, q)

	spew.Dump("q", q)

	resp, err := gd.Stuff.Get(fmt.Sprintf("%s&q=%s", gd.BasePath, url.QueryEscape(q)), map[string]string{}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result GoogleDorkingSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return result.Items, nil
}

// calculateStringMatch calculates how well a string matches the query
func (gd *GoogleDorking) calculateStringMatch(text, query string, params domain.GoogleDorkingSearchParams) float64 {
	if text == "" {
		return 0.0
	}

	if !params.CaseSensitive {
		text = strings.ToLower(text)
	}

	// Exact match gets highest score
	if text == query {
		return 1.0
	}

	// Contains match
	if strings.Contains(text, query) {
		// Calculate frequency and position bonus
		frequency := float64(strings.Count(text, query))
		position := float64(strings.Index(text, query))
		length := float64(len(text))

		// Position bonus (earlier is better)
		positionBonus := 1.0 - (position / length)

		// Frequency bonus (more occurrences is better, but with diminishing returns)
		frequencyBonus := math.Log(frequency + 1)

		return (0.7 + positionBonus*0.2 + frequencyBonus*0.1)
	}

	// Fuzzy match using Levenshtein distance
	distance := gd.levenshteinDistance(text, query)
	maxLen := math.Max(float64(len(text)), float64(len(query)))

	if maxLen == 0 {
		return 0.0
	}

	similarity := 1.0 - (float64(distance) / maxLen)

	// Only consider fuzzy matches above a threshold
	if similarity > 0.6 {
		return similarity * 0.5
	}

	return 0.0
}

// calculateKeywordMatch calculates keyword matching score
func (gd *GoogleDorking) calculateKeywordMatch(keywords []string, query string, params domain.GoogleDorkingSearchParams) float64 {
	if len(keywords) == 0 {
		return 0.0
	}

	score := 0.0
	queryWords := strings.Fields(query)

	for _, keyword := range keywords {
		keywordLower := keyword
		if !params.CaseSensitive {
			keywordLower = strings.ToLower(keyword)
		}

		for _, queryWord := range queryWords {
			queryWordLower := queryWord
			if !params.CaseSensitive {
				queryWordLower = strings.ToLower(queryWord)
			}

			if keywordLower == queryWordLower {
				score += 1.0
			} else if strings.Contains(keywordLower, queryWordLower) {
				score += 0.7
			} else if strings.Contains(queryWordLower, keywordLower) {
				score += 0.7
			}
		}
	}

	return score / float64(len(keywords))
}

// hasExactMatch checks if the result has an exact match
func (gd *GoogleDorking) hasExactMatch(result domain.GoogleDorkingResult, query string, params domain.GoogleDorkingSearchParams) bool {
	texts := []string{result.Title, result.Description, result.URL}

	for _, text := range texts {
		if !params.CaseSensitive {
			text = strings.ToLower(text)
		}

		if text == query {
			return true
		}
	}

	return false
}

// levenshteinDistance calculates the Levenshtein distance between two strings
func (gd *GoogleDorking) levenshteinDistance(s1, s2 string) int {
	r1, r2 := []rune(s1), []rune(s2)
	rows := len(r1) + 1
	cols := len(r2) + 1

	d := make([][]int, rows)
	for i := range d {
		d[i] = make([]int, cols)
	}

	for i := 1; i < rows; i++ {
		d[i][0] = i
	}
	for j := 1; j < cols; j++ {
		d[0][j] = j
	}

	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}
			d[i][j] = min(d[i-1][j]+1, d[i][j-1]+1, d[i-1][j-1]+cost)
		}
	}

	return d[rows-1][cols-1]
}

// min returns the minimum of three integers
func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}

// Implement DomainConnector[domain.GoogleDorkingResult] for GoogleDorking

// ProcessData processes a Google Docking result
func (gd *GoogleDorking) ProcessData(data domain.GoogleDorkingResult) (domain.GoogleDorkingResult, error) {
	if err := gd.ValidateData(data); err != nil {
		return domain.GoogleDorkingResult{}, err
	}
	return gd.TransformData(data), nil
}

// ValidateData validates a Google Docking result
func (gd *GoogleDorking) ValidateData(data domain.GoogleDorkingResult) error {
	if data.URL == "" {
		return fmt.Errorf("URL is required")
	}
	if data.Title == "" {
		return fmt.Errorf("title is required")
	}
	if data.Relevance < 0 || data.Relevance > 1 {
		return fmt.Errorf("relevance must be between 0 and 1")
	}
	return nil
}

// TransformData transforms a Google Docking result
func (gd *GoogleDorking) TransformData(data domain.GoogleDorkingResult) domain.GoogleDorkingResult {
	transformed := data
	transformed.URL = strings.TrimSpace(data.URL)
	transformed.Title = strings.TrimSpace(data.Title)
	transformed.Description = strings.TrimSpace(data.Description)

	// Clean and normalize keywords
	cleanedKeywords := make([]string, 0, len(data.Keywords))
	for _, keyword := range data.Keywords {
		cleaned := strings.TrimSpace(keyword)
		if cleaned != "" {
			cleanedKeywords = append(cleanedKeywords, cleaned)
		}
	}
	transformed.Keywords = cleanedKeywords

	return transformed
}

// GetDataByCategory extracts data by keyword category
func (gd *GoogleDorking) GetDataByCategory(data domain.GoogleDorkingResult, category domain.KeywordCategory) []string {
	switch category {
	// case domain.KeywordCategoryCompanyName:
	// 	return gd.extractCompanyNames(data)
	// case domain.KeywordCategoryPersonName:
	// 	return gd.extractPersonNames(data)
	// case domain.KeywordCategoryAddress:
	// 	return gd.extractAddresses(data)
	case domain.KeywordCategorySocialMedia:
		return gd.extractSocialMedia(data)
	default:
		return []string{}
	}
}

// extractCompanyNames extracts potential company names from the result
func (gd *GoogleDorking) extractCompanyNames(data domain.GoogleDorkingResult) []string {
	companies := []string{}
	text := data.Title + " " + data.Description

	// Simple heuristic: look for capitalized words that might be company names
	words := strings.Fields(text)
	for i, word := range words {
		if len(word) > 2 && unicode.IsUpper(rune(word[0])) {
			// Check if it's followed by common company suffixes
			if i+1 < len(words) {
				nextWord := strings.ToLower(words[i+1])
				if nextWord == "inc" || nextWord == "corp" || nextWord == "llc" ||
					nextWord == "ltd" || nextWord == "co" || nextWord == "company" {
					companies = append(companies, word+" "+words[i+1])
				}
			}
			// Add standalone capitalized words
			companies = append(companies, word)
		}
	}

	return companies
}

// extractPersonNames extracts potential person names from the result
func (gd *GoogleDorking) extractPersonNames(data domain.GoogleDorkingResult) []string {
	names := []string{}
	text := data.Title + " " + data.Description

	// Simple heuristic: look for patterns like "First Last" or "Mr. Last"
	words := strings.Fields(text)
	for i, word := range words {
		if len(word) > 1 && unicode.IsUpper(rune(word[0])) {
			// Check for title prefixes
			if word == "Mr." || word == "Mrs." || word == "Ms." || word == "Dr." {
				if i+1 < len(words) {
					names = append(names, word+" "+words[i+1])
				}
			}
			// Check for first name + last name pattern
			if i+1 < len(words) {
				nextWord := words[i+1]
				if len(nextWord) > 1 && unicode.IsUpper(rune(nextWord[0])) {
					names = append(names, word+" "+nextWord)
				}
			}
		}
	}

	return names
}

// extractAddresses extracts potential addresses from the result
func (gd *GoogleDorking) extractAddresses(data domain.GoogleDorkingResult) []string {
	addresses := []string{}
	text := data.Title + " " + data.Description

	// Look for common address patterns
	words := strings.Fields(text)
	for i, word := range words {
		// Look for street numbers
		if len(word) > 0 && unicode.IsDigit(rune(word[0])) {
			address := word
			// Collect following words that might be part of the address
			for j := i + 1; j < len(words) && j < i+5; j++ {
				nextWord := words[j]
				if strings.Contains(nextWord, ",") || strings.Contains(nextWord, ".") {
					address += " " + nextWord
					break
				}
				address += " " + nextWord
			}
			addresses = append(addresses, address)
		}
	}

	return addresses
}

// extractSocialMedia extracts social media handles and URLs
func (gd *GoogleDorking) extractSocialMedia(data domain.GoogleDorkingResult) []string {
	social := []string{}
	text := data.Title + " " + data.Description + " " + data.URL

	// Loo	 for social media patterns
	socialPatterns := []string{"@", "twitter.com", "facebook.com", "instagram.com", "linkedin.com", "youtube.com"}

	for _, pattern := range socialPatterns {
		if strings.Contains(strings.ToLower(text), pattern) {
			social = append(social, pattern)
		}
	}

	return social
}

// GetSearchableKeywordCategories returns the categories that can be searched
func (gd *GoogleDorking) GetSearchableKeywordCategories() []domain.KeywordCategory {
	return []domain.KeywordCategory{
		domain.KeywordCategoryCompanyName,
		domain.KeywordCategoryPersonName,
	}
}

// GetFoundKeywordCategories returns the categories that can be found in results
func (gd *GoogleDorking) GetFoundKeywordCategories() []domain.KeywordCategory {
	return []domain.KeywordCategory{}
}

// Advanced search methods

// SearchWithFilters performs a search with advanced filtering options
func (gd *GoogleDorking) SearchWithFilters(query string, filters map[string]interface{}) ([]domain.GoogleDorkingResult, error) {
	params := domain.GoogleDorkingSearchParams{
		Query:        query,
		MaxResults:   10,
		MinRelevance: 0.1,
	}

	// Apply filters
	if maxResults, ok := filters["max_results"].(int); ok {
		params.MaxResults = maxResults
	}
	if minRelevance, ok := filters["min_relevance"].(float64); ok {
		params.MinRelevance = minRelevance
	}
	if exactMatch, ok := filters["exact_match"].(bool); ok {
		params.ExactMatch = exactMatch
	}
	if caseSensitive, ok := filters["case_sensitive"].(bool); ok {
		params.CaseSensitive = caseSensitive
	}
	if includeKeywords, ok := filters["include_keywords"].([]string); ok {
		params.IncludeKeywords = includeKeywords
	}
	if excludeKeywords, ok := filters["exclude_keywords"].([]string); ok {
		params.ExcludeKeywords = excludeKeywords
	}

	return gd.SearchWithParams(params)
}

// GetSearchSuggestions returns search suggestions based on the query
func (gd *GoogleDorking) GetSearchSuggestions(query string) ([]string, error) {
	if query == "" {
		return []string{}, nil
	}

	// Simple suggestion generation based on common patterns
	suggestions := []string{
		query + " company",
		query + " person",
		query + " address",
		"about " + query,
		query + " information",
		query + " details",
	}

	return suggestions, nil
}

// GetSearchStatistics returns statistics about the search results
func (gd *GoogleDorking) GetSearchStatistics(results []domain.GoogleDorkingResult) map[string]interface{} {
	if len(results) == 0 {
		return map[string]interface{}{
			"total_results":     0,
			"average_relevance": 0.0,
			"max_relevance":     0.0,
			"min_relevance":     0.0,
		}
	}

	totalRelevance := 0.0
	maxRelevance := results[0].Relevance
	minRelevance := results[0].Relevance

	for _, result := range results {
		totalRelevance += result.Relevance
		if result.Relevance > maxRelevance {
			maxRelevance = result.Relevance
		}
		if result.Relevance < minRelevance {
			minRelevance = result.Relevance
		}
	}

	return map[string]interface{}{
		"total_results":     len(results),
		"average_relevance": totalRelevance / float64(len(results)),
		"max_relevance":     maxRelevance,
		"min_relevance":     minRelevance,
	}
}

// Fluent API examples and helper functions

// QuickSearch provides a simple one-liner search
func QuickSearch(query string) ([]domain.GoogleDorkingResult, error) {
	return NewGoogleDorkingBuilder().Query(query).Build()
}

// AdvancedSearch provides a more complex search with multiple parameters
func AdvancedSearch(query string, maxResults int, minRelevance float64) ([]domain.GoogleDorkingResult, error) {
	return NewGoogleDorkingBuilder().
		Query(query).
		MaxResults(maxResults).
		MinRelevance(minRelevance).
		Build()
}

// ExactSearch performs an exact match search
func ExactSearch(query string) ([]domain.GoogleDorkingResult, error) {
	return NewGoogleDorkingBuilder().
		Query(query).
		ExactMatch(true).
		Build()
}

// CaseSensitiveSearch performs a case-sensitive search
func CaseSensitiveSearch(query string) ([]domain.GoogleDorkingResult, error) {
	return NewGoogleDorkingBuilder().
		Query(query).
		CaseSensitive(true).
		Build()
}

// FilteredSearch performs a search with keyword filtering
func FilteredSearch(query string, includeKeywords, excludeKeywords []string) ([]domain.GoogleDorkingResult, error) {
	return NewGoogleDorkingBuilder().
		Query(query).
		IncludeKeywords(includeKeywords...).
		ExcludeKeywords(excludeKeywords...).
		Build()
}
