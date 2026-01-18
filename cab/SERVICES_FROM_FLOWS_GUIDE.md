# Deriving Services from Flows: A Complete Guide

## 🎯 Overview

Services are **groupings of related operations** that work together to fulfill a business capability. They should have **one clear responsibility** and be defined through **interfaces** for extensibility.

---

## Step 1: Extract Operations from Flows

### Example: Request Ride Flow

```
1. Rider provides: pickup, drop, vehicle type
2. System validates rider (no active rides, account active)
3. System calculates fare estimate
4. System finds available drivers nearby
5. System sends ride request to drivers
6. Driver accepts
7. System creates ride record
8. System generates OTP
9. System notifies rider with driver details
10. System starts tracking driver location
```

### Extract Operations:

From this flow, identify **verbs** (actions):

- ✅ **validate** rider
- ✅ **calculate** fare
- ✅ **find** drivers
- ✅ **send** notification (to drivers)
- ✅ **create** ride
- ✅ **generate** OTP
- ✅ **notify** rider
- ✅ **track** location

---

## Step 2: Group Operations by Responsibility

### Principle: Single Responsibility Principle (SRP)

**Each service should have ONE reason to change.**

Ask yourself: "What is this service responsible for?"

### Grouping Strategy:

#### ❌ Bad Grouping (Too Broad)

```
RideManagementService:
- validate rider
- calculate fare
- find drivers
- create ride
- generate OTP
- notify users
- track location
```

**Problem**: Too many responsibilities! Changes to fare calculation affect ride creation, notifications, etc.

#### ✅ Good Grouping (Focused)

```
RideService:
- validate rider
- create ride
- update ride state

FareService:
- calculate fare

MatchingService:
- find drivers

NotificationService:
- notify rider
- notify driver

LocationService:
- track location
- get location

AuthenticationService:
- generate OTP
- validate OTP
```

**Better**: Each service has one clear responsibility!

---

## Step 3: Identify Service Responsibilities

### How to Identify Responsibilities:

Ask: **"What business capability does this service provide?"**

| Service                   | Responsibility         | Operations                                        |
| ------------------------- | ---------------------- | ------------------------------------------------- |
| **RideService**           | Manage ride lifecycle  | Create, start, end, cancel, update state          |
| **FareService**           | Calculate pricing      | Calculate fare, apply surge, apply discounts      |
| **MatchingService**       | Match drivers to rides | Find available drivers, rank drivers              |
| **NotificationService**   | Send notifications     | Notify rider, notify driver, send alerts          |
| **LocationService**       | Track locations        | Update location, get location, track in real-time |
| **AuthenticationService** | Verify identity        | Generate OTP, validate OTP                        |
| **PaymentService**        | Process payments       | Charge rider, refund, process payment             |

### Key Insight:

**One responsibility = One reason to change**

- If fare calculation logic changes → Only FareService changes
- If notification channel changes → Only NotificationService changes
- If matching algorithm changes → Only MatchingService changes

---

## Step 4: Define Service Interfaces (Contracts)

### Why Interfaces?

1. **Abstraction**: Hide implementation details
2. **Extensibility**: Easy to swap implementations
3. **Testability**: Easy to mock for testing
4. **Flexibility**: Multiple implementations can coexist

### Interface Design Principles:

#### 1. **Clear Input/Output**

```go
// ❌ Bad: Unclear what it does
type FareService interface {
    Calculate(data interface{}) interface{}
}

// ✅ Good: Clear inputs and outputs
type FareCalculator interface {
    CalculateFare(
        pickup *Location,
        drop *Location,
        vehicleType VehicleType,
    ) (*Fare, error)
}
```

#### 2. **Single Responsibility in Interface**

```go
// ❌ Bad: Too many responsibilities
type RideService interface {
    CreateRide(...) error
    CalculateFare(...) (*Fare, error)  // Should be in FareService
    NotifyRider(...) error              // Should be in NotificationService
    TrackLocation(...) error            // Should be in LocationService
}

// ✅ Good: One responsibility
type RideService interface {
    RequestRide(riderID string, pickup, drop *Location, vehicleType VehicleType) (*Ride, error)
    AcceptRide(driverID string, rideID string) error
    StartRide(rideID string, otp string) error
    EndRide(rideID string) error
    CancelRide(rideID string, userID string) error
    GetRide(rideID string) (*Ride, error)
}
```

#### 3. **Dependency on Abstractions**

