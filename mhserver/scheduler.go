package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorhill/cronexpr"
)

type schedule struct {
	expression *cronexpr.Expression
	cmd        string
	timer      *time.Timer
}

var filename string

var todos []schedule
var cases []reflect.SelectCase

func scheduler() {
	filename = "crontab"
	loadSchedule(filename)

}

/*

 func looper() {
	 err = watcher.Add("/tmp/foo")
	 if err != nil {
	     log.Fatal(err)
	 }


    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                continue
            }
            log.Println("event:", event)
            if event.Op&fsnotify.Write == fsnotify.Write {
                log.Println("modified file:", event.Name)
            }
        case err, ok := <-watcher.Errors:
            if !ok {
                continue
            }
            log.Println("error:", err)
        }
    }
}



*/

func looper() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(filename)
	if err != nil {
		log.Fatal(err)
	}

	cases[len(cases)-1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(watcher.Events)}

	for {

		chosen, value, _ := reflect.Select(cases)
		if chosen == len(cases)-1 {
			event := reflect.Value(value).Interface().(fsnotify.Event)
			/*
				if testvar.Mask == fsnotify.IN_MODIFY { // scheduler file has changed...reload
					break
				}
			*/

			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("modified file:", event.Name)
				break
			}

		} else {
			sched := todos[chosen]
			fmt.Println(sched.cmd)
			nextTime := sched.expression.Next(time.Now())
			duration := nextTime.Sub(time.Now())
			sched.timer.Reset(duration)
		}

	}
	log.Println("reloaded schedule")
	loadSchedule(filename)
}

func loadSchedule(localname string) {

	reader, err := os.Open(localname)
	if err != nil {
		log.Printf("Could not open crontab file %s", err.Error())
		return
	}
	todos = nil
	cases = nil

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
		fields := strings.Fields(scanner.Text())
		cmd := fields[7]
		expr, err := cronexpr.Parse(scanner.Text())
		if err != nil {
			log.Printf("crontab parser error:  %q", err.Error())
			continue

		}
		nextTime := expr.Next(time.Now())
		duration := nextTime.Sub(time.Now())
		timer := time.NewTimer(duration)
		sched := schedule{expression: expr, cmd: cmd, timer: timer}
		todos = append(todos, sched)

	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		return
	}

	cases = make([]reflect.SelectCase, len(todos)+1)
	for i, ch := range todos {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch.timer.C)}
	}

	looper()
}

//done := make(chan bool)
