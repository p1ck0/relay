package cli

import (
    "github.com/fatih/color"
    "github.com/urfave/cli/v2"
    "fmt"
    "os"
    "log"
)

//CliApp - launch cli
func CliApp(port, servers string) {
    app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Value:       "8888",
				Usage:       "the port on which the server will run",
				Aliases:     []string{"p"},
				Destination: &port,
			},
			&cli.StringFlag{
				Name:        "conn",
				Value:       "",
				Usage:       "the port on which the server will run",
				Aliases:     []string{"c"},
				Destination: &servers,
			},
		},
		Action: func(c *cli.Context) error {
            colorgrenn := color.New(color.FgBlack).Add(color.BgHiCyan)
            logoPrint()
			colorgrenn.Printf("*:*:*& USES PORT " + port + " &*:*:*")
			fmt.Printf("\n")
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}