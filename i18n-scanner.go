package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var langArg string
var funcKeyword string
var scanDir string
var msgFile string

// Total number of messages found
var msgCount int

type File map[string]Messages
type Messages map[string]string

// Save messages to path
func (f File) Save(path string) error {
	raw, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, raw, 0666)
	return err
}

// Load messages from path. All messages will be flushed before load
func (f File) Load(path string) error {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(raw, &f)
	return err
}

// Add a new message. If the message translation is exists, the file is inact
func (f File) Add(msg string) {
	for _, v := range f {
		_, ok := v[msg]
		if ok != true {
			v[msg] = ""
		}
	}

}

func (f File) AddLocale(locale string) {
	_, ok := f[locale]
	if ok == false {
		f[locale] = make(Messages)
	}
}

func main() {
	flag.StringVar(&langArg, "l", "en,zh", "Language names")
	flag.StringVar(&funcKeyword, "k", "_", "Translate function name, This arg will be a part of regular expression, if there are special chars you have to escape it manually")
	flag.StringVar(&scanDir, "d", "", "Directory to scan")
	flag.StringVar(&msgFile, "m", "messages.json", "Message file")

	flag.Usage = func() {
		fmt.Printf("i18n-scanner (https://github.com/greensea/i18n-scanner)\n\n")
		fmt.Printf("Usage: \n  i18n-scanner [OPTIONS] -d <DIR TO SCAN>\n\n")
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if scanDir == "" {
		flag.Usage()
		fmt.Printf("\nError: missing -d argument\n")
		os.Exit(0)
	}

	Scan()
}

func Scan() {
	f := make(File)
	f.Load(msgFile)

	langs := strings.Split(langArg, ",")
	for _, lang := range langs {
		f.AddLocale(lang)
	}

	ScanDir(f, scanDir)
	fmt.Printf("Found %d messages\n", msgCount)

	f.Save(msgFile)
	fmt.Printf("Messages json file saved to %s\n", msgFile)
}

func ScanDir(f File, path string) {
	entry, err := os.ReadDir(path)
	if err != nil {
		log.Printf("Unable to scan %s: %v\n", path, err)
		return
	}

	for _, v := range entry {
		subpath := fmt.Sprintf("%s/%s", path, v.Name())
		if v.IsDir() {
			ScanDir(f, subpath)
		} else {
			ScanFile(f, subpath)
		}
	}
}

func ScanFile(f File, path string) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Unable to read file %s: %v\n", path, err)
		return
	}

	fmt.Printf("Parsing %s\n", path)

	msgs := Parse(string(raw), funcKeyword)
	for _, v := range msgs {
		f.Add(v)
		msgCount++
	}
}

func Parse(raw string, funcName string) []string {
	var ret []string

	// Match _('bla')
	exp1 := fmt.Sprintf(`%s\([ ]*"(.+)"[ ]*\)`, funcName)
	r1, err := regexp.Compile(exp1)
	if err != nil {
		fmt.Printf("Unable to compile regular expression %s\nCheck if there are special chars in function name parameter", exp1)
		os.Exit(0)
	}

	// Match _("bla")
	exp2 := fmt.Sprintf(`%s\([ ]*'(.+)'[ ]*\)`, funcName)
	r2, err := regexp.Compile(exp2)
	if err != nil {
		fmt.Printf("Unable to compile regular expression %s\nCheck if there are special chars in function name parameter", exp1)
		os.Exit(0)
	}

	ret1 := r1.FindAllStringSubmatch(raw, -1)
	ret2 := r2.FindAllStringSubmatch(raw, -1)

	for _, v := range ret1 {
		if len(v) >= 2 {
			ret = append(ret, v[1])
		}
	}

	for _, v := range ret2 {
		if len(v) >= 2 {
			ret = append(ret, v[1])
		}
	}

	return ret
}
