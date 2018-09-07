# godwmstatus

## Summary

This is a remake of [schachmat's "gods"](https://github.com/schachmat/gods).
It is a dwm status script that displays date/time, cpu usage, memory
consumption, and network transfer speeds.

## Dependencies

A working Go environment and the xsetroot binary.

## Usage

To install, run

	go get github.com/anieri/godwmstatus

Then add the following line to your `.xinitrc` or wherever you start dwm, but
before actually starting dwm:

	$GOPATH/bin/godwmstatus &

## License

As denoted in the original source:

"THE BEER-WARE LICENSE" (Revision 42):
<teichm@in.tum.de> wrote this file. As long as you retain this notice you
can do whatever you want with this stuff. If we meet some day, and you think
this stuff is worth it, you can buy me a beer in return Markus Teich
