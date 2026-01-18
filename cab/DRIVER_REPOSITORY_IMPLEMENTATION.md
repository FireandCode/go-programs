# DriverRepository Implementation Guide

## 🎯 Overview

This guide shows how to implement `DriverRepository` with different storage backends, following the Repository Pattern.

---

## 📋 Repository Interface

### DriverRepository Interface

```go
package repository

import (
    "errors"
    "time"
)

// DriverRepository defines the interface for driver data access
type DriverRepository interface {
    // Basic CRUD operations
    Create(driver *Driver) error
    GetByID(driverID string) (*Driver, error)
    Update(driver *Driver) error
    Delete(driverID string) error
    
    // Query operations
    GetAvailableDrivers(
        location *Location,
        vehicleType VehicleType,
        radius float64,
    ) ([]*Driver, error)
    
    GetDriversNearLocation(
        location *Location,
        radius float64,
    ) ([]*Driver, error)
    
    GetDriversByVehicleType(vehicleType VehicleType) ([]*Driver, error)
    GetAllDrivers() ([]*Driver, error)
    
    // Status operations
    UpdateDriverStatus(driverID string, isAvailable bool) error
    UpdateDriverLocation(driverID string, location *Location) error
}

// Common errors
var (
    ErrDriverNotFound    = errors.New("driver not found")
    ErrDriverExists      = errors.New("driver already exists")
    ErrInvalidLocation   = errors.New("invalid location")
)
```

---

## 🗄️ Implementation 1: In-Memory Repository (For Testing/Development)

### In-Memory Implementation

```go
package repository

import (
    "sync"
    "math"
)

// InMemoryDriverRepository implements DriverRepository using in-memory storage
type InMemoryDriverRepository struct {
    drivers map[string]*Driver
    mu      sync.RWMutex // For thread safety
}

// NewInMemoryDriverRepository creates a new in-memory repository
func NewInMemoryDriverRepository() *InMemoryDriverRepository {
    return &InMemoryDriverRepository{
        drivers: make(map[string]*Driver),
    }
}

// Create adds a new driver
func (r *InMemoryDriverRepository) Create(driver *Driver) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.drivers[driver.ID]; exists {
        return ErrDriverExists
    }
    
    r.drivers[driver.ID] = driver
    return nil
}

// GetByID retrieves a driver by ID
func (r *InMemoryDriverRepository) GetByID(driverID string) (*Driver, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    driver, exists := r.drivers[driverID]
    if !exists {
        return nil, ErrDriverNotFound
    }
    
    return driver, nil
}

// Update updates an existing driver
func (r *InMemoryDriverRepository) Update(driver *Driver) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.drivers[driver.ID]; !exists {
        return ErrDriverNotFound
    }
    
    r.drivers[driver.ID] = driver
    return nil
}

// Delete removes a driver
func (r *InMemoryDriverRepository) Delete(driverID string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.drivers[driverID]; !exists {
        return ErrDriverNotFound
    }
    
    delete(r.drivers, driverID)
    return nil
}

// GetAvailableDrivers finds available drivers matching criteria
func (r *InMemoryDriverRepository) GetAvailableDrivers(
    location *Location,
    vehicleType VehicleType,
    radius float64,
) ([]*Driver, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    var result []*Driver
    
    for _, driver := range r.drivers {
        // Check if driver is available
        if !driver.IsAvailable || driver.CurrentRideID != "" {
            continue
        }
        
        // Check vehicle type match
        if driver.Vehicle.Type != vehicleType {
            continue
        }
        
        // Check if driver is within radius
        distance := calculateDistance(location, driver.Location)
        if distance <= radius {
            result = append(result, driver)
        }
    }
    
    return result, nil
}

// GetDriversNearLocation finds all drivers within radius (regardless of availability)
func (r *InMemoryDriverRepository) GetDriversNearLocation(
    location *Location,
    radius float64,
) ([]*Driver, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    var result []*Driver
    
    for _, driver := range r.drivers {
        distance := calculateDistance(location, driver.Location)
        if distance <= radius {
            result = append(result, driver)
        }
    }
    
    return result, nil
}

// GetDriversByVehicleType finds all drivers with specific vehicle type
func (r *InMemoryDriverRepository) GetDriversByVehicleType(
    vehicleType VehicleType,
) ([]*Driver, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    var result []*Driver
    
    for _, driver := range r.drivers {
        if driver.Vehicle.Type == vehicleType {
            result = append(result, driver)
        }
    }
    
    return result, nil
}

// GetAllDrivers returns all drivers
func (r *InMemoryDriverRepository) GetAllDrivers() ([]*Driver, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    result := make([]*Driver, 0, len(r.drivers))
    for _, driver := range r.drivers {
        result = append(result, driver)
    }
    
    return result, nil
}

// UpdateDriverStatus updates driver availability
func (r *InMemoryDriverRepository) UpdateDriverStatus(
    driverID string,
    isAvailable bool,
) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    driver, exists := r.drivers[driverID]
    if !exists {
        return ErrDriverNotFound
    }
    
    driver.IsAvailable = isAvailable
    return nil
}

// UpdateDriverLocation updates driver's current location
func (r *InMemoryDriverRepository) UpdateDriverLocation(
    driverID string,
    location *Location,
) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    driver, exists := r.drivers[driverID]
    if !exists {
        return ErrDriverNotFound
    }
    
    driver.Location = location
    driver.Location.Timestamp = time.Now()
    return nil
}

// Helper function to calculate distance between two locations (Haversine formula)
func calculateDistance(loc1, loc2 *Location) float64 {
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

## 🗄️ Implementation 2: SQL Database Repository

### SQL Implementation (PostgreSQL Example)

```go
package repository

