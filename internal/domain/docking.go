package domain

import (
	"fmt"
	"insightful-intel/internal/stuff"
	"math"
	"net/url"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/gocolly/colly"
)

var _ DomainConnector[GoogleDockingResult] = &GoogleDocking{}

var FRAUD_KEYWORDS = []string{
	"fraude",
	"estafa",
	"denuncia",
	"engaño",
	"irregular",
	"contrato",
	"pagos",
	"retrasos",
	"robo",
}

// GoogleDocking represents a Google Docking string search connector
type GoogleDocking struct {
	Stuff    stuff.Stuff
	BasePath string
	PathMap  stuff.PathMap
}

// NewGoogleDockingDomain creates a new Google Docking domain instance
func NewGoogleDockingDomain() GoogleDocking {
	return GoogleDocking{
		BasePath: "https://html.duckduckgo.com/html/",
		Stuff:    *stuff.NewStuff(),
	}
}

// GoogleDockingResult represents a search result from Google Docking
type GoogleDockingResult struct {
	ID                   ID        `json:"id"`
	DomainSearchResultID ID        `json:"domainSearchResultId"`
	URL                  string    `json:"url"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	Relevance            float64   `json:"relevance"`
	Rank                 int       `json:"rank"`
	Keywords             []string  `json:"keywords"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

// GoogleDockingSearchParams holds parameters for Google Docking search
type GoogleDockingSearchParams struct {
	Query           string   `json:"query"`
	MaxResults      int      `json:"max_results"`
	MinRelevance    float64  `json:"min_relevance"`
	ExactMatch      bool     `json:"exact_match"`
	CaseSensitive   bool     `json:"case_sensitive"`
	IncludeKeywords []string `json:"include_keywords"`
	ExcludeKeywords []string `json:"exclude_keywords"`
}

// GoogleDockingBuilder provides a fluent interface for building Google Docking searches
type GoogleDockingBuilder struct {
	params GoogleDockingSearchParams
	gd     *GoogleDocking
}

// NewGoogleDockingBuilder creates a new Google Docking builder
func NewGoogleDockingBuilder() *GoogleDockingBuilder {
	return &GoogleDockingBuilder{
		params: GoogleDockingSearchParams{
			MaxResults:    10,
			MinRelevance:  0.1,
			ExactMatch:    false,
			CaseSensitive: false,
		},
		gd: func() *GoogleDocking {
			gd := NewGoogleDockingDomain()
			return &gd
		}(),
	}
}

// Query sets the search query
func (b *GoogleDockingBuilder) Query(query string) *GoogleDockingBuilder {
	b.params.Query = query
	return b
}

// MaxResults sets the maximum number of results
func (b *GoogleDockingBuilder) MaxResults(max int) *GoogleDockingBuilder {
	if max > 0 {
		b.params.MaxResults = max
	}
	return b
}

// MinRelevance sets the minimum relevance threshold
func (b *GoogleDockingBuilder) MinRelevance(min float64) *GoogleDockingBuilder {
	if min >= 0 && min <= 1 {
		b.params.MinRelevance = min
	}
	return b
}

// ExactMatch enables exact matching
func (b *GoogleDockingBuilder) ExactMatch(exact bool) *GoogleDockingBuilder {
	b.params.ExactMatch = exact
	return b
}

// CaseSensitive enables case-sensitive search
func (b *GoogleDockingBuilder) CaseSensitive(caseSensitive bool) *GoogleDockingBuilder {
	b.params.CaseSensitive = caseSensitive
	return b
}

// IncludeKeywords adds keywords that must be present
func (b *GoogleDockingBuilder) IncludeKeywords(keywords ...string) *GoogleDockingBuilder {
	b.params.IncludeKeywords = append(b.params.IncludeKeywords, keywords...)
	return b
}

// ExcludeKeywords adds keywords to exclude
func (b *GoogleDockingBuilder) ExcludeKeywords(keywords ...string) *GoogleDockingBuilder {
	b.params.ExcludeKeywords = append(b.params.ExcludeKeywords, keywords...)
	return b
}

// Build executes the search and returns results
func (b *GoogleDockingBuilder) Build() ([]GoogleDockingResult, error) {
	if b.params.Query == "" {
		return nil, fmt.Errorf("query is required")
	}
	return b.gd.SearchWithParams(b.params)
}

// BuildWithStats executes the search and returns results with statistics
func (b *GoogleDockingBuilder) BuildWithStats() ([]GoogleDockingResult, map[string]interface{}, error) {
	results, err := b.Build()
	if err != nil {
		return nil, nil, err
	}

	stats := b.gd.GetSearchStatistics(results)
	return results, stats, nil
}

// GetDomainType returns the domain type for Google Docking
func (*GoogleDocking) GetDomainType() DomainType {
	return DomainTypeGoogleDocking
}

// Search performs a Google Docking string search
func (gd *GoogleDocking) Search(query string) ([]GoogleDockingResult, error) {
	params := GoogleDockingSearchParams{
		Query:         query,
		MaxResults:    10,
		MinRelevance:  0.1,
		ExactMatch:    false,
		CaseSensitive: false,
	}
	return gd.SearchWithParams(params)
}

// SearchWithParams performs a Google Docking search with custom parameters
func (gd *GoogleDocking) SearchWithParams(params GoogleDockingSearchParams) ([]GoogleDockingResult, error) {
	if params.Query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	///
	q := params.Query

	if len(params.IncludeKeywords) > 0 {
		quotedElements := make([]string, len(params.IncludeKeywords))
		for i, s := range params.IncludeKeywords {
			quotedElements[i] = s // Add double quotes to each element
		}
		result := strings.Join(quotedElements, " OR ")

		q = fmt.Sprintf("%s intext:%s", q, result)
	}

	if len(params.ExcludeKeywords) > 0 {
		quotedElements := make([]string, len(params.ExcludeKeywords))
		for i, s := range params.ExcludeKeywords {
			quotedElements[i] = "-\"" + s + "\"" // Add double quotes to each element
		}
		result := strings.Join(quotedElements, " OR ")

		q = fmt.Sprintf("%s intext:(%s)", q, result)
	}

	URL, err := url.Parse(gd.BasePath)
	if err != nil {
		return nil, fmt.Errorf("wrong google url")
	}

	c := colly.NewCollector(
		colly.AllowedDomains(URL.Hostname()),
	)
	c.Async = true

	var news []GoogleDockingResult

	c.OnHTML("#links div.result", func(e *colly.HTMLElement) {
		article := GoogleDockingResult{}

		article.URL = e.ChildAttr("a.result__a", "href")
		article.Title = e.ChildText("a.result__a")
		article.Description = e.ChildText("a.result__snippet")

		news = append(news, article)

	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",             // Apply to all domains
		Parallelism: 2,               // Allow 2 concurrent requests
		Delay:       5 * time.Second, // Delay between requests to the same domain
	})

	c.Visit(fmt.Sprintf("%s?q=%s", gd.BasePath, url.QueryEscape(q)))

	c.Wait() // Wait for all requests to complete

	// Simulate Google search results (in a real implementation, this would call Google's API)

	// // Apply string matching and ranking
	// rankedResults := gd.rankResults(results, params)

	// // Filter by minimum relevance
	// filteredResults := gd.filterByRelevance(rankedResults, params.MinRelevance)

	return news, nil
}

// generateMockResults generates mock search results for demonstration
func (gd *GoogleDocking) generateMockResults(query string, maxResults int) []GoogleDockingResult {
	// In a real implementation, this would make HTTP requests to Google's search API
	mockResults := []GoogleDockingResult{
		{
			URL:         "https://example.com/page1",
			Title:       "Example Page 1 - " + query,
			Description: "This is a description containing " + query + " and other relevant information.",
			Keywords:    []string{query, "example", "page"},
		},
		{
			URL:         "https://example.com/page2",
			Title:       "Another Example - " + strings.ToUpper(query),
			Description: "Another description with " + query + " mentioned multiple times for better relevance.",
			Keywords:    []string{query, "another", "example"},
		},
		{
			URL:         "https://example.com/page3",
			Title:       "Related Content",
			Description: "This page discusses topics related to " + query + " and provides additional context.",
			Keywords:    []string{"related", "content", query},
		},
		{
			URL:         "https://example.com/page4",
			Title:       "Unrelated Page",
			Description: "This page doesn't contain the search term and should have low relevance.",
			Keywords:    []string{"unrelated", "page"},
		},
	}

	// Limit results to maxResults
	if len(mockResults) > maxResults {
		mockResults = mockResults[:maxResults]
	}

	return mockResults
}

// rankResults ranks search results based on relevance scoring
func (gd *GoogleDocking) rankResults(results []GoogleDockingResult, params GoogleDockingSearchParams) []GoogleDockingResult {
	for i := range results {
		results[i].Relevance = gd.calculateRelevance(results[i], params)
	}

	// Sort by relevance (highest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Relevance > results[j].Relevance
	})

	// Assign ranks
	for i := range results {
		results[i].Rank = i + 1
	}

	return results
}

