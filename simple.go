package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"simple/simpl"
)

var usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "%s <input file> \n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	in := flag.Arg(0)
	if in == "" {
		usage()
		return
	}

	infile, err := os.Open(in)
	if err != nil {
		log.Fatalf("could not open file '%v' with err: %v\n", in, err)
	}

	finfo, err := infile.Stat()
	if err != nil {
		log.Fatalf("could not stat file '%v' with err: %v\n", in, err)
	}

	if finfo.IsDir() {
		log.Fatalf("'%v' is a directory\n", in)
	}

	l := simpl.Lexer{In: infile}
	tokens, errors := l.Lex()
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println("ERROR:", err)
		}
		verb := "was"
		e := "error"
		if len(errors) > 1 {
			verb = "were"
			e += "s"
		}
		log.Fatalf("there %s %v %s lexing '%v'", verb, len(errors), e, in)
	}

	p := simpl.Parser{Tokens: tokens}
	p.Parse()

	i := simpl.NewInterpreter(&p.Lines, os.Stdout)
	i.Interpret()
}
