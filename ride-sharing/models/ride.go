package models

// User represents a user in the system
type User struct {
	Name          string    // Immutable, hence value type
	Gender        string    // Immutable, hence value type
	Age           int       // Immutable, hence value type
	Vehicle       []*Vehicle // Slice of pointers to vehicles as a user may own multiple vehicles
	RidesOffered  int       // Count, hence value type
	RidesTaken    int       // Count, hence value type
}

// Vehicle represents a vehicle owned by a user
type Vehicle struct {
	Model       string // Immutable, hence value type
	NumberPlate string // Immutable, hence value type
}

// Ride represents a shared ride
type Ride struct {
	Driver         *User    // Pointer, as the driver is a large struct and shared
	Vehicle        *Vehicle // Pointer, as the vehicle is part of the driver
	Origin         string   // Immutable, hence value type
	Destination    string   // Immutable, hence value type
	AvailableSeats int      // Mutable but small, so value type is fine
	IsActive       bool     // Boolean status, hence value type
}
