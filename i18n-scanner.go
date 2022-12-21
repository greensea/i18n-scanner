package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var langArg string
var funcKeyword string
var scanDir string
var msgFile string

// Total number of messages found
var msgCount int

type File struct {
	Data      map[string]Messages
	msgOrder  int
	msgOrders map[string]int
}

type Messages map[string]string

func NewFile() *File {
	f := File{}
	f.Data = make(map[string]Messages)
	f.msgOrders = make(map[string]int)
	return &f
}

// Save messages to path
func (f *File) Save(path string) error {
	raw, err := f.MarshalJSON()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, raw, 0666)
	return err
}

func (f *File) MarshalJSON() ([]byte, error) {
	// 1. Load keys
	var locales []string
	var msgs []string

	for k := range f.Data {
		locales = append(locales, k)

		for k2 := range f.Data[k] {
			msgs = append(msgs, k2)
		}
	}

	// 2. Sort keys and make them unique
	msgs = Unique(msgs)
	sort.Slice(locales, func(i, j int) bool {
		return strings.Compare(locales[i], locales[j]) < 0
	})

	buf := &bytes.Buffer{}
	buf.Write([]byte("{\n"))

	for k1, localeName := range locales {
		lToken := EscapeString(localeName)
		buf.Write([]byte(fmt.Sprintf("  %s: {\n", lToken)))

		// Sort again, put un-translate messages at top
		sort.Slice(msgs, func(i, j int) bool {
			// Sort by un-translate
			msgi, _ := f.Data[localeName][msgs[i]]
			msgj, _ := f.Data[localeName][msgs[j]]
			ilen := len(msgi)
			jlen := len(msgj)

			if ilen == 0 && jlen > 0 {
				return true
			} else if ilen > 0 && jlen == 0 {
				return false
			}

			// Sort by scan order
			ival, ok1 := f.msgOrders[msgs[i]]
			jval, ok2 := f.msgOrders[msgs[j]]
			if ok1 == false || ok2 == false {
				return strings.Compare(msgs[i], msgs[j]) < 0
			}

			// Sort by original messages string
			if ival == jval {
				return strings.Compare(msgs[i], msgs[j]) < 0
			}

			return ival < jval

		})

		for k2, msgName := range msgs {
			msg, _ := f.Data[localeName][msgName]
			nameToken := EscapeString(msgName)
			msgToken := EscapeString(msg)
			buf.Write([]byte(fmt.Sprintf("    %s: %s", nameToken, msgToken)))

			if k2 != len(msgs)-1 {
				buf.Write([]byte(",\n"))
			} else {
				buf.Write([]byte("\n"))
			}
		}

		if k1 != len(locales)-1 {
			buf.Write([]byte("  },\n"))
		} else {
			buf.Write([]byte("  }\n"))
		}

	}
	buf.Write([]byte("}"))

	return buf.Bytes(), nil
}

// Load messages from path. All messages will be flushed before load
func (f *File) Load(path string) error {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(raw, &f.Data)
	return err
}

// Add a new message. If the message translation is exists, the file is inact
func (f *File) Add(msg string) {
	for _, v := range f.Data {
		_, ok := v[msg]
		if ok != true {
			v[msg] = ""
		}
	}

	_, ok := f.msgOrders[msg]
	if ok == false {
		f.msgOrder++
		f.msgOrders[msg] = f.msgOrder
	}

}

func (f *File) AddLocale(locale string) {
	_, ok := f.Data[locale]
	if ok == false {
		f.Data[locale] = make(Messages)
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
	f := NewFile()
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

func ScanDir(f *File, path string) {
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

func ScanFile(f *File, path string) {
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

	// Match _("bla")
	exp1 := fmt.Sprintf(`%s\([ ]*"(.+)"[ ]*[\,\)]`, funcName)
	r1, err := regexp.Compile(exp1)
	if err != nil {
		fmt.Printf("Unable to compile regular expression %s\nCheck if there are special chars in function name parameter", exp1)
		os.Exit(0)
	}

	// Match _('bla')
	exp2 := fmt.Sprintf(`%s\([ ]*'(.+)'[ ]*[\,\)]`, funcName)
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

func Unique(ss []string) []string {
	u := make(map[string]struct{})
	var ret []string
	for _, v := range ss {
		_, ok := u[v]
		if ok != true {
			ret = append(ret, v)
			u[v] = struct{}{}
		}
	}

	return ret
}

// Escape a string to JSON format
func EscapeString(s string) string {
	ret := strconv.Quote(s)
	return ret

}
