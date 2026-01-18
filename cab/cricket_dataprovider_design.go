package main

/*
==========================================
CRICKET DATA PROVIDER API - LOW LEVEL DESIGN
==========================================

REQUIREMENTS:
- DataProvider API with three sources: cricbuzz, cricinfo, skysports
- Some methods are common to all sources
- Some methods are specific to individual sources
- API should return response based on the source sent in the request

DESIGN APPROACH:
- Strategy Pattern: Each source implements common interface
- Factory Pattern: Create appropriate source based on input
- Interface Segregation: Separate interfaces for common and source-specific methods
- Adapter Pattern: Unified API layer that routes to appropriate source

==========================================
INTERFACES & STRUCTS
==========================================
*/

// SourceType represents the cricket data source
type SourceType string

const (
	SourceCricbuzz  SourceType = "cricbuzz"
	SourceCricinfo  SourceType = "cricinfo"
	SourceSkysports SourceType = "skysports"
)

// CommonDataProvider defines methods common to all cricket data sources
type CommonDataProvider interface {
	// GetMatchDetails retrieves basic match information
	GetMatchDetails(matchID string) (*MatchResponse, error)

	// GetLiveScore retrieves current live score
	GetLiveScore(matchID string) (*ScoreResponse, error)

	// GetTeamInfo retrieves team information
	GetTeamInfo(teamID string) (*TeamResponse, error)

	// GetPlayerInfo retrieves player information
	GetPlayerInfo(playerID string) (*PlayerResponse, error)

	// GetUpcomingMatches retrieves list of upcoming matches
	GetUpcomingMatches() (*MatchesResponse, error)

	// GetSourceType returns the source type identifier
	GetSourceType() SourceType
}

// CricbuzzProvider defines methods specific to Cricbuzz source
type CricbuzzProvider interface {
	CommonDataProvider

	// GetCommentary retrieves ball-by-ball commentary (Cricbuzz specific)
	GetCommentary(matchID string, over int) (*CommentaryResponse, error)

	// GetMatchPredictions retrieves match predictions (Cricbuzz specific)
	GetMatchPredictions(matchID string) (*PredictionsResponse, error)
}

// CricinfoProvider defines methods specific to Cricinfo source
type CricinfoProvider interface {
	CommonDataProvider

	// GetDetailedStats retrieves comprehensive statistics (Cricinfo specific)
	GetDetailedStats(matchID string) (*StatsResponse, error)

	// GetHistoricalData retrieves historical match data (Cricinfo specific)
	GetHistoricalData(teamID string, year int) (*HistoricalResponse, error)
}

// SkysportsProvider defines methods specific to Skysports source
type SkysportsProvider interface {
	CommonDataProvider

	// GetVideoHighlights retrieves video highlights (Skysports specific)
	GetVideoHighlights(matchID string) (*HighlightsResponse, error)

	// GetNewsArticles retrieves cricket news articles (Skysports specific)
	GetNewsArticles(limit int) (*NewsResponse, error)
}

// Structs for each source implementation
// These structs will implement the respective interfaces

// CricbuzzDataProvider implements CricbuzzProvider
type CricbuzzDataProvider struct {
	// Configuration fields
	apiKey     string
	baseURL   string
	timeout   int
	// Add other Cricbuzz-specific fields
}

// CricinfoDataProvider implements CricinfoProvider
type CricinfoDataProvider struct {
	// Configuration fields
	apiKey     string
	baseURL   string
	timeout   int
	// Add other Cricinfo-specific fields
}

// SkysportsDataProvider implements SkysportsProvider
type SkysportsDataProvider struct {
	// Configuration fields
	apiKey     string
	baseURL   string
	timeout   int
	// Add other Skysports-specific fields
}

// Response structs (common response types)
type MatchResponse struct {
	MatchID      string
	Team1        string
	Team2        string
	Venue        string
	Date         string
	Status       string
	Source       SourceType
	// Add other common fields
}

type ScoreResponse struct {
	MatchID      string
	Team1Score   string
	Team2Score   string
	Overs        string
	CurrentOver  string
	Source       SourceType
	// Add other common fields
}

type TeamResponse struct {
	TeamID       string
	TeamName     string
	Players      []string
	Source       SourceType
	// Add other common fields
}

type PlayerResponse struct {
	PlayerID     string
	Name         string
	Role         string
	Source       SourceType
	// Add other common fields
}

