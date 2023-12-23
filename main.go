package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
)

func GetSummarizeResponse(client gpt3.Client, ctx context.Context, quesiton string) {
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			quesiton,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	fmt.Printf("\n")
}

func main() {
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	// Command-line flags
	inputFile := flag.String("input", "", "Path to the input text file")
	flag.Parse()

	// Validate input file flag
	if *inputFile == "" {
		fmt.Println("Please provide a valid path to the input file using -input flag")
		os.Exit(1)
	}

	// Read input file
	file, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal("Error opening input file:", err)
	}
	defer file.Close()

	// Use io.ReadAll to read the content
	inputContent, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading input file:", err)
	}
	msgPrefix := "Break down the steps of the code between ``` and explain each part in a way that is easy for someone without a technical background to understand\n```\n"
	msgSuffix := "\n```"
	msg := msgPrefix + string(inputContent) + msgSuffix

	GetSummarizeResponse(client, ctx, msg)

}
