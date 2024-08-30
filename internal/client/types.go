package client

type Memo struct {
	Name    string   `json:"name"`
	ID      string   `json:"uid"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type ListMemosResponse struct {
	Memos         []Memo `json:"memos"`
	NextPageToken string `json:"nextPageToken"`
}
