package main

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

/*
Requirements:
- adEvent can be of type impression or click
- get all the adEvent of id ad_id sorted in created date order
- find for an ad_id if there are X impressions in last Y days without a click.

Extensibility Features:
- Interface-based design for storage, ID generation, and queries
- Event handlers/hooks for custom behavior
- Configurable options
- Strategy pattern for queries
*/

// ==================== Interfaces for Extensibility ====================

// IDGenerator generates unique IDs for events
type IDGenerator interface {
	GenerateID() int64
}

// EventStorage defines the interface for storing and retrieving events
type EventStorage interface {
	// AddEvent adds an event to storage
	AddEvent(adID int64, event *Event) error
	// GetEvents returns all events for an ad ID sorted by timestamp
	GetEvents(adID int64) ([]*Event, error)
	// GetEventsInRange returns events within a time range
	GetEventsInRange(adID int64, from, to time.Time) ([]*Event, error)
	// RegisterAd registers a new ad ID
	RegisterAd(adID int64) error
	// Exists checks if an ad ID exists
	Exists(adID int64) bool
}

// QueryStrategy defines the interface for custom query strategies
type QueryStrategy interface {
	// Execute executes the query and returns results
	Execute(storage EventStorage, adID int64, params interface{}) (interface{}, error)
}

// EventHandler is called when events are added
type EventHandler func(adID int64, event *Event) error

// TimeSource provides the current time (useful for testing)
type TimeSource interface {
	Now() time.Time
}

// ==================== Configuration ====================

// ManagerConfig holds configuration for AdEventManager
type ManagerConfig struct {
	IDGenerator    IDGenerator
	Storage        EventStorage
	TimeSource     TimeSource
	EventHandlers  []EventHandler
	AutoRegister   bool // Auto-register ad IDs on first event
}

// DefaultConfig returns a default configuration
func DefaultConfig() *ManagerConfig {
	return &ManagerConfig{
		IDGenerator:   NewAtomicIDGenerator(),
		Storage:       NewInMemoryStorage(),
		TimeSource:    NewSystemTimeSource(),
		EventHandlers: []EventHandler{},
		AutoRegister:  true,
	}
}

// ==================== Event Types ====================

// EventType represents the type of ad event
type EventType int

const (
	IMPRESSION EventType = iota
	CLICK
)

// String returns the string representation of EventType
func (e EventType) String() string {
	switch e {
	case IMPRESSION:
		return "IMPRESSION"
	case CLICK:
		return "CLICK"
	default:
		return "UNKNOWN"
	}
}

// EventTypeRegistry allows registering custom event types
type EventTypeRegistry struct {
	types map[string]EventType
	mu    sync.RWMutex
}

// NewEventTypeRegistry creates a new event type registry
func NewEventTypeRegistry() *EventTypeRegistry {
	reg := &EventTypeRegistry{
		types: make(map[string]EventType),
	}
	// Register default types
	reg.Register("IMPRESSION", IMPRESSION)
	reg.Register("CLICK", CLICK)
	return reg
}

// Register registers a new event type
func (r *EventTypeRegistry) Register(name string, eventType EventType) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.types[name] = eventType
}

// Get retrieves an event type by name
func (r *EventTypeRegistry) Get(name string) (EventType, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	eventType, exists := r.types[name]
	return eventType, exists
}

// ==================== Event ====================

// Event represents a single ad event (impression or click)
type Event struct {
	ID        int64
	AdID      int64
	Type      EventType
	Timestamp time.Time
}

// NewEvent creates a new Event with the given parameters
func NewEvent(adID int64, eventType EventType, timestamp time.Time, idGen IDGenerator) *Event {
	return &Event{
		ID:        idGen.GenerateID(),
		AdID:      adID,
		Type:      eventType,
		Timestamp: timestamp,
	}
}

// ==================== Storage Implementations ====================

// InMemoryStorage is an in-memory implementation of EventStorage
type InMemoryStorage struct {
	adEvents map[int64]*AdEventList
	mu       sync.RWMutex
}

// NewInMemoryStorage creates a new in-memory storage
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		adEvents: make(map[int64]*AdEventList),
	}
}

