package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Run() {
	for {
		cmds := map[string]any{
			".help": displayHelp,
			".clear": clearScreen,
		}

		scanner := bufio.NewScanner(os.Stdin)
		printPrompt()

		for scanner.Scan() {
			input := cleanInput(scanner.Text())
			if cmd, exists := cmds[input]; exists {
				cmd.(func()) ()
			} else if strings.EqualFold(".exit", input) || strings.EqualFold(".quit", input) {
				return
			} else {
				handleCmd(input)
			}
			printPrompt()
		}

		// exit on io.EOF. scanner.Err() returns nil on io.EOF character.
		if err := scanner.Err(); err == nil {
			fmt.Println()
			return
		}
	}
}

func printPrompt() {
	fmt.Print("booktracker> ")
}

func printUnknown(cmd string) {
	fmt.Println(cmd, ": command not found")
}

func displayHelp() {
	fmt.Println("List of commands:")
	fmt.Println(".help		- Show available commands")
	fmt.Println(".clear		- Clear the screen")
	fmt.Println(".exit | .quit	- Close the connection")

}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func handleInvalidCmd(cmd string) {
	defer printUnknown(cmd)
}

func handleCmd(cmd string) {
	handleInvalidCmd(cmd)
}

func cleanInput(input string) string {
	output := strings.TrimSpace(input)
	output = strings.ToLower(output)
	return output
}
