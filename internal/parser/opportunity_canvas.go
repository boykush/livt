package parser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/boykush/livt/internal/domain"
	"gopkg.in/yaml.v3"
)

// canvasYAML is the ten boxes as authored. Each box is a list because a canvas
// box holds stickies, not a paragraph: keeping them apart is what lets the
// board render one card per idea, the way it sat in the room.
type canvasYAML struct {
	SolutionIdeas      []string `yaml:"solution-ideas"`
	Problems           []string `yaml:"problems"`
	UsersAndCustomers  []string `yaml:"users-and-customers"`
	SolutionsToday     []string `yaml:"solutions-today"`
	BusinessChallenges []string `yaml:"business-challenges"`
	UserValue          []string `yaml:"user-value"`
	UserMetrics        []string `yaml:"user-metrics"`
	AdoptionStrategy   []string `yaml:"adoption-strategy"`
	BusinessImpact     []string `yaml:"business-impact"`
	Budget             []string `yaml:"budget"`
}

// opportunityCanvasYAML nests the boxes under canvas: so the file has room for
// fields that describe the canvas rather than sit in it — ubiquitous today.
type opportunityCanvasYAML struct {
	Canvas     canvasYAML `yaml:"canvas"`
	Ubiquitous []string   `yaml:"ubiquitous"`
}

func ParseOpportunityCanvas(path string) (*domain.OpportunityCanvas, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw opportunityCanvasYAML
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	key := domain.OpportunityKey{Value: strings.TrimSuffix(filepath.Base(path), ".yaml")}
	c := raw.Canvas
	return &domain.OpportunityCanvas{
		OpportunityKey:     key,
		SolutionIdeas:      c.SolutionIdeas,
		Problems:           c.Problems,
		UsersAndCustomers:  c.UsersAndCustomers,
		SolutionsToday:     c.SolutionsToday,
		BusinessChallenges: c.BusinessChallenges,
		UserValue:          c.UserValue,
		UserMetrics:        c.UserMetrics,
		AdoptionStrategy:   c.AdoptionStrategy,
		BusinessImpact:     c.BusinessImpact,
		Budget:             c.Budget,
		Ubiquitous:         raw.Ubiquitous,
	}, nil
}