import (
    "database/sql"
    "time"
    
    _ "github.com/lib/pq" // PostgreSQL driver
)

// SQLDriverRepository implements DriverRepository using PostgreSQL
type SQLDriverRepository struct {
    db *sql.DB
}

// NewSQLDriverRepository creates a new SQL repository
func NewSQLDriverRepository(db *sql.DB) *SQLDriverRepository {
    return &SQLDriverRepository{db: db}
}

// Create inserts a new driver
func (r *SQLDriverRepository) Create(driver *Driver) error {
    query := `
        INSERT INTO drivers (
            id, name, phone, email, 
            latitude, longitude, address,
            vehicle_type, vehicle_id, license_plate,
            is_available, current_ride_id, rating, total_rides,
            created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
    `
    
    _, err := r.db.Exec(
        query,
        driver.ID,
        driver.Name,
        driver.Phone,
        driver.Email,
        driver.Location.Latitude,
        driver.Location.Longitude,
        driver.Location.Address,
        driver.Vehicle.Type,
        driver.Vehicle.ID,
        driver.Vehicle.LicensePlate,
        driver.IsAvailable,
        driver.CurrentRideID,
        driver.Rating,
        driver.TotalRides,
        time.Now(),
        time.Now(),
    )
    
    return err
}

// GetByID retrieves a driver by ID
func (r *SQLDriverRepository) GetByID(driverID string) (*Driver, error) {
    query := `
        SELECT 
            id, name, phone, email,
            latitude, longitude, address,
            vehicle_type, vehicle_id, license_plate,
            is_available, current_ride_id, rating, total_rides,
            created_at, updated_at
        FROM drivers
        WHERE id = $1
    `
    
    var driver Driver
    var locationTimestamp time.Time
    
    err := r.db.QueryRow(query, driverID).Scan(
        &driver.ID,
        &driver.Name,
        &driver.Phone,
        &driver.Email,
        &driver.Location.Latitude,
        &driver.Location.Longitude,
        &driver.Location.Address,
        &driver.Vehicle.Type,
        &driver.Vehicle.ID,
        &driver.Vehicle.LicensePlate,
        &driver.IsAvailable,
        &driver.CurrentRideID,
        &driver.Rating,
        &driver.TotalRides,
        &driver.CreatedAt,
        &locationTimestamp,
    )
    
    if err == sql.ErrNoRows {
        return nil, ErrDriverNotFound
    }
    if err != nil {
        return nil, err
    }
    
    driver.Location.Timestamp = locationTimestamp
    return &driver, nil
}

