package itype

import "flag"

func PaseFlag(F *Flag) {

	flag.StringVar(&F.FileName, "f", "", "file name")
	flag.BoolVar(&F.Debug, "d", false, "debug mode")
	flag.Parse()

}
