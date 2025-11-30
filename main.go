package main

import (
	"errors"
	"fmt"
	_ "net/http/pprof"
)


type Cache struct {
	data map[string]interface{}
	msg chan string
}

func (c *Cache) Get(key string) (interface{}, error) {
	val, ok := c.data[key]
	if !ok {
		return nil, errors.New("key is not present: GET")
	}

	return val, nil
}

func (c *Cache) Set(key string , value interface{}) error {
	c.data[key] = value
	return nil
}

func (c *Cache) Delete(key string) error {
	if _, exists := c.data[key]; exists {
		delete(c.data, key)
		return nil
	}
	
	return errors.New("key is not present : DELETE")
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
		msg: make(chan string, 500),
	}
}

func modify(nums *[]int) {
	(*nums)[0] = 42
	*nums = append(*nums, 35)
}

type Handler func(c *Cache, args ...interface{})

func handlerGet(c *Cache, args ...interface{})  {
	key := args[0].(string)
	val, err := c.Get(key)
	if err !=nil {
		c.msg <-  err.Error()
		return
	}
	c.msg <-  val.(string)
}

func handlerSet(c *Cache, args ...interface{})  {
	key := args[0].(string)
	value := args[1]
	err := c.Set(key, value)
	if err !=nil {
		c.msg <-  err.Error()
	}

}

func handlerDelete(c *Cache, args ...interface{})  {
	key := args[0].(string)
	err := c.Delete(key)
	if err !=nil {
	c.msg <-  err.Error()
	}
}

type Status int 

const (
	ACCEPTED Status = iota
	REJECTED
)

var x Status

func someValue()  {
	x = ACCEPTED
	fmt.Println(x)
}

func main() {
	someValue()
	y := REJECTED
	if y == REJECTED{
		fmt.Println("you are not some")
	}
}


/*
File Storage sytem with version control and file metadata

functional requirements
1. User can create/edit/delete/view files
2. User can see the file metadata
3. version can see the version_history and switch to a previous version

Non Functional Requirements
1. Data should be persistent
2. Handle huge amount of data.

File 
- id
- name
- data
- file_metadata_id

file_metadata
- id
- size
- created_at
- created_by
- version_no.
- updated_at
- accessType(Public, private)

file_user
- file_id
- user_id

User 
- id
- name
- created_at
- role(Admin, client)

version
- id
- file_id
- data
- version_no.

CreateFile(name, data, user_id)
UpdateFile(name, data, user_id)
ViewFile(name,user_id)
DeleteFile(name, user_id)
ViewVersionHistory(name, user_id)
SwitchToPrevVersion(name, version_no, user_id)
ViewMetadata(name, user_id)
EditMetadata(name, user_id, metadata) (edit name, accessType) 

*/