func (s *InMemoryStorage) RegisterAd(adID int64) error {
	if adID <= 0 {
		return errors.New("adID must be positive")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.adEvents[adID]; exists {
		return fmt.Errorf("adID %d already registered", adID)
	}

	s.adEvents[adID] = NewAdEventList(adID)
	return nil
}

func (s *InMemoryStorage) Exists(adID int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.adEvents[adID]
	return exists
}

func (s *InMemoryStorage) AddEvent(adID int64, event *Event) error {
	s.mu.RLock()
	adEventList, exists := s.adEvents[adID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("adID %d not found", adID)
	}

	adEventList.insertEvent(event)
	return nil
}

func (s *InMemoryStorage) GetEvents(adID int64) ([]*Event, error) {
	s.mu.RLock()
	adEventList, exists := s.adEvents[adID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("adID %d not found", adID)
	}

	return adEventList.getEvents(), nil
}

func (s *InMemoryStorage) GetEventsInRange(adID int64, from, to time.Time) ([]*Event, error) {
	s.mu.RLock()
	adEventList, exists := s.adEvents[adID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("adID %d not found", adID)
	}

	return adEventList.getEventsInRange(from, to), nil
}

// AdEventList manages events for a single ad ID
type AdEventList struct {
	AdID      int64
	events    []*Event // sorted by timestamp (ascending)
	mu        sync.RWMutex
	lastClick *Event // optimization: track last click for faster queries
}

// NewAdEventList creates a new AdEventList for the given ad ID
func NewAdEventList(adID int64) *AdEventList {
	return &AdEventList{
		AdID:   adID,
		events: make([]*Event, 0),
	}
}

// insertEvent inserts an event in sorted order by timestamp
func (ael *AdEventList) insertEvent(event *Event) {
	ael.mu.Lock()
	defer ael.mu.Unlock()

	// Find insertion point using binary search
	insertPos := sort.Search(len(ael.events), func(i int) bool {
		return ael.events[i].Timestamp.After(event.Timestamp) ||
			ael.events[i].Timestamp.Equal(event.Timestamp)
	})

	// Insert at the found position
	ael.events = append(ael.events, nil)
	copy(ael.events[insertPos+1:], ael.events[insertPos:])
	ael.events[insertPos] = event

	// Update last click if this is a click
	if event.Type == CLICK {
		if ael.lastClick == nil || event.Timestamp.After(ael.lastClick.Timestamp) {
			ael.lastClick = event
		}
	}
}

// getEvents returns all events sorted by timestamp (ascending)
func (ael *AdEventList) getEvents() []*Event {
	ael.mu.RLock()
	defer ael.mu.RUnlock()

	result := make([]*Event, len(ael.events))
	copy(result, ael.events)
	return result
}

// getEventsInRange returns events within the specified time range
func (ael *AdEventList) getEventsInRange(from, to time.Time) []*Event {
	ael.mu.RLock()
	defer ael.mu.RUnlock()

	// Binary search for start position
	startIdx := sort.Search(len(ael.events), func(i int) bool {
		return !ael.events[i].Timestamp.Before(from)
	})

	// Binary search for end position
	endIdx := sort.Search(len(ael.events), func(i int) bool {
		return ael.events[i].Timestamp.After(to)
	})

	if startIdx >= endIdx {
		return []*Event{}
	}

	result := make([]*Event, endIdx-startIdx)
	copy(result, ael.events[startIdx:endIdx])
	return result
}

// ==================== ID Generator Implementations ====================

// AtomicIDGenerator generates IDs using atomic operations
type AtomicIDGenerator struct {
	counter int64
}

// NewAtomicIDGenerator creates a new atomic ID generator
func NewAtomicIDGenerator() *AtomicIDGenerator {
	return &AtomicIDGenerator{}
}

func (g *AtomicIDGenerator) GenerateID() int64 {
	return atomic.AddInt64(&g.counter, 1)
}

// ==================== Time Source Implementations ====================

// SystemTimeSource uses the system time
type SystemTimeSource struct{}

// NewSystemTimeSource creates a new system time source
func NewSystemTimeSource() *SystemTimeSource {
	return &SystemTimeSource{}
}

func (s *SystemTimeSource) Now() time.Time {
	return time.Now()
}

// MockTimeSource allows setting custom time for testing
type MockTimeSource struct {
	currentTime time.Time
	mu          sync.RWMutex
}

// NewMockTimeSource creates a new mock time source
func NewMockTimeSource(initialTime time.Time) *MockTimeSource {
	return &MockTimeSource{currentTime: initialTime}
}

func (m *MockTimeSource) Now() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentTime
}

