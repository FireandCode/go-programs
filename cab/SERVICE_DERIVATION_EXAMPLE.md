# Step-by-Step: Deriving Services from Flows

## 🎯 Complete Walkthrough with Example

Let's walk through the **Request Ride** flow step-by-step and see how services emerge.

---

## Step 1: Write the Flow in Detail

### Request Ride Flow

```
1. Rider calls: RequestRide(pickup, drop, vehicleType)
2. System validates: Rider exists, no active rides, account is active
3. System calculates: Fare estimate based on distance, time, vehicle type
4. System finds: Available drivers within 5km of pickup location
5. System filters: Drivers with matching vehicle type
6. System ranks: Drivers by distance and rating
7. System sends: Ride request notification to top 3 drivers
8. Driver accepts: AcceptRide(driverID, rideID)
9. System creates: Ride record with state = ACCEPTED
10. System generates: 4-digit OTP
11. System stores: OTP with expiration (5 minutes)
12. System notifies: Rider with driver details and OTP
13. System starts: Real-time location tracking for driver
```

---

## Step 2: Extract Operations (Verbs)

### Highlight all actions:

```
1. Rider calls: RequestRide(pickup, drop, vehicleType)
2. System **validates**: Rider exists, no active rides, account is active
3. System **calculates**: Fare estimate based on distance, time, vehicle type
4. System **finds**: Available drivers within 5km of pickup location
5. System **filters**: Drivers with matching vehicle type
6. System **ranks**: Drivers by distance and rating
7. System **sends**: Ride request notification to top 3 drivers
8. Driver accepts: AcceptRide(driverID, rideID)
9. System **creates**: Ride record with state = ACCEPTED
10. System **generates**: 4-digit OTP
11. System **stores**: OTP with expiration (5 minutes)
12. System **notifies**: Rider with driver details and OTP
13. System **starts**: Real-time location tracking for driver
```

### Operations List:

- ✅ **validate** rider
- ✅ **calculate** fare
- ✅ **find** drivers
- ✅ **filter** drivers
- ✅ **rank** drivers
- ✅ **send** notification (to drivers)
- ✅ **create** ride
- ✅ **generate** OTP
- ✅ **store** OTP
- ✅ **notify** rider
- ✅ **start** tracking

---

## Step 3: Group by Responsibility

### Ask: "What is this operation responsible for?"

| Operation                   | Responsibility                | Group                 |
| --------------------------- | ----------------------------- | --------------------- |
| validate rider              | Ensure rider can request ride | **Ride Management**   |
| create ride                 | Create and manage ride record | **Ride Management**   |
| calculate fare              | Determine pricing             | **Pricing**           |
| find drivers                | Locate available drivers      | **Matching**          |
| filter drivers              | Filter by criteria            | **Matching**          |
| rank drivers                | Rank by algorithm             | **Matching**          |
| send notification (drivers) | Communicate with drivers      | **Communication**     |
| notify rider                | Communicate with rider        | **Communication**     |
| generate OTP                | Create authentication code    | **Authentication**    |
| store OTP                   | Persist authentication data   | **Authentication**    |
| start tracking              | Monitor location              | **Location Tracking** |

### Grouped Operations:

#### Group 1: Ride Management

- validate rider
- create ride

#### Group 2: Pricing

- calculate fare

#### Group 3: Matching

- find drivers
- filter drivers
- rank drivers

#### Group 4: Communication

- send notification (drivers)
- notify rider

#### Group 5: Authentication

- generate OTP
- store OTP

#### Group 6: Location Tracking

- start tracking

---

## Step 4: Name the Services

### Based on Responsibility:

| Group             | Responsibility         | Service Name              |
| ----------------- | ---------------------- | ------------------------- |
| Ride Management   | Manage ride lifecycle  | **RideService**           |
| Pricing           | Calculate fares        | **FareService**           |
| Matching          | Match drivers to rides | **MatchingService**       |
| Communication     | Send notifications     | **NotificationService**   |
| Authentication    | Verify identity        | **AuthenticationService** |
| Location Tracking | Track locations        | **LocationService**       |

---

## Step 5: Define Service Interfaces

### For Each Service, Define What It Does:

#### 1. RideService

**Responsibility**: Orchestrate ride lifecycle

