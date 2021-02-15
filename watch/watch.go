package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	script := "script.json"
	if len(os.Args) > 1 {
		script = os.Args[1]
		report("using %s as script file...\n", script)
	}

	data, err := ioutil.ReadFile(script)
	if err != nil {
		panic(err)
	}

	var list []interface{}
	if err := json.Unmarshal(data, &list); err != nil {
		panic(err)
	}

	for _, task := range list {
		dict, converts := task.(map[string]interface{})
		if converts {
			process(dict)
		} else {
			panic("wrong format")
		}
	}
}

func process(dict map[string]interface{}) {
	name := dict["name"].(string)
	action := dict["action"].(string)

	report("checking: %s\n", name)

	from := dict["url"].(string)
	content, err := GetContent(from)
	if err != nil {
		panic(err)
	}

	content = trimPrefix(content, dict["after"])
	content = trimSuffix(content, dict["before"])

	if missing, found := dict["missing"].(string); found && len(missing) > 0 {
		if !strings.Contains(content, missing) {
			report("FOUND: %s\n", name)
			alert(name, action)
		}
	} else if present, found := dict["present"].(string); found && len(present) > 0 {
		if strings.Contains(content, present) {
			report("FOUND: %s\n", name)
			alert(name, action)
		}
	} else {
		report("ERROR: no missing not present entry is found...\n")
	}
}

func alert(name, address string) {
	sendEmail([]string{address}, name)
}

func trimPrefix(original string, key interface{}) string {
	if key == nil {
		return original
	}
	var steps []string
	switch actual := key.(type) {
	case string:
		steps = []string{actual}
	case []interface{}:
		for _, one := range actual {
			if txt, converts := one.(string); converts && len(txt) > 0 {
				steps = append(steps, txt)
			}
		}
	default:
		report("unhandled type: %v\n", actual)
	}

	for _, step := range steps {
		if found := strings.Index(original, step); len(step) > 0 && found > 0 {
			original = original[found+len(step):]
		}
	}

	original = strings.Trim(original, "\r\n\t ")
	return original
}

func trimSuffix(original string, key interface{}) string {
	if key == nil {
		return original
	}
	var steps []string
	switch actual := key.(type) {
	case string:
		steps = []string{actual}
	case []interface{}:
		for _, one := range actual {
			if txt, converts := one.(string); converts && len(txt) > 0 {
				steps = append(steps, txt)
			}
		}
	default:
		report("unhandled type: %v\n", actual)
	}

	for _, step := range steps {
		if found := strings.Index(original, step); len(step) > 0 && found > 0 {
			original = original[:found]
		}
	}

	original = strings.Trim(original, "\r\n\t ")
	return original
}

func report(format string, arg ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, format, arg)
}
