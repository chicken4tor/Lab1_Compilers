package main

import (
	_ "bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func LexIt(input, output *os.File) {
	buffer, _ := io.ReadAll(input)
	code := string(buffer)
	for code != "" {
		code = strings.TrimSpace(code)
		for _, v := range Types {
			ok, str, err := v.IsThisStringYourType(code)
			if !ok {
				continue
			}
			if err != nil {
				output.WriteString(err.Error() + "\n")
				strings.Replace(code, str, "", 1)
				break
			}

			if !strings.HasPrefix(code, str) {
				continue
			}

			qwe, err := v.PostProcessingFunc(str)
			if err != nil {
				output.WriteString(fmt.Sprintf("ERROR\t %v\n", err))
			} else {
				_, err = output.WriteString(fmt.Sprintf("%v\t %v\n", v.Name, qwe))
				if err != nil {
					log.Printf("fatal: %v", err)
				}
			}

			code = strings.Replace(code, str, "", 1)
			break
		}
	}
}