```go
type RideService interface {
    // RequestRide orchestrates the entire flow
    RequestRide(riderID string, pickup, drop *Location, vehicleType VehicleType) (*Ride, error)

    // AcceptRide handles driver acceptance
    AcceptRide(driverID string, rideID string) error

    // Other ride operations...
    StartRide(rideID string, otp string) error
    EndRide(rideID string) error
    CancelRide(rideID string, userID string) error
}
```

**Why this interface?**

- Clear responsibility: Manage rides
- All methods relate to ride lifecycle
- One reason to change: Ride business rules

#### 2. FareService

**Responsibility**: Calculate pricing

```go
type FareCalculator interface {
    // Calculate fare for a ride
    CalculateFare(pickup, drop *Location, vehicleType VehicleType) (*Fare, error)

    // Calculate final fare (after ride completion)
    CalculateFinalFare(ride *Ride) (*Fare, error)
}
```

**Why this interface?**

- Single responsibility: Pricing
- Clear inputs/outputs
- One reason to change: Pricing rules

#### 3. MatchingService

**Responsibility**: Find and rank drivers

```go
type MatchingService interface {
    // Find available drivers
    FindAvailableDrivers(pickup *Location, vehicleType VehicleType, radius float64) ([]*Driver, error)

    // Find best driver (includes ranking)
    FindBestDriver(pickup *Location, vehicleType VehicleType) (*Driver, error)
}
```

**Why this interface?**

- Single responsibility: Driver matching
- Encapsulates finding, filtering, ranking
- One reason to change: Matching algorithm

**Implementation with DriverRepository:**

```go
// MatchingService uses DriverRepository directly ✅
type matchingService struct {
    driverRepo DriverRepository  // ✅ Direct access to driver data
    strategy   MatchingStrategy
}

func (s *matchingService) FindAvailableDrivers(
    pickup *Location,
    vehicleType VehicleType,
    radius float64,
) ([]*Driver, error) {
    // Step 1: Get drivers from repository (data access)
    drivers, err := s.driverRepo.GetDriversNearLocation(pickup, radius)
    if err != nil {
        return nil, err
    }

    // Step 2: Filter by availability and vehicle type (business logic)
    availableDrivers := s.filterAvailableDrivers(drivers, vehicleType)

    // Step 3: Rank drivers using strategy (business logic)
    rankedDrivers := s.strategy.RankDrivers(availableDrivers, pickup)

    return rankedDrivers, nil
}

func (s *matchingService) filterAvailableDrivers(
    drivers []*Driver,
    vehicleType VehicleType,
) []*Driver {
    var filtered []*Driver
    for _, driver := range drivers {
        if driver.IsAvailable &&
           driver.Vehicle.Type == vehicleType &&
           driver.CurrentRideID == "" {
            filtered = append(filtered, driver)
        }
    }
    return filtered
}
```

**Why DriverRepository?**

- MatchingService needs driver data to fulfill its matching responsibility
- Repository provides data access (no business logic)
- Direct access is clear and efficient
- No unnecessary coupling through other services

#### 4. NotificationService

**Responsibility**: Send notifications

```go
type NotificationService interface {
    // Notify rider
    NotifyRider(riderID string, message string, data map[string]interface{}) error

    // Notify driver
    NotifyDriver(driverID string, message string, data map[string]interface{}) error
}
```

**Why this interface?**

- Single responsibility: Communication
- Reusable across all flows
- One reason to change: Notification channels

#### 5. AuthenticationService

**Responsibility**: Verify identity

```go
type AuthenticationService interface {
    // Generate OTP for ride
    GenerateOTP(rideID string) (string, error)

    // Validate OTP
    ValidateOTP(rideID string, otp string) (bool, error)
}
```

**Why this interface?**

- Single responsibility: Authentication
- Reusable for other auth flows
- One reason to change: Auth mechanism

#### 6. LocationService

**Responsibility**: Track locations

```go
type LocationService interface {
    // Update location
    UpdateLocation(userID string, location *Location) error

    // Get current location
    GetLocation(userID string) (*Location, error)

    // Track location in real-time
    TrackLocation(userID string) (<-chan *Location, error)
}
```

**Why this interface?**

- Single responsibility: Location tracking
- Reusable for any user tracking
- One reason to change: Tracking mechanism

---

## Step 6: Show How Services Work Together

### RideService Implementation (Orchestrator)

