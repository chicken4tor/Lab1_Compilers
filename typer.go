package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

var Types []Type

type Type struct {
	Name               string
	Regexps            []*regexp.Regexp
	AdditionalInfo     string
	PreProcessingFunc  func(str string) (string, error)
	PostProcessingFunc func(str string) (string, error)
}

func defaultPostprocessingFunc(str string) (string, error) {
	return str, nil
}

func defaultPreprocessingFunc(str string) (string, error) { return str, nil }

func NewType() Type {
	return Type{
		PostProcessingFunc: defaultPostprocessingFunc,
		PreProcessingFunc:  defaultPreprocessingFunc,
	}
}

func (t Type) Print() {
	for _, v := range t.Regexps {
		fmt.Printf("%s\t%v\n", t.Name, v.String())
	}
	fmt.Print("\n")
}

func (t Type) IsThisStringYourType(str string) (bool, string, error) {
	var err error
	if t.PreProcessingFunc != nil {
		_, err = t.PreProcessingFunc(str)
		if err != nil {
			return false, str, err
		}
		str, _ = t.PreProcessingFunc(str)
	}
	for _, v := range t.Regexps {
		sttr := v.FindString(str)
		if sttr != "" {
			return true, sttr, nil
		}
	}
	return false, "", errors.New("kavo")
}

func CreateRegexps(str ...string) []*regexp.Regexp {
	rgxs := make([]*regexp.Regexp, 0, len(Keywords()))
	for _, v := range str {
		rgx, err := regexp.Compile(v)
		if err != nil {
			log.Fatal(err)
		}
		rgxs = append(rgxs, rgx)
	}
	return rgxs
}

func Keywords() []string {
	return []string{
		`class`,
		`else`,
		`false`,
		`fi`,
		`if`,
		`in`,
		`inherits`,
		`isvoid`,
		`let`,
		`loop`,
		`pool`,
		`then`,
		`while`,
		`case`,
		`esac`,
		`new`,
		`of`,
		`not`,
		`true`,
	}
}

func init() {
	Types = make([]Type, 0)
	keywords := NewType()
	keywords.Name = "KEYWORD"
	keywords.Regexps = CreateRegexps(Keywords()...)

	integer := NewType()
	integer.Name = "INTEGER"
	integer.Regexps = CreateRegexps(`\d`)

	classIdentifier := NewType()
	classIdentifier.Name = "CLASS_IDENTIFIER"
	classIdentifier.Regexps = CreateRegexps(`[A-Z]\w*`)

	objectIdentifier := NewType()
	objectIdentifier.Name = "OBJECT_IDENTIFIER"
	objectIdentifier.Regexps = CreateRegexps(`[a-z]\w*`)

	self := NewType()
	self.Name = "SELF"
	self.Regexps = CreateRegexps(`self`)

	selfType := NewType()
	selfType.Name = "SELF_TYPE"
	selfType.Regexps = CreateRegexps(`SELF_TYPE`)

	String := NewType()
	String.Name = "STRING"
	String.Regexps = CreateRegexps(`"(.|\(\/\n\))*?"`)
	String.PreProcessingFunc = func(str string) (string, error) {
		strings.ReplaceAll(str, "\\b", "\b")
		strings.ReplaceAll(str, "\\t", "\t")
		strings.ReplaceAll(str, "\\n", "\n")
		strings.ReplaceAll(str, "\\f", "\f")
		strings.ReplaceAll(str, "\\", "")
		return str, nil

	}
	String.PostProcessingFunc = func(str string) (string, error) {
		if len(str) > 128 {
			return str, errors.New("string is too long")
		}
		return str, nil
	}

	unterminatedStringError := NewType()
	unterminatedStringError.Name = "ERROR"
	unterminatedStringError.Regexps = CreateRegexps(`".*\n`)
	unterminatedStringError.PostProcessingFunc = func(str string) (string, error) {
		return "Unterminated string", nil
	}

	oneLineComment := NewType()
	oneLineComment.Name = "ONE_LINE_COMMENT"
	oneLineComment.Regexps = CreateRegexps(`--[^\n]*`)
	oneLineComment.PreProcessingFunc = nil

	multilineComment := NewType()
	multilineComment.Name = "MULTI_LINE_COMMENT"
	multilineComment.Regexps = CreateRegexps(`(\\()([\*])+(.|\n)*([\*])(\\))`)

	errorEOFinComment := NewType()
	errorEOFinComment.Name = "ERROR"
	errorEOFinComment.Regexps = CreateRegexps(`\(\\()([\*])+(.|\n)*`)
	errorEOFinComment.PostProcessingFunc = func(str string) (string, error) {
		return "EOF in comment", nil
	}

	whiteSpace := NewType()
	whiteSpace.Name = "WHITE_SPACE"
	whiteSpace.Regexps = CreateRegexps(`[\n\f\r\t\ ]*`)

	operator := NewType()
	operator.Name = "OPERATOR"
	operator.Regexps = CreateRegexps("\\.", "@", "~", "isvoid", "\\*", "/", "\\+", "-", "<=", "<", "=", "not", "<-")

	punctuation := NewType()
	punctuation.Name = "PUNCTUATION"
	punctuation.Regexps = CreateRegexps(",", ":", ";", "\\(", "\\)", "{", "}")

	errorDot := NewType()
	errorDot.Name = "ERROR"
	errorDot.Regexps = CreateRegexps(".")

	Types = append(
		Types,
		keywords,
		integer,
		classIdentifier,
		objectIdentifier,
		self,
		selfType,
		String,
		unterminatedStringError,
		oneLineComment,
		multilineComment,
		errorEOFinComment,
		whiteSpace,
		operator,
		punctuation,
		errorDot,
	)
	for _, v := range Types {
		if v.PostProcessingFunc == nil {
			v.PostProcessingFunc = defaultPostprocessingFunc
		}
		if v.PreProcessingFunc == nil {
			v.PreProcessingFunc = func(str string) (string, error) { return str, nil }
		}
	}
	for _, v := range Types {
		v.Print()
	}
}
