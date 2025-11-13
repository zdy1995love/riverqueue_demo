package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"zdy/worker/addone"
	"zdy/worker/addthree"
	"zdy/worker/multiplytwo"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

// Config represents the configuration structure
type Config struct {
	RiverDatabaseURL string `json:"river_database_url"`
}

// loadConfig loads configuration from file
func loadConfig() (*Config, error) {
	data, err := os.ReadFile("setting/config_DEV.jsonc")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func main() {
	ctx := context.Background()

	fmt.Println("🚀 River Queue Task Insertion Example")
	fmt.Println("================================")

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	dbPool, err := pgxpool.New(ctx, config.RiverDatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbPool.Close()

	fmt.Println("✅ Database connected successfully")

	// Create River client (for task insertion only)
	riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{})
	if err != nil {
		log.Fatal("Failed to create River client:", err)
	}

	fmt.Println("\n📝 Inserting test tasks...")

	// Task 1: AddOne(5)
	job1, err := riverClient.Insert(ctx, addone.Args{Number: 5}, nil)
	if err != nil {
		log.Printf("❌ Failed to insert task: %v", err)
	} else {
		fmt.Printf("✅ Task inserted [ID: %d] AddOne(5)\n", job1.Job.ID)
	}

	// Task 2: MultiplyTwo(10) - will auto-trigger AddThree(20)
	job2, err := riverClient.Insert(ctx, multiplytwo.Args{Number: 10}, nil)
	if err != nil {
		log.Printf("❌ Failed to insert task: %v", err)
	} else {
		fmt.Printf("✅ Task inserted [ID: %d] MultiplyTwo(10) → will trigger AddThree(20)\n", job2.Job.ID)
	}

	// Task 3: AddOne(100)
	job3, err := riverClient.Insert(ctx, addone.Args{Number: 100}, nil)
	if err != nil {
		log.Printf("❌ Failed to insert task: %v", err)
	} else {
		fmt.Printf("✅ Task inserted [ID: %d] AddOne(100)\n", job3.Job.ID)
	}

	// Task 4: MultiplyTwo(7) - will auto-trigger AddThree(14)
	job4, err := riverClient.Insert(ctx, multiplytwo.Args{Number: 7}, nil)
	if err != nil {
		log.Printf("❌ Failed to insert task: %v", err)
	} else {
		fmt.Printf("✅ Task inserted [ID: %d] MultiplyTwo(7) → will trigger AddThree(14)\n", job4.Job.ID)
	}

	// Task 5: Direct AddThree insertion
	job5, err := riverClient.Insert(ctx, addthree.Args{Number: 50}, nil)
	if err != nil {
		log.Printf("❌ Failed to insert task: %v", err)
	} else {
		fmt.Printf("✅ Task inserted [ID: %d] AddThree(50)\n", job5.Job.ID)
	}

	fmt.Println("\n✅ All tasks inserted successfully")
	fmt.Println("💡 Check the River Queue main program output to see task execution")
}
