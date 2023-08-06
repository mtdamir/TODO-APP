package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type Task struct {
	ID         int
	Title      string
	DueDate    string
	CategoryID int
	IsDone     bool
	UserID     int
}

type Category struct {
	ID     int
	Title  string
	Color  string
	UserID int
}

// Storage Layer
// userStorage is Global variable
var userStorage []User
var authenthicatedUser *User

var categoryStorage []Category
var taskStorage []Task

const userStoragePath = "user.txt"

func main() {

	// load user storage from file
	loadUserStorageFromFile()

	fmt.Println("Hello TODO app")

	command := flag.String("command", "no command", "command to run")
	flag.Parse()

	// This is a loop for the runCommand function that iterates over all its values.
	for {
		runCommand(*command)

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Please enter another command")
		scanner.Scan()
		*command = scanner.Text()
	}

}

func runCommand(command string) {
	// zero value equal nil so we compare with nil
	if command != "register-user" && command != "exit" && authenthicatedUser == nil {
		login()
		if authenthicatedUser == nil {
			return
		}
	}

	switch command {
	case "create-task":
		createTask()
	case "create-category":
		createCategory()
	case "register-user":
		registerUser()
	case "list-task":
		listTask()
	case "login":
		login()
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("command is not valid", command)
	}
}

func createTask() {
	// We use this line below to get input from the user
	scanner := bufio.NewScanner(os.Stdin)
	var title, category, duedate string

	fmt.Println("please enter the task title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("please enter the task category id")
	scanner.Scan()
	category = scanner.Text()

	categoryID, err := strconv.Atoi(category)
	if err != nil {
		fmt.Printf("category id is not valid integer. %v\n", err)

		return
	}

	isFound := false
	for _, c := range categoryStorage {
		if c.ID == categoryID && c.UserID == authenthicatedUser.ID {
			isFound = true

			break
		}
	}

	if !isFound {
		fmt.Printf("category-id is not found\n")

		return
	}

	fmt.Println("please enter the task date")
	scanner.Scan()
	duedate = scanner.Text()

	task := Task{
		ID:         len(taskStorage) + 1,
		Title:      title,
		DueDate:    duedate,
		CategoryID: categoryID,
		IsDone:     false,
		UserID:     authenthicatedUser.ID,
	}

	taskStorage = append(taskStorage, task)

	fmt.Println("task:", title, duedate, category)
}

func createCategory() {
	scanner := bufio.NewScanner(os.Stdin)
	var title, color string

	fmt.Println("please enter the category title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("please enter the category color")
	scanner.Scan()
	color = scanner.Text()

	c := Category{
		ID:     len(categoryStorage) + 1,
		Title:  title,
		Color:  color,
		UserID: authenthicatedUser.ID,
	}

	categoryStorage = append(categoryStorage, c)

	fmt.Println("catagory:", title, color)
}

func registerUser() {
	scanner := bufio.NewScanner(os.Stdin)
	var id, name, email, password string

	fmt.Println("please enter the  name")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("please enter the  email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter the  password")
	scanner.Scan()
	password = scanner.Text()

	id = email

	fmt.Println("user:", id, name, email, password)

	user := User{
		ID:       len(userStorage) + 1,
		Name:     name,
		Email:    email,
		Password: password,
	}

	userStorage = append(userStorage, user)

	var file *os.File

	file, err := os.OpenFile(userStoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("can't create or open file", err)

		return
	}

	data := fmt.Sprintf("id: %d, name: %s, email: %s, password: %s\n", user.ID, user.Name,
		user.Email, user.Password)

	var b = []byte(data)

	numberOfWrittenBytes, wErr := file.Write(b)
	if wErr != nil {
		fmt.Printf("can't write to the file %v\n", wErr)

		return
	}

	fmt.Println("numberOfWrittenBytes", numberOfWrittenBytes)

	file.Close()

}

func login() {
	fmt.Println("login process")
	scanner := bufio.NewScanner(os.Stdin)
	var email, password string

	fmt.Println("please enter the email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter the password")
	scanner.Scan()
	password = scanner.Text()

	// get the email and password from the client
	for _, user := range userStorage {
		if user.Email == email && user.Password == password {
			authenthicatedUser = &user

			break
		}
	}

	if authenthicatedUser == nil {
		fmt.Println("the email or password is not correct")

		return
	}

}

func listTask() {
	for _, task := range taskStorage {
		if task.UserID == authenthicatedUser.ID {
			fmt.Println(task)
		}
	}
}

func loadUserStorageFromFile() {
	file, err := os.Open(userStoragePath)
	if err != nil {
		fmt.Println("can't open the file", err)
	}

	var data = make([]byte, 1024)
	_, oErr := file.Read(data)
	if oErr != nil {
		fmt.Println("can't read from the file", oErr)
	}

	var dataStr = string(data)

	userSlice := strings.Split(dataStr, "\n")
	fmt.Println("userSlice", len(userSlice))
	for _, u := range userSlice {
		if u == "" {
			continue
		}

		var user = User{}

		userFields := strings.Split(u, ",")
		for _, field := range userFields {
			values := strings.Split(field, ": ")
			if len(values) != 2 {
				fmt.Println("field is not valid, skipping...", len(values))

				continue
			}
			fieldName := strings.ReplaceAll(values[0], " ", "")
			fieldValue := values[1]

			switch fieldName {
			case "id":
				id, err := strconv.Atoi(fieldValue)
				if err != nil {
					fmt.Println("strconv error", err)

					return
				}
				user.ID = id
			case "name":
				user.Name = fieldValue
			case "email":
				user.Email = fieldValue
			case "password":
				user.Password = fieldValue
			}

		}

		fmt.Printf("user: %+v\n", user)
	}

	//fmt.Println(data)
}
