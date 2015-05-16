package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type advice struct {
	funktion string
	code     string
}

func prep() {
	out, err := exec.Command("mkdir", "-p", "/tmp/coffee").CombinedOutput()
	if err != nil {
		log.Println("fuck" + err.Error())
	}
	fmt.Printf("%s\n", out)
}

func build() {
	fmt.Println("building")
	buildstr := "ENVZ=$PWD && cd /tmp/coffee && /usr/local/bin/go build && cp coffee $ENVZ/."
	out, err := exec.Command("bash", "-c", buildstr).CombinedOutput()
	if err != nil {
		log.Println("fuck" + err.Error())
	}
	fmt.Printf("%s\n", out)
	fmt.Println("building done")

}

func hasAdvice(a []advice, l string) advice {
	for i := 0; i < len(a); i++ {
		if strings.Contains(l, a[i].funktion) {
			return a[i]
		}
	}

	return advice{}
}

func transform(a []advice) {
	file, err := os.Open("main.go")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	out := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()

		avize := hasAdvice(a, l)
		if avize.funktion != "" {
			out += l + "\n" + avize.code + "\n"
		} else {
			out += l + "\n"
		}

	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	f, err := os.Create("/tmp/coffee/main.go")
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	b, err := f.WriteString(out)
	fmt.Println(b)
	if err != nil {
		log.Println(err)
	}
}

func grab_aspects() []advice {
	file, err := os.Open("main.goa")
	if err != nil {
		log.Println(err)
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
			cur_advice = a
		} else if strings.Contains(l, "}") {
			results = append(results, cur_advice)
		} else {
			cur_advice.code += l
		}

		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	return results
}

func main() {
	fmt.Println("goa v0.1")

	if len(os.Args) > 1 {
		prep()
		advice := grab_aspects()
		transform(advice)
		build()

	} else {
		fmt.Println("invalid")
	}
}
