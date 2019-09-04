package tool

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/j75689/slack-bot/tool/valuechain"
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
	// calculate points
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

	// assign paramter
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

// ReplaceVariables replace variables to content
func ReplaceVariables(content []byte, variables map[string]interface{}) []byte {
	r, _ := regexp.Compile("\\$\\{(.*?)\\}")
	reply := string(content)

	for _, match := range r.FindAllStringSubmatch(reply, -1) {
		var (
			name  string
			tag   string
			value string
		)
		nameReg, _ := regexp.Compile(`##.*:(.*?)##`)
		tagReg, _ := regexp.Compile(`##(.*?):.*##`)

		nameMatch := nameReg.FindStringSubmatch(match[1])
		tagMatch := tagReg.FindStringSubmatch(match[1])

		// find variable
		if len(nameMatch) > 0 {
			name = nameMatch[1]
		} else {
			name = match[1]
		}

		if v := GetParamterValue(name, variables); v != nil {
			value = fmt.Sprintf("%v", v)
		} else {
			value = name
		}

		// process tag
		if len(tagMatch) > 1 {
			tag = tagMatch[1]
			value = valuechain.Execute(tag, value)
		}

		// replace
		reply = strings.Replace(reply, match[0], value, -1)
	}

	return []byte(reply)
}

// GetParamterValue by Layer
func GetParamterValue(path string, data map[string]interface{}) interface{} {
	layer := strings.Split(path, ".")
	return getValue(layer, data)
}

func getValue(layer []string, data interface{}) interface{} {
	if len(layer) <= 0 || data == nil {
		return data
	}

	switch data.(type) {
	case map[string]interface{}:
		data = (data.(map[string]interface{}))[layer[0]]
	default:
		return data
	}

	return getValue(layer[1:len(layer)], data)
}
