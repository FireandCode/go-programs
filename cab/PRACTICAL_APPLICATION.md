# Practical Application: Uber-like Service Design

## Applying the Framework to Your Service

### Step 1: Refined Entity Design

```go
// Core Entities with proper structure

// User (Base entity)
type User struct {
    ID       string
    Name     string
    Phone    string
    Email    string
    Location *Location
}

// Rider extends User
type Rider struct {
    User
    Rating      float64
    ActiveRides []string // Ride IDs
}

// Driver extends User
type Driver struct {
    User
    Vehicle         *Vehicle
    Rating          float64
    IsAvailable     bool
    CurrentRideID   string
    TotalRides      int
}

// Ride (Core entity)
type Ride struct {
    ID            string
    RiderID       string
    DriverID      string
    VehicleType   VehicleType
    PickupPoint   *Location
    DropPoint     *Location
    State         RideState
    Fare          *Fare
    OTP           string
    StartTime     time.Time
    EndTime       time.Time
    CreatedAt     time.Time
}

// Vehicle
type Vehicle struct {
    ID          string
    Type        VehicleType
    LicensePlate string
    Model       string
    Capacity    int
}

// Location
type Location struct {
    Latitude  float64
    Longitude float64
    Address   string
    Timestamp time.Time
}

// Fare
type Fare struct {
    BaseFare      float64
    DistanceFare  float64
    TimeFare      float64
    SurgeMultiplier float64
    TotalFare    float64
    Currency     string
}
```

### Step 2: Service Interfaces (Abstraction)

```go
// Service interfaces for extensibility

type FareCalculator interface {
    CalculateFare(pickup, drop *Location, vehicleType VehicleType) (*Fare, error)
}

type MatchingStrategy interface {
    FindDriver(ride *Ride, availableDrivers []*Driver) (*Driver, error)
}

type LocationTracker interface {
    UpdateLocation(userID string, location *Location) error
    GetLocation(userID string) (*Location, error)
    TrackLocation(userID string) (<-chan *Location, error)
}

type NotificationService interface {
    NotifyRider(riderID string, message string) error
    NotifyDriver(driverID string, message string) error
}
```

### Step 3: Service Implementation Structure

```go
// RideService - Core orchestration
type RideService struct {
    rideRepo          RideRepository
    fareCalculator    FareCalculator
    matchingStrategy  MatchingStrategy
    locationTracker   LocationTracker
    notificationService NotificationService
    paymentService    PaymentService
}

// Methods:
// - RequestRide(riderID, pickup, drop, vehicleType) (*Ride, error)
// - AcceptRide(driverID, rideID) error
// - StartRide(rideID, otp) error
// - EndRide(rideID) error
// - CancelRide(rideID, userID) error

// FareService - Fare calculation
type FareService struct {
    strategies map[VehicleType]FareCalculator
    surgeService SurgePricingService
}

// MatchingService - Driver-rider matching
type MatchingService struct {
    strategy MatchingStrategy
    driverRepo DriverRepository
}

// LocationService - Real-time tracking
type LocationService struct {
    tracker LocationTracker
    cache   LocationCache
}

// NotificationService - Push notifications
type NotificationService struct {
    channels []NotificationChannel // SMS, Push, Email
}
```

### Step 4: Design Patterns Applied

#### Strategy Pattern for Fare Calculation

```go
type BaseFareStrategy struct {
    baseRate float64
    perKmRate float64
    perMinRate float64
}

type SurgeFareStrategy struct {
    baseStrategy FareCalculator
    multiplier   float64
}

type DiscountFareStrategy struct {
    baseStrategy FareCalculator
    discountPercent float64
}
```

#### Strategy Pattern for Matching

```go
type NearestDriverStrategy struct {
    maxDistance float64 // km
}

type RatingBasedStrategy struct {
    minRating float64
    maxDistance float64
}
```

#### State Pattern for Ride Lifecycle

```go
type RideState interface {
    Start(ride *Ride) error
    End(ride *Ride) error
    Cancel(ride *Ride) error
}

type RequestedState struct{}
type AcceptedState struct{}
type StartedState struct{}
type CompletedState struct{}
type CancelledState struct{}
```

### Step 5: Complete Flow Implementation

#### Request Ride Flow (Detailed)

