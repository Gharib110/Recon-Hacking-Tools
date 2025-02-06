package main

import (
	"bufio"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func parse(filename string) (map[string]string, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	records := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ",", 2)
		if len(parts) < 2 {
			return records, fmt.Errorf("invalid line: %s", line)
		}
		records[parts[0]] = parts[1]
	}
	log.Println("records set to: ")
	for k, v := range records {
		log.Printf("%s -> %s\n", k, v)
	}

	return records, scanner.Err()
}

func main() {
	records, err := parse("proxy.config")
	if err != nil {
		log.Fatal(err)
	}

	var recordLock sync.RWMutex
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		if len(r.Question) < 1 {
			dns.HandleFailed(w, r)
			return
		}

		name := r.Question[0].Name
		parts := strings.Split(name, ".")
		if len(parts) > 1 {
			name = strings.Join(parts[len(parts)-2:], ".")
		}
		recordLock.RLock()
		match, ok := records[name]
		recordLock.RUnlock()

		if !ok {
			dns.HandleFailed(w, r)
			return
		}
		respp, err := dns.Exchange(r, match)
		if err != nil {
			dns.HandleFailed(w, r)
			return
		}
		err = w.WriteMsg(respp)
		if err != nil {
			dns.HandleFailed(w, r)
			return
		}

	})

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGUSR1)

		for sig := range sigs {
			switch sig {
			case syscall.SIGUSR1:
				log.Println("SIGUSR1")
				recordsUpdate, err := parse("proxy.config")
				if err != nil {
					log.Fatal(err)
				} else {
					recordLock.Lock()
					records = recordsUpdate
					recordLock.Unlock()
				}
			}
		}
	}()

	log.Fatal(dns.ListenAndServe(":53", "udp", nil))
}
