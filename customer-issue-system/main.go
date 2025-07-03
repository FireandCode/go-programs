package main

import (
	"fmt"
)

// --- AgentObserver (Observer Pattern) ---
type AgentObserver interface {
	Update(issue *Issue)
}

// --- Agent ---
type Agent struct {
	ID        string
	Name      string
	Available bool
	CurrentLoad int
	Priority  int
}

func (a *Agent) SetAvailable(available bool) {
	a.Available = available
}
func (a *Agent) IncrementLoad() {
	a.CurrentLoad++
}
func (a *Agent) DecrementLoad() {
	if a.CurrentLoad > 0 {
		a.CurrentLoad--
	}
}
// Observer update method
func (a *Agent) Update(issue *Issue) {
	fmt.Printf("Agent %s notified: Assigned to issue %s\n", a.Name, issue.ID)
}

// --- Issue ---
type IssueStatus string

const (
	StatusOpen      IssueStatus = "Open"
	StatusAssigned  IssueStatus = "Assigned"
	StatusResolved  IssueStatus = "Resolved"
)

type Issue struct {
	ID             string
	Description    string
	Priority       int
	AssignedAgent  *Agent
	Status         IssueStatus
}

func (i *Issue) AssignTo(agent *Agent) {
	i.AssignedAgent = agent
	i.Status = StatusAssigned
}
func (i *Issue) SetStatus(status IssueStatus) {
	i.Status = status
}

// --- AssignmentStrategy (Strategy Pattern) ---
type AssignmentStrategy interface {
	Assign(agents []*Agent, issue *Issue) *Agent
}

// --- PriorityBasedStrategy ---
type PriorityBasedStrategy struct{}

func (s *PriorityBasedStrategy) Assign(agents []*Agent, issue *Issue) *Agent {
	var selected *Agent
	for _, agent := range agents {
		if agent.Available {
			if selected == nil || agent.Priority > selected.Priority {
				selected = agent
			}
		}
	}
	return selected
}

// --- IssueManager (Subject/Manager) ---
type IssueManager struct {
	agents    []*Agent
	issues   []*Issue
	strategy AssignmentStrategy
}

func NewIssueManager() *IssueManager {
	return &IssueManager{
		agents:    []*Agent{},
		issues:   []*Issue{},
		strategy: &PriorityBasedStrategy{}, // default
	}
}

func (im *IssueManager) RegisterAgent(agent *Agent) {
	im.agents = append(im.agents, agent)
}

func (im *IssueManager) SetStrategy(strategy AssignmentStrategy) {
	im.strategy = strategy
}

func (im *IssueManager) AddIssue(issue *Issue) {
	im.issues = append(im.issues, issue)
	im.AssignIssue(issue)
}

// AssignIssue uses the current strategy to assign the issue and notifies the agent
func (im *IssueManager) AssignIssue(issue *Issue) {
	agent := im.strategy.Assign(im.agents, issue)
	if agent != nil {
		issue.AssignTo(agent)
		agent.IncrementLoad()
		agent.SetAvailable(false) // Mark as unavailable for demo
		agent.Update(issue) // Observer notification
		fmt.Printf("Issue %s assigned to agent %s\n", issue.ID, agent.Name)
	} else {
		fmt.Printf("No available agent for issue %s\n", issue.ID)
	}
}

// --- Main: Demonstrate the flow ---
func main() {
	manager := NewIssueManager()

	// Register agents
	agent1 := &Agent{ID: "A1", Name: "Alice", Available: true, Priority: 2}
	agent2 := &Agent{ID: "A2", Name: "Bob", Available: true, Priority: 5}
	agent3 := &Agent{ID: "A3", Name: "Charlie", Available: true, Priority: 3}
	manager.RegisterAgent(agent1)
	manager.RegisterAgent(agent2)
	manager.RegisterAgent(agent3)

	// Set assignment strategy (PriorityBased)
	manager.SetStrategy(&PriorityBasedStrategy{})

	// Add an issue
	issue := &Issue{ID: "ISSUE-101", Description: "Customer cannot login", Priority: 1, Status: StatusOpen}
	manager.AddIssue(issue)

	// Add another issue (should go to next highest priority available agent)
	issue2 := &Issue{ID: "ISSUE-102", Description: "Payment failed", Priority: 2, Status: StatusOpen}
	manager.AddIssue(issue2)
} 