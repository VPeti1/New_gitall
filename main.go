package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func sudoCheck() {
	if syscall.Getuid() != 0 {
		fmt.Println("Please rerun the program as root!")
		os.Exit(1)
	}
}

func main() {
	// Check the number of arguments
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Printf("Usage: gitall <init|config|print|git> [git command]\n")
		os.Exit(1)
	}
	sudoCheck()
	filename := "/usr/gitall.db"
	switch os.Args[1] {
	case "init":
		// Handle the "init" argument
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer file.Close()
		stat, err := file.Stat()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		if stat.Size() == 0 {
			fmt.Println("Created", filename)
		} else {
			fmt.Println(filename, "already exists.")
			fmt.Println("Do you want to reinit the database?")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input == "y" || input == "yes" {
				file.Truncate(0)
				fmt.Println("Reinitialized", filename)
			}
		}
	case "config":
		// Handle the "config" argument
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer file.Close()
		fi, err := file.Stat()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		if fi.Size() > 0 {
			fmt.Println("Data already exists in", filename)
			return
		}
		fmt.Println("How many paths do you want?")
		reader := bufio.NewReader(os.Stdin)
		pathStr, _ := reader.ReadString('\n')
		pathStr = strings.TrimSpace(pathStr)
		lineInt, err := strconv.Atoi(pathStr)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		for i := 1; i <= lineInt; i++ {
			fmt.Printf("Path %d? ", i)
			path, _ := reader.ReadString('\n')
			path = strings.TrimSpace(path)
			if path == "exit" {
				break
			}
			fmt.Fprintln(file, path)
		}
		fmt.Println("Data saved to", filename)
	case "print":
		// Handle the "print" argument
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		count := 0
		for scanner.Scan() && count < 15 {
			fmt.Println(scanner.Text())
			count++
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "git":
		// Handle the "git" argument
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer file.Close()
		cmd := os.Args[1]
		fmt.Println("Your command: git", cmd)
		fmt.Println("Do you want to continue?")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if !(input == "y" || input == "yes") {
			os.Exit(1)
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			directory := scanner.Text()
			err := os.Chdir(directory)
			if err != nil {
				fmt.Println("Error: failed to change directory to", directory)
				os.Exit(1)
			}
			fmt.Println("Success: changed directory to", directory)
			cmd := exec.Command("git", os.Args[2:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				fmt.Println("Error:", err)
			}
		}
	default:
		fmt.Println("Invalid argument!")
		os.Exit(1)
	}
}