```go
// ✅ Good: Depend on interfaces, not concrete types
type RideService struct {
    fareCalculator    FareCalculator      // Interface
    matchingService   MatchingService     // Interface
    notificationService NotificationService // Interface
    locationService   LocationService     // Interface
    rideRepo         RideRepository       // Interface
}

// ❌ Bad: Depend on concrete implementations
type RideService struct {
    fareCalculator    *FareService        // Concrete type
    matchingService   *MatchingService    // Concrete type
}
```

---

## Step 5: Complete Example - Uber-like Service

### Flow: Request Ride

```
1. Rider provides: pickup, drop, vehicle type
2. Validate rider
3. Calculate fare
4. Find available drivers
5. Notify drivers
6. Driver accepts
7. Create ride
8. Generate OTP
9. Notify rider
10. Start tracking
```

### Services Derived:

#### 1. RideService (Orchestrator)

**Responsibility**: Orchestrate ride lifecycle

```go
type RideService interface {
    RequestRide(riderID string, pickup, drop *Location, vehicleType VehicleType) (*Ride, error)
    AcceptRide(driverID string, rideID string) error
    StartRide(rideID string, otp string) error
    EndRide(rideID string) error
    CancelRide(rideID string, userID string) error
}

// Implementation uses other services
type rideService struct {
    rideRepo          RideRepository
    fareCalculator    FareCalculator
    matchingService   MatchingService
    notificationService NotificationService
    locationService   LocationService
    authService       AuthenticationService
}

func (s *rideService) RequestRide(riderID string, pickup, drop *Location, vehicleType VehicleType) (*Ride, error) {
    // 1. Validate rider
    rider, err := s.rideRepo.GetRider(riderID)
    if err != nil {
        return nil, err
    }
    if rider.HasActiveRide() {
        return nil, ErrRiderHasActiveRide
    }

    // 2. Calculate fare
    fare, err := s.fareCalculator.CalculateFare(pickup, drop, vehicleType)
    if err != nil {
        return nil, err
    }

    // 3. Find available drivers
    drivers, err := s.matchingService.FindAvailableDrivers(pickup, vehicleType)
    if err != nil {
        return nil, err
    }

    // 4. Notify drivers
    for _, driver := range drivers {
        s.notificationService.NotifyDriver(driver.ID, "New ride request")
    }

    // 5. Create ride (waiting for driver acceptance)
    ride := &Ride{
        RiderID:     riderID,
        PickupPoint: pickup,
        DropPoint:   drop,
        VehicleType: vehicleType,
        Fare:        fare,
        State:       RideStateRequested,
    }

    err = s.rideRepo.Create(ride)
    return ride, err
}
```

#### 2. FareService

**Responsibility**: Calculate pricing

```go
type FareCalculator interface {
    CalculateFare(pickup, drop *Location, vehicleType VehicleType) (*Fare, error)
    CalculateFinalFare(ride *Ride) (*Fare, error)
}

type fareService struct {
    strategies map[VehicleType]FareStrategy
    surgeService SurgePricingService
}

func (s *fareService) CalculateFare(pickup, drop *Location, vehicleType VehicleType) (*Fare, error) {
    strategy := s.strategies[vehicleType]
    baseFare := strategy.CalculateBaseFare(pickup, drop)

    surgeMultiplier := s.surgeService.GetSurgeMultiplier(pickup)

    return &Fare{
        BaseFare: baseFare,
        SurgeMultiplier: surgeMultiplier,
        TotalFare: baseFare * surgeMultiplier,
    }, nil
}
```

#### 3. MatchingService

**Responsibility**: Match drivers to rides

```go
type MatchingService interface {
    FindAvailableDrivers(pickup *Location, vehicleType VehicleType) ([]*Driver, error)
    FindBestDriver(pickup *Location, vehicleType VehicleType) (*Driver, error)
}

type matchingService struct {
    driverRepo DriverRepository
    strategy   MatchingStrategy
}

func (s *matchingService) FindAvailableDrivers(pickup *Location, vehicleType VehicleType) ([]*Driver, error) {
    return s.driverRepo.GetAvailableDrivers(pickup, vehicleType, 5.0) // 5km radius
}
```

#### 4. NotificationService

**Responsibility**: Send notifications

