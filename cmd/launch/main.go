package main

import (
	"flag"
	"fmt"
	"launch/internal/config"
	"launch/internal/logger"
	"launch/internal/project"
	"os"
)

func main() {
	verbose := flag.Bool("v", false, "verbose output")
	flag.Parse()
	logger.Init(*verbose)

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: dev [-v] <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println("  init              - Initialize config")
		fmt.Println("  list              - List projects")
		fmt.Println("  add <name> <path> - Add a project")
		fmt.Println("  open <name> <profile> - Open a project")
		os.Exit(1)
	}

	switch args[0] {
	case "init":
		cfg := &config.Config{Projects: make(map[string]config.Project)}
		if err := cfg.Save(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created config at %s\n", config.DefaultPath())

	case "list":
		projects, err := project.List()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		if len(projects) == 0 {
			fmt.Println("No projects found.")
			return
		}

		fmt.Println("Projects:")
		for _, p := range projects {
			fmt.Printf("- %s\n", p.Name)
		}

	case "add":
		if len(os.Args) < 4 {
			fmt.Println("Usage: dev add <name> <path>")
			os.Exit(1)
		}

		name := os.Args[2]
		path := os.Args[3]

		err := project.Add(name, path)
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(1)
		}

		fmt.Println("Project added successfully!")

	case "open":
		if len(os.Args) < 4 {
			fmt.Println("Usage: dev open <project> <profile>")
			os.Exit(1)
		}

		err := project.Open(args[1], args[2])
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
	}
}
