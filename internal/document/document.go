package document

import (
	"encoding/json"
	"fmt"
)

type Passage struct {
	Text string `json:"text"`
}

type Document struct {
	// List of passages that compromise the document, ranked by topk relevance to query
	Passages  []Passage `json:"passages"`
	Reference Reference `json:"reference"`
}

// Reference represents where the document came from, so it can be referenced or cited later
type Reference struct {
	Source Source `json:"source"`
	URL    string `json:"url"`
}

// Source represents what type of corpus the document came from
type Source int

const (
	Web Source = iota
	Personal
)

func (s Source) String() string {
	return [...]string{"web", "personal"}[s]
}

func (s Source) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Source) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	switch str {
	case "web":
		*s = Web
	case "personal":
		*s = Personal
	default:
		return fmt.Errorf("invalid source value: %s", str)
	}

	return nil
}
