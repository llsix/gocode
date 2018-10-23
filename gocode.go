package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
)

var (
	g_is_server   = flag.Bool("s", false, "run a server instead of a client")
	g_format      = flag.String("f", "nice", "output format (vim | emacs | nice | csv | csv-with-package | json)")
	g_input       = flag.String("in", "", "use this file instead of stdin input")
	g_sock        = create_sock_flag("sock", "socket type (unix | tcp)")
	g_addr        = flag.String("addr", "127.0.0.1:37372", "address for tcp socket")
	g_debug       = flag.Bool("debug", false, "enable server-side debug mode")
	g_profile     = flag.Int("profile", 0, "port on which to expose profiling information for pprof; 0 to disable profiling")
	g_daemon_name = flag.String("daemon", "gocode-daemon-x", "unix socket daemon prefix name")
)

func get_socket_filename() string {
	user := os.Getenv("USER")
	if user == "" {
		user = "all"
	}
	return filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", *g_daemon_name, user))
}

func show_usage() {
	fmt.Fprintf(os.Stderr,
		"Usage: %s [-s] [-f=<format>] [-in=<path>] [-sock=<type>] [-addr=<addr>]\n"+
			"       <command> [<args>]\n\n",
		os.Args[0])
	fmt.Fprintf(os.Stderr,
		"Flags:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr,
		"\nCommands:\n"+
			"  autocomplete [<path>] <offset>     main autocompletion command\n"+
			"  close                              close the gocode daemon\n"+
			"  drop-cache                         drop gocode daemon's cache\n"+
			"  options                            list config options (extended)\n"+
			"  set [<name> [<value>]]             list or set config options\n"+
			"  status                             gocode daemon status report\n"+
			"")
}

func main() {
	flag.Usage = show_usage
	flag.Parse()

	var retval int
	if *g_is_server {
		go func() {
			if *g_profile <= 0 {
				return
			}
			addr := fmt.Sprintf("localhost:%d", *g_profile)
			// Use the following commands to profile the binary:
			// go tool pprof http://localhost:6060/debug/pprof/profile   # 30-second CPU profile
			// go tool pprof http://localhost:6060/debug/pprof/heap      # heap profile
			// go tool pprof http://localhost:6060/debug/pprof/block     # goroutine blocking profile
			// See http://blog.golang.org/profiling-go-programs for more info.
			log.Printf("enabling  profiler on %s", addr)
			log.Print(http.ListenAndServe(addr, nil))
		}()
		retval = do_server()
	} else {
		retval = do_client()
	}
	os.Exit(retval)
}
