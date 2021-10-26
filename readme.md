## FallingElk 
[![License](https://img.shields.io/badge/license-MIT-_red.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/schwartz1375/fallingelk)](https://goreportcard.com/report/github.com/schwartz1375/fallingelk)

FallingElk is a golang (recommend v1.15, v1.13 min) cross platform webshell.

## Dependencies
* [Go](https://golang.org/) 1.13+
* GNU make
* [goversioninfo](https://github.com/josephspurrier/goversioninfo) (needed to add the windows resource infromation to the binary)
* osslsigncode (needed to sign the windows binary) if you don't care about signing you can comment these lines out in the Makefile.

## How to build
To get the POC working you shouldn't have to modify anything.  You may need to uncomment the log.Fatal lines in the main.go file for troubleshooting, these statements are commented out for threat emulation purposes to excerise the blue IR and RE teams.  Note, you may also need to run "make depends" first to generate the generate self signed certificate. 

## Usage
After the binary has been executed on the target you interact with via the following:

    https://<FQDN or IPaddr>/?cmd=whoami

## OPSEC
Things to consider for windows:
* ```You are embedding credentials in the binary being deployed, take precautions!```
* The name of the binary
* The icon
* The values in the ./resource/verioninfo.json file especially "OriginalFilename" - this should match the binary name set in the Makefile
* The fact this uses a self signed certificate 
