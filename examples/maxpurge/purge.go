package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/codegangsta/cli"
	"github.com/jmervine/go-maxcdn"
)

var alias, token, secret, zone, file string
var start time.Time

func init() {

	// Override cli's default help template
	cli.AppHelpTemplate = `Usage: {{.Name}} [arguments...]
Options:
   {{range .Flags}}{{.}}
   {{end}}
`

	app := cli.NewApp()

	app.Name = "maxpurge"
	app.Version = "0.0.1"

	cli.HelpPrinter = helpPrinter

	app.Flags = []cli.Flag{
		cli.StringFlag{"alias, a", "", "[required] consumer alias"},
		cli.StringFlag{"token, t", "", "[required] consumer token"},
		cli.StringFlag{"secret, s", "", "[required] consumer secret"},
		cli.StringFlag{"zone, z", "", "[required] zone to be purged"},
		cli.StringFlag{"file, f", "", "cached file to be purged"},
	}

	app.Action = func(c *cli.Context) {
		alias = ensureArg(c.String("alias"), "ALIAS", c)
		token = ensureArg(c.String("token"), "TOKEN", c)
		secret = ensureArg(c.String("secret"), "SECRET", c)
		zone = ensureArg(c.String("zone"), "ZONE", c)
		file = c.String("file")
	}

	app.Run(os.Args)

	start = time.Now()
}

func main() {
	max := maxcdn.NewMaxCDN(alias, token, secret)

	i, err := strconv.ParseInt(zone, 0, 64)
	check(err)

	zoneid := int(i)

	var response *maxcdn.GenericResponse
	if file != "" {
		response, err = max.PurgeFile(zoneid, file)
	} else {
		response, err = max.PurgeZone(zoneid)
	}
	check(err)

	if response.Code == 200 {
		fmt.Printf("Purge successful after: %v.\n", time.Since(start))
	}
}

func ensureArg(arg string, key string, c *cli.Context) string {
	if arg == "" {
		if env := os.Getenv(key); env != "" {
			return env
		}
		cli.ShowAppHelp(c)
	}
	return arg
}

func check(err error) {
	if err != nil {
		fmt.Printf("%v.\n\nPurge failed after %v.\n", err, time.Since(start))
		os.Exit(2)
	}
}

// Replace cli's default help printer with cli's default help printer
// plus an exit at the end.
func helpPrinter(templ string, data interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	t := template.Must(template.New("help").Parse(templ))
	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
	w.Flush()
	os.Exit(0)
}