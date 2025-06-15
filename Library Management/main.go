package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
Library
- users map[id]User
- books map[id]Book
- searchBooks map[title/author/category][]Book
- tickets map[id]Ticket
-
User struct
id
name
role - (enum)
- map[ticketId]Ticket

role -
ADMIN
LIBRARIAN
MEMBER

Book
- id
- title
- author
- category
- isAvailable

Ticket
- Id
- bookId
- userId
- bookDate
- returnDate
- dueAmount


Library
- users map[id]User
- books map[id]Book
- searchBooks map[title/author/category][]Book
- tickets map[id]Ticket


*/

type Library struct {
	users map[int]*User
	books map[int]*Book
	searchBooks map[string][]Book
	tickets map[int]*Ticket
}

type UserRole int 

const (
	ADMIN UserRole = iota
	LIBRARIAN
	MEMBER
)


type Category string

const (
	SPORTS   Category = "SPORTS"
	SCIENCE  Category = "SCIENCE"
	FICTION  Category = "FICTION"
)



type Book struct {
	id int
	title string
	author string
	category Category
	Issued bool
}

type Ticket struct {
	id int
	bookId int
	userId int
	issueDate int
	duration int
	actualReturnDate int
	overDue int
}

type User struct {
	id int
	role UserRole
	name string
	tickets map[int]Ticket
}

var userId = 0
var bookId = 0
var ticketId =0
var currentTime = 1000
/*
- showBorrowedBooks
- ShowUserActivity()
*/

func generateRandomID() int {
	return rand.Intn(10000000)
}

func (l *Library) ShowUserActivity(userId int) {
	for _, v := range l.users[userId].tickets {
		fmt.Println("User issued a book: bookId", v.bookId)
		if v.actualReturnDate != 0 {
			fmt.Println("User retured the book on ", v.actualReturnDate)
		}
	
	}
}

func (l *Library) showBorrowedBooks() {
	for k, v := range l.books {
		if v.Issued == true {
			fmt.Println("Book Id", k , "is Borrowed")
		}

	}
}


func (l *Library) BookAvailability() {
	for k, v := range l.books {
		fmt.Println("Book Id", k)
		fmt.Println("Book Title", v.title)
		fmt.Println("Book Availability", !(v.Issued))
	}
}

func (l *Library) PayTheDue(userId, ticketId int) {
	ticket := l.tickets[ticketId]
	book := l.books[ticket.bookId]
	amountDue := generateOverDueAmount(*ticket, book.category)

	ticket.overDue = amountDue
	ticket.actualReturnDate = currentTime
	l.users[userId].tickets[ticketId] = *ticket
}

func(l *Library) IssueABook(bookId int, userId int)  {
	//update the book

	tickets := l.users[userId].tickets
	for _, v := range tickets {
		if v.issueDate + v.duration < currentTime && v.overDue == 0 {
			fmt.Println("user has over due amount, book issue failed")
			return 
		}
	}
	book := *l.books[bookId]
	book.Issued = true

	book = l.UpdateABook(book)
	ticketId++
	ticket := Ticket{
		id: ticketId,
		bookId: book.id,
		userId: userId,
		issueDate: 100,
		duration: 200,
	}
	(l.tickets[ticketId]) = &ticket
	(l.users[userId]).tickets[ticketId] = ticket
}

func (l *Library) UpdateABook(book Book) Book {
	l.RemoveBook(book.id)
	return l.AddBook(book)

}

func (l *Library) RemoveBook(bookId int) {
	author := l.books[bookId].author
	category := l.books[bookId].category
	title := l.books[bookId].title

	delete(l.books, bookId)

	for i, v := range l.searchBooks[author] {
		if v.id == bookId {
			l.searchBooks[author] = append(l.searchBooks[author][:i], l.searchBooks[author][i+1:]...)
		}
	}
	for i, v := range l.searchBooks[title] {
		if v.id == bookId {
			l.searchBooks[title] = append(l.searchBooks[title][:i], l.searchBooks[title][i+1:]...)
		}
	}
	for i, v := range l.searchBooks[string(category)] {
		if v.id == bookId {
			l.searchBooks[string(category)] = append(l.searchBooks[string(category)][:i], l.searchBooks[string(category)][i+1:]...)
		}
	}
}

func (l*Library) AddBook(book Book) Book {
	bookId++
	book.id = bookId
	l.books[bookId] = &book
	l.searchBooks[book.author] = append(l.searchBooks[book.author], book)
	l.searchBooks[book.title] = append(l.searchBooks[book.title], book)
	l.searchBooks[string(book.category)] = append(l.searchBooks[string(book.category)], book)

	fmt.Println("BOOK is added ", bookId)
	return book
}

func(l *Library) SearchBook(author string, category Category, title string) []Book{
	result := []Book{}

	result = append(result, l.searchBooks[author]...)
	result = append(result, l.searchBooks[string(category)]...)
	result = append(result, l.searchBooks[title]...)

	return result
}


func generateOverDueAmount(ticket Ticket, category Category) int {
	amountDue := (currentTime- ticket.issueDate - ticket.duration)* getPrice(category)

	return amountDue
}

func getPrice(category Category) int {
	switch category {
	case SCIENCE:
		return 100
	case SPORTS:
		return 50
	case FICTION:
		return 80
	default:
		fmt.Println("Category doesn't have a price defined")
	}
	return 0
}

func (l *Library) RemoveUser(userId int) {
	
	delete(l.users, userId)
}

func (l *Library) AddUser(user User) {
	userId++	
	user.id = userId
	user.tickets = make(map[int]Ticket)

	l.users[user.id] = &user
}

func NewLibrary(users []User) Library  {
	library := Library{}
	library.users = make(map[int]*User)
	library.books = make(map[int]*Book)
	library.searchBooks = make(map[string][]Book)
	library.tickets = make(map[int]*Ticket)
	for _, user := range users {
		library.AddUser(user)
	}
	return library
}



func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println(generateRandomID())
	fmt.Println(generateRandomID())
	fmt.Println(generateRandomID())
	fmt.Println(generateRandomID())
	fmt.Println(generateRandomID())
	fmt.Println(generateRandomID())
	users := []User{
		{role: ADMIN, name: "Pradeep"},
		{role: LIBRARIAN, name :"Rahul"},	
		{role: MEMBER, name: "Chinmay"},
		{role: MEMBER, name: "Rohit"},
		{role: MEMBER, name: "Nitish"},
	}
	
	library := NewLibrary(users)
	book := Book{
		title  : "Treasure Island",
		author : "Unknown",
		category : FICTION,
	}
	book2 := Book{
		title  : "Leaving Time",
		author : "some one",
		category : FICTION,
	}
	book3 := Book{
		title  : "NCERT",
		author : "Govt.",
		category : SCIENCE,
	}
	book4 := Book{
		title  : "Just Play",
		author : "Viral Kohli.",
		category : SPORTS,
	}
	library.AddBook(book)
	library.AddBook(book2)
	library.AddBook(book3)
	library.AddBook(book4)

	books := library.SearchBook("Just Play", "", "")
	fmt.Println(books)

	library.IssueABook(2, 4)

	library.showBorrowedBooks()
	library.ShowUserActivity(4)
	fmt.Println(library.tickets[1].overDue)
	fmt.Println(library.users[4].tickets)
	library.PayTheDue(4, 1)
	fmt.Println(library.tickets[1].overDue)
	fmt.Println(library.users[4].tickets)

}