package main

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
)

// Data Models

type Difficulty string

const (
	Easy   Difficulty = "easy"
	Medium Difficulty = "medium"
	Hard   Difficulty = "hard"
)

type SolveInfo struct {
	UserID    string
	TimeTaken int // in seconds
}

type Problem struct {
	ID          string
	Description string
	Tag         string
	Difficulty  Difficulty
	Score       int
	Likes       int32
	SolvedBy    map[string]SolveInfo // userID -> SolveInfo
	AvgSolveTime float64
	TotalSolveTime int
	NumSolves int
	mu sync.Mutex
}

type User struct {
	ID            string
	Name          string
	Department    string
	SolvedProblems map[string]SolveInfo // problemID -> SolveInfo
	mu            sync.Mutex
}

// Service Layers

type ProblemService struct {
	problems              map[string]*Problem // problemID -> Problem
	problemsByTag         map[string][]*Problem
	problemsByDifficulty  map[Difficulty][]*Problem
	mu                    sync.RWMutex
}

type UserService struct {
	users map[string]*User // userID -> User
	mu    sync.RWMutex
}

type RecommendationService struct {
	problemService *ProblemService
	userService    *UserService
}

// Method signatures for services (to be implemented)

// ProblemService methods
func (ps *ProblemService) AddProblem(p *Problem) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if _, exists := ps.problems[p.ID]; exists {
		return fmt.Errorf("problem with ID %s already exists", p.ID)
	}
	p.SolvedBy = make(map[string]SolveInfo)
	ps.problems[p.ID] = p
	ps.problemsByTag[p.Tag] = append(ps.problemsByTag[p.Tag], p)
	ps.problemsByDifficulty[p.Difficulty] = append(ps.problemsByDifficulty[p.Difficulty], p)
	return nil
}

func (ps *ProblemService) FetchProblems(filter map[string]string, sortBy string) []*Problem {
	var filtered []*Problem
	if tag, ok := filter["tag"]; ok {
		filtered = append(filtered, ps.problemsByTag[tag]...)
	} else if diff, ok := filter["difficulty"]; ok {
		filtered = append(filtered, ps.problemsByDifficulty[Difficulty(diff)]...)
	} else {
		for _, p := range ps.problems {
			filtered = append(filtered, p)
		}
	}
	// If both tag and difficulty are present, take intersection
	if tag, ok := filter["tag"]; ok {
		if diff, ok2 := filter["difficulty"]; ok2 {
			var intersection []*Problem
			m := make(map[string]bool)
			for _, p := range ps.problemsByDifficulty[Difficulty(diff)] {
				m[p.ID] = true
			}
			for _, p := range ps.problemsByTag[tag] {
				if m[p.ID] {
					intersection = append(intersection, p)
				}
			}
			filtered = intersection
		}
	}
	if sortBy == "score" {
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].Score > filtered[j].Score
		})
	}
	return filtered
}

func (ps *ProblemService) LikeProblem(problemID string) error {

	problem, exists := ps.problems[problemID]
	if !exists {
		return fmt.Errorf("problem with ID %s does not exist", problemID)
	}
	atomic.AddInt32(&problem.Likes, 1)
	return nil
}

func (ps *ProblemService) GetTopNProblems(tag string, n int) []*Problem {
	tagged := ps.problemsByTag[tag]
	sort.Slice(tagged, func(i, j int) bool {
		return tagged[i].Likes > tagged[j].Likes
	})
	if n > len(tagged) {
		n = len(tagged)
	}
	return tagged[:n]
}

// UserService methods
func (us *UserService) AddUser(u *User) error {
	us.mu.Lock()
	defer us.mu.Unlock()
	if _, exists := us.users[u.ID]; exists {
		return fmt.Errorf("user with ID %s already exists", u.ID)
	}
	u.SolvedProblems = make(map[string]SolveInfo)
	us.users[u.ID] = u
	return nil
}

func (us *UserService) SolveProblem(userID, problemID string, timeTaken int, ps *ProblemService) error {
	us.mu.RLock()
	user, userExists := us.users[userID]
	us.mu.RUnlock()
	if !userExists {
		return fmt.Errorf("user with ID %s does not exist", userID)
	}
	ps.mu.RLock()
	problem, problemExists := ps.problems[problemID]
	ps.mu.RUnlock()
	if !problemExists {
		return fmt.Errorf("problem with ID %s does not exist", problemID)
	}
	user.mu.Lock()
	defer user.mu.Unlock()
	problem.mu.Lock()
	defer problem.mu.Unlock()
	if _, alreadySolved := user.SolvedProblems[problemID]; alreadySolved {
		return fmt.Errorf("user %s already solved problem %s", userID, problemID)
	}
	info := SolveInfo{UserID: userID, TimeTaken: timeTaken}
	user.SolvedProblems[problemID] = info
	problem.SolvedBy[userID] = info
	// Update average time
	problem.TotalSolveTime += timeTaken
	problem.NumSolves++
	problem.AvgSolveTime = float64(problem.TotalSolveTime) / float64(problem.NumSolves)
	return nil
}

func (us *UserService) FetchSolvedProblems(userID string, ps *ProblemService) []*Problem {
	us.mu.RLock()
	user, userExists := us.users[userID]
	us.mu.RUnlock()
	if !userExists {
		return nil
	}
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	problems := []*Problem{}
	for pid := range user.SolvedProblems {
		if p, ok := ps.problems[pid]; ok {
			problems = append(problems, p)
		}
	}
	return problems
}

