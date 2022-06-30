package my

import (
	"github.com/jinzhu/gorm"
)

// User model
type User struct {
	gorm.Model
	Account  string
	Name     string
	Password string
	Message  string
}

type Memo struct {
	gorm.Model
	Address string
	Title   string
	Message string
	UserId  int
	ListId  int
}

type List struct {
	gorm.Model
	UserId  int
	Name    string
	Message string
}

type Comment struct {
	gorm.Model
	UserId  int
	MemoId  int
	Message string
}

type CommentJoin struct {
	Comment
	User
	Memo
}