```go
type rideService struct {
    // Dependencies (all interfaces!)
    rideRepo          RideRepository
    fareCalculator    FareCalculator
    matchingService   MatchingService
    notificationService NotificationService
    locationService   LocationService
    authService       AuthenticationService
}

func (s *rideService) RequestRide(
    riderID string,
    pickup, drop *Location,
    vehicleType VehicleType,
) (*Ride, error) {

    // Step 2: Validate rider
    rider, err := s.rideRepo.GetRider(riderID)
    if err != nil {
        return nil, err
    }
    if rider.HasActiveRide() {
        return nil, ErrRiderHasActiveRide
    }

    // Step 3: Calculate fare
    fare, err := s.fareCalculator.CalculateFare(pickup, drop, vehicleType)
    if err != nil {
        return nil, err
    }

    // Steps 4-6: Find, filter, rank drivers
    drivers, err := s.matchingService.FindAvailableDrivers(pickup, vehicleType, 5.0)
    if err != nil {
        return nil, err
    }

    // Step 7: Notify drivers
    for _, driver := range drivers[:3] { // Top 3
        s.notificationService.NotifyDriver(
            driver.ID,
            "New ride request",
            map[string]interface{}{
                "pickup": pickup,
                "drop": drop,
                "fare": fare,
            },
        )
    }

    // Step 9: Create ride
    ride := &Ride{
        RiderID:     riderID,
        PickupPoint: pickup,
        DropPoint:   drop,
        VehicleType: vehicleType,
        Fare:        fare,
        State:       RideStateRequested,
    }
    err = s.rideRepo.Create(ride)
    if err != nil {
        return nil, err
    }

    // Ride is created, waiting for driver acceptance
    // Steps 10-13 happen in AcceptRide()

    return ride, nil
}

func (s *rideService) AcceptRide(driverID string, rideID string) error {
    // Get ride
    ride, err := s.rideRepo.GetByID(rideID)
    if err != nil {
        return err
    }

    // Update ride
    ride.DriverID = driverID
    ride.State = RideStateAccepted

    // Step 10-11: Generate and store OTP
    otp, err := s.authService.GenerateOTP(rideID)
    if err != nil {
        return err
    }
    ride.OTP = otp

    err = s.rideRepo.Update(ride)
    if err != nil {
        return err
    }

    // Step 12: Notify rider
    s.notificationService.NotifyRider(
        ride.RiderID,
        "Driver assigned",
        map[string]interface{}{
            "driverID": driverID,
            "otp": otp,
        },
    )

    // Step 13: Start tracking
    s.locationService.TrackLocation(driverID)

    return nil
}
```

---

## Step 7: Verify Single Responsibility

### Check Each Service:

#### ✅ RideService

- **Responsibility**: Orchestrate ride lifecycle
- **Changes when**: Ride business rules change
- **Does NOT change when**:
  - Fare calculation changes (FareService)
  - Matching algorithm changes (MatchingService)
  - Notification channel changes (NotificationService)

#### ✅ FareService

- **Responsibility**: Calculate pricing
- **Changes when**: Pricing rules change
- **Does NOT change when**:
  - Ride creation changes (RideService)
  - Driver matching changes (MatchingService)

#### ✅ MatchingService

- **Responsibility**: Match drivers
- **Changes when**: Matching algorithm changes
- **Does NOT change when**:
  - Ride state changes (RideService)
  - Notification changes (NotificationService)

---

## Visual Summary

```
Request Ride Flow
│
├─ Step 2: Validate
│  └─> RideService (business rule)
│
├─ Step 3: Calculate Fare
│  └─> FareService (pricing)
│
├─ Steps 4-6: Find/Filter/Rank Drivers
│  └─> MatchingService (matching)
│
├─ Step 7: Notify Drivers
│  └─> NotificationService (communication)
│
├─ Step 9: Create Ride
│  └─> RideService (orchestration)
│
├─ Steps 10-11: Generate/Store OTP
│  └─> AuthenticationService (auth)
│
├─ Step 12: Notify Rider
│  └─> NotificationService (communication)
│
└─ Step 13: Start Tracking
   └─> LocationService (tracking)
```

---

## Key Takeaways

1. **Extract verbs** from flows → These are operations
2. **Group by responsibility** → What business capability?
3. **One service = one responsibility** → One reason to change
4. **Define interfaces** → Clear contracts, extensible
5. **Orchestrate in main service** → RideService coordinates others
6. **Depend on interfaces** → Not concrete implementations

**Remember**: Services should be **focused**, **testable**, and **composable**!
