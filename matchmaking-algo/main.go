package main

/*
Matchmaking System – Key Points

Goal: Pair players into matches based on similar skill and low latency for fair and smooth gameplay.

Queue Management: Players enter a matchmaking pool; system continuously searches for suitable matches.

Skill Matching: Use player ratings/skill levels to ensure balanced matches.

Latency Matching: Consider network proximity/ping to reduce lag.

Match Size: Support different formats (1v1, team matches, battle royale).

Trade-offs: Balance match quality vs. wait time; expand skill/latency range if no match is found quickly.

Edge Cases: Handle disconnects, parties/friends, timeouts, and smurfing/fake skills.

Output: Matched players are assigned to a game room/server for gameplay.

Structs

MatchMaking
-> workers
-> pool
	-> globalList []Players
	-> levels map[string][]Players
-> gameRooms []GameRoom

Player
-> Name
-> level
-> latency
-> isPlaying
-> isSearching

GameRoom
-> []Players
*/

