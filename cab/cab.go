package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// ============================================================================
// ENTITIES
// ============================================================================

// VehicleType represents the type of vehicle
type VehicleType string

const (
	VehicleTypeCar  VehicleType = "CAR"
	VehicleTypeBike VehicleType = "BIKE"
	VehicleTypeSUV  VehicleType = "SUV"
	VehicleTypeAuto VehicleType = "AUTO"
)

// Location represents a geographical location
type Location struct {
	Latitude  float64
	Longitude float64
	Address   string
	Timestamp time.Time
}

// CalculateDistance calculates distance between two locations using Haversine formula
// Returns distance in kilometers
func CalculateDistance(loc1, loc2 *Location) float64 {
	const earthRadiusKm = 6371.0

	lat1Rad := loc1.Latitude * math.Pi / 180
	lat2Rad := loc2.Latitude * math.Pi / 180
	deltaLat := (loc2.Latitude - loc1.Latitude) * math.Pi / 180
	deltaLon := (loc2.Longitude - loc1.Longitude) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// Vehicle represents a vehicle
type Vehicle struct {
	ID          string
	Type        VehicleType
	LicensePlate string
	Model       string
	Capacity    int
	Year        int
}

// RideState represents the state of a ride
type RideState string

const (
	RideStateRequested RideState = "REQUESTED"
	RideStateAccepted  RideState = "ACCEPTED"
	RideStateStarted   RideState = "STARTED"
	RideStateCompleted RideState = "COMPLETED"
	RideStateCancelled RideState = "CANCELLED"
)

// Fare represents fare calculation
type Fare struct {
	BaseFare      float64
	DistanceFare  float64
	TimeFare      float64
	SurgeMultiplier float64
	TotalFare     float64
}

// Ride represents a ride in the system
type Ride struct {
	ID          string
	RiderID     string
	DriverID    string
	PickupPoint *Location
	DropPoint   *Location
	VehicleType VehicleType
	Fare        *Fare
	State       RideState
	OTP         string
	StartTime   time.Time
	EndTime     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Driver represents a driver in the system
type Driver struct {
	ID            string
	Name          string
	Phone         string
	Email         string
	Location      *Location
	Vehicle       *Vehicle
	IsAvailable   bool
	CurrentRideID string
	Rating        float64
	TotalRides    int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// HasActiveRide checks if driver has an active ride
func (d *Driver) HasActiveRide() bool {
	return d.CurrentRideID != ""
}

// IsWithinRadius checks if driver is within radius of a location
func (d *Driver) IsWithinRadius(location *Location, radiusKm float64) bool {
	if d.Location == nil {
		return false
	}
	distance := CalculateDistance(d.Location, location)
	return distance <= radiusKm
}

// Rider represents a rider in the system
type Rider struct {
	ID            string
	Name          string
	Phone         string
	Email         string
	CurrentRideID string
	Rating        float64
	TotalRides    int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// HasActiveRide checks if rider has an active ride
func (r *Rider) HasActiveRide() bool {
	return r.CurrentRideID != ""
}

// ============================================================================
// REPOSITORIES (Data Access Layer)
// ============================================================================

// RideRepository interface for ride data access
type RideRepository interface {
	Create(ride *Ride) error
	GetByID(rideID string) (*Ride, error)
	Update(ride *Ride) error
	GetByRiderID(riderID string) ([]*Ride, error)
	GetByDriverID(driverID string) ([]*Ride, error)
}

// DriverRepository interface for driver data access
type DriverRepository interface {
	Create(driver *Driver) error
	GetByID(driverID string) (*Driver, error)
	Update(driver *Driver) error
	GetAvailableDrivers() ([]*Driver, error)
	GetDriversNearLocation(location *Location, radiusKm float64) ([]*Driver, error)
}

// RiderRepository interface for rider data access
type RiderRepository interface {
	Create(rider *Rider) error
	GetByID(riderID string) (*Rider, error)
	Update(rider *Rider) error
}

// In-memory implementations (for demonstration)
type inMemoryRideRepository struct {
	rides map[string]*Ride
}

func NewInMemoryRideRepository() RideRepository {
	return &inMemoryRideRepository{
		rides: make(map[string]*Ride),
	}
}

func (r *inMemoryRideRepository) Create(ride *Ride) error {
	r.rides[ride.ID] = ride
	return nil
}

func (r *inMemoryRideRepository) GetByID(rideID string) (*Ride, error) {
	ride, exists := r.rides[rideID]
	if !exists {
		return nil, errors.New("ride not found")
	}
	return ride, nil
}

func (r *inMemoryRideRepository) Update(ride *Ride) error {
	if _, exists := r.rides[ride.ID]; !exists {
		return errors.New("ride not found")
	}
	ride.UpdatedAt = time.Now()
	r.rides[ride.ID] = ride
	return nil
}

func (r *inMemoryRideRepository) GetByRiderID(riderID string) ([]*Ride, error) {
	var rides []*Ride
	for _, ride := range r.rides {
		if ride.RiderID == riderID {
			rides = append(rides, ride)
		}
	}
	return rides, nil
}

func (r *inMemoryRideRepository) GetByDriverID(driverID string) ([]*Ride, error) {
	var rides []*Ride
	for _, ride := range r.rides {
		if ride.DriverID == driverID {
			rides = append(rides, ride)
		}
	}
	return rides, nil
}

type inMemoryDriverRepository struct {
	drivers map[string]*Driver
}

func NewInMemoryDriverRepository() DriverRepository {
	return &inMemoryDriverRepository{
		drivers: make(map[string]*Driver),
	}
}

func (r *inMemoryDriverRepository) Create(driver *Driver) error {
	r.drivers[driver.ID] = driver
	return nil
}

func (r *inMemoryDriverRepository) GetByID(driverID string) (*Driver, error) {
	driver, exists := r.drivers[driverID]
	if !exists {
		return nil, errors.New("driver not found")
	}
	return driver, nil
}

func (r *inMemoryDriverRepository) Update(driver *Driver) error {
	if _, exists := r.drivers[driver.ID]; !exists {
		return errors.New("driver not found")
	}
	driver.UpdatedAt = time.Now()
	r.drivers[driver.ID] = driver
	return nil
}

func (r *inMemoryDriverRepository) GetAvailableDrivers() ([]*Driver, error) {
	var drivers []*Driver
	for _, driver := range r.drivers {
		if driver.IsAvailable && !driver.HasActiveRide() {
			drivers = append(drivers, driver)
		}
	}
	return drivers, nil
}

func (r *inMemoryDriverRepository) GetDriversNearLocation(location *Location, radiusKm float64) ([]*Driver, error) {
	var drivers []*Driver
	for _, driver := range r.drivers {
		if driver.IsAvailable && !driver.HasActiveRide() && driver.IsWithinRadius(location, radiusKm) {
			drivers = append(drivers, driver)
		}
	}
	return drivers, nil
}

type inMemoryRiderRepository struct {
	riders map[string]*Rider
}

func NewInMemoryRiderRepository() RiderRepository {
	return &inMemoryRiderRepository{
		riders: make(map[string]*Rider),
	}
}

func (r *inMemoryRiderRepository) Create(rider *Rider) error {
	r.riders[rider.ID] = rider
	return nil
}

func (r *inMemoryRiderRepository) GetByID(riderID string) (*Rider, error) {
	rider, exists := r.riders[riderID]
	if !exists {
		return nil, errors.New("rider not found")
	}
	return rider, nil
}

func (r *inMemoryRiderRepository) Update(rider *Rider) error {
	if _, exists := r.riders[rider.ID]; !exists {
		return errors.New("rider not found")
	}
	rider.UpdatedAt = time.Now()
	r.riders[rider.ID] = rider
	return nil
}

// ============================================================================
// SERVICES (Business Logic Layer)
// ============================================================================

// FareCalculator interface (Strategy Pattern)
type FareCalculator interface {
	CalculateFare(pickup, drop *Location, vehicleType VehicleType) (*Fare, error)
	CalculateFinalFare(ride *Ride) (*Fare, error)
}

// MatchingStrategy interface (Strategy Pattern)
type MatchingStrategy interface {
	FindDrivers(pickup *Location, vehicleType VehicleType, radiusKm float64, driverRepo DriverRepository) ([]*Driver, error)
}

// FareService implements fare calculation
type fareService struct {
	baseFares map[VehicleType]float64
}

func NewFareService() FareCalculator {
	return &fareService{
		baseFares: map[VehicleType]float64{
			VehicleTypeCar:  50.0,
			VehicleTypeBike: 30.0,
			VehicleTypeSUV:  80.0,
			VehicleTypeAuto: 40.0,
		},
	}
}

func (s *fareService) CalculateFare(pickup, drop *Location, vehicleType VehicleType) (*Fare, error) {
	distance := CalculateDistance(pickup, drop)
	baseFare := s.baseFares[vehicleType]
	distanceFare := distance * 10.0 // 10 per km
	timeFare := 0.0 // Simplified, would calculate based on estimated time
	surgeMultiplier := 1.0 // Simplified, would calculate based on demand

	totalFare := (baseFare + distanceFare + timeFare) * surgeMultiplier

	return &Fare{
		BaseFare:       baseFare,
		DistanceFare:   distanceFare,
		TimeFare:       timeFare,
		SurgeMultiplier: surgeMultiplier,
		TotalFare:      totalFare,
	}, nil
}

func (s *fareService) CalculateFinalFare(ride *Ride) (*Fare, error) {
	if ride.StartTime.IsZero() || ride.EndTime.IsZero() {
		return nil, errors.New("ride not completed")
	}

	duration := ride.EndTime.Sub(ride.StartTime)
	timeFare := duration.Minutes() * 1.0 // 1 per minute
	distanceFare := CalculateDistance(ride.PickupPoint, ride.DropPoint) * 10.0
	baseFare := s.baseFares[ride.VehicleType]

	totalFare := baseFare + distanceFare + timeFare

	return &Fare{
		BaseFare:       baseFare,
		DistanceFare:   distanceFare,
		TimeFare:       timeFare,
		SurgeMultiplier: 1.0,
		TotalFare:      totalFare,
	}, nil
}

// NearestDriverStrategy implements nearest driver matching
type nearestDriverStrategy struct{}

func NewNearestDriverStrategy() MatchingStrategy {
	return &nearestDriverStrategy{}
}

func (s *nearestDriverStrategy) FindDrivers(pickup *Location, vehicleType VehicleType, radiusKm float64, driverRepo DriverRepository) ([]*Driver, error) {
	drivers, err := driverRepo.GetDriversNearLocation(pickup, radiusKm)
	if err != nil {
		return nil, err
	}

	// Filter by vehicle type and sort by distance
	var matchingDrivers []*Driver
	for _, driver := range drivers {
		if driver.Vehicle != nil && driver.Vehicle.Type == vehicleType {
			matchingDrivers = append(matchingDrivers, driver)
		}
	}

	// Sort by distance (simplified - would use proper sorting)
	return matchingDrivers, nil
}

// MatchingService finds available drivers
type MatchingService interface {
	FindAvailableDrivers(pickup *Location, vehicleType VehicleType, radiusKm float64) ([]*Driver, error)
}

type matchingService struct {
	driverRepo      DriverRepository
	matchingStrategy MatchingStrategy
}

func NewMatchingService(driverRepo DriverRepository, strategy MatchingStrategy) MatchingService {
	return &matchingService{
		driverRepo:      driverRepo,
		matchingStrategy: strategy,
	}
}

func (s *matchingService) FindAvailableDrivers(pickup *Location, vehicleType VehicleType, radiusKm float64) ([]*Driver, error) {
	return s.matchingStrategy.FindDrivers(pickup, vehicleType, radiusKm, s.driverRepo)
}

// NotificationService handles notifications
type NotificationService interface {
	NotifyRider(riderID string, message string, data map[string]interface{}) error
	NotifyDriver(driverID string, message string, data map[string]interface{}) error
}

type notificationService struct{}

func NewNotificationService() NotificationService {
	return &notificationService{}
}

func (s *notificationService) NotifyRider(riderID string, message string, data map[string]interface{}) error {
	fmt.Printf("[NOTIFICATION] Rider %s: %s - %+v\n", riderID, message, data)
	return nil
}

func (s *notificationService) NotifyDriver(driverID string, message string, data map[string]interface{}) error {
	fmt.Printf("[NOTIFICATION] Driver %s: %s - %+v\n", driverID, message, data)
	return nil
}

// LocationService tracks locations
type LocationService interface {
	TrackLocation(userID string) error
	StopTracking(userID string) error
	GetLocation(userID string) (*Location, error)
}

type locationService struct {
	locations map[string]*Location
}

func NewLocationService() LocationService {
	return &locationService{
		locations: make(map[string]*Location),
	}
}

func (s *locationService) TrackLocation(userID string) error {
	fmt.Printf("[LOCATION] Started tracking location for user: %s\n", userID)
	return nil
}

func (s *locationService) StopTracking(userID string) error {
	fmt.Printf("[LOCATION] Stopped tracking location for user: %s\n", userID)
	delete(s.locations, userID)
	return nil
}

func (s *locationService) GetLocation(userID string) (*Location, error) {
	location, exists := s.locations[userID]
	if !exists {
		return nil, errors.New("location not found")
	}
	return location, nil
}

// AuthenticationService handles OTP generation and validation
type AuthenticationService interface {
	GenerateOTP(rideID string) (string, error)
	ValidateOTP(rideID string, otp string) (bool, error)
}

type authenticationService struct {
	otps map[string]string
}

func NewAuthenticationService() AuthenticationService {
	return &authenticationService{
		otps: make(map[string]string),
	}
}

func (s *authenticationService) GenerateOTP(rideID string) (string, error) {
	otp := fmt.Sprintf("%04d", rand.Intn(10000))
	s.otps[rideID] = otp
	fmt.Printf("[AUTH] Generated OTP for ride %s: %s\n", rideID, otp)
	return otp, nil
}

func (s *authenticationService) ValidateOTP(rideID string, otp string) (bool, error) {
	storedOTP, exists := s.otps[rideID]
	if !exists {
		return false, errors.New("OTP not found")
	}
	isValid := storedOTP == otp
	if isValid {
		delete(s.otps, rideID)
	}
	return isValid, nil
}

// PaymentService processes payments
type PaymentService interface {
	ProcessPayment(riderID string, amount float64) error
}

type paymentService struct{}

func NewPaymentService() PaymentService {
	return &paymentService{}
}

func (s *paymentService) ProcessPayment(riderID string, amount float64) error {
	fmt.Printf("[PAYMENT] Processing payment of %.2f for rider %s\n", amount, riderID)
	// Simulate payment processing
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("[PAYMENT] Payment successful for rider %s\n", riderID)
	return nil
}

// ============================================================================
// ORCHESTRATOR SERVICE (RideService)
// ============================================================================

var (
	ErrRiderHasActiveRide = errors.New("rider already has an active ride")
	ErrNoDriversAvailable = errors.New("no drivers available")
	ErrInvalidRideState   = errors.New("invalid ride state")
	ErrInvalidOTP         = errors.New("invalid OTP")
)

// RideService orchestrates all services for ride management
type RideService interface {
	RequestRide(riderID string, pickup, drop *Location, vehicleType VehicleType) (*Ride, error)
	AcceptRide(driverID string, rideID string) error
	StartRide(rideID string, otp string) error
	EndRide(rideID string) error
	CancelRide(rideID string, userID string) error
	GetRide(rideID string) (*Ride, error)
}

type rideService struct {
	// Repositories
	rideRepo  RideRepository
	riderRepo RiderRepository
	driverRepo DriverRepository

	// Services (dependencies)
	fareCalculator     FareCalculator
	matchingService    MatchingService
	notificationService NotificationService
	locationService    LocationService
	authService        AuthenticationService
	paymentService     PaymentService
}

// NewRideService creates a new RideService with all dependencies
func NewRideService(
	rideRepo RideRepository,
	riderRepo RiderRepository,
	driverRepo DriverRepository,
	fareCalculator FareCalculator,
	matchingService MatchingService,
	notificationService NotificationService,
	locationService LocationService,
	authService AuthenticationService,
	paymentService PaymentService,
) RideService {
	return &rideService{
		rideRepo:          rideRepo,
		riderRepo:         riderRepo,
		driverRepo:        driverRepo,
		fareCalculator:    fareCalculator,
		matchingService:   matchingService,
		notificationService: notificationService,
		locationService:   locationService,
		authService:       authService,
		paymentService:    paymentService,
	}
}

// generateRideID generates a unique ride ID
func generateRideID() string {
	return fmt.Sprintf("RIDE-%d", time.Now().UnixNano())
}

// RequestRide implements the request ride flow
func (s *rideService) RequestRide(
	riderID string,
	pickup, drop *Location,
	vehicleType VehicleType,
) (*Ride, error) {
	// Step 1: Validate rider
	rider, err := s.riderRepo.GetByID(riderID)
	if err != nil {
		return nil, fmt.Errorf("rider not found: %w", err)
	}

	if rider.HasActiveRide() {
		return nil, ErrRiderHasActiveRide
	}

	// Step 2: Calculate fare estimate
	fare, err := s.fareCalculator.CalculateFare(pickup, drop, vehicleType)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate fare: %w", err)
	}

	// Step 3: Find available drivers
	drivers, err := s.matchingService.FindAvailableDrivers(pickup, vehicleType, 5.0)
	if err != nil {
		return nil, fmt.Errorf("failed to find drivers: %w", err)
	}

	if len(drivers) == 0 {
		return nil, ErrNoDriversAvailable
	}

	// Step 4: Create ride record
	ride := &Ride{
		ID:          generateRideID(),
		RiderID:     riderID,
		PickupPoint: pickup,
		DropPoint:   drop,
		VehicleType: vehicleType,
		Fare:        fare,
		State:       RideStateRequested,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.rideRepo.Create(ride)
	if err != nil {
		return nil, fmt.Errorf("failed to create ride: %w", err)
	}

	// Step 5: Update rider's current ride
	rider.CurrentRideID = ride.ID
	s.riderRepo.Update(rider)

	// Step 6: Notify drivers
	for i, driver := range drivers {
		if i >= 3 { // Top 3 drivers
			break
		}
		s.notificationService.NotifyDriver(
			driver.ID,
			"New ride request",
			map[string]interface{}{
				"rideID":  ride.ID,
				"pickup":  pickup.Address,
				"drop":    drop.Address,
				"fare":    fare.TotalFare,
			},
		)
	}

	return ride, nil
}

// AcceptRide implements the accept ride flow
func (s *rideService) AcceptRide(driverID string, rideID string) error {
	// Step 1: Get ride
	ride, err := s.rideRepo.GetByID(rideID)
	if err != nil {
		return err
	}

	if ride.State != RideStateRequested {
		return ErrInvalidRideState
	}

	// Step 2: Validate driver
	driver, err := s.driverRepo.GetByID(driverID)
	if err != nil {
		return fmt.Errorf("driver not found: %w", err)
	}

	if !driver.IsAvailable || driver.HasActiveRide() {
		return errors.New("driver not available")
	}

	// Step 3: Update ride with driver
	ride.DriverID = driverID
	ride.State = RideStateAccepted

	// Step 4: Generate OTP
	otp, err := s.authService.GenerateOTP(rideID)
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}
	ride.OTP = otp

	// Step 5: Update driver status
	driver.CurrentRideID = rideID
	driver.IsAvailable = false
	s.driverRepo.Update(driver)

	// Step 6: Save ride
	err = s.rideRepo.Update(ride)
	if err != nil {
		return err
	}

	// Step 7: Notify rider
	s.notificationService.NotifyRider(
		ride.RiderID,
		"Driver assigned",
		map[string]interface{}{
			"driverID": driverID,
			"driverName": driver.Name,
			"otp":      otp,
		},
	)

	// Step 8: Start tracking driver location
	s.locationService.TrackLocation(driverID)

	return nil
}

// StartRide implements the start ride flow
func (s *rideService) StartRide(rideID string, otp string) error {
	// Step 1: Get ride
	ride, err := s.rideRepo.GetByID(rideID)
	if err != nil {
		return err
	}

	if ride.State != RideStateAccepted {
		return ErrInvalidRideState
	}

	// Step 2: Validate OTP
	isValid, err := s.authService.ValidateOTP(rideID, otp)
	if err != nil {
		return err
	}

	if !isValid {
		return ErrInvalidOTP
	}

	// Step 3: Update ride state
	ride.State = RideStateStarted
	ride.StartTime = time.Now()

	err = s.rideRepo.Update(ride)
	if err != nil {
		return err
	}

	// Step 4: Start tracking both locations
	s.locationService.TrackLocation(ride.DriverID)
	s.locationService.TrackLocation(ride.RiderID)

	// Step 5: Notify rider
	s.notificationService.NotifyRider(
		ride.RiderID,
		"Ride started",
		map[string]interface{}{
			"rideID": rideID,
		},
	)

	return nil
}

// EndRide implements the end ride flow
func (s *rideService) EndRide(rideID string) error {
	// Step 1: Get ride
	ride, err := s.rideRepo.GetByID(rideID)
	if err != nil {
		return err
	}

	if ride.State != RideStateStarted {
		return ErrInvalidRideState
	}

	// Step 2: Calculate final fare
	finalFare, err := s.fareCalculator.CalculateFinalFare(ride)
	if err != nil {
		return fmt.Errorf("failed to calculate final fare: %w", err)
	}
	ride.Fare = finalFare

	// Step 3: Update ride state
	ride.State = RideStateCompleted
	ride.EndTime = time.Now()

	err = s.rideRepo.Update(ride)
	if err != nil {
		return err
	}

	// Step 4: Process payment
	err = s.paymentService.ProcessPayment(ride.RiderID, finalFare.TotalFare)
	if err != nil {
		return fmt.Errorf("payment failed: %w", err)
	}

	// Step 5: Update driver and rider status
	driver, _ := s.driverRepo.GetByID(ride.DriverID)
	if driver != nil {
		driver.CurrentRideID = ""
		driver.IsAvailable = true
		driver.TotalRides++
		s.driverRepo.Update(driver)
	}

	rider, _ := s.riderRepo.GetByID(ride.RiderID)
	if rider != nil {
		rider.CurrentRideID = ""
		rider.TotalRides++
		s.riderRepo.Update(rider)
	}

	// Step 6: Stop location tracking
	s.locationService.StopTracking(ride.DriverID)
	s.locationService.StopTracking(ride.RiderID)

	// Step 7: Notify both parties
	s.notificationService.NotifyRider(
		ride.RiderID,
		"Ride completed",
		map[string]interface{}{
			"fare": finalFare.TotalFare,
		},
	)

	s.notificationService.NotifyDriver(
		ride.DriverID,
		"Ride completed",
		map[string]interface{}{
			"fare": finalFare.TotalFare,
		},
	)

	return nil
}

// CancelRide cancels a ride
func (s *rideService) CancelRide(rideID string, userID string) error {
	ride, err := s.rideRepo.GetByID(rideID)
	if err != nil {
		return err
	}

	if ride.State == RideStateCompleted || ride.State == RideStateCancelled {
		return ErrInvalidRideState
	}

	// Update ride state
	ride.State = RideStateCancelled
	ride.UpdatedAt = time.Now()
	s.rideRepo.Update(ride)

	// Update driver and rider status
	if ride.DriverID != "" {
		driver, _ := s.driverRepo.GetByID(ride.DriverID)
		if driver != nil {
			driver.CurrentRideID = ""
			driver.IsAvailable = true
			s.driverRepo.Update(driver)
		}
	}

	rider, _ := s.riderRepo.GetByID(ride.RiderID)
	if rider != nil {
		rider.CurrentRideID = ""
		s.riderRepo.Update(rider)
	}

	// Notify parties
	s.notificationService.NotifyRider(ride.RiderID, "Ride cancelled", nil)
	if ride.DriverID != "" {
		s.notificationService.NotifyDriver(ride.DriverID, "Ride cancelled", nil)
	}

	return nil
}

// GetRide retrieves a ride by ID
func (s *rideService) GetRide(rideID string) (*Ride, error) {
	return s.rideRepo.GetByID(rideID)
}

// ============================================================================
// EXAMPLE USAGE - Main Function
// ============================================================================

func main() {
	fmt.Println("=== Cab Booking System Demo ===\n")

	// Step 1: Initialize repositories
	rideRepo := NewInMemoryRideRepository()
	riderRepo := NewInMemoryRiderRepository()
	driverRepo := NewInMemoryDriverRepository()

	// Step 2: Create sample data
	// Create a rider
	rider := &Rider{
		ID:        "RIDER-001",
		Name:      "John Doe",
		Phone:     "+1234567890",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	riderRepo.Create(rider)

	// Create drivers
	driver1 := &Driver{
		ID:          "DRIVER-001",
		Name:        "Alice Smith",
		Phone:       "+1987654321",
		Email:       "alice@example.com",
		IsAvailable: true,
		Location: &Location{
			Latitude:  28.6139, // Delhi coordinates
			Longitude: 77.2090,
			Address:   "Connaught Place, New Delhi",
			Timestamp: time.Now(),
		},
		Vehicle: &Vehicle{
			ID:          "VEH-001",
			Type:        VehicleTypeCar,
			LicensePlate: "DL-01-AB-1234",
			Model:       "Honda City",
			Capacity:    4,
			Year:        2020,
		},
		Rating:     4.8,
		TotalRides: 150,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	driverRepo.Create(driver1)

	driver2 := &Driver{
		ID:          "DRIVER-002",
		Name:        "Bob Johnson",
		Phone:       "+1555555555",
		Email:       "bob@example.com",
		IsAvailable: true,
		Location: &Location{
			Latitude:  28.6140,
			Longitude: 77.2091,
			Address:   "Near Connaught Place, New Delhi",
			Timestamp: time.Now(),
		},
		Vehicle: &Vehicle{
			ID:          "VEH-002",
			Type:        VehicleTypeCar,
			LicensePlate: "DL-01-CD-5678",
			Model:       "Maruti Swift",
			Capacity:    4,
			Year:        2021,
		},
		Rating:     4.9,
		TotalRides: 200,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	driverRepo.Create(driver2)

	fmt.Println("✓ Created rider and drivers\n")

	// Step 3: Create services
	fareService := NewFareService()
	matchingStrategy := NewNearestDriverStrategy()
	matchingService := NewMatchingService(driverRepo, matchingStrategy)
	notificationService := NewNotificationService()
	locationService := NewLocationService()
	authService := NewAuthenticationService()
	paymentService := NewPaymentService()

	// Step 4: Create orchestrator service (RideService)
	rideService := NewRideService(
		rideRepo,
		riderRepo,
		driverRepo,
		fareService,
		matchingService,
		notificationService,
		locationService,
		authService,
		paymentService,
	)

	fmt.Println("✓ Initialized all services\n")

	// Step 5: Demonstrate complete ride flow
	fmt.Println("--- Step 1: Request Ride ---")
	pickup := &Location{
		Latitude:  28.6139,
		Longitude: 77.2090,
		Address:   "Connaught Place, New Delhi",
		Timestamp: time.Now(),
	}
	drop := &Location{
		Latitude:  28.5355,
		Longitude: 77.3910,
		Address:   "Gurgaon, Haryana",
		Timestamp: time.Now(),
	}

	ride, err := rideService.RequestRide(rider.ID, pickup, drop, VehicleTypeCar)
	if err != nil {
		fmt.Printf("❌ Error requesting ride: %v\n", err)
		return
	}
	fmt.Printf("✓ Ride requested successfully!\n")
	fmt.Printf("  Ride ID: %s\n", ride.ID)
	fmt.Printf("  State: %s\n", ride.State)
	fmt.Printf("  Estimated Fare: ₹%.2f\n\n", ride.Fare.TotalFare)

	// Step 6: Driver accepts ride
	fmt.Println("--- Step 2: Driver Accepts Ride ---")
	err = rideService.AcceptRide(driver1.ID, ride.ID)
	if err != nil {
		fmt.Printf("❌ Error accepting ride: %v\n", err)
		return
	}
	fmt.Printf("✓ Driver %s accepted the ride\n\n", driver1.Name)

	// Get updated ride to see OTP
	ride, _ = rideService.GetRide(ride.ID)
	fmt.Printf("  OTP: %s\n", ride.OTP)
	fmt.Printf("  State: %s\n\n", ride.State)

	// Step 7: Start ride
	fmt.Println("--- Step 3: Start Ride ---")
	err = rideService.StartRide(ride.ID, ride.OTP)
	if err != nil {
		fmt.Printf("❌ Error starting ride: %v\n", err)
		return
	}
	fmt.Printf("✓ Ride started successfully!\n")
	ride, _ = rideService.GetRide(ride.ID)
	fmt.Printf("  State: %s\n", ride.State)
	fmt.Printf("  Start Time: %s\n\n", ride.StartTime.Format(time.RFC3339))

	// Simulate ride duration
	fmt.Println("--- Ride in progress... ---")
	time.Sleep(2 * time.Second)

	// Step 8: End ride
	fmt.Println("\n--- Step 4: End Ride ---")
	err = rideService.EndRide(ride.ID)
	if err != nil {
		fmt.Printf("❌ Error ending ride: %v\n", err)
		return
	}
	fmt.Printf("✓ Ride completed successfully!\n")
	ride, _ = rideService.GetRide(ride.ID)
	fmt.Printf("  State: %s\n", ride.State)
	fmt.Printf("  Final Fare: ₹%.2f\n", ride.Fare.TotalFare)
	fmt.Printf("  End Time: %s\n\n", ride.EndTime.Format(time.RFC3339))

	fmt.Println("=== Demo Complete ===")
}
