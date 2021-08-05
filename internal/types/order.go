package types

type OrderItem struct {
	ID          string   `json:"id"`
	Quantity    int      `json:"quantity"`
	Size        string   `json:"size_id"`
	Ingredients []string `json:"ingredient_ids"`
	Extras      []string `json:"extra_ids"`
}

type OrderRequest struct {
	ID    string      `json:"id"`
	POS   string      `json:"pos"`
	Items []OrderItem `json:"items"`
}
