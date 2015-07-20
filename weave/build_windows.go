package weave

func init() {
	oscmdenv = &osCmdInfo{
		shell:  []string{"cmd", "/e:on", "/c"},
		mkdir:  "mkdir",
		rmdir:  "rmdir /s /q",
		pwd:    "cd",
		cp:     "copy",
		xcp:    "xcopy /e /q",
		exeext: ".exe",
		envsep: ";",
	}
}
