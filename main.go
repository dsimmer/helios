package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func inter(booler bool) int {
	if booler {
		return 1
	}
	return 0
}

// Init saves the required CD workaround to the bash so it can be called. Will need sudo permissions to write to it
func Init() {
	path := string(os.PathSeparator) + "etc" + string(os.PathSeparator) + "bash.bashrc"
	data, err := ioutil.ReadFile(path)
	check(err)

	stringData := string(data)
	prefix := `# Generated function that allows helios to use cd`
	// todo fix this bash func
	content := `
function hc() {
	hset=0
	anyset=0
	for opt in "$@"
	do
		echo $opt
		case $opt in
			-h) hset=1 ;;
			-*) anyset=1 ;;
		esac
	done
	if [[ (! ("$#" -eq 1 && $anyset -eq 0 && $hset -eq 0)) && $hset -eq 0 ]]; then
		helios "$@"
	else
		cdResult=$(helios "$@")
		cd "$cdResult"
	fi
}

export -f hc
`
	suffix := `# End of generated helios function`
	var newFile string
	if strings.Contains(stringData, prefix) {
		from := strings.Index(stringData, prefix)
		to := strings.LastIndex(stringData, suffix)
		newFile = stringData[:from] + prefix + content + stringData[to:]
	} else {
		newFile = stringData + "\n" + prefix + content + suffix + "\n"
	}
	ioutil.WriteFile(path, []byte(newFile), 0666)

	ex, err := os.Executable()
	check(err)
	exPath := filepath.Dir(ex)
	if _, err := os.Stat(exPath + string(os.PathSeparator) + "helios_history"); os.IsNotExist(err) {
		f, err := os.Create(exPath + string(os.PathSeparator) + "helios_history")
		check(err)
		check(f.Close())
		os.Chmod(exPath+string(os.PathSeparator)+"helios_history", 0666)
	}
	if _, err := os.Stat(exPath + string(os.PathSeparator) + "helios_notes.yml"); os.IsNotExist(err) {
		f, err := os.Create(exPath + string(os.PathSeparator) + "helios_notes.yml")
		check(err)
		_, err = f.WriteString(`linux:
  - chmod -R 777 dir

js:
  - === != ==

go:
  - check(err) good practice

clojurescript:
  - (.. object -property -nestedproperty)

clojure:
  - macros dont parse inputs

elixir:
  - No error checking`)
		check(err)
		check(f.Close())
		os.Chmod(exPath+string(os.PathSeparator)+"helios_notes.yml", 0666)
	}
	if _, err := os.Stat(exPath + string(os.PathSeparator) + "helios_achievements.yml"); os.IsNotExist(err) {
		f, err := os.Create(exPath + string(os.PathSeparator) + "helios_achievements.yml")
		check(err)
		_, err = f.WriteString(`Achievements to date:
`)
		check(err)
		check(f.Close())
		os.Chmod(exPath+string(os.PathSeparator)+"helios_achievements.yml", 0666)
	}
	if _, err := os.Stat(exPath + string(os.PathSeparator) + "helios_settings"); os.IsNotExist(err) {
		f, err := os.Create(exPath + string(os.PathSeparator) + "helios_settings")
		check(err)
		check(f.Close())
		os.Chmod(exPath+string(os.PathSeparator)+"helios_settings", 0666)
	}
	if _, err := os.Stat(exPath + string(os.PathSeparator) + "helios_favourites"); os.IsNotExist(err) {
		f, err := os.Create(exPath + string(os.PathSeparator) + "helios_favourites")
		check(err)
		check(f.Close())
		os.Chmod(exPath+string(os.PathSeparator)+"helios_favourites", 0666)
		var empty map[string]string
		saveFavourites(empty)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func SaveNote(category string, line string) {
	ex, err := os.Executable()
	check(err)
	exPath := filepath.Dir(ex)
	data, err := ioutil.ReadFile(exPath + string(os.PathSeparator) + "helios_notes.yml")
	check(err)
	dataString := string(data)
	if strings.Index(dataString, category+":") > -1 {
		newString := strings.SplitAfter(dataString, category+":")
		dataString = strings.Join([]string{newString[0], "\n  - " + line, newString[1]}, "")
	} else {
		dataString = dataString + "\n\n" + category + ":" + "\n" + "  - " + line
	}

	err = ioutil.WriteFile(exPath+string(os.PathSeparator)+"helios_notes.yml", []byte(dataString), 0777)
	check(err)
}

// todo write an import method so you can import and sort achievements
// todo stats on your achievements
func SaveAchievement(line string) {
	ex, err := os.Executable()
	check(err)
	exPath := filepath.Dir(ex)
	f, err := os.OpenFile(exPath+string(os.PathSeparator)+"helios_achievements.yml", os.O_APPEND|os.O_WRONLY, 0600)
	check(err)
	f.WriteString("\n" + time.Now().String() + ": " + line)
}

//todo fuzzy search

func GrepNote(category string, line string) {
	regexer := regexp.MustCompile(`  -.*?` + line + `.*?\n`)
	ex, err := os.Executable()
	check(err)
	exPath := filepath.Dir(ex)
	data, err := ioutil.ReadFile(exPath + string(os.PathSeparator) + "helios_notes.yml")
	check(err)
	dataString := string(data)
	if category != "" {
		newString := strings.SplitAfter(dataString, category+":")
		newString2 := strings.Split(newString[1], "\n\n")
		matches := regexer.FindAll([]byte(newString2[0]), -1)
		for _, res := range matches {
			fmt.Println(strings.TrimLeft(string(res), "  -"))
		}
	} else {
		matches := regexer.FindAll([]byte(dataString), -1)
		for _, res := range matches {
			fmt.Println(strings.TrimLeft(string(res), "  -"))
		}
	}
}

func addToHistory(line string) {
	ex, err := os.Executable()
	check(err)
	exPath := filepath.Dir(ex)
	f, err := os.OpenFile(exPath+string(os.PathSeparator)+"helios_history", os.O_APPEND|os.O_WRONLY, 0600)
	check(err)

	defer f.Close()

	if _, err = f.WriteString(line + "\n"); err != nil {
		panic(err)
	}
}

//todo this should probably use a regex instead?
func SearchSettings() {
	ex, err := os.Executable()
	check(err)
	exPath := filepath.Dir(ex)
	data, err := ioutil.ReadFile(exPath + string(os.PathSeparator) + "helios_settings")
	check(err)
	splitData := strings.Split(string(data), "#!")
	result := ""
	for i, d := range splitData {
		if i%2 == 1 {
			result = result + "\n" + d
		}
	}
	fmt.Println(result)
}

// todo parse this to import and recreate programs
func addToSettings(line string) {
	ex, err := os.Executable()
	check(err)
	exPath := filepath.Dir(ex)
	f, err := os.OpenFile(exPath+string(os.PathSeparator)+"helios_settings", os.O_APPEND|os.O_WRONLY, 0600)
	check(err)

	defer f.Close()

	if _, err = f.WriteString(line); err != nil {
		panic(err)
	}
}

// SaveScript saves the provided script to the filename provided in an executable directory (default bin)
func SaveScript(args []string) {
	if len(args) > 2 {
		err := errors.New("Incorrect number of arguments (>2)")
		panic(err)
	}

	addToSettings("#!" + args[0] + "#!/bin/sh\n" + args[1])
	ex, err := os.Executable()
	check(err)
	exPath := filepath.Dir(ex)
	content := []byte("#!/bin/sh\n" + args[1])
	err = ioutil.WriteFile(exPath+string(os.PathSeparator)+args[0], content, 0777)
	check(err)
}

// ExportAll exports all the saved scripts, history and favourite to a file
func ExportAll(args []string) {
	if len(args) > 1 {
		err := errors.New("Incorrect number of arguments (>1)")
		panic(err)
	}
	ex, err := os.Executable()
	check(err)
	exPath := filepath.Dir(ex)
	data, err := ioutil.ReadFile(exPath + string(os.PathSeparator) + "helios_settings")
	check(err)

	err = ioutil.WriteFile(args[0], data, 0666)
	check(err)
}

// ImportAll imports all the saved scripts, history and favourite from a file generated by ExportAll
func ImportAll(args []string) {
	if len(args) > 1 {
		err := errors.New("Incorrect number of arguments (>1)")
		panic(err)
	}
	data, err := ioutil.ReadFile(args[0])
	check(err)
	ex, err := os.Executable()
	check(err)
	exPath := filepath.Dir(ex)
	err = os.Remove(exPath + string(os.PathSeparator) + "helios_settings")
	check(err)
	addToSettings(string(data))
}

func saveFavourites(favourites map[string]string) {
	ex, err := os.Executable()
	check(err)
	exPath := filepath.Dir(ex)

	b := new(bytes.Buffer)

	e := gob.NewEncoder(b)
	err = e.Encode(favourites)
	check(err)

	err = ioutil.WriteFile(exPath+string(os.PathSeparator)+"helios_favourites", b.Bytes(), 0666)
	check(err)
}

func loadFavourites() map[string]string {
	ex, err := os.Executable()
	check(err)
	exPath := filepath.Dir(ex)

	data, err := os.Open(exPath + string(os.PathSeparator) + "helios_favourites")
	defer data.Close()

	var decodedMap map[string]string
	d := gob.NewDecoder(data)

	err = d.Decode(&decodedMap)
	check(err)
	return decodedMap
}

//todo test this works
// CD improves the editors regular functionality with a favourite and regex serch option. Regex also searches favourites
func CD(fav bool, regex bool, args []string) {
	if len(args) > 2 || (!fav && (len(args) > 1)) {
		err := errors.New("Incorrect number of arguments (>2 or !favourite and >1)")
		panic(err)
	}
	var favName string
	search := args[0]
	if fav {
		favName = args[0]
		search = args[1]
	}
	result := search
	if regex {
		//todo search through all files for anything with that path, preference for later matches
	}
	favourites := loadFavourites()
	if _, ok := favourites[search]; ok {
		result = favourites[search]
	}

	addToHistory(result)

	// Print result for input into cd
	fmt.Println(result)

	if fav {
		favourites[favName] = result
		saveFavourites(favourites)
	}
}

// History shows your previous helios commands and allows you to jump to previous directories or commands
func History(args []string) {
	if len(args) > 1 {
		err := errors.New("Incorrect number of arguments (>1)")
		panic(err)
	} else if len(args) == 1 {
		ex, err := os.Executable()
		check(err)
		exPath := filepath.Dir(ex)
		data, err := ioutil.ReadFile(exPath + string(os.PathSeparator) + "helios_history")
		check(err)
		// todo Goto line in history and return it for the cd bash command
		fmt.Println(data)
	} else {
		ex, err := os.Executable()
		check(err)
		exPath := filepath.Dir(ex)
		data, err := ioutil.ReadFile(exPath + string(os.PathSeparator) + "helios_history")
		check(err)
		fmt.Println(data)
	}
}

func main() {
	snPtr := flag.Bool("sn", false, "Save note (arg 1: category, arg2: line)")
	saPtr := flag.Bool("sa", false, "Save achievement, automatically timestamped (arg 1: line)")
	gnPtr := flag.Bool("gn", false, "Get note (via grep) (arg 1: category, arg2: line) ||  (arg1: line)")
	sPtr := flag.Bool("s", false, "Define and save a script (arg 1: name, arg2: program, one string, auto prefixed with #!/bin/sh\\n and saved as a file). Don't forget to add help via bash 'if [ \"$1\" == \"-h\" ]; then\necho \"Usage: `basename $0` [somestuff]\"\nexit 0\nfi")
	slPtr := flag.Bool("sl", false, "List all available scripts (noargs)")
	ePtr := flag.Bool("e", false, "Export your settings and scripts (arg 1: path to save)")

	iPtr := flag.Bool("i", false, "Import your settings and scripts (arg 1: path to import)")

	fPtr := flag.Bool("f", false, "TODO: Favourite a directory, also works with the r option. Automatically navigates there. (arg 1: favourite name, arg 2: search text) || (arg 1: search text)")
	rPtr := flag.Bool("r", false, "TODO: Regex search")

	hPtr := flag.Bool("h", false, "History of navigation in helios. Additional argument navigates to that item (noargs) || (arg 1: history line item to navigate to)")

	initPtr := flag.Bool("init", false, "Init, required for CD functions to work properly. Creates storage files")

	flag.Parse()
	frPtr := *fPtr || *rPtr
	counter := inter(*sPtr) + inter(*ePtr) + inter(*iPtr) + inter(frPtr) + inter(*hPtr) + inter(*initPtr) + inter(*snPtr) + inter(*gnPtr)
	if counter >= 2 {
		err := errors.New("Incorrect combination of arguments (>=2)")
		panic(err)
	}

	if *initPtr {
		Init()
	}
	if *gnPtr {
		if flag.NArg() > 1 {
			GrepNote(flag.Arg(0), flag.Arg(1))
		} else {
			GrepNote("", flag.Arg(0))
		}
	}
	if *snPtr {
		SaveNote(flag.Arg(0), flag.Arg(1))
	}
	if *saPtr {
		SaveAchievement(flag.Arg(0))
	}
	if *sPtr {
		SaveScript(flag.Args())
	}
	if *slPtr {
		if flag.NArg() == 0 {
			SearchSettings()
		} else {
			err := errors.New("Too many arguments for -sl")
			panic(err)
		}
	}
	if *ePtr {
		ExportAll(flag.Args())
	}
	if *iPtr {
		ImportAll(flag.Args())
	}
	if frPtr || counter == 0 {
		CD(*fPtr, *rPtr, flag.Args())
	}
	if *hPtr {
		History(flag.Args())
	}
}