type MatchesResponse struct {
	Matches      []MatchResponse
	Total        int
	Source       SourceType
	// Add other common fields
}

// Source-specific response structs
type CommentaryResponse struct {
	MatchID      string
	Over         int
	Commentary   []string
	Source       SourceType
}

type PredictionsResponse struct {
	MatchID      string
	Predictions  map[string]float64
	Source       SourceType
}

type StatsResponse struct {
	MatchID      string
	Statistics   map[string]interface{}
	Source       SourceType
}

type HistoricalResponse struct {
	TeamID       string
	Year         int
	Matches      []MatchResponse
	Source       SourceType
}

type HighlightsResponse struct {
	MatchID      string
	VideoURLs    []string
	Source       SourceType
}

type NewsResponse struct {
	Articles     []Article
	Total        int
	Source       SourceType
}

type Article struct {
	ID          string
	Title       string
	Content     string
	PublishedAt string
}

// DataProviderFactory creates appropriate data provider based on source type
type DataProviderFactory struct {
	providers map[SourceType]CommonDataProvider
}

// NewDataProviderFactory creates a new factory instance
func NewDataProviderFactory() *DataProviderFactory {
	return &DataProviderFactory{
		providers: make(map[SourceType]CommonDataProvider),
	}
}

// RegisterProvider registers a data provider for a specific source
func (f *DataProviderFactory) RegisterProvider(source SourceType, provider CommonDataProvider) {
	f.providers[source] = provider
}

// GetProvider returns the data provider for the given source
func (f *DataProviderFactory) GetProvider(source SourceType) (CommonDataProvider, error) {
	provider, exists := f.providers[source]
	if !exists {
		return nil, ErrSourceNotSupported
	}
	return provider, nil
}

// DataProviderService is the main service that routes requests to appropriate source
type DataProviderService struct {
	factory *DataProviderFactory
}

// NewDataProviderService creates a new service instance
func NewDataProviderService(factory *DataProviderFactory) *DataProviderService {
	return &DataProviderService{
		factory: factory,
	}
}

// GetMatchDetails routes to appropriate source and returns match details
func (s *DataProviderService) GetMatchDetails(source SourceType, matchID string) (*MatchResponse, error) {
	provider, err := s.factory.GetProvider(source)
	if err != nil {
		return nil, err
	}
	return provider.GetMatchDetails(matchID)
}

// GetLiveScore routes to appropriate source and returns live score
func (s *DataProviderService) GetLiveScore(source SourceType, matchID string) (*ScoreResponse, error) {
	provider, err := s.factory.GetProvider(source)
	if err != nil {
		return nil, err
	}
	return provider.GetLiveScore(matchID)
}

// GetTeamInfo routes to appropriate source and returns team info
func (s *DataProviderService) GetTeamInfo(source SourceType, teamID string) (*TeamResponse, error) {
	provider, err := s.factory.GetProvider(source)
	if err != nil {
		return nil, err
	}
	return provider.GetTeamInfo(teamID)
}

// GetPlayerInfo routes to appropriate source and returns player info
func (s *DataProviderService) GetPlayerInfo(source SourceType, playerID string) (*PlayerResponse, error) {
	provider, err := s.factory.GetProvider(source)
	if err != nil {
		return nil, err
	}
	return provider.GetPlayerInfo(playerID)
}

// GetUpcomingMatches routes to appropriate source and returns upcoming matches
func (s *DataProviderService) GetUpcomingMatches(source SourceType) (*MatchesResponse, error) {
	provider, err := s.factory.GetProvider(source)
	if err != nil {
		return nil, err
	}
	return provider.GetUpcomingMatches()
}

// Source-specific methods that need type assertion

// GetCommentary retrieves commentary from Cricbuzz source
func (s *DataProviderService) GetCommentary(matchID string, over int) (*CommentaryResponse, error) {
	provider, err := s.factory.GetProvider(SourceCricbuzz)
	if err != nil {
		return nil, err
	}
	cricbuzzProvider, ok := provider.(CricbuzzProvider)
	if !ok {
		return nil, ErrInvalidProviderType
	}
	return cricbuzzProvider.GetCommentary(matchID, over)
}

// GetMatchPredictions retrieves predictions from Cricbuzz source
func (s *DataProviderService) GetMatchPredictions(matchID string) (*PredictionsResponse, error) {
	provider, err := s.factory.GetProvider(SourceCricbuzz)
	if err != nil {
		return nil, err
	}
	cricbuzzProvider, ok := provider.(CricbuzzProvider)
	if !ok {
		return nil, ErrInvalidProviderType
	}
	return cricbuzzProvider.GetMatchPredictions(matchID)
}

