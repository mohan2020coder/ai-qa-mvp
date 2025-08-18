package store

import (
	"ai-qa-mvp/orchestrator/internal/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	SaveRun(context.Context, *types.Run) error
	GetRun(context.Context, string) (*types.Run, error)
	ListRuns(context.Context) ([]*types.Run, error)
}

type mongoStore struct{ c *mongo.Collection }

func NewStore(ctx context.Context) (Store, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return nil, errors.New("MONGO_URI not set")
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "aiqa" // default
	}

	colName := os.Getenv("MONGO_COLLECTION")
	if colName == "" {
		colName = "runs" // default
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &mongoStore{c: client.Database(dbName).Collection(colName)}, nil
}

func (m *mongoStore) SaveRun(ctx context.Context, run *types.Run) error {
	_, err := m.c.ReplaceOne(ctx, map[string]any{"_id": run.ID}, run, options.Replace().SetUpsert(true))
	return err
}

func (m *mongoStore) GetRun(ctx context.Context, id string) (*types.Run, error) {
	var run types.Run
	if err := m.c.FindOne(ctx, map[string]any{"_id": id}).Decode(&run); err != nil {
		return nil, err
	}
	return &run, nil
}

func (m *mongoStore) ListRuns(ctx context.Context) ([]*types.Run, error) {
	cur, err := m.c.Find(ctx, map[string]any{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []*types.Run
	for cur.Next(ctx) {
		var r types.Run
		if err := cur.Decode(&r); err != nil {
			return nil, err
		}
		out = append(out, &r)
	}
	return out, nil
}

type fileStore struct{ base string }

func NewFileStore(base string) Store { return &fileStore{base: base} }

func (f *fileStore) SaveRun(ctx context.Context, run *types.Run) error {
	runsDir := os.Getenv("RUNS_DIR")
	if runsDir == "" {
		runsDir = "runs" // default
	}

	d := filepath.Join(f.base, runsDir, run.ID)
	if err := os.MkdirAll(d, 0o755); err != nil {
		return err
	}
	b, _ := json.MarshalIndent(run, "", "  ")
	return os.WriteFile(filepath.Join(d, "run.json"), b, 0o644)
}

func (f *fileStore) GetRun(ctx context.Context, id string) (*types.Run, error) {
	runsDir := os.Getenv("RUNS_DIR")
	if runsDir == "" {
		runsDir = "runs" // default
	}

	p := filepath.Join(f.base, runsDir, id, "run.json")
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var run types.Run
	if err := json.Unmarshal(b, &run); err != nil {
		return nil, err
	}
	return &run, nil
}

func (f *fileStore) ListRuns(ctx context.Context) ([]*types.Run, error) {
	runsDir := os.Getenv("RUNS_DIR")
	if runsDir == "" {
		runsDir = "runs" // default
	}

	root := filepath.Join(f.base, runsDir)
	es, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return []*types.Run{}, nil
		}
		return nil, err
	}

	var out []*types.Run
	for _, e := range es {
		p := filepath.Join(root, e.Name(), "run.json")
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		var r types.Run
		if err := json.Unmarshal(b, &r); err != nil {
			continue
		}
		out = append(out, &r)
	}
	return out, nil
}

func RunDir(base, id string) string {
	runsDir := os.Getenv("RUNS_DIR")
	if runsDir == "" {
		runsDir = "runs"
	}
	return filepath.Join(base, runsDir, id)
}

func ArtifactURL(id, file string) string {
	artifactBase := os.Getenv("ARTIFACT_BASE_URL")
	if artifactBase == "" {
		artifactBase = "/artifacts/runs" // default
	}
	return fmt.Sprintf("%s/%s/%s", artifactBase, id, file)
}
