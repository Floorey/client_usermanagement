package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("What would you like to do?")
		fmt.Println("1. Create user")
		fmt.Println("2. Query user")
		fmt.Println("3. Quit")

		fmt.Print("Choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			createUser(reader)
		case "2":
			getUser(reader)
		case "3":
			fmt.Println("Application is quitting.")
			return
		default:
			fmt.Println("Invalid choice.")
		}
	}
}

func createUser(reader *bufio.Reader) {
	// Authorization check: Only users with admin privileges can create new users
	if !isAdmin() {
		fmt.Println("Only users with admin privileges can create new users.")
		return
	}

	fmt.Println("Enter user data:")
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("Role: ")
	role, _ := reader.ReadString('\n')
	role = strings.TrimSpace(role)

	user := User{
		Username: username,
		Email:    email,
		Password: password,
		Role:     role,
	}

	// Encode the user object to JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error encoding user:", err)
		return
	}

	// Send an HTTP request to the server to create the user
	resp, err := http.Post("http://localhost:8080/add-user", "application/json", strings.NewReader(string(userJSON)))

	if err != nil {
		fmt.Println("Error creating user:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("User created successfully.")
}

func getUser(reader *bufio.Reader) {
	// Authorization check: Only users with admin privileges can query users
	if !isAdmin() {
		fmt.Println("Only users with admin privileges can query users.")
		return
	}

	fmt.Print("Enter username to query: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	// Send an HTTP request to the server to get the user details
	resp, err := http.Get("http://localhost:8080/get-user?username=" + username)
	if err != nil {
		fmt.Println("Error querying user:", err)
		return
	}
	defer resp.Body.Close()

	// Check if the status code is OK
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error retrieving user. Status code:", resp.StatusCode)
		return
	}

	// Decode the JSON response object
	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	// Print the user details
	fmt.Println("User data:", user)
}

func isAdmin() bool {
	// Here you can implement the authorization check
	// For example, you could check the role of the current user
	// and return true if the user is an admin, otherwise return false.
	return true // Temporarily return true here to demonstrate functionality
}
