# S57 parser cli tool

A very simple cli tool for parsing folders of S57 data using ogr2ogr by GDAL.
Can do folders of ENC files (.000) and shape files (.shp).
Can also polygonise BSB (.kap) files.

```
NAME:
   S57 parser CLI - A CLI tool for parsing ENC files. Can also polygonise BSB files

USAGE:
   s57parse [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR:
   MortenOJ

COMMANDS:
     enc, e      enc [ROOT_DIR] [LAYERNAME] [-s | --simplify]
     bsb, b      bsb [ROOT_DIR]
     shape, shp  shp [ROOT_DIR] [-s | --simplify]
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
