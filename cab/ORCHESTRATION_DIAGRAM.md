# Service Orchestration Visual Guide

## 🎯 Who Orchestrates What?

### Simple Answer:
**RideService orchestrates all other services!**

---

## 📊 Visual Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Request Layer                     │
│  POST /api/rides/request                                    │
│  POST /api/rides/accept                                     │
│  POST /api/rides/start                                      │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Handler Layer                            │
│  - RideHandler.RequestRide()                               │
│  - RideHandler.AcceptRide()                                 │
│  - Validates HTTP input                                     │
│  - Calls RideService                                        │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              ORCHESTRATOR: RideService                      │
│  ┌─────────────────────────────────────────────────────┐  │
│  │ RequestRide()                                        │  │
│  │  1. Validate rider                                    │  │
│  │  2. Calculate fare → FareService                     │  │
│  │  3. Find drivers → MatchingService                    │  │
│  │  4. Create ride → RideRepository                     │  │
│  │  5. Notify drivers → NotificationService             │  │
│  └─────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐  │
│  │ AcceptRide()                                          │  │
│  │  1. Get ride                                          │  │
│  │  2. Generate OTP → AuthenticationService            │  │
│  │  3. Update ride                                       │  │
│  │  4. Notify rider → NotificationService               │  │
│  │  5. Start tracking → LocationService                 │  │
│  └─────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐  │
│  │ StartRide()                                          │  │
│  │  1. Validate OTP → AuthenticationService            │  │
│  │  2. Update ride state                                │  │
│  │  3. Start tracking → LocationService                 │  │
│  │  4. Notify rider → NotificationService               │  │
│  └─────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐  │
│  │ EndRide()                                            │  │
│  │  1. Calculate final fare → FareService               │  │
│  │  2. Process payment → PaymentService                 │  │
│  │  3. Stop tracking → LocationService                  │  │
│  │  4. Notify both → NotificationService               │  │
│  └─────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────┘
                         │
         ┌───────────────┼───────────────┬───────────────┐
         │               │               │               │
         ▼               ▼               ▼               ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│   Fare       │ │  Matching    │ │ Notification │ │  Location    │
│   Service    │ │  Service     │ │  Service    │ │  Service     │
│              │ │              │ │              │ │              │
│ Calculate   │ │ Find         │ │ Notify       │ │ Track        │
│ Fare        │ │ Drivers       │ │ Users        │ │ Location     │
└──────┬───────┘ └──────┬───────┘ └──────┬───────┘ └──────┬───────┘
       │                │                │                │
       │                │                │                │
       ▼                ▼                │                ▼
┌──────────────┐ ┌──────────────┐       │       ┌──────────────┐
│ Pricing      │ │ Driver       │       │       │ Location     │
│ Repository   │ │ Repository   │       │       │ Repository   │
│ (if needed)  │ │              │       │       │              │
└──────────────┘ └──────────────┘       │       └──────────────┘
                                        │
                                        ▼
                               ┌──────────────┐
                               │ Notification │
                               │ Channels     │
                               │ (SMS/Push)   │
                               └──────────────┘
```

---

## 🔄 Request Ride Flow (Step by Step)

```
1. HTTP Request
   POST /api/rides/request
   {
     "rider_id": "rider-123",
     "pickup": {...},
     "drop": {...}
   }
   │
   ▼
2. RideHandler.RequestRide()
   - Parse JSON
   - Validate input
   │
   ▼
3. RideService.RequestRide()  ← ORCHESTRATOR STARTS
   │
   ├─→ 3a. RiderRepository.GetByID()
   │      (Check if rider exists, has active ride)
   │
   ├─→ 3b. FareService.CalculateFare()
   │      (Calculate fare estimate)
   │
   ├─→ 3c. MatchingService.FindAvailableDrivers()
   │      │
   │      └─→ DriverRepository.GetDriversNearLocation()
   │          (Get drivers within 5km)
   │
   ├─→ 3d. RideRepository.Create()
   │      (Save ride record)
   │
   └─→ 3e. NotificationService.NotifyDriver()
          (Send ride request to top 3 drivers)
   │
   ▼
4. Return Ride object to Handler
   │
   ▼
5. Handler returns JSON response
   {
     "ride_id": "ride-123",
     "fare": 25.50,
     "status": "requested"
   }
```

---

## 🎭 Role of Each Component

### 1. **Handler (HTTP Layer)**
- **Role**: Entry point, HTTP concerns
- **Does**: Parse requests, validate input, return responses
- **Does NOT**: Business logic

### 2. **RideService (Orchestrator)**
- **Role**: Coordinates all services
- **Does**: 
  - Implements complete business flows
  - Calls other services in sequence
  - Manages state transitions
  - Handles errors
- **Does NOT**: Direct data access (uses repositories)

### 3. **Specialized Services**
- **Role**: Do one thing well
- **Examples**:
  - FareService: Calculate fares
  - MatchingService: Find drivers
  - NotificationService: Send notifications
  - LocationService: Track locations
- **Does NOT**: Call each other directly

### 4. **Repositories**
- **Role**: Data access only
- **Does**: Query, create, update, delete
- **Does NOT**: Business logic

---

## 🔗 Service Dependencies

### Dependency Graph:

```
RideService (Orchestrator)
    │
    ├─→ FareService (independent)
    │
    ├─→ MatchingService
    │   └─→ DriverRepository
    │
    ├─→ NotificationService (independent)
    │
    ├─→ LocationService (independent)
    │
    ├─→ AuthenticationService (independent)
    │
    ├─→ PaymentService (independent)
    │
    ├─→ RideRepository
    │
    └─→ RiderRepository
```

### Key Points:
- ✅ RideService depends on all other services
- ✅ Other services are independent (don't know about each other)
- ✅ Services depend on repositories (for data)
- ✅ No circular dependencies

---

## 💡 Why This Design?

### Benefits:

1. **Single Point of Coordination**
   - RideService knows the complete flow
   - Easy to understand and maintain

2. **Service Independence**
   - Services can be tested independently
   - Easy to swap implementations

3. **Clear Responsibilities**
   - Each service has one job
   - Orchestrator coordinates them

4. **Scalability**
   - Services can be deployed separately
   - Can scale independently

---

## 🎯 Quick Reference

### Who Does What?

| Component | Role | Example |
|-----------|------|---------|
| **Handler** | HTTP entry point | Parse JSON, call service |
| **RideService** | Orchestrator | Coordinate all services |
| **FareService** | Calculate fare | Base + distance + time |
| **MatchingService** | Find drivers | Query, filter, rank |
| **NotificationService** | Send notifications | SMS, Push, Email |
| **LocationService** | Track locations | Update, get, stream |
| **Repository** | Data access | Query database |

### Flow Direction:

```
HTTP → Handler → RideService → [Other Services] → Repositories
```

**Remember**: RideService is the conductor of the orchestra! 🎼