// GetDetailedStats retrieves detailed stats from Cricinfo source
func (s *DataProviderService) GetDetailedStats(matchID string) (*StatsResponse, error) {
	provider, err := s.factory.GetProvider(SourceCricinfo)
	if err != nil {
		return nil, err
	}
	cricinfoProvider, ok := provider.(CricinfoProvider)
	if !ok {
		return nil, ErrInvalidProviderType
	}
	return cricinfoProvider.GetDetailedStats(matchID)
}

// GetHistoricalData retrieves historical data from Cricinfo source
func (s *DataProviderService) GetHistoricalData(teamID string, year int) (*HistoricalResponse, error) {
	provider, err := s.factory.GetProvider(SourceCricinfo)
	if err != nil {
		return nil, err
	}
	cricinfoProvider, ok := provider.(CricinfoProvider)
	if !ok {
		return nil, ErrInvalidProviderType
	}
	return cricinfoProvider.GetHistoricalData(teamID, year)
}

// GetVideoHighlights retrieves video highlights from Skysports source
func (s *DataProviderService) GetVideoHighlights(matchID string) (*HighlightsResponse, error) {
	provider, err := s.factory.GetProvider(SourceSkysports)
	if err != nil {
		return nil, err
	}
	skysportsProvider, ok := provider.(SkysportsProvider)
	if !ok {
		return nil, ErrInvalidProviderType
	}
	return skysportsProvider.GetVideoHighlights(matchID)
}

// GetNewsArticles retrieves news articles from Skysports source
func (s *DataProviderService) GetNewsArticles(limit int) (*NewsResponse, error) {
	provider, err := s.factory.GetProvider(SourceSkysports)
	if err != nil {
		return nil, err
	}
	skysportsProvider, ok := provider.(SkysportsProvider)
	if !ok {
		return nil, ErrInvalidProviderType
	}
	return skysportsProvider.GetNewsArticles(limit)
}

// Error definitions
var (
	ErrSourceNotSupported = &ProviderError{Message: "source not supported"}
	ErrInvalidProviderType = &ProviderError{Message: "invalid provider type"}
)

type ProviderError struct {
	Message string
}

func (e *ProviderError) Error() string {
	return e.Message
}

/*
==========================================
DESIGN PATTERNS USED
==========================================

1. Strategy Pattern:
   - CommonDataProvider interface defines common strategy
   - Each source (Cricbuzz, Cricinfo, Skysports) implements the strategy differently

2. Factory Pattern:
   - DataProviderFactory creates appropriate provider based on source type
   - Centralized creation logic

3. Adapter Pattern:
   - DataProviderService acts as adapter/unified interface
   - Routes requests to appropriate source implementation

4. Interface Segregation:
   - CommonDataProvider for common methods
   - Source-specific interfaces (CricbuzzProvider, CricinfoProvider, SkysportsProvider)
   - Clients only depend on interfaces they use

==========================================
USAGE FLOW
==========================================

1. Initialize Factory:
   factory := NewDataProviderFactory()
   factory.RegisterProvider(SourceCricbuzz, &CricbuzzDataProvider{...})
   factory.RegisterProvider(SourceCricinfo, &CricinfoDataProvider{...})
   factory.RegisterProvider(SourceSkysports, &SkysportsDataProvider{...})

2. Create Service:
   service := NewDataProviderService(factory)

3. Use Common Methods:
   matchDetails, err := service.GetMatchDetails(SourceCricbuzz, "match123")

4. Use Source-Specific Methods:
   commentary, err := service.GetCommentary("match123", 5)
   stats, err := service.GetDetailedStats("match123")
   highlights, err := service.GetVideoHighlights("match123")

==========================================
EXTENSIBILITY
==========================================

- To add a new source:
  1. Create new SourceType constant
  2. Create new provider interface extending CommonDataProvider
  3. Create struct implementing the new interface
  4. Register in factory

- To add new common methods:
  1. Add method to CommonDataProvider interface
  2. Implement in all provider structs
  3. Add routing method in DataProviderService

- To add source-specific methods:
  1. Add method to specific provider interface
  2. Implement in corresponding struct
  3. Add routing method in DataProviderService with type assertion

*/

