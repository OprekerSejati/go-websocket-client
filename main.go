package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/websocket"
)  
  
func main() {  
	// Server WebSocket URL
	serverURL := os.Getenv("SERVER_URL")  
       if serverURL == "" {  
           log.Fatal("SERVER_URL environment variable is not set")  
       }   
  
	// Connect to the WebSocket server  
    conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)  
    if err != nil {  
        log.Fatalf("Failed to connect to server: %v", err)  
    }  
    defer conn.Close()  
  
    done := make(chan os.Signal, 1)  
    signal.Notify(done, os.Interrupt)  
    var room string  
  
    // Menu (login/register)  
    reader := bufio.NewReader(os.Stdin)  
    log.Print("Do you want to login or register? (login/register): ")  
    action, _ := reader.ReadString('\n')  
  
    // Periksa apakah action kosong  
    if len(action) > 0 {  
        action = action[:len(action)-1] // Menghapus newline  
    } else {  
        log.Println("No action provided. Exiting.")  
        return  
    }  
  
    log.Print("Enter your username: ")  
    username, _ := reader.ReadString('\n')  
    username = strings.TrimSpace(username) // Menghapus spasi di awal dan akhir  
  
    log.Print("Enter your password: ")  
    password, _ := reader.ReadString('\n')  
    password = strings.TrimSpace(password) // Menghapus spasi di awal dan akhir  
  
    // Send login or register request to the server  
    err = conn.WriteJSON(map[string]string{  
        "action":   action,  
        "username": username,  
        "password": password,  
    })  
    if err != nil {  
        log.Fatalf("Failed to send request: %v", err)  
    }  
  
    // Handle server response  
    var response map[string]string  
    err = conn.ReadJSON(&response)  
    if err != nil {  
        log.Println("Error reading response:", err)  
        return  
    }  
    if msg, ok := response["error"]; ok {  
        log.Println("Error:", msg)  
        return  
    }  
    log.Println(response["message"])  
  
    if response["message"] == "Registration successful" {  
        log.Println("Please Re-Open the program and login")  
        os.Exit(0)  
    }  
  
    // choose of action (join/dm)  
    if action == "login" {  
        for {  
            log.Print("Do you want to join a room or send a DM? (join/dm): ")  
            choice, _ := reader.ReadString('\n')  
            choice = strings.TrimSpace(choice) // Menghapus spasi di awal dan akhir  
  
            if choice == "join" {  
                log.Print("Enter room name: ")  
                room, _ := reader.ReadString('\n')  
                room = strings.TrimSpace(room) // Menghapus spasi di awal dan akhir  
  
                // Send join room request  
                err = conn.WriteJSON(map[string]string{  
                    "action":   "join",  
                    "room":     room,  
                    "username": username,  
                })  
                if err != nil {  
                    log.Println("Error joining room:", err)  
                    continue  
                }  
                log.Printf("Joined room: %s", room)  
                break  
            } else if choice == "dm" {  
                log.Print("Enter username to DM: ")  
                dmUser, _ := reader.ReadString('\n')  
                dmUser = strings.TrimSpace(dmUser) // Menghapus spasi di awal dan akhir  
  
                log.Print("Enter your message: ")  
                message, _ := reader.ReadString('\n')  
                message = strings.TrimSpace(message) // Menghapus spasi di awal dan akhir  
  
                // Send DM request  
                err = conn.WriteJSON(map[string]string{  
                    "action":  "dm",  
                    "to":      dmUser,  
                    "message": message,  
                })  
                if err != nil {  
                    log.Println("Error sending DM:", err)  
                    continue  
                }  
                log.Printf("DM sent to %s: %s", dmUser, message)  
            } else {  
                log.Println("Invalid choice. Please enter 'join' or 'dm'.")  
            }  
        }  
    }   
  
	go func() {  
		for {  
			var msg map[string]string  
			err := conn.ReadJSON(&msg)  
			if err != nil {  
				log.Println("Error reading message:", err)  
				return  
			}  

			if msg["action"] == "dm" && msg["to"] == username {
				log.Printf("DM received from %s: %s", msg["from"], msg["message"])
			}
		}  
	}()
	

	go func() {  
		for {  
			log.Print("Enter message (or /dm username message): ")  
			input, _ := reader.ReadString('\n')  
			input = input[:len(input)-1] 
  
			// Check if the input starts with "/dm"  
			if len(input) > 3 && input[:3] == "/dm" {  
				// Remove the "/dm " part  
				input = input[4:] 
  
				// Split the input into username and message  
				parts := strings.SplitN(input, " ", 2) // Memisahkan menjadi 2 bagian  
				if len(parts) < 2 {  
					log.Println("Invalid DM format. Use /dm username message.")  
					continue  
				}  
  
				dmUser := parts[0]   // Username  
				message := parts[1]   // Message  
  
				// Send DM request  
				err = conn.WriteJSON(map[string]string{  
					"action":  "dm",  
					"to":      dmUser,  
					"message": message,  
				})  
				if err != nil {  
					log.Println("Error sending DM:", err)  
					return  
				}  
				log.Printf("DM sent to %s: %s", dmUser, message)  
			} else {  
				// Send message to the server  
				err = conn.WriteJSON(map[string]string{  
					"room":     room,  
					"username": username,  
					"message":  input,  
				})  
				if err != nil {  
					log.Println("Error sending message:", err)  
					return  
				}  
			}  
		}  
	}()  
  
	// Wait for OS signal to terminate  
	<-done  
	log.Println("Closing connection...")  
}  
