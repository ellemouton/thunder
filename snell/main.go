package main

import (
        "flag"
        "html/template"
        "log"
)

var templates = template.Must(template.ParseGlob("assets/*.html"))

func main() {
        flag.Parse()

        s, err := newState()
        if err != nil {
                log.Fatalf("newState: %s", err)
        }


        startRouter(s)
}