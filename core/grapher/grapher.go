package grapher

func GetShortestPath(graph Graph) ([]string, int) {
	if path, distance := Dijkstra(graph, "S", "P"); path != nil {
		return path, distance
	}
	
	return nil, 0
}