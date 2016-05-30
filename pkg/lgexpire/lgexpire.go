// Package lgexpire implements log file cleanup of old lg (glog) logs. This
// package's API's might change, escpecially when it comes to what's considered
// and which errors are returned.
package lgexpire

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Expire  .
type Expire struct {
	LogDir   string   // The directory where the log files are located
	Programs []string // the programs to consider.
	Rules    []Rule   // Additional rules for removal

	logFiles []logFile
}

// Rule for log expiery, all specified fields are considered together.
type Rule struct {
	Level string        // applies to all levels if not specified
	Age   time.Duration // keep logs newer than this
	Count uint          // keep maxium amount of logs
}

var defaultRule = Rule{
	Age:   30 * 24 * time.Hour,
	Count: 30,
}

type logFile struct {
	Filename string // full path to file
	// TODO Filesize int // log file size
	Program  string    // program name
	Host     string    // hostname
	Username string    // username
	Level    string    // level
	Time     time.Time // timestamp
	Pid      uint64    // program pid
	Ext      string    // extra file extension (.gz, etc..)
}

var noLogFile = logFile{}
var validExts = map[string]bool{
	"gz": true,
}
var allLevels = []string{"INFO", "WARNING", "ERROR", "FATAL"}
var validLevels = make(map[string]bool, 0)

func init() {
	for _, v := range allLevels {
		validLevels[v] = true
	}
}

func (r *Expire) Run() error {
	if len(r.Programs) == 0 {
		return fmt.Errorf("Programs is empty")
	}
	if r.LogDir == "" {
		r.LogDir = "/tmp" // TODO: do it the same way as lg does it
	}

	fs, err := filepath.Glob(r.LogDir + "/*") // TODO: nothing gained by using globs
	if err != nil {
		return fmt.Errorf("file error: %s", err)
	}
	if len(fs) == 0 {
		return fmt.Errorf("no files found")
	}
	var logFiles []logFile
	for _, f := range fs {

		lf, err := parseLogFileName(f)
		if err == nil {
			logFiles = append(logFiles, lf)
		}
	}
	r.logFiles = logFiles

	var errors []error
	for _, name := range r.Programs {
		err := r.clean(name)
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("%v", errors)
	}
	return nil

}

var ErrNotLgFile = errors.New("non a lg log file name")

func parseLogFileName(filename string) (logFile, error) {
	basename := filepath.Base(filename)
	fields := strings.Split(basename, ".")

	var ext string
	{
		if validExts[fields[len(fields)-1]] {
			ext = fields[len(fields)-1]
			fields = fields[:len(fields)-1]
		}
	}

	if len(fields) < 7 {
		return noLogFile, ErrNotLgFile
	}

	if fields[len(fields)-4] != "log" {
		return noLogFile, ErrNotLgFile
	}

	var pid uint64
	{
		var err error
		pidStr := fields[len(fields)-1]
		pid, err = strconv.ParseUint(pidStr, 10, 64)
		if err != nil {
			log.Printf("invalid pid %s in %s: %v", pidStr, filename, err)
		}
	}

	var timestamp time.Time
	{
		var err error
		timestamp, err = time.Parse("20060102-150405", fields[len(fields)-2])
		if err != nil {
			return noLogFile, fmt.Errorf("invalid date: %s", err)
		}
	}

	var level string
	{
		level = fields[len(fields)-3]
		if !validLevels[level] {
			return noLogFile, fmt.Errorf("%s is not a supprted log level", level)
		}
	}

	v := logFile{
		Filename: filename,
		Program:  strings.Join(fields[0:len(fields)-6], "."),
		Host:     fields[len(fields)-6],
		Username: fields[len(fields)-5],
		Level:    level,
		Time:     timestamp,
		Ext:      ext,
		Pid:      pid,
	}

	return v, nil
}

func (r Expire) clean(program string) error {

	programLogfiles := make(map[string][]logFile, 0)
	for level, _ := range validLevels {
		programLogfiles[level] = make([]logFile, 0)
	}

	for _, lf := range r.logFiles {
		if program == lf.Program {
			programLogfiles[lf.Level] = append(programLogfiles[lf.Level], lf)
		}
	}
	if len(programLogfiles) < 1 {
		log.Printf("found no log files for %s", program)
		return nil
	}
	for _, v := range programLogfiles {
		sort.Sort(sort.Reverse(logfilesByTime(v)))

	}

	filesToDelete := make(map[string]bool, 0)
	now := time.Now()
	for _, rule := range r.Rules {
		levels := allLevels
		if rule.Level != "" {
			levels = []string{rule.Level}
		}
	levels:
		for _, level := range levels {
			if len(programLogfiles[level]) < 1 {
				continue levels
			}
			var keptLogs []logFile
			keptLogs = append(keptLogs, programLogfiles[level]...)

			if rule.Count != 0 {
				if uint(len(keptLogs)) > rule.Count {
					keptLogs = keptLogs[:rule.Count]
				}
			}

			if rule.Age != 0 {
				var ageFiltered []logFile
				for _, v := range keptLogs {
					if now.Before(v.Time.Add(rule.Age)) {
						ageFiltered = append(ageFiltered, v)
					}
				}
				if len(ageFiltered) < 1 {
					ageFiltered = append(ageFiltered, keptLogs[0])
				}
				keptLogs = ageFiltered
			}
			keep := make(map[string]bool, 0)
			for _, v := range keptLogs {
				keep[v.Filename] = true

			}
			for _, v := range programLogfiles[level] {
				if !keep[v.Filename] {
					filesToDelete[v.Filename] = true
				}
			}
		}
	}
	for filename, _ := range filesToDelete {
		err := os.Remove(filename)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}
