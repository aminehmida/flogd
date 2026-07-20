package matcher

import (
	"fmt"
	"regexp"
	"sync"

	"time"

	"github.com/rs/zerolog/log"
)

type matchTracker struct {
	count  int
	oldest time.Time
}

func Monitor(regex string, count int, interval_s int, lineInput <-chan string, matchOutput chan<- string, wg *sync.WaitGroup) {
	// store number of matches for each unique_name
	matches := make(map[string]matchTracker)
	matchCount := 0

	// compile the regex once, up front, rather than on every line
	re := regexp.MustCompile(regex)

	// range read lines from lineInput
	for line := range lineInput {
		log.Debug().Msgf("match map size: %d", len(matches))
		log.Debug().Msg("Line received: " + line)
		if wg != nil {
			wg.Add(1)
		}

		// Delete matches that timed out
		for k := range matches {
			if time.Since(matches[k].oldest) > time.Duration(interval_s)*time.Second {
				log.Debug().Msgf("Deleting match because timeout: %s - %s, %v > %v.", k, matches[k].oldest, time.Since(matches[k].oldest), time.Duration(interval_s)*time.Second)
				delete(matches, k)
			}
		}
		// match line against regex
		results := re.FindStringSubmatch(line)

		if len(results) > 1 {
			log.Debug().Msg("Matched with capture group")
			if count == 1 {
				matchOutput <- results[1]
			} else if item, ok := matches[results[1]]; ok {
				if item.count < count-1 {
					matches[results[1]] = matchTracker{
						count:  item.count + 1,
						oldest: matches[results[1]].oldest,
					}
					log.Debug().Msg(fmt.Sprintf("Matched group: %s, count incremented: %d", results[1], item.count+1))
				} else {
					log.Debug().Msg("Matched with capture group and count reached!")
					matchOutput <- results[1]
					delete(matches, results[1])
				}
			} else {
				log.Debug().Msg("Matched with capture group added to matches map")
				matches[results[1]] = matchTracker{1, time.Now()}
				log.Debug().Msgf("Time now: %v", matches[results[1]].oldest)
			}
		} else if len(results) == 1 {
			log.Debug().Msg("Matched with no capture group")
			if count == 1 {
				log.Debug().Msg("Matched with no capture group and count 1 reached")
				matchOutput <- results[0]
			} else if matchCount < count-1 {
				log.Debug().Msg("Matched with no capture group and count not reached")
				matchCount++
			} else {
				log.Debug().Msg("Matched with no capture group and count reached")
				matchOutput <- results[0]
				matchCount = 0
			}
		}
		if wg != nil {
			wg.Done()
		}
	}
}
