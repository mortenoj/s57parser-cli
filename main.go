package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/urfave/cli"
)

// Command struct
type Command struct {
	Name string
	Args []string
}

var app = cli.NewApp()

func main() {
	info()
	commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func info() {
	app.Name = "S57 parser CLI"
	app.Usage = "A CLI tool for parsing ENC files. Can also polygonise BSB files"
	app.Author = "MortenOJ"
	app.Version = "1.0.0"
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name:    "enc",
			Aliases: []string{"e"},
			Usage:   "enc [ROOT_DIR] [LAYERNAME] [-s | --simplify]",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "simplify, s",
				},
			},
			Action: parseENCDir,
		},
		{
			Name:    "bsb",
			Aliases: []string{"b"},
			Usage:   "bsb [ROOT_DIR]",
			Action:  parseBSBDir,
		},
		{
			Name:    "shape",
			Aliases: []string{"shp"},
			Usage:   "shp [ROOT_DIR] [-s | --simplify]",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "simplify, s",
				},
			},
			Action: parseSHPDir,
		},
	}
}

func parseSHPDir(c *cli.Context) error {
	if len(c.Args()) < 1 || len(c.Args()) > 2 {
		fmt.Println("\n\t\t\tToo few or too many arguments called")
		return errors.New("\tUse --help | -h to see usage")
	}

	simplify := false

	if contains(c.Args(), "-s") || contains(c.Args(), "--simplify") {
		simplify = true
	}

	rootdir := c.Args().Get(0)
	outdir := "./shp_output"
	os.MkdirAll(outdir, os.ModePerm)

	libRegEx, err := regexp.Compile(`([a-zA-Z0-9\s_\\.\-\(\):])+(.shp)$`)
	if err != nil {
		return err
	}

	err = filepath.Walk(rootdir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !libRegEx.MatchString(info.Name()) {
			return err
		}
		fmt.Printf("\rParsing file: %s...", info.Name())

		outFile := fmt.Sprintf("%s/%s.json", outdir, info.Name())
		cmd := &Command{
			Name: "ogr2ogr",
			Args: []string{"-f", "GeoJSON", outFile, path},
		}

		if simplify {
			cmd.Args = append([]string{"-simplify", "0.125"}, cmd.Args...)
		}

		return cmd.exec()
	})

	if err != nil {
		return err
	}

	fmt.Println("\nStored GeoJSON files in: ", outdir)

	return nil
}

func parseBSBDir(c *cli.Context) error {
	if len(c.Args()) < 1 || len(c.Args()) > 2 {
		fmt.Println("\n\t\t\tToo few or too many arguments called")
		return errors.New("\tUse --help | -h to see usage")
	}

	rootdir := c.Args().Get(0)
	outdir := "./bsb_output"
	os.MkdirAll(outdir, os.ModePerm)

	libRegEx, err := regexp.Compile(`([a-zA-Z0-9\s_\\.\-\(\):])+(.kap)$`)
	if err != nil {
		return err
	}

	err = filepath.Walk(rootdir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !libRegEx.MatchString(info.Name()) {
			return err
		}

		outFile := fmt.Sprintf("%s/%s.json", outdir, info.Name())
		cmd := &Command{
			Name: "gdal_polygonize.py",
			Args: []string{path, "-f", "GeoJSON", outFile},
		}
		fmt.Printf("\rParsing %s...", info.Name())

		return cmd.exec()
	})

	if err != nil {
		return err
	}

	return nil
}

func parseENCDir(c *cli.Context) error {
	if len(c.Args()) < 1 || len(c.Args()) > 2 {
		fmt.Println("\n\t\t\tToo few or too many arguments called")
		return errors.New("\tUse --help | -h to see usage")
	}

	simplify := false

	if c.Bool("s") {
		simplify = true
	}

	rootdir := c.Args().Get(0)
	outdir := "./enc_output"
	os.MkdirAll(outdir, os.ModePerm)

	layername := "LNDARE"
	if c.Args().Get(1) != "" {
		layername = c.Args().Get(1)
	}
	layerFolder := fmt.Sprintf("%s/%s", outdir, layername)
	os.MkdirAll(layerFolder, os.ModePerm)

	// Match for *.000 files
	libRegEx, err := regexp.Compile(`([a-zA-Z0-9\s_\\.\-\(\):])+(.000)$`)
	if err != nil {
		return err
	}

	err = filepath.Walk(rootdir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !libRegEx.MatchString(info.Name()) {
			return err
		}

		outFile := fmt.Sprintf("%s/%s/%s.json", outdir, layername, info.Name())
		cmd := &Command{
			Name: "ogr2ogr",
			Args: []string{"-f", "GeoJSON", outFile, path, layername},
		}

		if simplify {
			cmd.Args = append([]string{"-simplify", "0.125"}, cmd.Args...)
		}

		fmt.Printf("\rParsing layer: %s from file: %s...", layername, info.Name())
		return cmd.exec()
	})

	if err != nil {
		return err
	}

	fmt.Println("\nStored GeoJSON files in: ", outdir)
	return nil
}

func (cmd *Command) exec() error {
	if log, err := exec.Command(cmd.Name, cmd.Args...).Output(); err != nil {
		fmt.Println(fmt.Sprintf("\nError running %s %s: %s \n", cmd.Name, cmd.Args, string(log)))
	}
	return nil
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
