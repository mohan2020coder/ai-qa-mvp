package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	// DB connection string (update to your settings)
	connStr := "postgres://postgres:postgres@localhost:5432/workflowdb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer db.Close()

	// Step 1: Insert workflow
	var workflowID string
	err = db.QueryRow(
		`INSERT INTO workflows (name) VALUES ($1) RETURNING id`,
		"Shopping App Builder",
	).Scan(&workflowID)
	if err != nil {
		log.Fatal("failed to insert workflow:", err)
	}
	fmt.Println("Workflow ID:", workflowID)

	// Step 2: Define fixed node UUIDs
	nodeProduct := "f433d8f4-e90c-445d-bbb5-0cc545883c8f"
	nodeCode := "22695631-c485-46c0-831b-0f7c1d83d190"
	nodeDesign := "841219b1-d443-444b-995e-217d2a64a5da"

	// Step 3: Insert workflow nodes
	nodes := []struct {
		id       string
		nodeType string
		config   string
	}{
		{nodeProduct, "llm", `{"name":"ProductAgent","max_tokens":512,"prompt":"Generate product requirements"}`},
		{nodeCode, "llm", `{"name":"CodeAgent","max_tokens":1024,"prompt":"Generate Go backend code"}`},
		{nodeDesign, "llm", `{"name":"DesignAgent","max_tokens":512,"prompt":"Suggest UI/UX design"}`},
	}

	for _, n := range nodes {
		_, err := db.Exec(
			`INSERT INTO workflow_nodes (id, workflow_id, type, config) VALUES ($1, $2, $3, $4::jsonb)`,
			n.id, workflowID, n.nodeType, n.config,
		)
		if err != nil {
			log.Fatal("failed to insert node:", err)
		}
	}
	fmt.Println("Inserted nodes:", nodeProduct, nodeCode, nodeDesign)

	// Step 4: Insert workflow edges
	edges := []struct {
		source string
		target string
	}{
		{nodeProduct, nodeCode},
		{nodeProduct, nodeDesign},
	}

	for _, e := range edges {
		_, err := db.Exec(
			`INSERT INTO workflow_edges (id, workflow_id, source_node_id, target_node_id, condition) 
			 VALUES ($1, $2, $3, $4, NULL)`,
			uuid.New().String(), workflowID, e.source, e.target,
		)
		if err != nil {
			log.Fatal("failed to insert edge:", err)
		}
	}
	fmt.Println("Inserted edges")
}