// Update updates an existing driver
func (r *SQLDriverRepository) Update(driver *Driver) error {
    query := `
        UPDATE drivers SET
            name = $2,
            phone = $3,
            email = $4,
            latitude = $5,
            longitude = $6,
            address = $7,
            vehicle_type = $8,
            vehicle_id = $9,
            license_plate = $10,
            is_available = $11,
            current_ride_id = $12,
            rating = $13,
            total_rides = $14,
            updated_at = $15
        WHERE id = $1
    `
    
    result, err := r.db.Exec(
        query,
        driver.ID,
        driver.Name,
        driver.Phone,
        driver.Email,
        driver.Location.Latitude,
        driver.Location.Longitude,
        driver.Location.Address,
        driver.Vehicle.Type,
        driver.Vehicle.ID,
        driver.Vehicle.LicensePlate,
        driver.IsAvailable,
        driver.CurrentRideID,
        driver.Rating,
        driver.TotalRides,
        time.Now(),
    )
    
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return ErrDriverNotFound
    }
    
    return nil
}

// Delete removes a driver
func (r *SQLDriverRepository) Delete(driverID string) error {
    query := `DELETE FROM drivers WHERE id = $1`
    
    result, err := r.db.Exec(query, driverID)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return ErrDriverNotFound
    }
    
    return nil
}

