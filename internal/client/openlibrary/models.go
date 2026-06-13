package openlibrary

type SearchResponse struct {
	Docs []SearchDoc `json:"docs"`
}

type SearchDoc struct {
	Key        string   `json:"key"`
	Title      string   `json:"title"`
	AuthorName []string `json:"author_name"`
	CoverID    int64    `json:"cover_i"`
}

type workResponse struct {
	Title  string  `json:"title"`
	Covers []int64 `json:"covers"`
}