```
1. RiderService.RequestRide()
   ├─ Validate rider (no active rides, account active)
   ├─ Create Ride entity (state: REQUESTED)
   ├─ FareService.CalculateFare()
   │  ├─ Get base fare for vehicle type
   │  ├─ Calculate distance
   │  ├─ Calculate estimated time
   │  ├─ Apply surge pricing if needed
   │  └─ Return total fare
   ├─ MatchingService.FindDriver()
   │  ├─ Get available drivers near pickup
   │  ├─ Filter by vehicle type
   │  ├─ Apply matching strategy (nearest/rating-based)
   │  └─ Return best match
   ├─ NotificationService.NotifyDriver()
   │  └─ Send ride request to driver
   └─ Return Ride with fare estimate
```

#### Accept Ride Flow

```
1. DriverService.AcceptRide()
   ├─ Validate driver (available, no active ride)
   ├─ RideService.UpdateRideState(ACCEPTED)
   ├─ Generate OTP
   ├─ NotificationService.NotifyRider()
   │  └─ Send driver details + OTP
   └─ LocationService.StartTracking(driverID)
```

### Step 6: Extensibility Examples

#### Adding New Vehicle Type

```go
// Just implement VehicleType enum and FareCalculator strategy
// No changes to RideService needed!

type VehicleType string
const (
    Car VehicleType = "CAR"
    Bike VehicleType = "BIKE"
    SUV VehicleType = "SUV"
    Auto VehicleType = "AUTO" // NEW - easy to add
)

type AutoFareStrategy struct {
    // Implement FareCalculator interface
}
```

#### Adding New Matching Strategy

```go
// Just implement MatchingStrategy interface
// No changes to MatchingService needed!

type PremiumMatchingStrategy struct {
    // Match only premium drivers to premium riders
}
```

#### Adding New Notification Channel

```go
// Just add to NotificationService channels
// No changes to other services needed!

type WhatsAppChannel struct {
    // Implement NotificationChannel interface
}
```

### Step 7: Error Handling

```go
type RideError struct {
    Code    string
    Message string
    Cause   error
}

var (
    ErrRiderHasActiveRide = &RideError{Code: "RIDER_ACTIVE", Message: "Rider has an active ride"}
    ErrNoDriversAvailable = &RideError{Code: "NO_DRIVERS", Message: "No drivers available"}
    ErrInvalidOTP = &RideError{Code: "INVALID_OTP", Message: "OTP verification failed"}
    ErrRideNotFound = &RideError{Code: "RIDE_NOT_FOUND", Message: "Ride not found"}
)
```

### Step 8: Repository Pattern (Data Access)

```go
type RideRepository interface {
    Create(ride *Ride) error
    GetByID(id string) (*Ride, error)
    Update(ride *Ride) error
    GetActiveRidesByRider(riderID string) ([]*Ride, error)
    GetActiveRidesByDriver(driverID string) ([]*Ride, error)
}

type DriverRepository interface {
    GetAvailableDrivers(vehicleType VehicleType, location *Location, radius float64) ([]*Driver, error)
    GetByID(id string) (*Driver, error)
    Update(driver *Driver) error
}
```

---

## 🎯 Key Design Decisions

### 1. Why Service Layer?

- **Separation of Concerns**: Business logic separate from data access
- **Testability**: Easy to mock dependencies
- **Reusability**: Services can be used by different controllers

### 2. Why Interfaces?

- **Extensibility**: Easy to swap implementations
- **Testability**: Easy to create mocks
- **Flexibility**: Multiple strategies can coexist

### 3. Why Repository Pattern?

- **Abstraction**: Hide database details
- **Testability**: Easy to mock data access
- **Flexibility**: Can switch databases easily

### 4. Why State Pattern for Ride?

- **Clear State Transitions**: Prevents invalid state changes
- **Business Logic**: State-specific behavior encapsulated
- **Maintainability**: Easy to add new states

---

## 📊 Comparison: Before vs After Framework

### Before (Your Initial Design):

- Entities identified ✓
- Functions listed ✓
- Services grouped ✓
- **Missing**: Interfaces, patterns, error handling, extensibility

### After (With Framework):

- ✅ Clear entity relationships
- ✅ Service interfaces for abstraction
- ✅ Design patterns applied
- ✅ Error handling strategy
- ✅ Extensibility built-in
- ✅ Complete flow design
- ✅ Repository pattern for data access

---

## 🚀 Next Steps

1. **Implement interfaces first** - Define contracts
2. **Implement concrete classes** - Build actual functionality
3. **Add error handling** - Handle edge cases
4. **Write unit tests** - Test each service
5. **Refactor** - Improve based on feedback

Remember: **Design is iterative**. Start simple, then add complexity as needed.