func (us *UserService) GetLeader(ps *ProblemService) *User {
	us.mu.RLock()
	defer us.mu.RUnlock()
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	var leader *User
	maxScore := -1
	for _, user := range us.users {
		totalScore := 0
		for pid := range user.SolvedProblems {
			if p, ok := ps.problems[pid]; ok {
				totalScore += p.Score
			}
		}
		if totalScore > maxScore {
			maxScore = totalScore
			leader = user
		}
	}
	return leader
}

// RecommendationService methods
func (rs *RecommendationService) RecommendProblems(userID, lastSolvedProblemID string) []*Problem {
	user, userExists := rs.userService.users[userID]
	if !userExists {
		return nil
	}
	lastProblem, probExists := rs.problemService.problems[lastSolvedProblemID]
	if !probExists {
		return nil
	}
	tag := lastProblem.Tag
	candidates := []*Problem{}
	for _, p := range rs.problemService.problems {
		if _, solved := user.SolvedProblems[p.ID]; solved {
			continue
		}
		if p.Tag == tag {
			candidates = append(candidates, p)
		}
	}
	// If less than 5, fill with other unsolved problems
	if len(candidates) < 5 {
		for _, p := range rs.problemService.problems {
			if _, solved := user.SolvedProblems[p.ID]; solved {
				continue
			}
			found := false
			for _, c := range candidates {
				if c.ID == p.ID {
					found = true
					break
				}
			}
			if !found {
				candidates = append(candidates, p)
			}
			if len(candidates) >= 5 {
				break
			}
		}
	}
	// Sort by score descending
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})
	if len(candidates) > 5 {
		candidates = candidates[:5]
	}
	return candidates
}

// Constructors for services
func NewProblemService() *ProblemService {
	return &ProblemService{
		problems:             make(map[string]*Problem),
		problemsByTag:        make(map[string][]*Problem),
		problemsByDifficulty: make(map[Difficulty][]*Problem),
	}
}

func NewUserService() *UserService {
	return &UserService{users: make(map[string]*User)}
}

// Add helper to Problem
func (p *Problem) Stats() (numUsers int, avgTime float64) {
	numUsers = len(p.SolvedBy)
	if numUsers == 0 {
		return 0, 0
	}
	
	return numUsers, p.AvgSolveTime
}

type LeaderboardEntry struct {
	User  *User
	Score int
}

type Leaderboard struct {
	Entries []*LeaderboardEntry
}

type Platform struct {
	Problems        *ProblemService
	Users           *UserService
	Recommendations *RecommendationService
	Leaderboard     *Leaderboard
}

func NewPlatform() *Platform {
	ps := NewProblemService()
	us := NewUserService()
	rs := &RecommendationService{problemService: ps, userService: us}
	return &Platform{
		Problems:        ps,
		Users:           us,
		Recommendations: rs,
		Leaderboard:     &Leaderboard{},
	}
}

func (p *Platform) UpdateLeaderboard() {
	p.Users.mu.RLock()
	defer p.Users.mu.RUnlock()
	p.Problems.mu.RLock()
	defer p.Problems.mu.RUnlock()

	var entries []*LeaderboardEntry
	for _, user := range p.Users.users {
		total := 0
		for pid := range user.SolvedProblems {
			if prob, ok := p.Problems.problems[pid]; ok {
				total += prob.Score
			}
		}
		entries = append(entries, &LeaderboardEntry{User: user, Score: total})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})
	p.Leaderboard.Entries = entries
}

func main() {
	platform := NewPlatform()

	// Add two users
	user1 := &User{ID: "u1", Name: "Alice", Department: "Engineering"}
	user2 := &User{ID: "u2", Name: "Bob", Department: "Science"}
	if err := platform.Users.AddUser(user1); err != nil {
		fmt.Println("Error adding user1:", err)
	}
	if err := platform.Users.AddUser(user2); err != nil {
		fmt.Println("Error adding user2:", err)
	}

	// Add two problems
	problem1 := &Problem{ID: "p1", Description: "FizzBuzz", Tag: "math", Difficulty: Easy, Score: 10}
	problem2 := &Problem{ID: "p2", Description: "Palindrome", Tag: "string", Difficulty: Easy, Score: 15}
	if err := platform.Problems.AddProblem(problem1); err != nil {
		fmt.Println("Error adding problem1:", err)
	}
	if err := platform.Problems.AddProblem(problem2); err != nil {
		fmt.Println("Error adding problem2:", err)
	}

	// Each user solves a different problem
	if err := platform.Users.SolveProblem("u1", "p1", 120, platform.Problems); err != nil {
		fmt.Println("Error solving problem1 by user1:", err)
	}
	if err := platform.Users.SolveProblem("u2", "p2", 90, platform.Problems); err != nil {
		fmt.Println("Error solving problem2 by user2:", err)
	}

	// Update and print leaderboard
	platform.UpdateLeaderboard()
	fmt.Println("Leaderboard:")
	for i, entry := range platform.Leaderboard.Entries {
		fmt.Printf("%d. %s (%s) - %d points\n", i+1, entry.User.Name, entry.User.Department, entry.Score)
	}

	// Print stats for both problems
	p1 := platform.Problems.problems["p1"]
	p2 := platform.Problems.problems["p2"]
	fmt.Printf("Problem: %s | Solved by %d user(s), Avg time: %.2fs\n", p1.Description, p1.NumSolves, p1.AvgSolveTime)
	fmt.Printf("Problem: %s | Solved by %d user(s), Avg time: %.2fs\n", p2.Description, p2.NumSolves, p2.AvgSolveTime)

	// Recommendation for user1 after solving problem1
	recommended := platform.Recommendations.RecommendProblems("u1", "p1")
	fmt.Println("Recommended problems for Alice after solving 'FizzBuzz':")
	for _, p := range recommended {
		fmt.Printf("- %s: %s (Tag: %s, Score: %d)\n", p.ID, p.Description, p.Tag, p.Score)
	}
}