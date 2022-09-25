package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	//"io/ioutil"
	"log"

	//"strings"

	"time"

	checki2p "github.com/eyedeekay/checki2cp"
	checki2pcontrol "github.com/eyedeekay/checki2cp/controlcheck"
	fcw "github.com/eyedeekay/go-fpw"
	goi2pbrowser "github.com/eyedeekay/go-i2pbrowser"
	"github.com/eyedeekay/go-i2pcontrol"
	"github.com/eyedeekay/i2p-traymenu/icon"
	toopiexec "github.com/eyedeekay/toopie.html/import"

	"fyne.io/systray"
)

var usage = `i2p-traymenu
===========

Tray interface to monitor and manage I2P router service. Basically, a
tray i2pcontrol client. Also has an embedded IRC client.

`

//        -block default:false

var (
	host     = flag.String("host", "localhost", "Host of the i2pcontrol and SAM interfaces")
	port     = flag.String("port", "7657", "Port of the i2pcontrol interface")
	dir      = flag.String("dir", defaultDir(), "Path to the configuration directory")
	path     = flag.String("path", "jsonrpc", "Path to the i2pcontrol interface")
	password = flag.String("password", "itoopie", "Password for the i2pcontrol interface")
	shelp    = flag.Bool("h", false, "Show the help message")
	lhelp    = flag.Bool("help", false, "Show the help message")
)

func defaultDir() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	if strings.Contains(exe, "plugins/i2p-traymenu") {
		return filepath.Dir(exe)
	}
	// if the path to me is the I2P plugin directory, then use the plugin directory as the default directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// if the working directory is the home directory, then use a default directory inside the I2P directory
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	if home == wd {
		return filepath.Join(home, ".i2p/plugins/i2p-traymenu")
	}
	return wd
}

func profileDir() string {
	return filepath.Join(*dir, "i2p.profile.firefox")
}

func browse(url string) {
	profilePath, err := goi2pbrowser.UnpackBase(profileDir())
	if err != nil {
		log.Println(err)
		return
	}
	FIREFOX, ERROR := fcw.BasicFirefox(profilePath, false, url)
	if ERROR != nil {
		log.Println(ERROR)
		return
	}
	defer FIREFOX.Close()
	<-FIREFOX.Done()
}

func main() {
	flag.Parse()
	if *shelp || *lhelp {
		fmt.Printf(usage)
		flag.PrintDefaults()
		return
	}
	onExit := func() {
		log.Println("Exiting now.")
	}

	systray.Run(onReady, onExit)
}

func onReady() {

	systray.SetTemplateIcon(icon.Icon, icon.Icon)
	systray.SetTitle("I2P Controller")
	systray.SetTooltip("Freestanding Invisisble Internet Router Control Appliance")

	mStartOrig := systray.AddMenuItem("Start I2P", "Start the I2P Service")
	mStopOrig := systray.AddMenuItem("Stop I2P", "Stop the I2P Service")
	mRestartOrig := systray.AddMenuItem("Restart I2P", "Restart the I2P Service")
	systray.AddSeparator()
	mBrowseOrig := systray.AddMenuItem("Launch an I2P Browser", "Start an available browser, configured for I2P")
	subMenuTop := systray.AddMenuItem("I2P Applications", "I2P Applications")
	smConsole := subMenuTop.AddSubMenuItem("I2P Router Console", "Go to the I2P config page")
	smTorrent := subMenuTop.AddSubMenuItem("Bittorrent", "Manage your Bittorrent Client")
	smEmail := subMenuTop.AddSubMenuItem("Mail", "Send and Recieve email")
	smServices := subMenuTop.AddSubMenuItem("Hidden Services Mangager", "Set up and tear down tunnels")
	smDNS := subMenuTop.AddSubMenuItem("Address Book", "Store contact addresses")
	mStatOrig := systray.AddMenuItem("I2P Router Stats", "View I2P Router Console Statistics")
	systray.AddSeparator()
	mQuitOrig := systray.AddMenuItem("Close Tray", "Close the tray app, but don't shutdown the router")
	mWarnOrig := systray.AddMenuItem("I2P is Running but I2PControl is Not available.\nEnable jsonrpc on your I2P router.", "Warn the user if functionality is limited.")
	sub := true

	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
	}()
	refreshStart := func() {
		ok, err := checki2p.CheckI2PIsRunning()
		if err != nil {
			log.Fatalln("I2P failed to start", err)
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
				<-subMenuTop.ClickedCh
				if sub {
					smConsole.Show()
					smTorrent.Show()
					smEmail.Show()
					smServices.Show()
					smDNS.Show()
					t := sub
					sub = !t
				} else {
					smConsole.Hide()
					smTorrent.Hide()
					smEmail.Hide()
					smServices.Hide()
					smDNS.Hide()
					t := sub
					sub = !t
				}
			}()

			go func() {
				<-smConsole.ClickedCh
				go browse("http://127.0.0.1:7657/console")
			}()

			go func() {
				<-smTorrent.ClickedCh
				go browse("http://127.0.0.1:7657/i2psnark/")
			}()

			go func() {
				<-smEmail.ClickedCh
				go browse("http://127.0.0.1:7657/susimail/")
			}()

			go func() {
				<-smServices.ClickedCh
				go browse("http://127.0.0.1:7657/i2ptunnel/")
			}()

			go func() {
				<-smDNS.ClickedCh
				go browse("http://127.0.0.1:7657/susidns/")

			}()

			go func() {
				<-mBrowseOrig.ClickedCh
				log.Println("Launching an I2P Browser")
				go browse("http://127.0.0.1:7657")
			}()

			go func() {
				<-mStatOrig.ClickedCh
				log.Println("Launching toopie.html")
				go toopiexec.Run()
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
			mWarnOrig.Show()
		}

		if ok {
			mWarnOrig.Hide()
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
