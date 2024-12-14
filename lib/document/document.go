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
	Passages     []Passage     `json:"passages"`
	Title        string        `json:"title"`
	Corpus       Corpus        `json:"corpus"`
	WebReference *WebReference `json:"webReference"` // Not present when Corpus is Personal
}

// WebReference represents where the document came from, so it can be referenced or cited later
type WebReference struct {
	Title         string `json:"title"`
	Link          string `json:"link"`
	DisplayedLink string `json:"displayedLink"`
	Blurb         string `json:"blurb"` // Blurb is intended to be a short summary, snippet, or preview of the full text
	Date          string `json:"date"`
	Author        string `json:"author"`
	Favicon       string `json:"favicon"`
	Thumbnail     string `json:"thumbnail"`
	APISource     string
}

// Corpus represents where the document came from
type Corpus int

const (
	Web Corpus = iota
	Personal
)

func (s Corpus) String() string {
	return [...]string{"web", "personal"}[s]
}

func (s Corpus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Corpus) UnmarshalJSON(data []byte) error {
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
