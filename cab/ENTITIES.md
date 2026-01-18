# Entity Definitions for DriverRepository

## 🎯 Entity Structures

These are the entity definitions used by `DriverRepository` and other services.

---

## Driver Entity

```go
package entity

import "time"

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
```

---

## Location Entity

```go
package entity

import "time"

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
```

---

## Vehicle Entity

```go
package entity

// VehicleType represents the type of vehicle
type VehicleType string

const (
    VehicleTypeCar  VehicleType = "CAR"
    VehicleTypeBike VehicleType = "BIKE"
    VehicleTypeSUV  VehicleType = "SUV"
    VehicleTypeAuto VehicleType = "AUTO"
)

// Vehicle represents a vehicle
type Vehicle struct {
    ID          string
    Type        VehicleType
    LicensePlate string
    Model       string
    Capacity    int
    Year        int
}
```

---

## Complete Entity Package

```go
package entity

import (
    "math"
    "time"
)

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

// Vehicle represents a vehicle
type Vehicle struct {
    ID          string
    Type        VehicleType
    LicensePlate string
    Model       string
    Capacity    int
    Year        int
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
```

---

## Usage in Repository

The repository would import and use these entities:

```go
package repository

import "your-project/entity"

// DriverRepository uses entity.Driver
type DriverRepository interface {
    Create(driver *entity.Driver) error
    GetByID(driverID string) (*entity.Driver, error)
    // ...
}

// Or you can alias for convenience
type (
    Driver   = entity.Driver
    Location = entity.Location
    Vehicle  = entity.Vehicle
    VehicleType = entity.VehicleType
)
```

---

## Notes

1. **Separation of Concerns**: Entities are separate from repository implementation
2. **Domain Logic**: Entities can have methods (like `HasActiveRide()`)
3. **Value Objects**: Location and Vehicle are value objects
4. **Type Safety**: Using `VehicleType` enum for type safety

