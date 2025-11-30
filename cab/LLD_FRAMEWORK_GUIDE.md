# Low-Level Design (LLD) Framework Guide

## 🎯 Systematic Approach to LLD Problems

### Phase 1: Requirements Analysis & Clarification

#### 1.1 Understand the Problem

- **Read carefully**: Understand what the system should do
- **Ask clarifying questions**:
  - Scale (number of users, requests per second)
  - Constraints (time, space, consistency requirements)
  - Edge cases (what if driver rejects? what if rider cancels mid-ride?)
  - Non-functional requirements (availability, latency, etc.)

#### 1.2 Identify Core Use Cases

List all the actions users can perform:

- **Rider**: Request ride, cancel ride, track driver, pay
- **Driver**: Accept/reject ride, update location, start/end ride
- **System**: Match driver-rider, calculate fare, track locations

#### 1.3 Define Actors

Who interacts with the system?

- Primary actors: Rider, Driver
- Secondary actors: Admin, Payment Gateway, Maps Service

---

### Phase 2: Flow Design & Use Case Analysis

**Why Flows First?** Understanding how the system works helps identify what entities, state, and data you actually need. This prevents over-engineering and ensures you build what's required.

**Benefits of Flow-First Approach:**

- ✅ **Accurate Entity Identification**: You identify entities based on what's actually needed in flows, not assumptions
- ✅ **Discover Hidden Requirements**: Flows reveal edge cases and state transitions you might miss
- ✅ **Better Relationships**: Understanding interactions helps define entity relationships correctly
- ✅ **Focused Design**: You only build what's needed, avoiding unnecessary complexity
- ✅ **Natural Service Boundaries**: Services emerge naturally from flow steps

#### 2.1 Design Core User Flows

Walk through each use case step-by-step:

**Request Ride Flow:**

```
1. Rider provides pickup & drop location, vehicle type
2. System calculates fare estimate
3. System finds available drivers nearby
4. System sends ride request to drivers
5. Driver accepts/rejects
6. If accepted: System creates ride, generates OTP, notifies rider
7. System starts tracking driver location
```

**Start Ride Flow:**

```
1. Driver arrives at pickup location
2. Driver enters OTP
3. System validates OTP
4. System updates ride state to STARTED
5. System starts tracking both locations
6. System notifies rider
```

**End Ride Flow:**

```
1. Driver arrives at drop location
2. Driver ends ride
3. System calculates final fare
4. System processes payment
5. System prompts for ratings
6. System updates ride state to COMPLETED
```

#### 2.2 Identify Data Needs from Flows

From the flows, identify what data you need to track:

- **Rider data**: ID, location, active rides
- **Driver data**: ID, location, vehicle, availability, current ride
- **Ride data**: ID, rider, driver, pickup/drop locations, state, fare, OTP
- **Location data**: Coordinates, timestamp
- **Fare data**: Base fare, distance, time, surge multiplier, total

#### 2.3 Identify State Transitions

From flows, identify state changes:

- **Ride states**: REQUESTED → ACCEPTED → STARTED → COMPLETED/CANCELLED
- **Driver states**: AVAILABLE → ASSIGNED → IN_RIDE → AVAILABLE
- **Rider states**: IDLE → WAITING_FOR_DRIVER → IN_RIDE → IDLE

#### 2.4 Identify Operations Needed

From flows, identify what operations are required:

- Calculate fare
- Find matching drivers
- Send notifications
- Track locations
- Validate OTP
- Process payment

---

### Phase 3: Entity Identification & Modeling

**Now that you understand the flows, identify entities more accurately.**

#### 3.1 Derive Entities from Flows

Based on the flows, identify entities:

- **User** (abstract) → **Rider**, **Driver** (need to track location, state)
- **Ride** (needs pickup/drop, state, fare, OTP)
- **Vehicle** (needed for matching and fare calculation)
- **Location** (needed for tracking and matching)
- **Fare** (needed for pricing)
- **Payment** (needed for transaction)

#### 3.2 Define Entity Relationships

From flows, understand relationships:

- Rider has many Rides (1:N) - rider can have multiple ride history
- Driver has many Rides (1:N) - driver can have multiple rides
- Ride has one Vehicle (1:1) - each ride uses one vehicle
- Ride has one Fare (1:1) - each ride has one fare calculation
- Ride has Pickup and Drop locations (1:2) - two locations per ride

#### 3.3 Design Entity Attributes

For each entity, define attributes based on flow requirements:

- **Identity**: Unique identifiers (ID)
- **State**: Current status (derived from state transitions in flows)
- **Properties**: What's needed in the flows (location for tracking, OTP for verification, etc.)

---

### Phase 4: Service Layer Design

**Now that you know the flows and entities, design services that orchestrate them.**

#### 4.1 Identify Services from Flows

Group operations from flows into services:

- **RideService**: Orchestrates ride lifecycle (create, start, end, cancel)
- **MatchingService**: Finds drivers (from "find available drivers" step)
- **FareService**: Calculates fare (from "calculate fare" step)
- **LocationService**: Tracks locations (from "track location" step)
- **NotificationService**: Sends notifications (from "notify" steps)
- **PaymentService**: Processes payment (from "process payment" step)
- **AuthenticationService**: Validates OTP (from "validate OTP" step)

#### 4.2 Define Service Interfaces

Each service should have:

- Clear input/output contracts
- Well-defined responsibilities
- Minimal dependencies

---

### Phase 5: Design Patterns & Principles

#### 4.1 Key Design Patterns for LLD

**Strategy Pattern** (for algorithms):

- `FareCalculationStrategy` interface
  - `BaseFareStrategy`
  - `SurgePricingStrategy`
  - `DiscountStrategy`
- `MatchingStrategy` interface
  - `NearestDriverStrategy`
  - `RatingBasedStrategy`

**Observer Pattern** (for notifications):

- Subject: Ride state changes
- Observers: Rider, Driver, Admin

**Factory Pattern** (for object creation):

- `VehicleFactory` → creates Car, Bike, SUV based on type

**State Pattern** (for ride lifecycle):

- Ride states: Requested → Accepted → Started → Completed/Cancelled

**Repository Pattern** (for data access):

- `RideRepository`, `DriverRepository`, `RiderRepository`

#### 4.2 SOLID Principles

**S - Single Responsibility**

- Each class/service does one thing well
- Example: `FareService` only calculates fare, doesn't send notifications

**O - Open/Closed**

- Open for extension, closed for modification
- Example: Add new vehicle type without changing existing code

**L - Liskov Substitution**

- Subtypes must be substitutable for base types
- Example: Any `Vehicle` implementation works with `RideService`

**I - Interface Segregation**

- Clients shouldn't depend on interfaces they don't use
- Example: `Rider` doesn't need `AcceptRide()` method

**D - Dependency Inversion**

- Depend on abstractions, not concretions
- Example: `RideService` depends on `FareCalculator` interface, not concrete implementation

---

### Phase 6: Extensibility Considerations

#### 6.1 Make it Extensible

- **Interface-based design**: Depend on interfaces, not implementations
- **Configuration over code**: Make algorithms configurable
- **Plugin architecture**: Easy to add new features
- **Event-driven**: Decouple services using events

#### 6.2 Future Enhancements (Design for Change)

- Multiple payment methods
- Ride sharing (pool)
- Scheduled rides
- Different pricing models
- Multiple languages
- Different vehicle types
- Promotional codes
- Loyalty programs

---

### Phase 7: Error Handling & Edge Cases

#### 7.1 Common Edge Cases

- Driver rejects ride → Find another driver
- Rider cancels before pickup → Refund policy
- Driver doesn't arrive → Timeout, reassign
- OTP mismatch → Retry mechanism
- Payment failure → Retry or alternative payment
- Network issues → Retry with exponential backoff

#### 7.2 Error Handling Strategy

- Use custom exceptions/errors
- Return Result objects (success/error)
- Log errors appropriately
- Provide meaningful error messages

---

### Phase 8: Testing Strategy

#### 8.1 Unit Tests

- Test each service in isolation
- Mock dependencies
- Test edge cases

#### 8.2 Integration Tests

- Test service interactions
- Test complete flows

#### 8.3 Test Coverage

- Happy paths
- Error scenarios
- Edge cases
- Boundary conditions

---

## 📋 Checklist for LLD Problems

### Before Coding:

- [ ] Requirements clarified
- [ ] **Flows designed (step-by-step for each use case)**
- [ ] **Data needs identified from flows**
- [ ] **State transitions identified from flows**
- [ ] Entities and relationships defined (derived from flows)
- [ ] Services identified with clear responsibilities (from flow steps)
- [ ] Design patterns chosen
- [ ] Edge cases considered

### During Implementation:

- [ ] SOLID principles followed
- [ ] Interfaces used for abstraction
- [ ] Error handling implemented
- [ ] Code is readable and maintainable
- [ ] Extensibility considered

### After Implementation:

- [ ] Code reviewed
- [ ] Unit tests written
- [ ] Integration tests written
- [ ] Documentation updated

---

## 🎓 Key Takeaways

1. **Start with requirements** - Don't jump to code
2. **Identify entities first** - Nouns become classes
3. **Group operations into services** - Verbs become methods
4. **Use design patterns** - They solve common problems
5. **Design for change** - Make it extensible
6. **Handle errors gracefully** - Edge cases matter
7. **Test your design** - Think about how to test it

---

## 🔄 Iterative Refinement

LLD is iterative:

1. Start with a simple design
2. Identify gaps
3. Refactor and improve
4. Add complexity only when needed
5. Keep it simple until you can't

Remember: **Perfect is the enemy of good**. Start simple, then refine.
