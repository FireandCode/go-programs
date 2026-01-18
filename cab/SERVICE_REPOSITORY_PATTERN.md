# Service-Repository Pattern: When Should Services Use Repositories?

## 🎯 The Question

**Should MatchingService use DriverRepository directly, or go through another service?**

**Short Answer**: ✅ **Yes, MatchingService can and should use DriverRepository directly.**

---

## 📊 Understanding the Layers

### Architecture Layers:

```
┌─────────────────────────────────────┐
│   Service Layer (Business Logic)   │
│  - RideService                      │
│  - MatchingService                  │
│  - FareService                      │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Repository Layer (Data Access)    │
│  - RideRepository                   │
│  - DriverRepository                 │
│  - RiderRepository                  │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Database/Storage Layer           │
└─────────────────────────────────────┘
```

---

## ✅ Correct Approach: MatchingService Uses DriverRepository

### Why This is Correct:

**MatchingService needs driver data to fulfill its responsibility (matching).**

```go
// ✅ GOOD: MatchingService uses DriverRepository
type MatchingService interface {
    FindAvailableDrivers(pickup *Location, vehicleType VehicleType, radius float64) ([]*Driver, error)
    FindBestDriver(pickup *Location, vehicleType VehicleType) (*Driver, error)
}

type matchingService struct {
    driverRepo DriverRepository  // ✅ Direct access to driver data
    strategy   MatchingStrategy
}

func (s *matchingService) FindAvailableDrivers(
    pickup *Location,
    vehicleType VehicleType,
    radius float64,
) ([]*Driver, error) {
    // MatchingService needs driver data for its matching responsibility
    // It's appropriate to access DriverRepository directly

    // Step 1: Get drivers from repository
    drivers, err := s.driverRepo.GetAvailableDrivers(pickup, vehicleType, radius)
    if err != nil {
        return nil, err
    }

    // Step 2: Apply matching strategy (ranking, filtering)
    rankedDrivers := s.strategy.RankDrivers(drivers, pickup)

    return rankedDrivers, nil
}
```

### Why This Works:

1. **Single Responsibility**: MatchingService is responsible for matching, which requires driver data
2. **Direct Data Access**: MatchingService needs to query drivers, so it uses DriverRepository
3. **No Business Logic in Repository**: Repository only handles data access, not business rules
4. **Clear Boundaries**: Each service accesses repositories for entities it needs

---

## ❌ Alternative (Not Recommended): Going Through Another Service

### Why NOT to go through DriverService:

```go
// ❌ BAD: MatchingService depends on DriverService
type matchingService struct {
    driverService DriverService  // ❌ Unnecessary indirection
    strategy      MatchingStrategy
}

func (s *matchingService) FindAvailableDrivers(...) ([]*Driver, error) {
    // Now MatchingService depends on DriverService
    // DriverService might have business logic we don't need
    // Creates unnecessary coupling
    drivers, err := s.driverService.GetAvailableDrivers(...)
    // ...
}
```

**Problems:**

1. **Unnecessary Coupling**: MatchingService depends on DriverService
2. **Indirection**: Extra layer without benefit
3. **Mixed Responsibilities**: DriverService might have other responsibilities
4. **Circular Dependencies Risk**: DriverService might need MatchingService

---

## 📋 Decision Framework: When Should a Service Use a Repository?

### ✅ Use Repository Directly When:

1. **Service needs data for its core responsibility**

   - MatchingService needs drivers → Use DriverRepository ✅
   - FareService needs pricing rules → Use PricingRepository ✅
   - LocationService needs location data → Use LocationRepository ✅

2. **Repository only provides data access (no business logic)**

   - Repository just queries/stores data
   - No complex business rules in repository

3. **Service needs to query/filter data**
   - MatchingService needs to query available drivers
   - RideService needs to query active rides

### ❌ Don't Use Repository When:

1. **Another service already provides the needed functionality**

   - If DriverService already has `GetAvailableDrivers()` with business logic
   - But this is rare - usually services orchestrate, repositories query

2. **You need to bypass business rules**
   - If you need raw data without validation
   - Usually indicates design issue

---

## 🎯 Complete Example: MatchingService Implementation

### Repository Interface:

```go
// DriverRepository - Data access only, no business logic
type DriverRepository interface {
    // Get available drivers (data query)
    GetAvailableDrivers(
        location *Location,
        vehicleType VehicleType,
        radius float64,
    ) ([]*Driver, error)

    // Get driver by ID
    GetByID(driverID string) (*Driver, error)

    // Update driver
    Update(driver *Driver) error

    // Get drivers by location (raw query)
    GetDriversNearLocation(location *Location, radius float64) ([]*Driver, error)
}
```

### MatchingService Implementation:

```go
// MatchingService - Business logic for matching
type MatchingService interface {
    FindAvailableDrivers(
        pickup *Location,
        vehicleType VehicleType,
        radius float64,
    ) ([]*Driver, error)

    FindBestDriver(
        pickup *Location,
        vehicleType VehicleType,
    ) (*Driver, error)
}

type matchingService struct {
    driverRepo DriverRepository    // ✅ Direct access
    strategy   MatchingStrategy    // Strategy for ranking
}

func (s *matchingService) FindAvailableDrivers(
    pickup *Location,
    vehicleType VehicleType,
    radius float64,
) ([]*Driver, error) {

    // Step 1: Get drivers from repository (data access)
    allDrivers, err := s.driverRepo.GetDriversNearLocation(pickup, radius)
    if err != nil {
        return nil, err
    }

    // Step 2: Filter by availability and vehicle type (business logic)
    availableDrivers := s.filterAvailableDrivers(allDrivers, vehicleType)

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
        // Business logic: Check if driver is available
        if driver.IsAvailable &&
           driver.Vehicle.Type == vehicleType &&
           driver.CurrentRideID == "" {
            filtered = append(filtered, driver)
        }
    }
    return filtered
}

func (s *matchingService) FindBestDriver(
    pickup *Location,
    vehicleType VehicleType,
) (*Driver, error) {
    drivers, err := s.FindAvailableDrivers(pickup, vehicleType, 5.0)
    if err != nil {
        return nil, err
    }

    if len(drivers) == 0 {
        return nil, ErrNoDriversAvailable
    }

    // Return top-ranked driver
    return drivers[0], nil
}
```

