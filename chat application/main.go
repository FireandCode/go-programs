package main

// import (
// 	"fmt"
// 	"sync"
// )

// func main() {
// 	fmt.Printf("%f is a good number", 1.55)
// }

// /*
// chat application with 1-1 and group chat feature
// functional requirements
// user can send/receive messages from direct/groups.
// user can create/delete/update groups
// user will notify when they relogin in to the system for all the pending messages
// user can edit/delete messages
// group will have access control

// non functional requirements
// - sent/receive APIs should have low latency.
// - message data should be persistent & uncorrupted.

// user
// -id
// - name
// - created_at
// - last_read_at

// group
// - id
// - name
// - owner_id
// - access(private, public, direct_message)

// group_users
// - user_id
// - group_id
// - role(Admin, member, owner)

// message
// - id
// - group_id
// - text
// - created_at
// - created_by(user_id)
// - updated_at

// message_user <= this is to check which users have seen a message
// - user_id
// - message_id
// - created_at

// */
// type AccessType int

// const (
// 	PUBLIC AccessType = iota
// 	PRIVATE
// 	DIRECT_MESSAGE
// )

// type User struct {
// 	id int
// 	name string
// 	createdAt int
// 	updatedAt int
// 	lastReadAt int
// 	groups map[int]*Group
// }

// type Group struct {
// 	id int
// 	name string
// 	createdAt int
// 	updatedAt int
// 	createdBy int
// 	accessType AccessType
// 	users map[int]*User
// 	messages map[int]*Message
// 	mu *sync.Mutex
// }

// type Message struct {
// 	id int
// 	createdAt int
// 	updatedAt int
// 	createdBy int
// 	text string
// }

// /*
// APIs
// sendAMessage
// ReceiveAMessage
// DeleteAMessage
// EditAMessage
// CreateGroup
// EditGroup
// DeleteGroup
// UpdateAccessToAGroup
// LogIn()
// SignUp()

// */

// type ChatApp struct {
// 	users map[int]*User
// 	groups map[int]*Group
// 	messages map[int]*Message
// 	createdAt int
// }
// func(ch *ChatApp) sendAMessage(userID int, text string, groupID int) {

// }

// func(ch *ChatApp) receiveMessage() {

// }

// func(ch *ChatApp) DeleteMessage() {

// }
// func(ch *ChatApp) EditMessage() {

// }
// func(ch *ChatApp) CreateGroup() {

// }
// func(ch *ChatApp) EditGroup() {

// }
// func(ch *ChatApp) DeleteGroup() {

// }
// func(ch *ChatApp) LogIn() {

// }
// func(ch *ChatApp) SignUp() {

// }

// func NewGroup() *Group {

// }

// func NewUser() *User {

// }

// func NewChatApp() ChatApp {

// }

// func main() {

// }
