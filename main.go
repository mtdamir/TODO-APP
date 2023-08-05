package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type Task struct {
	ID       int
	Title    string
	DueDate  string
	Category string
	IsDone   bool
	UserID   int
}

// Storage Layer
// userStorage is Global variable
var userStorage []User
var authenthicatedUser *User

var taskStorage []Task

func (u User) print() {
	fmt.Println("User:", u.ID, u.Email, u.Name)
}

func main() {
	fmt.Println("Hello TODO app")

	command := flag.String("command", "no command", "command to run")
	flag.Parse()

	// if there is user record with corresponding data alow  the user to continue

	for {
		runCommand(*command)

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Please enter another command")
		scanner.Scan()
		*command = scanner.Text()
	}

	// input: name

}

func runCommand(command string) {

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

	scanner := bufio.NewScanner(os.Stdin)
	var title, category, duedate string

	fmt.Println("please enter the task title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("please enter the task category")
	scanner.Scan()
	category = scanner.Text()

	fmt.Println("please enter the task date")
	scanner.Scan()
	duedate = scanner.Text()

	task := Task{
		ID:       len(taskStorage) + 1,
		Title:    title,
		DueDate:  duedate,
		Category: category,
		IsDone:   false,
		UserID:   authenthicatedUser.ID,
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
