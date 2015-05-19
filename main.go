package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	version = "v0.1"
)

var (
	// flog is a global logger
	flog *log.Logger
)

// adviceType returns a map of id to human expression of advice types
func adviceType() map[int]string {
	return map[int]string{
		1: "before",
		2: "after",
		3: "around",
	}
}

// advice has a function to wrap advice around and code for said
// function
type advice struct {
	funktion     string
	code         string
	adviceTypeId int
}

// buildDir determines what the root build dir is
func buildDir() string {
	out, err := exec.Command("bash", "-c", "pwd").CombinedOutput()
	if err != nil {
		flog.Println(err.Error())
	}

	return strings.TrimSpace(string(out))
}

// binName returns the expected bin name
func binName() string {
	s := buildDir()
	stuff := strings.Split(s, "/")
	return stuff[len(stuff)-1]
}

// prep prepares any tmp. build dirs
func prep() {
	fmt.Println("building" + tmpLocation())
	out, err := exec.Command("mkdir", "-p", tmpLocation()).CombinedOutput()
	if err != nil {
		flog.Println(err.Error())
	}
	fmt.Printf("%s\n", out)
}

// whichgo determines provides the full go path to the current go build
// tool
func whichGo() string {
	return "/usr/local/bin/go"
}

// tmpLocation returns the tmp build dir
func tmpLocation() string {
	return "/tmp" + buildDir()
}

// build does the actual compilation
// right nowe we piggy back off of 6g/8g
func build() {
	buildstr := "cd " + tmpLocation() + " && " + whichGo() + " build && cp " +
		binName() + " " + buildDir() + "/."

	fmt.Println(buildstr)
	_, err := exec.Command("bash", "-c", buildstr).CombinedOutput()
	if err != nil {
		flog.Println(err.Error())
	}

}

// hasAdvice returns advice if there is advice for a given line of
// source or returns empty advice
func hasAdvice(a []advice, l string) advice {
	for i := 0; i < len(a); i++ {
		if strings.Contains(l, a[i].funktion) {
			return a[i]
		}
	}

	return advice{}
}

// transform reads line by line over each src file and inserts advice
// where appropriate
func transform(a []advice) {
	curfile := "main.go"

	file, err := os.Open(curfile)
	if err != nil {
		flog.Println(err)
	}
	defer file.Close()

	out := ""

	// poor man's scope
	scope := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()

		avize := hasAdvice(a, l)
		if avize.funktion != "" {
			scope += 1

			fmt.Println("found advice:\t" + avize.funktion)

			// insert before
			if avize.adviceTypeId == 1 {
				out += l + "\n" + avize.code + "\n"
			} else {
				out += l + "\n"
			}

		} else {

			// dat scope
			if strings.Contains(l, "}") {
				scope -= 1
			}

			out += l + "\n"
		}

	}

	if err := scanner.Err(); err != nil {
		flog.Println(err)
	}

	fmt.Println("writing out" + tmpLocation() + "/" + curfile)

	f, err := os.Create(tmpLocation() + "/" + curfile)
	if err != nil {
		flog.Println(err)
	}

	defer f.Close()

	b, err := f.WriteString(out)
	fmt.Println(b)
	if err != nil {
		flog.Println(err)
	}
}

// grab_aspects looks for an aspect file for each file
// this seems lame and contrary to what we want...
// I'd go as far as to say that we want this to be cross-pkg w/in root?
//
// maybe the rule should be - aspects are valid for anything in a
// project root?
func grab_aspects() []advice {
	aspectsFile := "main.goa"

	file, err := os.Open(aspectsFile)
	if err != nil {
		flog.Println(err)
	}
	defer file.Close()

	results := []advice{}

	cur_advice := advice{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()

		if strings.Contains(l, "advice") {
			blah := strings.Split(l, "execution(\"")[1]
			shiz := strings.Split(blah, "\"")[0]
			fmt.Println("Method: " + shiz)

			a := advice{}
			a.funktion = shiz
			flog.Println("function:" + a.funktion)

			a.adviceTypeId = set_advice_type(l)

			flog.Println(a.adviceTypeId)

			cur_advice = a
		} else if strings.Contains(l, "}") {
			results = append(results, cur_advice)
		} else {
			cur_advice.code += l + "\n"
		}

		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		flog.Println(err)
	}

	fmt.Println(len(results))
	return results
}

func set_advice_type(l string) int {
	stuff := strings.Split(l, ": ")
	ostuff := strings.Split(stuff[1], " ")

	switch ostuff[0] {
	case "before":
		return 1
	case "after":
		return 2
	case "around":
		return 3
	}

	return -1
}

// main is the main point of entry for running goa
func main() {
	flog = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	fmt.Println("goa " + version)

	prep()
	advice := grab_aspects()
	transform(advice)
	build()

}
