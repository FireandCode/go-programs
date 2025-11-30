# LLD Quick Reference Checklist

## 🚀 5-Minute LLD Framework

### Step 1: Requirements (1 min)

- [ ] Read problem statement carefully
- [ ] List all use cases (what can users do?)
- [ ] Identify actors (who uses the system?)
- [ ] Note scale/constraints if mentioned

### Step 2: Flows (2 min) ⭐ **DO THIS FIRST!**

- [ ] Design step-by-step flow for each use case
- [ ] Identify what data is needed at each step
- [ ] Identify state transitions
- [ ] Note operations needed (calculate, find, notify, etc.)

### Step 3: Entities (1 min) - Derived from Flows

- [ ] Identify entities from data needs in flows
- [ ] Define relationships based on flow interactions
- [ ] List key attributes needed in flows

### Step 4: Services (30 sec) - Derived from Flows

- [ ] Group flow operations into services
- [ ] One service = one responsibility
- [ ] Define service interfaces (contracts)

### Step 5: Patterns (30 sec)

- [ ] Strategy → for algorithms (fare, matching)
- [ ] Observer → for notifications
- [ ] State → for lifecycle (ride states)
- [ ] Factory → for object creation
- [ ] Repository → for data access

---

## 📋 Detailed Checklist

### Phase 1: Analysis

- [ ] Requirements understood
- [ ] Use cases listed
- [ ] Actors identified
- [ ] Edge cases considered

### Phase 2: Flow Design ⭐ **CRITICAL STEP**

- [ ] Flows designed for each use case
- [ ] Data needs identified from flows
- [ ] State transitions identified
- [ ] Operations needed identified

### Phase 3: Design (Derived from Flows)

- [ ] Entities identified from flow data needs
- [ ] Relationships defined from flow interactions
- [ ] Services designed from flow operations
- [ ] Interfaces created
- [ ] Patterns applied

### Phase 3: Implementation

- [ ] SOLID principles followed
- [ ] Error handling added
- [ ] Edge cases handled
- [ ] Code is extensible

### Phase 4: Validation

- [ ] All use cases covered
- [ ] Flows are complete
- [ ] Extensibility demonstrated
- [ ] Code is testable

---

## 🎯 Common Patterns by Problem Type

### Matching/Assignment Problems

- **Strategy Pattern**: Different matching algorithms
- **Observer Pattern**: Notify when match found

### State Management Problems

- **State Pattern**: Object lifecycle
- **Command Pattern**: State transitions

### Calculation Problems

- **Strategy Pattern**: Different algorithms
- **Chain of Responsibility**: Multiple calculations

### Notification Problems

- **Observer Pattern**: Subscribe to events
- **Strategy Pattern**: Different notification channels

---

## 💡 Pro Tips

1. **Design flows first** ⭐ - Understand how system works before identifying entities
2. **Start with interfaces** - Define contracts first
3. **Use composition** - Prefer composition over inheritance
4. **Single Responsibility** - One class, one reason to change
5. **Dependency Injection** - Pass dependencies, don't create them
6. **Fail Fast** - Validate inputs early
7. **Return Results** - Use Result<T, Error> pattern
8. **Log Everything** - Debugging is easier with logs

---

## 🔍 Red Flags (Things to Avoid)

❌ **God Classes** - Classes doing too much
❌ **Tight Coupling** - Classes depending on concrete implementations
❌ **No Error Handling** - Assuming everything works
❌ **Hard-coded Values** - Magic numbers/strings
❌ **No Abstraction** - Direct dependencies everywhere
❌ **Mixed Concerns** - Business logic in data access layer

---

## ✅ Green Flags (Good Design)

✅ **Small, Focused Classes** - Single responsibility
✅ **Interface-based Design** - Depend on abstractions
✅ **Comprehensive Error Handling** - All edge cases covered
✅ **Configuration-driven** - Easy to change behavior
✅ **Testable** - Easy to write unit tests
✅ **Extensible** - Easy to add new features

---

## 📚 Quick Pattern Reference

### Strategy Pattern

**When**: Multiple algorithms for same task
**Example**: Fare calculation (base, surge, discount)

### Observer Pattern

**When**: Need to notify multiple objects
**Example**: Notify rider and driver on ride updates

### State Pattern

**When**: Object behavior changes with state
**Example**: Ride lifecycle (requested → accepted → started → completed)

### Factory Pattern

**When**: Complex object creation
**Example**: Create different vehicle types

### Repository Pattern

**When**: Need to abstract data access
**Example**: Database operations for rides, drivers

---

## 🎓 Remember

> **"Make it work, make it right, make it fast"** - Kent Beck

1. **Make it work** - Get basic functionality
2. **Make it right** - Apply design principles
3. **Make it fast** - Optimize if needed

**Start simple, refactor iteratively!**
