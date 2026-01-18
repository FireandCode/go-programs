# Service Orchestration: Who Coordinates All Services?

## 🎯 The Question

**Who orchestrates all these services? How do they work together?**

**Answer**: **RideService** acts as the **orchestrator** that coordinates all other services to fulfill business flows.

---

## 📊 Architecture Layers

### Complete Architecture:

```
┌─────────────────────────────────────────┐
│   API/Controller Layer (HTTP Handlers) │
│   - Handles HTTP requests               │
│   - Validates input                     │
│   - Calls services                      │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│   Orchestrator Service (RideService)    │
│   - Coordinates multiple services       │
│   - Implements business flows           │
│   - Manages transactions                │
└──────────────┬──────────────────────────┘
               │
       ┌───────┴───────┬───────────┬──────────┐
       │               │           │          │
┌──────▼──────┐ ┌─────▼─────┐ ┌───▼────┐ ┌───▼──────┐
│ Matching    │ │ Fare      │ │ Notif  │ │ Location │
│ Service     │ │ Service   │ │ Service│ │ Service  │
└──────┬──────┘ └─────┬─────┘ └───┬────┘ └───┬──────┘
       │               │           │          │
┌──────▼──────────────▼───────────▼──────────▼──────┐
│   Repository Layer (Data Access)                   │
│   - DriverRepository                                │
│   - RideRepository                                  │
│   - RiderRepository                                 │
└─────────────────────────────────────────────────────┘
```

---

## 🎭 Orchestrator Pattern

### RideService as Orchestrator

**RideService** is the **orchestrator** that:

1. Coordinates multiple services
2. Implements complete business flows
3. Manages the sequence of operations
4. Handles errors and rollbacks

---

## 🎯 Complete Orchestration Example

### RideService Implementation

```go
package service

import (
    "errors"
    "time"
)

// RideService orchestrates all services for ride management
type RideService interface {
    RequestRide(riderID string, pickup, drop *Location, vehicleType VehicleType) (*Ride, error)
    AcceptRide(driverID string, rideID string) error
    StartRide(rideID string, otp string) error
    EndRide(rideID string) error
    CancelRide(rideID string, userID string) error
}

type rideService struct {
    // Repositories
    rideRepo  RideRepository
    riderRepo RiderRepository

    // Services (dependencies)
    fareCalculator    FareCalculator
    matchingService   MatchingService
    notificationService NotificationService
    locationService   LocationService
    authService       AuthenticationService
    paymentService    PaymentService
}

// NewRideService creates a new RideService with all dependencies
func NewRideService(
    rideRepo RideRepository,
    riderRepo RiderRepository,
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
        fareCalculator:    fareCalculator,
        matchingService:   matchingService,
        notificationService: notificationService,
        locationService:   locationService,
        authService:       authService,
        paymentService:    paymentService,
    }
}
```

---

## 🔄 Request Ride Flow (Orchestration)

### Complete Orchestration Example

```go
func (s *rideService) RequestRide(
    riderID string,
    pickup, drop *Location,
    vehicleType VehicleType,
) (*Ride, error) {

    // Step 1: Validate rider (business rule)
    rider, err := s.riderRepo.GetByID(riderID)
    if err != nil {
        return nil, err
    }

    if rider.HasActiveRide() {
        return nil, ErrRiderHasActiveRide
    }

    // Step 2: Calculate fare (delegate to FareService)
    fare, err := s.fareCalculator.CalculateFare(pickup, drop, vehicleType)
    if err != nil {
        return nil, err
    }

    // Step 3: Find available drivers (delegate to MatchingService)
    drivers, err := s.matchingService.FindAvailableDrivers(pickup, vehicleType, 5.0)
    if err != nil {
        return nil, err
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
    }

    err = s.rideRepo.Create(ride)
    if err != nil {
        return nil, err
    }

    // Step 5: Notify drivers (delegate to NotificationService)
    for _, driver := range drivers[:3] { // Top 3 drivers
        s.notificationService.NotifyDriver(
            driver.ID,
            "New ride request",
            map[string]interface{}{
                "rideID":  ride.ID,
                "pickup":  pickup,
                "drop":    drop,
                "fare":    fare.TotalFare,
            },
        )
    }

    return ride, nil
}
```

---

## 🔄 Accept Ride Flow (Orchestration)

