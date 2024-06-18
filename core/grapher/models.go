package grapher

type Item struct {
	Node     string
	Priority int
	Index    int
}

type Edge struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Weight int    `json:"weight"`
}

type Graph struct {
	Nodes []string `json:"nodes"`
	Edges []Edge   `json:"edges"`
}

type NodeProject struct {
	ID        string  `json:"id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Traversed bool    `json:"-"`
}

type EdgeProject struct {
	ID     string `json:"id"`
	From   string `json:"from"`
	To     string `json:"to"`
	Weight int    `json:"weight"`
}

type GraphProject struct {
	Nodes []NodeProject `json:"nodes"`
	Edges []EdgeProject `json:"edges"`
}

type PriorityQueue []*Item
