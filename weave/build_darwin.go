package weave

func init() {
	oscmdenv = &osCmdInfo{
		shell:  []string{"/bin/sh", "-c"},
		mkdir:  "mkdir -p",
		rmdir:  "rm -rf",
		pwd:    "pwd",
		cp:     "cp",
		xcp:    "cp -r",
		envsep: ":",
	}
}
