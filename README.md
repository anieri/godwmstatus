# gods

## Summary

This is a clone of [schachmat's repo](https://github.com/schachmat/gods).
It is a dwm status script that displays date/time, cpu usage, memory 
consumption, and network transfer speeds.

## Dependencies

Only a working Go environment and the xsetroot binary is needed.
The dwm's [statuscolor patch](http://dwm.suckless.org/patches/statuscolors) 
is recommended.

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
