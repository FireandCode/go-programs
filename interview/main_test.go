package main

import (
	"testing"
)

func setup() *Platform {
	return NewPlatform()
}

func TestAddUserAndProblem(t *testing.T) {
	platform := setup()
	err := platform.Users.AddUser(&User{ID: "u1", Name: "Test", Department: "Dept"})
	if err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}
	err = platform.Problems.AddProblem(&Problem{ID: "p1", Description: "Desc", Tag: "tag", Difficulty: Easy, Score: 10})
	if err != nil {
		t.Fatalf("AddProblem failed: %v", err)
	}
}

func TestSolveAndFetchSolvedProblems(t *testing.T) {
	platform := setup()
	platform.Users.AddUser(&User{ID: "u1", Name: "Test", Department: "Dept"})
	platform.Problems.AddProblem(&Problem{ID: "p1", Description: "Desc", Tag: "tag", Difficulty: Easy, Score: 10})
	err := platform.Users.SolveProblem("u1", "p1", 100, platform.Problems)
	if err != nil {
		t.Fatalf("SolveProblem failed: %v", err)
	}
	solved := platform.Users.FetchSolvedProblems("u1", platform.Problems)
	if len(solved) != 1 || solved[0].ID != "p1" {
		t.Fatalf("FetchSolvedProblems failed, got: %v", solved)
	}
}

func TestFetchProblemsFilterSort(t *testing.T) {
	platform := setup()
	platform.Problems.AddProblem(&Problem{ID: "p1", Description: "A", Tag: "math", Difficulty: Easy, Score: 10})
	platform.Problems.AddProblem(&Problem{ID: "p2", Description: "B", Tag: "math", Difficulty: Hard, Score: 30})
	platform.Problems.AddProblem(&Problem{ID: "p3", Description: "C", Tag: "string", Difficulty: Easy, Score: 20})
	problems := platform.Problems.FetchProblems(map[string]string{"tag": "math"}, "score")
	if len(problems) != 2 || problems[0].ID != "p2" {
		t.Fatalf("FetchProblems filter/sort failed, got: %v", problems)
	}
}

func TestGetLeader(t *testing.T) {
	platform := setup()
	platform.Users.AddUser(&User{ID: "u1", Name: "A", Department: "D"})
	platform.Users.AddUser(&User{ID: "u2", Name: "B", Department: "D"})
	platform.Problems.AddProblem(&Problem{ID: "p1", Description: "A", Tag: "t", Difficulty: Easy, Score: 10})
	platform.Problems.AddProblem(&Problem{ID: "p2", Description: "B", Tag: "t", Difficulty: Easy, Score: 20})
	platform.Users.SolveProblem("u1", "p1", 100, platform.Problems)
	platform.Users.SolveProblem("u2", "p2", 100, platform.Problems)
	leader := platform.Users.GetLeader(platform.Problems)
	if leader == nil || leader.ID != "u2" {
		t.Fatalf("GetLeader failed, got: %v", leader)
	}
}

func TestLikeAndTopNProblems(t *testing.T) {
	platform := setup()
	platform.Problems.AddProblem(&Problem{ID: "p1", Description: "A", Tag: "t", Difficulty: Easy, Score: 10})
	platform.Problems.AddProblem(&Problem{ID: "p2", Description: "B", Tag: "t", Difficulty: Easy, Score: 20})
	platform.Problems.LikeProblem("p1")
	platform.Problems.LikeProblem("p2")
	platform.Problems.LikeProblem("p2")
	top := platform.Problems.GetTopNProblems("t", 1)
	if len(top) != 1 || top[0].ID != "p2" {
		t.Fatalf("GetTopNProblems failed, got: %v", top)
	}
}

func TestRecommendProblems(t *testing.T) {
	platform := setup()
	platform.Users.AddUser(&User{ID: "u1", Name: "A", Department: "D"})
	platform.Problems.AddProblem(&Problem{ID: "p1", Description: "A", Tag: "t", Difficulty: Easy, Score: 10})
	platform.Problems.AddProblem(&Problem{ID: "p2", Description: "B", Tag: "t", Difficulty: Easy, Score: 20})
	platform.Problems.AddProblem(&Problem{ID: "p3", Description: "C", Tag: "x", Difficulty: Easy, Score: 30})
	platform.Users.SolveProblem("u1", "p1", 100, platform.Problems)
	recs := platform.Recommendations.RecommendProblems("u1", "p1")
	found := false
	for _, rec := range recs {
		if rec.ID == "p2" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("RecommendProblems failed, expected p2 in recommendations, got: %v", recs)
	}
} 