package matcher

import (
	"regexp"
	"sync"

	// "github.com/google/go-dap"
	"time"
)

type matchTracker struct {
	count  int
	oldest time.Time
}

func Monitor(regex string, count int, interval_s int, lineInput <-chan string, matchOutput chan<- string, wg *sync.WaitGroup) {
	// store number of matches for each unique_name
	matches := make(map[string]matchTracker)
	matchCount := 0

	// range read lins from lineInput
	for line := range lineInput {
		if wg != nil {
			wg.Add(1)
		}
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
			if count == 1 {
				matchOutput <- results[1]
			} else if item, ok := matches[results[1]]; ok {
				if item.count < count-1 {
					matches[results[1]] = matchTracker{
						count: item.count + 1,
					}
				} else {
					matchOutput <- results[1]
					delete(matches, results[1])
				}
			} else {
				matches[results[1]] = matchTracker{1, time.Now()}
			}
		} else if len(results) == 1 {
			if count == 1 {
				matchOutput <- results[0]
			} else if matchCount < count-1 {
				matchCount++
			} else {
				matchOutput <- results[0]
				matchCount = 0
			}
		}
		if wg != nil {
			wg.Done()
		}
	}
}
