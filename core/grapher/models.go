package grapher

type Item struct {
	Node     string
	Priority int
	Index    int
}

type Graph struct {
	Nodes []string `json:"nodes"`
	Edges []Edge   `json:"edges"`
}

type Edge struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Weight int    `json:"weight"`
}

type PriorityQueue []*Item