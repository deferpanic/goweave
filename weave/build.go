package weave

import (
	"log"
	"os/exec"
	"strings"
)

// buildDir determines what the root build dir is
func (w *Weave) buildDir() string {
	out, err := exec.Command("bash", "-c", "pwd").CombinedOutput()
	if err != nil {
		w.flog.Println(err.Error())
	}

	return strings.TrimSpace(string(out))
}

func buildDir() string {
	out, err := exec.Command("bash", "-c", "pwd").CombinedOutput()
	if err != nil {
		log.Println(err.Error())
	}

	return strings.TrimSpace(string(out))
}

// binName returns the expected bin name
func (w *Weave) binName() string {
	s := w.buildDir()
	stuff := strings.Split(s, "/")
	return stuff[len(stuff)-1]
}

// whichgo determines provides the full go path to the current go build
// tool
func (w *Weave) whichGo() string {
	out, err := exec.Command("bash", "-c", "which go").CombinedOutput()
	if err != nil {
		w.flog.Println(err.Error())
	}

	return strings.TrimSpace(string(out))
}

// tmpLocation returns the tmp build dir
func (w *Weave) tmpLocation() string {
	out, err := exec.Command("bash", "-c", "echo $GOPATH").CombinedOutput()
	if err != nil {
		w.flog.Println(err.Error())
	}

	return string(out) + "src/_weave" + w.buildDir()
}

// tmpLocation returns the tmp build dir
func tmpLocation() string {
	out, err := exec.Command("bash", "-c", "echo $GOPATH").CombinedOutput()
	if err != nil {
		log.Println(err.Error())
	}

	return strings.TrimSpace(string(out)) + "/src/_weave/" + setBase()
}

// build does the actual compilation
// right nowe we piggy back off of 6g/8g
func (w *Weave) build() {
	buildstr := "cd " + w.buildLocation + " && " + w.whichGo() + " build && cp " +
		w.binName() + " " + w.buildDir() + "/."

	o, err := exec.Command("bash", "-c", buildstr).CombinedOutput()
	if err != nil {
		w.flog.Println(string(o))
	}

}

// prep prepares any tmp. build dirs
func (w *Weave) prep() {

	fstcmd := "mkdir -p " + w.buildLocation

	sndcmd := `find . -type d -exec mkdir -p "` + w.buildLocation + `/{}" \;`

	_, err := exec.Command("bash", "-c", fstcmd).CombinedOutput()
	if err != nil {
		w.flog.Println(err.Error())
	}

	_, err = exec.Command("bash", "-c", sndcmd).CombinedOutput()
	if err != nil {
		w.flog.Println(err.Error())
	}

}

// rootPkg returns the root package of a go build
// this is needed to determine whether or not sub-pkg imports need to be
// re-written - which is basically any project w/more than one folder
func (w *Weave) rootPkg() string {
	out, err := exec.Command("bash", "-c", "go list").CombinedOutput()
	if err != nil {
		w.flog.Println(err.Error())
	}

	return strings.TrimSpace(string(out))
}

// rootPkg returns the root package of a go build
// this is needed to determine whether or not sub-pkg imports need to be
// re-written - which is basically any project w/more than one folder
func setBase() string {
	out, err := exec.Command("bash", "-c", "go list").CombinedOutput()
	if err != nil {
		log.Println(err.Error())
	}

	return strings.TrimSpace(string(out))
}
