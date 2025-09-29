package parser

import (
	"bufio"
	"encoding/json"
	"math"
	"os"
	"path/filepath"

	"codex_usage_report/internal/model"
)

// FindSessionFiles walks through all folders inside baseDir and returns rollout-*.jsonl files
func FindSessionFiles(baseDir string) ([]string, error) {
	var files []string
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".jsonl" && len(info.Name()) > 8 && info.Name()[0:8] == "rollout-" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// ParseFile reads a .jsonl file and returns two timelines:
// - timelineFull: with all entries (including duplicates)
// - timelineClean: only entries when Primary/Secondary values changed
// It also returns the maximum total tokens found and the sum of last task tokens.
func ParseFile(filePath string, debug bool) ([]model.TimelineEntry, []model.TimelineEntry, int, int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil, 0, 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	var (
		timelineFull  []model.TimelineEntry
		timelineClean []model.TimelineEntry
		maxTotal      int
		sumLast       int
		lastP         = -1
		lastS         = -1
	)

	for scanner.Scan() {
		line := scanner.Bytes()
		var logLine model.LogLine
		if err := json.Unmarshal(line, &logLine); err != nil {
			continue
		}

		rl := logLine.Payload.RateLimits

		primaryRaw := rl.Primary.UsedPercent
		if primaryRaw == 0 {
			primaryRaw = rl.PrimaryUsedPercent
		}

		secondaryRaw := rl.Secondary.UsedPercent
		if secondaryRaw == 0 {
			secondaryRaw = rl.SecondaryUsedPercent
		}


		if primaryRaw == 0 && secondaryRaw == 0 {
			continue
		}

		p := toDisplayPercent(primaryRaw)
		s := toDisplayPercent(secondaryRaw)

		total := logLine.Payload.Info.TotalTokenUsage.TotalTokens
		last := logLine.Payload.Info.LastTokenUsage.TotalTokens
		if total > maxTotal {
			maxTotal = total
		}
		sumLast += last

		entry := model.TimelineEntry{
			Timestamp: logLine.Timestamp,
			Primary:   p,
			Secondary: s,
		}

		if debug {
			println("DEBUG:", entry.Timestamp, "primary=", entry.Primary, "secondary=", entry.Secondary)
		}

		// always add to full timeline
		timelineFull = append(timelineFull, entry)

		// only add to clean timeline if values changed
		if p != lastP || s != lastS {
			timelineClean = append(timelineClean, entry)
			lastP, lastS = p, s
		}
	}

	if err := scanner.Err(); err != nil {
		return timelineFull, timelineClean, maxTotal, sumLast, err
	}

	return timelineFull, timelineClean, maxTotal, sumLast, nil
}

func toDisplayPercent(value float64) int {
	if value <= 1 {
		value *= 100
	}
	return int(math.Round(value))
}
