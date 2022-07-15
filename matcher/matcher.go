package matcher

import (
	"regexp"

	// "github.com/google/go-dap"
	"time"
)

type matchTracker struct {
	count  int
	oldest time.Time
}

func Monitor(regex string, count int, interval_s int, lineInput <-chan string, matchOutput chan<- string) {
	// store number of matches for each unique_name
	matches := make(map[string]matchTracker)
	matchCount := 0

	// range read lins from lineInput
	for line := range lineInput {
		// fmt.Println("line:", line)
		// Delete matches that timed out
		for k := range matches {
			if time.Since(matches[k].oldest) > time.Duration(interval_s)*time.Second {
				// fmt.Println("deleting", k)
				delete(matches, k)
			}
		}
		// match line against regex
		var re = regexp.MustCompile(regex)
		results := re.FindStringSubmatch(line)

		if len(results) > 1 {
			// fmt.Println("match:", results[0], results[1])
			// if match, increment count
			if item, ok := matches[results[1]]; ok {
				// fmt.Println("found match for", results[1])
				if item.count < count-1 {
					matches[results[1]] = matchTracker{
						count: item.count + 1,
					}
					// fmt.Println("incrementing", results[1], matches[results[1]].count)
				} else {
					// if count is reached, do something
					// fmt.Println("match count reached", results[1])
					matchOutput <- results[1]
					// fmt.Println("Triggered on:", line, "with:", results[1])
					delete(matches, results[1])
				}
			} else {
				// fmt.Println("adding", results[1])
				matches[results[1]] = matchTracker{1, time.Now()}
			}
			// fmt.Println("==========================")
		} else if len(results) == 1 {
			if matchCount < count-1 {
				matchCount++
			} else {
				matchOutput <- results[0]
				// fmt.Println("Triggered on:", line)
				matchCount = 0
			}
		}
	}
}