```go
func (s *rideService) AcceptRide(driverID string, rideID string) error {
    // Step 1: Get ride
    ride, err := s.rideRepo.GetByID(rideID)
    if err != nil {
        return err
    }

    if ride.State != RideStateRequested {
        return ErrInvalidRideState
    }

    // Step 2: Update ride with driver
    ride.DriverID = driverID
    ride.State = RideStateAccepted

    // Step 3: Generate OTP (delegate to AuthenticationService)
    otp, err := s.authService.GenerateOTP(rideID)
    if err != nil {
        return err
    }
    ride.OTP = otp

    // Step 4: Update driver status
    // (In real system, you'd update driver's current ride)

    // Step 5: Save ride
    err = s.rideRepo.Update(ride)
    if err != nil {
        return err
    }

    // Step 6: Notify rider (delegate to NotificationService)
    s.notificationService.NotifyRider(
        ride.RiderID,
        "Driver assigned",
        map[string]interface{}{
            "driverID": driverID,
            "otp":      otp,
        },
    )

    // Step 7: Start tracking driver location (delegate to LocationService)
    s.locationService.TrackLocation(driverID)

    return nil
}
```

---

## 🔄 Start Ride Flow (Orchestration)

```go
func (s *rideService) StartRide(rideID string, otp string) error {
    // Step 1: Get ride
    ride, err := s.rideRepo.GetByID(rideID)
    if err != nil {
        return err
    }

    if ride.State != RideStateAccepted {
        return ErrInvalidRideState
    }

    // Step 2: Validate OTP (delegate to AuthenticationService)
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

    // Step 4: Start tracking both locations (delegate to LocationService)
    s.locationService.TrackLocation(ride.DriverID)
    s.locationService.TrackLocation(ride.RiderID)

    // Step 5: Notify rider (delegate to NotificationService)
    s.notificationService.NotifyRider(
        ride.RiderID,
        "Ride started",
        map[string]interface{}{
            "rideID": rideID,
        },
    )

    return nil
}
```

---

## 🔄 End Ride Flow (Orchestration)

```go
func (s *rideService) EndRide(rideID string) error {
    // Step 1: Get ride
    ride, err := s.rideRepo.GetByID(rideID)
    if err != nil {
        return err
    }

    if ride.State != RideStateStarted {
        return ErrInvalidRideState
    }

    // Step 2: Calculate final fare (delegate to FareService)
    finalFare, err := s.fareCalculator.CalculateFinalFare(ride)
    if err != nil {
        return err
    }
    ride.Fare = finalFare

    // Step 3: Update ride state
    ride.State = RideStateCompleted
    ride.EndTime = time.Now()

    err = s.rideRepo.Update(ride)
    if err != nil {
        return err
    }

    // Step 4: Process payment (delegate to PaymentService)
    err = s.paymentService.ProcessPayment(ride.RiderID, finalFare.TotalFare)
    if err != nil {
        return err
    }

    // Step 5: Stop location tracking
    s.locationService.StopTracking(ride.DriverID)
    s.locationService.StopTracking(ride.RiderID)

    // Step 6: Notify both parties (delegate to NotificationService)
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
```

---

## 🎮 Controller/Handler Layer

### HTTP Handler (Entry Point)

```go
package handler

import (
    "encoding/json"
    "net/http"
    "your-project/service"
)

type RideHandler struct {
    rideService service.RideService
}

func NewRideHandler(rideService service.RideService) *RideHandler {
    return &RideHandler{
        rideService: rideService,
    }
}

// RequestRide handles HTTP request for ride creation
func (h *RideHandler) RequestRide(w http.ResponseWriter, r *http.Request) {
    // Step 1: Parse request
    var req struct {
        RiderID     string  `json:"rider_id"`
        PickupLat   float64 `json:"pickup_lat"`
        PickupLon   float64 `json:"pickup_lon"`
        DropLat     float64 `json:"drop_lat"`
        DropLon     float64 `json:"drop_lon"`
        VehicleType string  `json:"vehicle_type"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Step 2: Create locations
    pickup := &service.Location{
        Latitude:  req.PickupLat,
        Longitude: req.PickupLon,
    }
    drop := &service.Location{
        Latitude:  req.DropLat,
        Longitude: req.DropLon,
    }

    // Step 3: Call orchestrator service
    ride, err := h.rideService.RequestRide(
        req.RiderID,
        pickup,
        drop,
        service.VehicleType(req.VehicleType),
    )

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Step 4: Return response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(ride)
}
```

---

## 🏗️ Complete Setup (Dependency Injection)

### Main Function - Wiring Everything Together

```go
package main

