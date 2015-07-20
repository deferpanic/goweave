package weave

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type osCmdInfo struct {
	shell  []string // Shell命令
	mkdir  string
	rmdir  string
	pwd    string
	cp     string
	xcp    string
	exeext string // 可执行文件扩展名
	envsep string // GOPATH环境变量中用来分隔不同路径的字符
}

var (
	oscmdenv *osCmdInfo
)

// buildDir determines what the root build dir is
func (w *Weave) buildDir() string {
	return buildDir()
}

func execShellCmd(cmd string) (out []byte, err error) {
	arglist := append(oscmdenv.shell, cmd)
	out, err = exec.Command(arglist[0], arglist[1:]...).CombinedOutput()
	if err != nil {
		log.Printf("\ncmd: %s\nerr: %v\n", cmd, err)
	}
	return
}

// buildDir determines what the root build dir is
func buildDir() string {
	out, _ := execShellCmd(oscmdenv.pwd)
	return strings.TrimSpace(string(out))
}

// binName returns the expected bin name
func (w *Weave) binName() string {
	s := w.buildDir()
	stuff := strings.Split(s, string(filepath.Separator))
	return stuff[len(stuff)-1] + oscmdenv.exeext
}

// whichgo determines provides the full go path to the current go build
// tool
func (w *Weave) whichGo() string {
	return "go"
}

// tmpLocation returns the tmp build dir
func (w *Weave) tmpLocation() string {
	return tmpLocation()
}

func tmpLocation() string {
	gopath := os.Getenv("GOPATH")
	if idx := strings.Index(gopath, oscmdenv.envsep); idx > 0 {
		gopath = gopath[:idx]
	}
	return filepath.Join(gopath, "src", "_weave", setBase())
}

// build does the actual compilation
// right nowe we piggy back off of 6g/8g
func (w *Weave) build() {

	idx := strings.Index(w.buildLocation, "_weave")
	weavedir := w.buildLocation[:idx+6]
	defer os.RemoveAll(weavedir)

	delbin := fmt.Sprintf("%s %s", oscmdenv.rmdir, filepath.Join(w.buildLocation, w.binName()))
	execShellCmd(delbin)

	buildstr := "cd " + w.buildLocation
	buildstr += " && " + w.whichGo() + " build "
	buildstr += fmt.Sprintf(" && %s %s %s", oscmdenv.cp, w.binName(), w.buildDir())

	o, err := execShellCmd(buildstr)
	if err != nil {
		w.flog.Println(string(o))
	}
}

// prep prepares any tmp. build dirs
func (w *Weave) prep() {

	// hacky dir prep
	os.RemoveAll(w.buildLocation)
	fstcmd := fmt.Sprintf("%s %s", oscmdenv.mkdir, w.buildLocation)

	// hack to get anything that might be ref'd in the env
	hackcmd := fmt.Sprintf("%s * %s", oscmdenv.xcp, w.buildLocation)

	_, err := execShellCmd(fstcmd)
	if err != nil {
		w.flog.Println(fstcmd, err.Error())
	}

	filepath.Walk(
		w.buildLocation,
		func(fname string, fi os.FileInfo, err error) error {
			if fi == nil || err != nil {
				return err
			} else if !fi.IsDir() || fi.Name() == "." || fi.Name() == ".." {
				return nil
			} else {
				mkdir := fmt.Sprintf("%s %s", oscmdenv.mkdir, filepath.Join(w.buildLocation, fi.Name()))
				execShellCmd(mkdir)
				return nil
			}
		},
	)

	_, err = execShellCmd(hackcmd)
	if err != nil {
		w.flog.Println(hackcmd, err.Error())
	}
}

// rootPkg returns the root package of a go build
// this is needed to determine whether or not sub-pkg imports need to be
// re-written - which is basically any project w/more than one folder
func (w *Weave) rootPkg() string {
	return setBase()
}

// rootPkg returns the root package of a go build
// this is needed to determine whether or not sub-pkg imports need to be
// re-written - which is basically any project w/more than one folder
func setBase() string {
	out, err := exec.Command("go", "list").CombinedOutput()
	if err != nil {
		log.Println(err.Error())
	}

	str := strings.TrimSpace(string(out))
	str = filepath.FromSlash(str)
	return str
}