// GetAvailableDrivers finds available drivers using SQL query
func (r *SQLDriverRepository) GetAvailableDrivers(
    location *Location,
    vehicleType VehicleType,
    radius float64,
) ([]*Driver, error) {
    // Using PostGIS for spatial queries (if available)
    // Or using Haversine formula in SQL
    query := `
        SELECT 
            id, name, phone, email,
            latitude, longitude, address,
            vehicle_type, vehicle_id, license_plate,
            is_available, current_ride_id, rating, total_rides,
            created_at, updated_at,
            (
                6371 * acos(
                    cos(radians($1)) * cos(radians(latitude)) *
                    cos(radians(longitude) - radians($2)) +
                    sin(radians($1)) * sin(radians(latitude))
                )
            ) AS distance
        FROM drivers
        WHERE is_available = true
          AND current_ride_id IS NULL
          AND vehicle_type = $3
          AND (
            6371 * acos(
                cos(radians($1)) * cos(radians(latitude)) *
                cos(radians(longitude) - radians($2)) +
                sin(radians($1)) * sin(radians(latitude))
            )
          ) <= $4
        ORDER BY distance
        LIMIT 50
    `
    
    rows, err := r.db.Query(
        query,
        location.Latitude,
        location.Longitude,
        vehicleType,
        radius,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var drivers []*Driver
    for rows.Next() {
        var driver Driver
        var locationTimestamp time.Time
        var distance float64
        
        err := rows.Scan(
            &driver.ID,
            &driver.Name,
            &driver.Phone,
            &driver.Email,
            &driver.Location.Latitude,
            &driver.Location.Longitude,
            &driver.Location.Address,
            &driver.Vehicle.Type,
            &driver.Vehicle.ID,
            &driver.Vehicle.LicensePlate,
            &driver.IsAvailable,
            &driver.CurrentRideID,
            &driver.Rating,
            &driver.TotalRides,
            &driver.CreatedAt,
            &locationTimestamp,
            &distance,
        )
        if err != nil {
            return nil, err
        }
        
        driver.Location.Timestamp = locationTimestamp
        drivers = append(drivers, &driver)
    }
    
    return drivers, rows.Err()
}

// GetDriversNearLocation finds all drivers within radius
func (r *SQLDriverRepository) GetDriversNearLocation(
    location *Location,
    radius float64,
) ([]*Driver, error) {
    query := `
        SELECT 
            id, name, phone, email,
            latitude, longitude, address,
            vehicle_type, vehicle_id, license_plate,
            is_available, current_ride_id, rating, total_rides,
            created_at, updated_at
        FROM drivers
        WHERE (
            6371 * acos(
                cos(radians($1)) * cos(radians(latitude)) *
                cos(radians(longitude) - radians($2)) +
                sin(radians($1)) * sin(radians(latitude))
            )
        ) <= $3
        ORDER BY (
            6371 * acos(
                cos(radians($1)) * cos(radians(latitude)) *
                cos(radians(longitude) - radians($2)) +
                sin(radians($1)) * sin(radians(latitude))
            )
        )
        LIMIT 100
    `
    
    rows, err := r.db.Query(query, location.Latitude, location.Longitude, radius)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var drivers []*Driver
    for rows.Next() {
        var driver Driver
        var locationTimestamp time.Time
        
        err := rows.Scan(
            &driver.ID,
            &driver.Name,
            &driver.Phone,
            &driver.Email,
            &driver.Location.Latitude,
            &driver.Location.Longitude,
            &driver.Location.Address,
            &driver.Vehicle.Type,
            &driver.Vehicle.ID,
            &driver.Vehicle.LicensePlate,
            &driver.IsAvailable,
            &driver.CurrentRideID,
            &driver.Rating,
            &driver.TotalRides,
            &driver.CreatedAt,
            &locationTimestamp,
        )
        if err != nil {
            return nil, err
        }
        
        driver.Location.Timestamp = locationTimestamp
        drivers = append(drivers, &driver)
    }
    
    return drivers, rows.Err()
}

// GetDriversByVehicleType finds drivers by vehicle type
func (r *SQLDriverRepository) GetDriversByVehicleType(
    vehicleType VehicleType,
) ([]*Driver, error) {
    query := `
        SELECT 
            id, name, phone, email,
            latitude, longitude, address,
            vehicle_type, vehicle_id, license_plate,
            is_available, current_ride_id, rating, total_rides,
            created_at, updated_at
        FROM drivers
        WHERE vehicle_type = $1
    `
    
    rows, err := r.db.Query(query, vehicleType)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var drivers []*Driver
    for rows.Next() {
        var driver Driver
        var locationTimestamp time.Time
        
        err := rows.Scan(
            &driver.ID,
            &driver.Name,
            &driver.Phone,
            &driver.Email,
            &driver.Location.Latitude,
            &driver.Location.Longitude,
            &driver.Location.Address,
            &driver.Vehicle.Type,
            &driver.Vehicle.ID,
            &driver.Vehicle.LicensePlate,
            &driver.IsAvailable,
            &driver.CurrentRideID,
            &driver.Rating,
            &driver.TotalRides,
            &driver.CreatedAt,
            &locationTimestamp,
        )
        if err != nil {
            return nil, err
        }
        
        driver.Location.Timestamp = locationTimestamp
        drivers = append(drivers, &driver)
    }
    
    return drivers, rows.Err()
}

// GetAllDrivers returns all drivers
func (r *SQLDriverRepository) GetAllDrivers() ([]*Driver, error) {
    query := `
        SELECT 
            id, name, phone, email,
            latitude, longitude, address,
            vehicle_type, vehicle_id, license_plate,
            is_available, current_ride_id, rating, total_rides,
            created_at, updated_at
        FROM drivers
    `
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var drivers []*Driver
    for rows.Next() {
        var driver Driver
        var locationTimestamp time.Time
        
        err := rows.Scan(
            &driver.ID,
            &driver.Name,
            &driver.Phone,
            &driver.Email,
            &driver.Location.Latitude,
            &driver.Location.Longitude,
            &driver.Location.Address,
            &driver.Vehicle.Type,
            &driver.Vehicle.ID,
            &driver.Vehicle.LicensePlate,
            &driver.IsAvailable,
            &driver.CurrentRideID,
            &driver.Rating,
            &driver.TotalRides,
            &driver.CreatedAt,
            &locationTimestamp,
        )
        if err != nil {
            return nil, err
        }
        
        driver.Location.Timestamp = locationTimestamp
        drivers = append(drivers, &driver)
    }
    
    return drivers, rows.Err()
}

// UpdateDriverStatus updates driver availability
func (r *SQLDriverRepository) UpdateDriverStatus(
    driverID string,
    isAvailable bool,
) error {
    query := `UPDATE drivers SET is_available = $2, updated_at = $3 WHERE id = $1`
    
    result, err := r.db.Exec(query, driverID, isAvailable, time.Now())
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return ErrDriverNotFound
    }
    
    return nil
}

