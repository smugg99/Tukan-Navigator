package datastorer

import (
	// "smuggr.xyz/tukan-navigator/core/datastorer"
	"smuggr.xyz/tukan-navigator/core/grapher"

	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"sort"
)

func SortAndHashGraph(graphProject grapher.GraphProject) (sortedJSON string, hash string, err error) {
    jsonData, err := json.Marshal(graphProject)
    if err != nil {
        return "", "", err
    }

    var sortedMap map[string]interface{}
    err = json.Unmarshal(jsonData, &sortedMap)
    if err != nil {
        return "", "", err
    }

    sortedKeys := make([]string, 0, len(sortedMap))
    for k := range sortedMap {
        sortedKeys = append(sortedKeys, k)
    }

    sort.Strings(sortedKeys)
    sortedJSONData := make(map[string]interface{})
    for _, k := range sortedKeys {
        sortedJSONData[k] = sortedMap[k]
    }

    sortedJSONBytes, err := json.Marshal(sortedJSONData)
    if err != nil {
        return "", "", err
    }
    sortedJSON = string(sortedJSONBytes)

    hashBytes := md5.Sum(sortedJSONBytes)
    hash = hex.EncodeToString(hashBytes[:]) // Convert hashBytes to hexadecimal string

    return sortedJSON, hash, nil
}

func StoreGraph(sortedJSON, hash string) error {
	graph := Graph{
		SortedJSON: sortedJSON,
		Hash:       hash,
	}
	result := DB.Create(&graph)
	return result.Error
}

func GetGraphByHash(hash string) (grapher.GraphProject, error) {
    var graph Graph
    result := DB.Where("hash = ?", hash).First(&graph)
    if result.Error != nil {
        return grapher.GraphProject{}, result.Error
    }

    var sortedData grapher.GraphProject
    err := json.Unmarshal([]byte(graph.SortedJSON), &sortedData)
    if err != nil {
        return grapher.GraphProject{}, err
    }

    return sortedData, nil
}