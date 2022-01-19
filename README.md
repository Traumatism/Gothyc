# Gothyc

### Minecraft port scanner

## Installation

1. `git clone https://github.com/traumatism/Gothyc.git`
2. `make install`

## Getting started

```gothyc --ports|-p <port range> --target|-t <CIDR> --threads <integer> --timeout <integer>```

### Examples

* `gothyc -p 25565 -t 51.79.0.0/16 --threads 1000 --timeout 5000`
* `gothyc -p 25560-25570,25580,25600-25605 -t 144.217.10.0/24 --threads 1000 --timeout 5000`

## TODO list

* Add support for every type of target

* Improve MOTD
    * text inside 'extra'

* Add silent mode

* Add support for target files
    * path
    * stdin

* Add support for custom output format
    * json
    * csv
    * markdown
    * html


## Utils

* https://www.ipconvertertools.com/iprange2cidr