// UpdateDriverLocation updates driver's location
func (r *SQLDriverRepository) UpdateDriverLocation(
    driverID string,
    location *Location,
) error {
    query := `
        UPDATE drivers 
        SET latitude = $2, longitude = $3, address = $4, updated_at = $5
        WHERE id = $1
    `
    
    result, err := r.db.Exec(
        query,
        driverID,
        location.Latitude,
        location.Longitude,
        location.Address,
        time.Now(),
    )
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return ErrDriverNotFound
    }
    
    return nil
}
```

---

## 📊 Database Schema (PostgreSQL)

### SQL Migration Script

```sql
-- Create drivers table
CREATE TABLE drivers (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    
    -- Location
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    address TEXT,
    
    -- Vehicle information
    vehicle_type VARCHAR(50) NOT NULL,
    vehicle_id VARCHAR(255) NOT NULL,
    license_plate VARCHAR(50) NOT NULL UNIQUE,
    
    -- Status
    is_available BOOLEAN DEFAULT true,
    current_ride_id VARCHAR(255),
    
    -- Statistics
    rating DECIMAL(3, 2) DEFAULT 0.0,
    total_rides INTEGER DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_drivers_location ON drivers(latitude, longitude);
CREATE INDEX idx_drivers_vehicle_type ON drivers(vehicle_type);
CREATE INDEX idx_drivers_available ON drivers(is_available) WHERE is_available = true;
CREATE INDEX idx_drivers_current_ride ON drivers(current_ride_id) WHERE current_ride_id IS NOT NULL;

-- For spatial queries (if using PostGIS)
-- CREATE EXTENSION IF NOT EXISTS postgis;
-- ALTER TABLE drivers ADD COLUMN location GEOGRAPHY(POINT, 4326);
-- CREATE INDEX idx_drivers_location_geo ON drivers USING GIST(location);
```

---

## 🧪 Usage Example

### Using the Repository

```go
package main

import (
    "database/sql"
    "log"
    
    "your-project/repository"
    _ "github.com/lib/pq"
)

func main() {
    // Connect to database
    db, err := sql.Open("postgres", "postgres://user:password@localhost/dbname?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Create repository
    driverRepo := repository.NewSQLDriverRepository(db)
    
    // Or use in-memory for testing
    // driverRepo := repository.NewInMemoryDriverRepository()
    
    // Create a driver
    driver := &repository.Driver{
        ID:   "driver-123",
        Name: "John Doe",
        Phone: "+1234567890",
        Email: "john@example.com",
        Location: &repository.Location{
            Latitude:  40.7128,
            Longitude: -74.0060,
            Address:   "New York, NY",
        },
        Vehicle: &repository.Vehicle{
            Type:        repository.VehicleTypeCar,
            ID:          "vehicle-123",
            LicensePlate: "ABC-123",
        },
        IsAvailable: true,
    }
    
    err = driverRepo.Create(driver)
    if err != nil {
        log.Fatal(err)
    }
    
    // Find available drivers
    pickup := &repository.Location{
        Latitude:  40.7580,
        Longitude: -73.9855,
    }
    
    drivers, err := driverRepo.GetAvailableDrivers(
        pickup,
        repository.VehicleTypeCar,
        5.0, // 5km radius
    )
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Found %d available drivers", len(drivers))
}
```

---

## 🎯 Key Points

### 1. **Repository Pattern Benefits**
- ✅ Abstracts data access layer
- ✅ Easy to swap implementations (in-memory → SQL → NoSQL)
- ✅ Testable (use in-memory for unit tests)
- ✅ Single responsibility (data access only)

### 2. **Implementation Considerations**
- **Thread Safety**: Use mutexes for in-memory implementation
- **Error Handling**: Return domain-specific errors
- **Performance**: Use indexes for spatial queries
- **Transactions**: Add transaction support if needed

### 3. **Extensibility**
- Easy to add caching layer
- Easy to add Redis for real-time location updates
- Easy to add MongoDB implementation
- Easy to add GraphQL data source

---

## ✅ Summary

The `DriverRepository` implementation:
- ✅ Provides clean interface for data access
- ✅ Supports multiple storage backends
- ✅ Handles spatial queries efficiently
- ✅ Follows Repository Pattern best practices
- ✅ Easy to test and extend

**Remember**: Repository = Data Access Only, No Business Logic! 🎯