// SetTime sets the current time
func (m *MockTimeSource) SetTime(t time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentTime = t
}

// ==================== Query Strategies ====================

// ImpressionsWithoutClickQueryParams holds parameters for the query
type ImpressionsWithoutClickQueryParams struct {
	X          int
	WithinDays int
}

// ImpressionsWithoutClickQuery implements the query for X impressions without click
type ImpressionsWithoutClickQuery struct{}

// Execute executes the impressions without click query
func (q *ImpressionsWithoutClickQuery) Execute(storage EventStorage, adID int64, params interface{}) (interface{}, error) {
	p, ok := params.(ImpressionsWithoutClickQueryParams)
	if !ok {
		return false, errors.New("invalid query parameters")
	}

	if p.X <= 0 || p.WithinDays <= 0 {
		return false, errors.New("x and withinDays must be positive")
	}

	// Calculate cutoff time
	cutoffTime := time.Now().AddDate(0, 0, -p.WithinDays)

	// Get events in the time range
	events, err := storage.GetEventsInRange(adID, cutoffTime, time.Now())
	if err != nil {
		return false, err
	}

	if len(events) == 0 {
		return false, nil
	}

	// Scan from oldest to newest (events are sorted ascending)
	// Count consecutive impressions, reset on any click
	impressionCount := 0
	for _, event := range events {
		if event.Type == CLICK {
			// Reset counter on any click
			impressionCount = 0
		} else if event.Type == IMPRESSION {
			impressionCount++
			if impressionCount >= p.X {
				return true, nil
			}
		}
	}

	return false, nil
}

// ==================== AdEventManager ====================

// AdEventManager manages ad events with extensible design
type AdEventManager struct {
	config *ManagerConfig
}

// NewAdEventManager creates a new AdEventManager with default config
func NewAdEventManager() *AdEventManager {
	return NewAdEventManagerWithConfig(DefaultConfig())
}

// NewAdEventManagerWithConfig creates a new AdEventManager with custom config
func NewAdEventManagerWithConfig(config *ManagerConfig) *AdEventManager {
	return &AdEventManager{
		config: config,
	}
}

// RegisterAd registers a new ad ID
func (aem *AdEventManager) RegisterAd(adID int64) error {
	return aem.config.Storage.RegisterAd(adID)
}

// AddEvent adds an event for the specified ad ID
func (aem *AdEventManager) AddEvent(adID int64, eventType EventType) error {
	if adID <= 0 {
		return errors.New("adID must be positive")
	}

	// Auto-register if enabled and ad doesn't exist
	if aem.config.AutoRegister && !aem.config.Storage.Exists(adID) {
		if err := aem.config.Storage.RegisterAd(adID); err != nil {
			return err
		}
	}

	// Create new event
	event := NewEvent(adID, eventType, aem.config.TimeSource.Now(), aem.config.IDGenerator)

	// Add to storage
	if err := aem.config.Storage.AddEvent(adID, event); err != nil {
		return err
	}

	// Call event handlers
	for _, handler := range aem.config.EventHandlers {
		if err := handler(adID, event); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: event handler error: %v\n", err)
		}
	}

	return nil
}

// GetAllEvents returns all events for an ad ID sorted by creation date (ascending)
func (aem *AdEventManager) GetAllEvents(adID int64) ([]*Event, error) {
	if adID <= 0 {
		return nil, errors.New("adID must be positive")
	}
	return aem.config.Storage.GetEvents(adID)
}

// GetEventsInRange returns events for an ad ID within the specified time range
func (aem *AdEventManager) GetEventsInRange(adID int64, from, to time.Time) ([]*Event, error) {
	if adID <= 0 {
		return nil, errors.New("adID must be positive")
	}
	if from.After(to) {
		return nil, errors.New("from time must be before to time")
	}
	return aem.config.Storage.GetEventsInRange(adID, from, to)
}