// calculateRelevance calculates the relevance score for a search result
func (gd *GoogleDocking) calculateRelevance(result GoogleDockingResult, params GoogleDockingSearchParams) float64 {
	score := 0.0
	query := params.Query

	if !params.CaseSensitive {
		query = strings.ToLower(query)
	}

	// Title relevance (highest weight)
	titleScore := gd.calculateStringMatch(result.Title, query, params)
	score += titleScore * 3.0

	// Description relevance (medium weight)
	descScore := gd.calculateStringMatch(result.Description, query, params)
	score += descScore * 2.0

	// URL relevance (lower weight)
	urlScore := gd.calculateStringMatch(result.URL, query, params)
	score += urlScore * 1.0

	// Keywords relevance
	keywordScore := gd.calculateKeywordMatch(result.Keywords, query, params)
	score += keywordScore * 1.5

	// Exact match bonus
	if params.ExactMatch && gd.hasExactMatch(result, query, params) {
		score += 2.0
	}

	// Include keywords bonus
	if len(params.IncludeKeywords) > 0 {
		includeBonus := gd.calculateIncludeKeywordsBonus(result, params.IncludeKeywords)
		score += includeBonus
	}

	// Exclude keywords penalty
	if len(params.ExcludeKeywords) > 0 {
		excludePenalty := gd.calculateExcludeKeywordsPenalty(result, params.ExcludeKeywords)
		score -= excludePenalty
	}

	// Normalize score to 0-1 range
	return math.Min(score/10.0, 1.0)
}

