package data_structure

import (
	"math"
	"sort"
	"strings"
)

type potential struct {
	dist float64
	path []ID
	code []ID
}

func (graph *Graph) Yen(source, destination ID, topK int) ([]float64, [][]ID, [][]ID, error) {
	var err error
	var i, j, k int
	var dijkstraDist map[ID]float64
	var dijkstraPrev map[ID]ID
	var existed bool
	var spurWeight float64
	var spurPath []ID
	var potentials []potential
	distTopK := make([]float64, topK)
	pathTopK := make([][]ID, topK)
	codeTopK := make([][]ID, topK)
	for i := 0; i < topK; i++ {
		distTopK[i] = math.Inf(1)
	}

	dijkstraDist, dijkstraPrev, err = graph.Dijkstra(source)
	if err != nil {
		return nil, nil, nil, err
	}
	distTopK[0] = dijkstraDist[destination]
	pathwithcode := getPath(dijkstraPrev, destination)
	pathTopK[0], codeTopK[0] = seperatePathWithCode(pathwithcode)

	for k = 1; k < topK; {
		for i = 0; i < len(pathTopK[k-1])-1; i++ {
			for j = 0; j < k; j++ {
				if isShareRootPath(pathTopK[j], pathTopK[k-1][:i+1]) {
					graph.DisableEdge(pathTopK[j][i], pathTopK[j][i+1])
				}
			}
			graph.DisablePath(pathTopK[k-1][:i])

			dijkstraDist, dijkstraPrev, _ = graph.Dijkstra(pathTopK[k-1][i])
			if dijkstraDist[destination] != math.Inf(1) {
				spurWeight = graph.GetPathWeight(pathTopK[k-1][:i+1]) + dijkstraDist[destination]

				pathwithcode := getPath(dijkstraPrev, destination)
				tempPath, tempCode := seperatePathWithCode(pathwithcode)
				tempCode = append(tempCode, codeTopK[k-1][:i]...)
				spurPath = mergePath(pathTopK[k-1][:i], tempPath)
				existed = false
				for _, each := range potentials {
					if isSamePath(each.path, spurPath) {
						existed = true
						break
					}
				}
				if !existed {
					potentials = append(potentials, potential{
						spurWeight,
						spurPath,
						tempCode,
					})
				}
			}

			graph.Reset()
		}

		if len(potentials) == 0 {
			break
		}
		sort.Slice(potentials, func(i, j int) bool {
			return potentials[i].dist < potentials[j].dist
		})

		if len(potentials) >= topK-k {
			for l := 0; k < topK; l++ {
				distTopK[k] = potentials[l].dist
				pathTopK[k] = potentials[l].path
				codeTopK[k] = potentials[l].code
				k++
			}
			break
		} else {
			distTopK[k] = potentials[0].dist
			pathTopK[k] = potentials[0].path
			codeTopK[k] = potentials[0].code
			potentials = potentials[1:]

			k++
		}
	}

	return distTopK, pathTopK, codeTopK, nil
}

func isShareRootPath(path, rootPath []ID) bool {
	if len(path) < len(rootPath) {
		return false
	}

	return isSamePath(path[:len(rootPath)], rootPath)
}
func seperatePathWithCode(arr []ID) (path, code []ID) {
	for _, val := range arr {
		splitArr := strings.Split(val.(string), "_")
		if len(splitArr) > 1 {
			code = append(code, splitArr[1])
		}
		path = append(path, splitArr[0])
	}
	return
}

func isSamePath(path1, path2 []ID) bool {
	if len(path1) != len(path2) {
		return false
	}

	for i := 0; i < len(path1); i++ {
		if path1[i] != path2[i] {
			return false
		}
	}

	return true
}

func mergePath(path1, path2 []ID) []ID {
	newPath := []ID{}
	newPath = append(newPath, path1...)
	newPath = append(newPath, path2...)

	return newPath
}

func GetJoinID(vall []ID) (res string) {
	for index, val := range vall {
		if index != 0 {
			res += "_"
		}
		res += val.(string)
	}
	return
}