```go
type NotificationService interface {
    NotifyRider(riderID string, message string, data map[string]interface{}) error
    NotifyDriver(driverID string, message string, data map[string]interface{}) error
}

type notificationService struct {
    channels []NotificationChannel
}

func (s *notificationService) NotifyRider(riderID string, message string, data map[string]interface{}) error {
    for _, channel := range s.channels {
        channel.Send(riderID, message, data)
    }
    return nil
}
```

#### 5. LocationService

**Responsibility**: Track and retrieve locations

```go
type LocationService interface {
    UpdateLocation(userID string, location *Location) error
    GetLocation(userID string) (*Location, error)
    TrackLocation(userID string) (<-chan *Location, error)
}

type locationService struct {
    cache LocationCache
    tracker LocationTracker
}

func (s *locationService) UpdateLocation(userID string, location *Location) error {
    return s.cache.Update(userID, location)
}
```

#### 6. AuthenticationService

**Responsibility**: Verify identity

```go
type AuthenticationService interface {
    GenerateOTP(rideID string) (string, error)
    ValidateOTP(rideID string, otp string) (bool, error)
}

type authService struct {
    otpRepo OTPRepository
}

func (s *authService) GenerateOTP(rideID string) (string, error) {
    otp := generateRandomOTP()
    s.otpRepo.Store(rideID, otp, time.Now().Add(5*time.Minute))
    return otp, nil
}
```

---

## Common Mistakes to Avoid

### ❌ Mistake 1: God Service (Too Many Responsibilities)

```go
// ❌ Bad: One service doing everything
type UberService struct {
    // Too many responsibilities!
}

func (s *UberService) DoEverything() {
    // Calculate fare
    // Find drivers
    // Send notifications
    // Track locations
    // Process payments
    // ...
}
```

**Fix**: Split into focused services

### ❌ Mistake 2: Anemic Services (No Business Logic)

```go
// ❌ Bad: Just a pass-through
type RideService struct {
    repo RideRepository
}

func (s *RideService) CreateRide(ride *Ride) error {
    return s.repo.Create(ride) // No validation, no business logic
}
```

**Fix**: Add business logic, validation, orchestration

### ❌ Mistake 3: Tight Coupling

```go
// ❌ Bad: Direct dependency on concrete types
type RideService struct {
    fareService *FareService  // Concrete type
}

// ✅ Good: Depend on interface
type RideService struct {
    fareCalculator FareCalculator  // Interface
}
```

### ❌ Mistake 4: Unclear Interface

```go
// ❌ Bad: Unclear what it does
type Service interface {
    Process(data interface{}) interface{}
}

// ✅ Good: Clear purpose
type FareCalculator interface {
    CalculateFare(pickup, drop *Location, vehicleType VehicleType) (*Fare, error)
}
```

---

## Decision Framework: When to Create a New Service?

### Ask These Questions:

1. **Does it have a distinct responsibility?**

   - Yes → New service
   - No → Part of existing service

2. **Will it change for different reasons?**

   - Yes → Separate service
   - No → Can be combined

3. **Can it be tested independently?**

   - Yes → Good candidate for separate service
   - No → Might be too coupled

4. **Is it reusable across different flows?**
   - Yes → Separate service
   - No → Might be flow-specific

### Example Decisions:

| Operation         | Service?               | Why?                                         |
| ----------------- | ---------------------- | -------------------------------------------- |
| Calculate fare    | ✅ FareService         | Distinct responsibility, reusable, testable  |
| Find drivers      | ✅ MatchingService     | Distinct algorithm, can change independently |
| Send notification | ✅ NotificationService | Different channels, reusable                 |
| Validate rider    | ❌ Part of RideService | Business rule for ride, not reusable         |
| Generate OTP      | ✅ AuthService         | Reusable for other auth flows                |

---

## Quick Reference: Service Identification

### From Flow Operations:

1. **Extract verbs** from flow steps
2. **Group by responsibility** (what business capability?)
3. **Check SRP** (one reason to change?)
4. **Define interface** (clear contract)
5. **Inject dependencies** (depend on interfaces)

### Service Naming:

- **Service suffix**: For orchestration (RideService)
- **Manager suffix**: For coordination (LocationManager)
- **Calculator suffix**: For calculations (FareCalculator)
- **Repository suffix**: For data access (RideRepository)

---

## Summary

1. **Extract operations** from flows (verbs)
2. **Group by responsibility** (one reason to change)
3. **Define interfaces** (clear contracts)
4. **Inject dependencies** (depend on abstractions)
5. **Keep focused** (one service = one responsibility)

**Remember**: Services should be **focused**, **testable**, and **extensible**!
