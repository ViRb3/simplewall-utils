package allow

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"github.com/TomOnTime/utfutil"
	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/unicode"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var processName string
var appendMode bool
var profilePath string
var logPath string

var Cmd = &cobra.Command{
	Use:   "allow",
	Short: "Adds allow rules for a process name by allowing its destination IPs from the packet log file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(); err != nil {
			log.Fatalln(err)
		}
	},
}

func run() error {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("%s doesn't exist, is logging enabled?", logPath))
	} else if err != nil {
		return err
	}
	log.Println("Processing...")

	// By default the file is in UTF16-LE BOM
	f, err := utfutil.OpenFile(logPath, utfutil.UTF8)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	var logEntries []*LogEntry
	if err := gocsv.Unmarshal(f, &logEntries); err != nil {
		return err
	}

	destIPs := map[string]bool{}
	for _, entry := range logEntries {
		if filepath.Base(entry.Path) == processName {
			remote := []string{entry.RemoteAddress, entry.RemotePort}
			for i := range remote {
				// Strip appended hostname if enabled in settings
				spaceIndex := strings.Index(remote[i], " ")
				if spaceIndex > 0 {
					remote[i] = remote[i][:spaceIndex]
				}
			}
			destIPs[remote[0]+":"+remote[1]] = true
		}
	}

	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("%s doesn't exist", profilePath))
	} else if err != nil {
		return err
	}

	// By default the file is in UTF16-LE BOM
	profileBytes, err := utfutil.ReadFile(profilePath, utfutil.UTF8)
	if err != nil {
		return err
	}

	var profile Profile
	if err := xml.Unmarshal(profileBytes, &profile); err != nil {
		return err
	}

	var removeIndexes []int
	for i, rule := range profile.RulesCustom.Item {
		if rule.Name != processName {
			continue
		}
		if appendMode {
			ips := strings.Split(rule.Rule, ";")
			for _, ip := range ips {
				destIPs[ip] = true
			}
		}
		removeIndexes = append(removeIndexes, i)
	}

	// Remove all old rules
	for i := len(removeIndexes) - 1; i >= 0; i-- {
		removeIndex := removeIndexes[i]
		profile.RulesCustom.Item = append(profile.RulesCustom.Item[:removeIndex], profile.RulesCustom.Item[removeIndex+1:]...)
	}

	// simplewall has a maximum limit of 255 characters in a rule. Anything after that is truncated.
	// https://github.com/henrypp/simplewall/issues/809
	var rules []string
	var rule string

	for ip, _ := range destIPs {
		newData := ";" + ip
		if len(rule)+len(newData) > 255 {
			rules = append(rules, rule[1:])
			rule = ""
		}
		rule += newData
	}
	rules = append(rules, rule[1:])

	for _, rule := range rules {
		profile.RulesCustom.Item = append(profile.RulesCustom.Item, RuleItem{
			Name:      processName,
			Rule:      rule,
			Protocol:  "6", // tcp
			IsEnabled: "1",
		})
	}

	newProfileBytes, err := xml.MarshalIndent(profile, "", "    ")
	if err != nil {
		return err
	}

	// Save back to default UTF16-LE BOM
	encoder := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder()
	newProfileBytes, err = encoder.Bytes(newProfileBytes)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(profilePath, newProfileBytes, 0644); err != nil {
		return err
	}
	log.Println("Done!")
	return nil
}

func init() {
	Cmd.PersistentFlags().StringVarP(&processName, "process-name", "n", "", "Process name to allow")
	if err := Cmd.MarkPersistentFlagRequired("process-name"); err != nil {
		log.Fatalln(err)
	}
	Cmd.PersistentFlags().BoolVarP(&appendMode, "append", "a", false, "Append to existing rules instead of overwriting")
	profilePathDefault := filepath.Join(os.Getenv("APPDATA"), "Henry++", "simplewall", "profile.xml")
	Cmd.PersistentFlags().StringVarP(&profilePath, "profile-path", "p", profilePathDefault, "Path to profile file")
	logPathDefault := filepath.Join(os.Getenv("USERPROFILE"), "simplewall.log")
	Cmd.PersistentFlags().StringVarP(&logPath, "log-path", "l", logPathDefault, "Path to log file")

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.LazyQuotes = true
		return r
	})
}