import (
    "database/sql"
    "log"
    "net/http"

    "your-project/handler"
    "your-project/repository"
    "your-project/service"

    _ "github.com/lib/pq"
)

func main() {
    // Step 1: Setup database
    db, err := sql.Open("postgres", "postgres://...")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Step 2: Create repositories
    driverRepo := repository.NewSQLDriverRepository(db)
    rideRepo := repository.NewSQLRideRepository(db)
    riderRepo := repository.NewSQLRiderRepository(db)

    // Step 3: Create services
    fareService := service.NewFareService()
    matchingService := service.NewMatchingService(driverRepo)
    notificationService := service.NewNotificationService()
    locationService := service.NewLocationService()
    authService := service.NewAuthenticationService()
    paymentService := service.NewPaymentService()

    // Step 4: Create orchestrator service (RideService)
    rideService := service.NewRideService(
        rideRepo,
        riderRepo,
        fareService,
        matchingService,
        notificationService,
        locationService,
        authService,
        paymentService,
    )

    // Step 5: Create handlers
    rideHandler := handler.NewRideHandler(rideService)

    // Step 6: Setup routes
    http.HandleFunc("/api/rides/request", rideHandler.RequestRide)
    http.HandleFunc("/api/rides/accept", rideHandler.AcceptRide)
    http.HandleFunc("/api/rides/start", rideHandler.StartRide)
    http.HandleFunc("/api/rides/end", rideHandler.EndRide)

    // Step 7: Start server
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 📊 Service Interaction Diagram

### Request Ride Flow:

```
HTTP Request
    │
    ▼
RideHandler.RequestRide()
    │
    ▼
RideService.RequestRide()  ← ORCHESTRATOR
    │
    ├─→ RiderRepository.GetByID()
    │
    ├─→ FareService.CalculateFare()
    │
    ├─→ MatchingService.FindAvailableDrivers()
    │   └─→ DriverRepository.GetDriversNearLocation()
    │
    ├─→ RideRepository.Create()
    │
    └─→ NotificationService.NotifyDriver()
```

---

## 🎯 Key Principles

### 1. **Single Orchestrator**

- **RideService** is the main orchestrator
- It coordinates all other services
- Other services don't call each other directly

### 2. **Service Independence**

- Services are independent (FareService, MatchingService, etc.)
- They don't know about each other
- They only know about their repositories

### 3. **Dependency Flow**

```
Handler → Orchestrator Service → Specialized Services → Repositories
```

### 4. **Clear Responsibilities**

- **Handler**: HTTP concerns (parsing, validation, response)
- **Orchestrator Service**: Business flow coordination
- **Specialized Services**: Specific business logic
- **Repositories**: Data access

---

## 🔄 Alternative: Event-Driven Orchestration

### Using Events (Advanced)

```go
// Event-driven approach
type RideService struct {
    eventBus EventBus
    // ...
}

func (s *rideService) RequestRide(...) (*Ride, error) {
    // Create ride
    ride := &Ride{...}
    s.rideRepo.Create(ride)

    // Publish event instead of direct call
    s.eventBus.Publish("ride.requested", RideRequestedEvent{
        RideID: ride.ID,
        Pickup: pickup,
        Drop:   drop,
    })

    // Other services subscribe to events
    // MatchingService listens to "ride.requested"
    // NotificationService listens to "driver.assigned"

    return ride, nil
}
```

---

## ✅ Summary

### Who Orchestrates?

**RideService** is the orchestrator that:

1. ✅ Coordinates all services
2. ✅ Implements complete business flows
3. ✅ Manages the sequence of operations
4. ✅ Handles errors and transactions

### Architecture Flow:

```
HTTP Request
    ↓
Handler (HTTP layer)
    ↓
RideService (Orchestrator)
    ↓
┌───┴───┬──────┬──────────┬──────────┐
│       │      │          │          │
Fare   Match  Notify   Location  Auth
Service Service Service Service  Service
    ↓       ↓      ↓        ↓        ↓
Repositories (Data Access)
```

### Key Points:

1. **RideService = Orchestrator** - Coordinates everything
2. **Other Services = Specialized** - Do one thing well
3. **Repositories = Data Access** - No business logic
4. **Handlers = Entry Point** - HTTP concerns only

**Remember**: One orchestrator (RideService) coordinates all specialized services! 🎯