---

## 🔄 Comparison: Different Approaches

### Approach 1: Direct Repository Access ✅ (Recommended)

```go
type matchingService struct {
    driverRepo DriverRepository
    strategy   MatchingStrategy
}

// MatchingService directly queries drivers
// Applies matching business logic
// Returns matched drivers
```

**Pros:**

- ✅ Clear responsibility
- ✅ Direct data access
- ✅ No unnecessary coupling
- ✅ Easy to test (mock repository)

**Cons:**

- None significant

---

### Approach 2: Through DriverService ❌ (Not Recommended)

```go
type matchingService struct {
    driverService DriverService  // ❌ Unnecessary
    strategy      MatchingStrategy
}

// MatchingService asks DriverService
// DriverService might have other responsibilities
// Creates coupling
```

**Pros:**

- None significant

**Cons:**

- ❌ Unnecessary indirection
- ❌ Creates coupling
- ❌ DriverService might have other responsibilities
- ❌ Harder to test

---

### Approach 3: Hybrid (When DriverService Has Business Logic)

```go
// If DriverService has complex business logic for "available drivers"
type DriverService interface {
    GetAvailableDrivers(criteria DriverCriteria) ([]*Driver, error)
    // This includes business logic like:
    // - Check driver status
    // - Check driver rating
    // - Check driver preferences
}

// Then MatchingService can use DriverService
type matchingService struct {
    driverService DriverService  // ✅ If it has business logic
    strategy      MatchingStrategy
}
```

**When to use:**

- DriverService has complex business logic for determining "available"
- Multiple services need the same business logic
- You want to centralize driver availability rules

---

## 📊 Service-Repository Mapping

### General Rule:

| Service             | Repositories It Uses            | Why?                           |
| ------------------- | ------------------------------- | ------------------------------ |
| **RideService**     | RideRepository, RiderRepository | Needs ride and rider data      |
| **MatchingService** | DriverRepository                | Needs driver data for matching |
| **FareService**     | PricingRepository (if exists)   | Needs pricing rules            |
| **LocationService** | LocationRepository              | Needs location data            |
| **PaymentService**  | PaymentRepository               | Needs payment data             |

### Key Principle:

**A service uses repositories for entities it needs to query/manage as part of its responsibility.**

---

## 🎓 Best Practices

### 1. Repository = Data Access Only

```go
// ✅ GOOD: Repository only does data access
type DriverRepository interface {
    GetByID(id string) (*Driver, error)
    GetAvailableDrivers(location *Location, radius float64) ([]*Driver, error)
    Create(driver *Driver) error
    Update(driver *Driver) error
}

// ❌ BAD: Repository has business logic
type DriverRepository interface {
    GetAvailableDrivers(...) ([]*Driver, error)
    // ❌ This should be in MatchingService:
    RankDriversByDistance(...) ([]*Driver, error)
}
```

### 2. Service = Business Logic

```go
// ✅ GOOD: Service has business logic
type MatchingService interface {
    FindAvailableDrivers(...) ([]*Driver, error)
    // Business logic: filtering, ranking, matching
}

// Implementation applies business rules
func (s *matchingService) FindAvailableDrivers(...) {
    // 1. Get data from repository
    drivers := s.driverRepo.GetDriversNearLocation(...)

    // 2. Apply business logic
    filtered := s.filterByAvailability(drivers)
    ranked := s.strategy.RankDrivers(filtered, ...)

    return ranked
}
```

### 3. One Service Can Use Multiple Repositories

```go
// ✅ GOOD: Service uses multiple repositories
type RideService struct {
    rideRepo   RideRepository    // For rides
    riderRepo  RiderRepository   // For riders
    driverRepo DriverRepository  // For drivers (if needed)
    // ...
}
```

### 4. Keep Repositories Focused

```go
// ✅ GOOD: Each repository handles one entity type
type DriverRepository interface {
    // Only driver-related queries
    GetByID(id string) (*Driver, error)
    GetAvailableDrivers(...) ([]*Driver, error)
}

// ❌ BAD: Repository handles multiple entities
type UberRepository interface {
    GetDriver(...) (*Driver, error)
    GetRider(...) (*Rider, error)
    GetRide(...) (*Ride, error)
    // Too broad!
}
```

---

## ✅ Summary

### For MatchingService:

**✅ YES, MatchingService should use DriverRepository directly**

**Why:**

1. MatchingService needs driver data for its matching responsibility
2. DriverRepository provides data access (no business logic)
3. Direct access is clear and efficient
4. No unnecessary coupling

**Structure:**

```go
type matchingService struct {
    driverRepo DriverRepository  // ✅ Direct access
    strategy   MatchingStrategy
}
```

**Remember:**

- **Repository** = Data access (queries, storage)
- **Service** = Business logic (matching, calculation, orchestration)
- **Service uses Repository** = Service needs data to do its job

This is the standard pattern! 🎯
