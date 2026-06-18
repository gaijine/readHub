package domain

type SearchBook struct {
	OpenLibraryID string
	Title         string
	Author        []string
	TotalPages    *int
	CoverURL      string
}
