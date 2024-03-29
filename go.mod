module github.com/eyedeekay/i2p-traymenu

go 1.19

require (
	fyne.io/systray v1.10.0
	github.com/eyedeekay/checki2cp v0.0.21
	github.com/eyedeekay/go-i2pbrowser v0.0.0-20230116041021-0efa770a0f26
	github.com/eyedeekay/go-i2pcontrol v0.1.4
	github.com/eyedeekay/toopie.html v0.0.0-20220731225754-04d68bd43b8d
	github.com/mitchellh/go-ps v1.0.0
)

require (
	github.com/artdarek/go-unzip v1.0.0 // indirect
	github.com/eyedeekay/go-fpw v0.0.5 // indirect
	github.com/eyedeekay/go-i2cp v0.0.0-20190716135428-6d41bed718b0 // indirect
	github.com/eyedeekay/i2pkeys v0.33.0 // indirect
	github.com/eyedeekay/sam3 v0.33.5 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/google/go-github v17.0.0+incompatible // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/tevino/abool v1.2.0 // indirect
	github.com/ybbus/jsonrpc/v2 v2.1.7 // indirect
	golang.org/x/net v0.0.0-20220923203811-8be639271d50 // indirect
	golang.org/x/sys v0.4.0 // indirect
)

replace github.com/artdarek/go-unzip v1.0.0 => github.com/eyedeekay/go-unzip v0.0.0-20220914222511-f2936bba53c2

replace github.com/eyedeekay/go-i2pbrowser => ../../../github.com/eyedeekay/go-i2pbrowser