// calculateStringMatch calculates how well a string matches the query
func (gd *GoogleDocking) calculateStringMatch(text, query string, params GoogleDockingSearchParams) float64 {
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
func (gd *GoogleDocking) calculateKeywordMatch(keywords []string, query string, params GoogleDockingSearchParams) float64 {
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
func (gd *GoogleDocking) hasExactMatch(result GoogleDockingResult, query string, params GoogleDockingSearchParams) bool {
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

// calculateIncludeKeywordsBonus calculates bonus for including required keywords
func (gd *GoogleDocking) calculateIncludeKeywordsBonus(result GoogleDockingResult, includeKeywords []string) float64 {
	bonus := 0.0
	text := strings.ToLower(result.Title + " " + result.Description)

	for _, keyword := range includeKeywords {
		if strings.Contains(text, strings.ToLower(keyword)) {
			bonus += 0.5
		}
	}

	return bonus
}

// calculateExcludeKeywordsPenalty calculates penalty for excluding unwanted keywords
func (gd *GoogleDocking) calculateExcludeKeywordsPenalty(result GoogleDockingResult, excludeKeywords []string) float64 {
	penalty := 0.0
	text := strings.ToLower(result.Title + " " + result.Description)

	for _, keyword := range excludeKeywords {
		if strings.Contains(text, strings.ToLower(keyword)) {
			penalty += 1.0
		}
	}

	return penalty
}

// filterByRelevance filters results by minimum relevance threshold
func (gd *GoogleDocking) filterByRelevance(results []GoogleDockingResult, minRelevance float64) []GoogleDockingResult {
	filtered := make([]GoogleDockingResult, 0, len(results))

	for _, result := range results {
		if result.Relevance >= minRelevance {
			filtered = append(filtered, result)
		}
	}

	return filtered
}

// levenshteinDistance calculates the Levenshtein distance between two strings
func (gd *GoogleDocking) levenshteinDistance(s1, s2 string) int {
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

// Implement DomainConnector[GoogleDockingResult] for GoogleDocking

// ProcessData processes a Google Docking result
func (gd *GoogleDocking) ProcessData(data GoogleDockingResult) (GoogleDockingResult, error) {
	if err := gd.ValidateData(data); err != nil {
		return GoogleDockingResult{}, err
	}
	return gd.TransformData(data), nil
}

// ValidateData validates a Google Docking result
func (gd *GoogleDocking) ValidateData(data GoogleDockingResult) error {
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
func (gd *GoogleDocking) TransformData(data GoogleDockingResult) GoogleDockingResult {
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
func (gd *GoogleDocking) GetDataByCategory(data GoogleDockingResult, category KeywordCategory) []string {
	switch category {
	case KeywordCategoryCompanyName:
		return gd.extractCompanyNames(data)
	case KeywordCategoryPersonName:
		return gd.extractPersonNames(data)
	case KeywordCategoryAddress:
		return gd.extractAddresses(data)
	case KeywordCategorySocialMedia:
		return gd.extractSocialMedia(data)
	default:
		return []string{}
	}
}

// extractCompanyNames extracts potential company names from the result
func (gd *GoogleDocking) extractCompanyNames(data GoogleDockingResult) []string {
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
func (gd *GoogleDocking) extractPersonNames(data GoogleDockingResult) []string {
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
func (gd *GoogleDocking) extractAddresses(data GoogleDockingResult) []string {
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
func (gd *GoogleDocking) extractSocialMedia(data GoogleDockingResult) []string {
	social := []string{}
	text := data.Title + " " + data.Description + " " + data.URL

	// Look for social media patterns
	socialPatterns := []string{"@", "twitter.com", "facebook.com", "instagram.com", "linkedin.com", "youtube.com"}

	for _, pattern := range socialPatterns {
		if strings.Contains(strings.ToLower(text), pattern) {
			social = append(social, pattern)
		}
	}

	return social
}

// GetSearchableKeywordCategories returns the categories that can be searched
func (gd *GoogleDocking) GetSearchableKeywordCategories() []KeywordCategory {
	return []KeywordCategory{
		KeywordCategoryPersonName,
	}
}

// GetFoundKeywordCategories returns the categories that can be found in results
func (gd *GoogleDocking) GetFoundKeywordCategories() []KeywordCategory {
	return []KeywordCategory{
		KeywordCategoryCompanyName,
		KeywordCategoryPersonName,
		KeywordCategoryAddress,
		KeywordCategorySocialMedia,
	}
}

// Advanced search methods

// SearchWithFilters performs a search with advanced filtering options
func (gd *GoogleDocking) SearchWithFilters(query string, filters map[string]interface{}) ([]GoogleDockingResult, error) {
	params := GoogleDockingSearchParams{
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
func (gd *GoogleDocking) GetSearchSuggestions(query string) ([]string, error) {
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
func (gd *GoogleDocking) GetSearchStatistics(results []GoogleDockingResult) map[string]interface{} {
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
func QuickSearch(query string) ([]GoogleDockingResult, error) {
	return NewGoogleDockingBuilder().Query(query).Build()
}

// AdvancedSearch provides a more complex search with multiple parameters
func AdvancedSearch(query string, maxResults int, minRelevance float64) ([]GoogleDockingResult, error) {
	return NewGoogleDockingBuilder().
		Query(query).
		MaxResults(maxResults).
		MinRelevance(minRelevance).
		Build()
}

// ExactSearch performs an exact match search
func ExactSearch(query string) ([]GoogleDockingResult, error) {
	return NewGoogleDockingBuilder().
		Query(query).
		ExactMatch(true).
		Build()
}

// CaseSensitiveSearch performs a case-sensitive search
func CaseSensitiveSearch(query string) ([]GoogleDockingResult, error) {
	return NewGoogleDockingBuilder().
		Query(query).
		CaseSensitive(true).
		Build()
}

// FilteredSearch performs a search with keyword filtering
func FilteredSearch(query string, includeKeywords, excludeKeywords []string) ([]GoogleDockingResult, error) {
	return NewGoogleDockingBuilder().
		Query(query).
		IncludeKeywords(includeKeywords...).
		ExcludeKeywords(excludeKeywords...).
		Build()
}