// HasXImpressionsWithoutClick checks if there are X impressions in the last Y days without a click
func (aem *AdEventManager) HasXImpressionsWithoutClick(adID int64, x int, withinDays int) (bool, error) {
	query := &ImpressionsWithoutClickQuery{}
	params := ImpressionsWithoutClickQueryParams{
		X:          x,
		WithinDays: withinDays,
	}
	result, err := query.Execute(aem.config.Storage, adID, params)
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

// ExecuteQuery executes a custom query strategy
func (aem *AdEventManager) ExecuteQuery(adID int64, strategy QueryStrategy, params interface{}) (interface{}, error) {
	return strategy.Execute(aem.config.Storage, adID, params)
}

// AddEventHandler adds an event handler that will be called when events are added
func (aem *AdEventManager) AddEventHandler(handler EventHandler) {
	aem.config.EventHandlers = append(aem.config.EventHandlers, handler)
}

func main() {
	// Example 1: Default usage
	fmt.Println("=== Example 1: Default Usage ===")
	manager := NewAdEventManager()

	// Register an ad ID
	adID := manager.config.IDGenerator.GenerateID()
	if err := manager.RegisterAd(adID); err != nil {
		fmt.Printf("Error registering ad: %v\n", err)
		return
	}

	fmt.Printf("Registered ad ID: %d\n", adID)

	// Add some events
	manager.AddEvent(adID, IMPRESSION)
	manager.AddEvent(adID, IMPRESSION)
	time.Sleep(1 * time.Second)
	manager.AddEvent(adID, CLICK)
	manager.AddEvent(adID, IMPRESSION)
	manager.AddEvent(adID, IMPRESSION)
	time.Sleep(1 * time.Second)
	manager.AddEvent(adID, IMPRESSION)
	manager.AddEvent(adID, IMPRESSION)
	manager.AddEvent(adID, IMPRESSION)
	time.Sleep(1 * time.Second)
	manager.AddEvent(adID, IMPRESSION)

	// Get all events
	events, err := manager.GetAllEvents(adID)
	if err != nil {
		fmt.Printf("Error getting events: %v\n", err)
		return
	}

	fmt.Println("\nAll events (sorted by timestamp):")
	for _, event := range events {
		fmt.Printf("  Event ID: %d, Type: %s, Time: %s\n",
			event.ID, event.Type.String(), event.Timestamp.Format(time.RFC3339))
	}

	// Check for 2 impressions without click in last 2 days
	hasImpressions, err := manager.HasXImpressionsWithoutClick(adID, 2, 2)
	if err != nil {
		fmt.Printf("Error checking impressions: %v\n", err)
		return
	}

	fmt.Printf("\nHas 2 impressions without click in last 2 days: %v\n", hasImpressions)

	// Example 2: Custom configuration with event handler
	fmt.Println("\n=== Example 2: Custom Configuration ===")
	customConfig := &ManagerConfig{
		IDGenerator:   NewAtomicIDGenerator(),
		Storage:       NewInMemoryStorage(),
		TimeSource:    NewSystemTimeSource(),
		EventHandlers: []EventHandler{},
		AutoRegister:  true,
	}

	customManager := NewAdEventManagerWithConfig(customConfig)

	// Add an event handler
	customManager.AddEventHandler(func(adID int64, event *Event) error {
		fmt.Printf("  [Handler] Event added: AdID=%d, Type=%s, ID=%d\n", adID, event.Type.String(), event.ID)
		return nil
	})

	adID2 := customManager.config.IDGenerator.GenerateID()
	customManager.AddEvent(adID2, IMPRESSION)
	customManager.AddEvent(adID2, CLICK)

	// Example 3: Using custom query strategy
	fmt.Println("\n=== Example 3: Custom Query Strategy ===")
	query := &ImpressionsWithoutClickQuery{}
	params := ImpressionsWithoutClickQueryParams{
		X:          3,
		WithinDays: 2,
	}
	result, err := customManager.ExecuteQuery(adID, query, params)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
	} else {
		fmt.Printf("Query result: %v\n", result)
	}
}
