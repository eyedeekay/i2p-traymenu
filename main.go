package main

import (
	"flag"
	"fmt"
	//	"io/ioutil"
	"log"
	//"strings"
	"time"

	"github.com/eyedeekay/checki2cp"
	"github.com/eyedeekay/checki2cp/controlcheck"
	"github.com/eyedeekay/go-i2pcontrol"
	"github.com/eyedeekay/i2p-traymenu/icon"
	"github.com/eyedeekay/i2pbrowser/import"
	"github.com/getlantern/systray"
)

var usage = `i2p-traymenu
===========

Tray interface to monitor and manage I2P router service. Basically, a
tray i2pcontrol client.

        -host default:"127.0.0.1"
        -port default:"7657"
        -path default:"jsonrpc"
        -password default:"itoopie"

Installation with go get

        go get -u github.com/eyedeekay/i2p-traymenu

`

//        -block default:false

var (
	host     = flag.String("host", "localhost", "Host of the i2pcontrol interface")
	port     = flag.String("port", "7657", "Port of the i2pcontrol interface")
	path     = flag.String("path", "jsonrpc", "Path to the i2pcontrol interface")
	password = flag.String("password", "itoopie", "Password for the i2pcontrol interface")
	shelp    = flag.Bool("h", false, "Show the help message")
	lhelp    = flag.Bool("help", false, "Show the help message")

//	block    = flag.Bool("block", false, "Block the terminal until the router is completely shut down")
)

func main() {
	flag.Parse()
	if *shelp || *lhelp {
		fmt.Printf(usage)
		return
	}

	onExit := func() {
		log.Println("Exiting now.")
	}

	systray.Run(onReady, onExit)
}

func onReady() {

	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("I2P Controller")
	systray.SetTooltip("Freestanding Invisisble Internet Router Control Appliance")
	systray.SetTemplateIcon(icon.Data, icon.Data)

	mStartOrig := systray.AddMenuItem("Start I2P", "Start the I2P Service")
	mStopOrig := systray.AddMenuItem("Stop I2P", "Stop the I2P Service")
	mRestartOrig := systray.AddMenuItem("Restart I2P", "Restart the I2P Service")
	mBrowseOrig := systray.AddMenuItem("Launch an I2P Browser", "Start an available browser, configured for I2P")
	mQuitOrig := systray.AddMenuItem("Close Tray", "Close the tray app, but don't shutdown the router")
	mWarnOrig := systray.AddMenuItem("I2P is Running but I2PControl is Not available.\nEnable jsonrpc on your I2P router.", "Warn the user if functionality is limited.")

	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
	}()

	refreshStart := func() {
		ok, err := checki2p.CheckI2PIsRunning()
		if err != nil {
		}
		if ok {
			mStartOrig.Hide()
			mBrowseOrig.Show()
		} else {
			mStartOrig.Show()
			mBrowseOrig.Hide()
		}
	}
	refreshStart()

	go func() {
		for {
			go func() {
				<-mStartOrig.ClickedCh
				checki2p.ConditionallyLaunchI2P()
			}()

			go func() {
				<-mBrowseOrig.ClickedCh
				log.Println("Launching an I2P Browser")
				go i2pbrowser.MainNoEmbeddedStuff()
			}()

			go func() {
				<-mStopOrig.ClickedCh
				log.Println("Beginning to shutdown I2P")
				i2pcontrol.ShutdownGraceful()
				refreshStart()
			}()

			go func() {
				<-mRestartOrig.ClickedCh
				log.Println("Beginning to restart I2P")
				i2pcontrol.RestartGraceful()
				refreshStart()
			}()

			time.Sleep(time.Second)
		}
	}()

	mWarnOrig.Hide()

	refreshMenu := func() {
		ok, err := checki2p.CheckI2PIsRunning()
		if err != nil {
			//mWarnOrig.Show()
		}

		if ok {
			mStopOrig.Show()
			mRestartOrig.Show()
			mBrowseOrig.Show()
		} else {
			mStopOrig.Hide()
			mRestartOrig.Hide()
			mBrowseOrig.Hide()
		}

		i2pcontrol.Initialize(*host, *port, *path)
		_, err = i2pcontrol.Authenticate(*password)
		if err != nil {
			mWarnOrig.Show()
			mStopOrig.Hide()
			mRestartOrig.Hide()
		}
		ok, err = checki2pcontrol.CheckI2PControlEcho(*host, *port, *path, "Will it blend?")
		if err != nil {
			mWarnOrig.Show()
			mStopOrig.Hide()
			mRestartOrig.Hide()
		}
		if ok {
			mWarnOrig.Hide()
		} else {
			mWarnOrig.Show()
			mStopOrig.Hide()
			mRestartOrig.Hide()
		}
	}
	refreshMenu()
	go func() {
		for {
			refreshMenu()
			log.Println("i2pcontrol check succeeded, sleeping for a while")
			time.Sleep(time.Minute)
		}
	}()

}
