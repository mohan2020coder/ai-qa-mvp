package workflow

import (
	"agentic-sdlc-api/internal/agents"
	"context"
	"log"
	"strings"
	"sync"
)

type Node struct {
	ID        string                `json:"id"`
	Agent     agents.AgentInterface `json:"-"`
	Config    map[string]any        `json:"config"`
	Next      []string              `json:"next,omitempty"`
	Condition string                `json:"condition,omitempty"`
}

type Workflow struct {
	ID    string
	Name  string
	Nodes map[string]*Node
}

func (wf *Workflow) Run(ctx context.Context, input string) (map[string]string, error) {
	results := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// define recursive runner
	var runNode func(n *Node, in string)
	runNode = func(n *Node, in string) {
		defer wg.Done()

		log.Printf("Starting node %s with input: %s", n.ID, in)

		output, err := n.Agent.Run(ctx, in)
		if err != nil {
			log.Printf("Error in node %s: %v", n.Agent.Name(), err)
			return
		}
		log.Printf("Node %s output: %s", n.ID, output)

		mu.Lock()
		results[n.ID] = output
		mu.Unlock()

		// go to next nodes
		for _, nextID := range n.Next {
			nextNode, ok := wf.Nodes[nextID]
			if !ok {
				continue
			}

			// condition check
			if nextNode.Condition != "" && !evaluateCondition(nextNode.Condition, output) {
				continue
			}

			wg.Add(1)
			go runNode(nextNode, output)
		}
	}

	// kick off execution at roots
	rootNodes := wf.getRootNodes()
	for _, n := range rootNodes {
		wg.Add(1)
		go runNode(n, input)
	}

	wg.Wait()
	return results, nil
}

// Simple evaluator for branching conditions
func evaluateCondition(cond, output string) bool {
	if cond == "" {
		return true
	}
	// example syntax: "contains:success"
	if len(cond) > 9 && cond[:9] == "contains:" {
		substr := cond[9:]
		return contains(output, substr)
	}
	return true
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}

// Find nodes with no incoming edges (roots)
func (wf *Workflow) getRootNodes() []*Node {
	incoming := make(map[string]bool)
	for _, n := range wf.Nodes {
		for _, next := range n.Next {
			incoming[next] = true
		}
	}
	var roots []*Node
	for id, n := range wf.Nodes {
		if !incoming[id] {
			roots = append(roots, n)
		}
	}
	return roots
}
