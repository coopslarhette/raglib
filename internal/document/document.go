package document

type Passage struct {
	Text string
}

type Document struct {
	// List of passages that compromise the document, ranked by topk relevance to query
	Passages []Passage
}
