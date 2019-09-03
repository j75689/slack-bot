package tool

import (
	"regexp"
	"strings"
	"sync"
)

// ResolveVariables assign cmd variable
// calculate point
// ${a} -> 0.5 point
//   a 	-> 1.0 point
func ResolveVariables(cmd string, pattern []string, variables *map[string]interface{}) {
	wg := &sync.WaitGroup{}
	points := make([]float64, len(pattern))
	r := regexp.MustCompile(`(\w+)|'(.*?)'|"(.*?)"|(\$\{.*?\})`)
	cmdarray := r.FindAllString(cmd, -1)
	for pointIdx, pat := range pattern {
		wg.Add(1)
		go func(pat string, pointIdx int, points *[]float64) {
			defer wg.Done()
			temp := r.FindAllString(pat, -1)
			if len(cmdarray) != len(temp) { // length not equal give 0 point
				(*points)[pointIdx] = 0
				return
			}

			point := 0.0
			for i, t := range temp {
				if t == cmdarray[i] {
					point += 1.0
					continue
				}
				if strings.HasPrefix(t, `${`) && strings.HasSuffix(t, `}`) {
					point += 0.5
				}
			}
			(*points)[pointIdx] = point

		}(pat, pointIdx, &points)
	}
	wg.Wait()
	maxPoint := 0.0
	var matchPattern []string
	for idx, point := range points {
		if point > maxPoint {
			matchPattern = r.FindAllString(pattern[idx], -1)
		}
	}

	for idx, pat := range matchPattern {
		if strings.HasPrefix(pat, `${`) && strings.HasSuffix(pat, `}`) {
			v := cmdarray[idx]
			if strings.HasPrefix(v, `'`) && strings.HasSuffix(v, `'`) {
				v = v[1 : len(v)-1]
			}
			if strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`) {
				v = v[1 : len(v)-1]
			}
			(*variables)[pat[2:len(pat)-1]] = v
		}
	}

}
