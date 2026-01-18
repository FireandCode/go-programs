package main

import (
	"fmt"

	"github.com/google/btree"
)

// ----------------------
// Domain objects
// ----------------------

type Player struct {
	ID    int64
	Score int64
}

// Ordering rule for B-Tree
func (p Player) Less(q btree.Item) bool {
	other := q.(Player)

	// Higher score = higher rank
	if p.Score != other.Score {
		return p.Score > other.Score
	}

	// Tie-break: higher ID first
	return p.ID > other.ID
}

// ----------------------
// Leaderboard
// ----------------------

type LeaderBoard struct {
	tree    *btree.BTree
	players map[int64]Player
}

func NewLeaderBoard() *LeaderBoard {
	return &LeaderBoard{
		tree:    btree.New(32), // degree
		players: make(map[int64]Player),
	}
}

// Add/update score
func (lb *LeaderBoard) Update(id, score int64) {
	// If exists, remove from tree first
	if old, exists := lb.players[id]; exists {
		lb.tree.Delete(old)
	}

	updated := Player{ID: id, Score: score}
	lb.players[id] = updated
	lb.tree.ReplaceOrInsert(updated)
}

// Top K players
func (lb *LeaderBoard) Top(k int) []Player {
	result := make([]Player, 0, k)

	lb.tree.Ascend(func(item btree.Item) bool {
		if len(result) == k {
			return false
		}
		result = append(result, item.(Player))
		return true
	})

	return result
}

// ----------------------
// Demo
// ----------------------

func main() {
	lb := NewLeaderBoard()

	lb.Update(1, 100)
	lb.Update(2, 200)
	lb.Update(3, 150)
	fmt.Println("Top 2:", lb.Top(2)) // [2=>200, 3=>150]

	lb.Update(1, 300) // score bump
	fmt.Println("Top 3:", lb.Top(3)) // [1=>300, 2=>200, 3=>150]

	lb.Update(4, 300) // tie in score → ID decides
	fmt.Println("Top 4:", lb.Top(4)) // [4=>300, 1=>300, 2=>200, 3=>150]
}
