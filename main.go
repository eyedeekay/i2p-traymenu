package main

import (
	"flag"
	"fmt"
	"net"
	//"io/ioutil"
	"log"
	"os"
	//"strings"
	"path/filepath"
	"time"

	"github.com/eyedeekay/checki2cp"
	"github.com/eyedeekay/checki2cp/controlcheck"
	"github.com/eyedeekay/di2prc/lib"
	"github.com/eyedeekay/go-i2pcontrol"
	"github.com/eyedeekay/i2p-traymenu/icon"
	"github.com/eyedeekay/i2p-traymenu/irc"
	"github.com/eyedeekay/i2pbrowser/import"

	Core "github.com/eyedeekay/opentracker"
	"github.com/eyedeekay/sam3/i2pkeys"
	"github.com/eyedeekay/toopie.html/import"
	"github.com/getlantern/systray"
	"github.com/vvampirius/retracker/core/common"
)

var usage = `i2p-traymenu
===========

Tray interface to monitor and manage I2P router service. Basically, a
tray i2pcontrol client. Also has an embedded IRC client.

`

//        -block default:false

var home, _ = os.UserHomeDir()

var (
	host     = flag.String("host", "localhost", "Host of the i2pcontrol and SAM interfaces")
	port     = flag.String("port", "7657", "Port of the i2pcontrol interface")
	dir      = flag.String("dir", filepath.Join(home, "i2p/opt/native-traymenu"), "Path to the configuration directory")
	path     = flag.String("path", "jsonrpc", "Path to the i2pcontrol interface")
	password = flag.String("password", "itoopie", "Password for the i2pcontrol interface")
	shelp    = flag.Bool("h", false, "Show the help message")
	lhelp    = flag.Bool("help", false, "Show the help message")
	debug    = flag.Bool("d", false, "Debug mode")
	ot       = flag.Bool("tracker", false, "Run an open torrent tracker")
	chat     = flag.Bool("irc", true, "Run an IRC client connected to I2P")
	age      = flag.Float64("a", 1800, "Keep 'n' minutes peer in memory")
	sam      = flag.String("sam", "7656", "Port of the SAMv3 interface, host must match i2pcontrol")

//	block    = flag.Bool("block", false, "Block the terminal until the router is completely shut down")
)

var di2prcln net.Listener

func main() {
	flag.Parse()
	if *shelp || *lhelp {
		fmt.Printf(usage)
		flag.PrintDefaults()
		return
	}
	if *chat {
		go trayirc.IRC(*dir)
	}
	if *ot {
		go tracker()
	}
	di2prcln = di2prc.Listen(*host+":"+*sam, "", "")
	onExit := func() {
		defer di2prcln.Close()
		log.Println("Exiting now.")
	}

	systray.Run(onReady, onExit)
}

func tracker() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	config := common.Config{
		Listen:  "127.0.0.1:80",
		Debug:   *debug,
		Age:     *age,
		XRealIP: false,
	}
	Core.New(&config)
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
	mIRC := subMenuTop.AddSubMenuItem("IRC Chat", "Talk to others on I2P IRC")
	mChatOrig := systray.AddMenuItem("Distributed Chat", "(Experimental) Distributed group-chat")
	mStatOrig := systray.AddMenuItem("I2P Router Stats", "View I2P Router Console Statistics")
	systray.AddSeparator()
	mQuitOrig := systray.AddMenuItem("Close Tray", "Close the tray app, but don't shutdown the router")
	mWarnOrig := systray.AddMenuItem("I2P is Running but I2PControl is Not available.\nEnable jsonrpc on your I2P router.", "Warn the user if functionality is limited.")
	sub := true

	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
	}()
	smConsole.Hide()
	smTorrent.Hide()
	smEmail.Hide()
	smServices.Hide()
	smDNS.Hide()
	mChatOrig.Hide()
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
				go i2pbrowser.MainNoEmbeddedStuff([]string{"--app", "http://127.0.0.1:7657/console"})
			}()

			go func() {
				<-smTorrent.ClickedCh
				go i2pbrowser.MainNoEmbeddedStuff([]string{"--app", "http://127.0.0.1:7657/i2psnark/"})
			}()

			go func() {
				<-smEmail.ClickedCh
				go i2pbrowser.MainNoEmbeddedStuff([]string{"--app", "http://127.0.0.1:7657/susimail/"})
			}()

			go func() {
				<-smServices.ClickedCh
				go i2pbrowser.MainNoEmbeddedStuff([]string{"--app", "http://127.0.0.1:7657/i2ptunnel/"})
			}()

			go func() {
				<-smDNS.ClickedCh
				go i2pbrowser.MainNoEmbeddedStuff([]string{"--app", "http://127.0.0.1:7657/susidns/"})
			}()

			go func() {
				<-mIRC.ClickedCh
				go i2pbrowser.MainNoEmbeddedStuff([]string{"--app", "http://127.0.0.1:7669/connect"})
			}()

			go func() {
				<-mBrowseOrig.ClickedCh
				log.Println("Launching an I2P Browser")
				go i2pbrowser.MainNoEmbeddedStuff(nil)
			}()

			go func() {
				<-mStatOrig.ClickedCh
				log.Println("Launching toopie.html")
				go toopiexec.Run()
			}()

			go func() {
				<-mChatOrig.ClickedCh
				log.Println("Launching di2prc")
				go i2pbrowser.MainNoEmbeddedStuff([]string{"about:blank", "http://" + di2prcln.Addr().(i2pkeys.I2PAddr).Base32()})
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
