package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
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
	host       = flag.String("host", "127.0.0.1", "Host of the i2pcontrol and SAM interfaces")
	port       = flag.String("port", consoleURLPort(), "Port of the i2pcontrol interface")
	dir        = flag.String("dir", defaultDir(), "Path to the configuration directory")
	path       = flag.String("path", "jsonrpc", "Path to the i2pcontrol interface")
	password   = flag.String("password", "itoopie", "Password for the i2pcontrol interface")
	routerconf = flag.String("client", routerConfig(), "path to the client.config file for the router console")
	shelp      = flag.Bool("h", false, "Show the help message")
	lhelp      = flag.Bool("help", false, "Show the help message")
	enableJson = flag.Bool("autoenable", true, "automatically enable the jsonrpc webapp(requires a router restart)")
)

var usability bool

func enableJsonRPC() {
	jsonconf := jsonConfig()
	if jsonconf != "" {
		log.Println("Ensuring jsonrpc.startOnLoad=true in", jsonconf)
		info, err := os.Stat(jsonconf)
		if err != nil {
			panic(err)
		}
		contents, err := ioutil.ReadFile(jsonconf)
		if err != nil {
			panic(err)
		}
		switched := strings.Replace(string(contents), "webapps.jsonrpc.startOnLoad=false", "webapps.jsonrpc.startOnLoad=true", 1)
		err = ioutil.WriteFile(jsonconf, []byte(switched), info.Mode())
		if err != nil {
			panic(err)
		}
	}
}

func jsonConfig() string {
	switch runtime.GOOS {
	case "windows":
		dir := filepath.Join(os.Getenv("LOCALAPPDATA"), "i2p")
		conf := filepath.Join(dir, "webapps.config")
		if _, err := os.Stat(conf); err == nil {
			return conf
		}
	case "linux":
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		dir := filepath.Join(userHomeDir, ".i2p")
		conf := filepath.Join(dir, "webapps.config")
		if _, err := os.Stat(conf); err == nil {
			return conf
		}
		dir = filepath.Join(userHomeDir, "i2p")
		conf = filepath.Join(dir, "webapps.config")
		if _, err := os.Stat(conf); err == nil {
			return conf
		}
	}
	return ""
}

func routerConfig() string {
	switch runtime.GOOS {
	case "windows":
		dir := filepath.Join(os.Getenv("LOCALAPPDATA"), "i2p")
		conf := filepath.Join(dir, "clients.config.d", "00-net.i2p.router.web.RouterConsoleRunner-clients.config")
		if _, err := os.Stat(conf); err == nil {
			return conf
		}
	case "linux":
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		dir := filepath.Join(userHomeDir, ".i2p")
		conf := filepath.Join(dir, "clients.config.d", "00-net.i2p.router.web.RouterConsoleRunner-clients.config")
		if _, err := os.Stat(conf); err == nil {
			return conf
		}
		dir = filepath.Join(userHomeDir, "i2p")
		conf = filepath.Join(dir, "clients.config.d", "00-net.i2p.router.web.RouterConsoleRunner-clients.config")
		if _, err := os.Stat(conf); err == nil {
			return conf
		}
	}
	return ""
}

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
	var profilePath string
	var err error
	if usability {
		profilePath, err = goi2pbrowser.UnpackUsability(profileDir())
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		profilePath, err = goi2pbrowser.UnpackBase(profileDir())
		if err != nil {
			log.Println(err)
			return
		}
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
	if *enableJson {
		enableJsonRPC()
	}
	onExit := func() {
		log.Println("Exiting now.")
	}

	systray.Run(onReady, onExit)
}

func usabilityMode() string {
	if !usability {
		return "Switch Browser to Usability Mode"
	}
	return "Switch Browser to Strict Mode"
}

func consoleURLPort() string {
	if *routerconf == "" {
		return "7657"
	}
	//clientApp.0.args=7657
	contents, err := ioutil.ReadFile(*routerconf)
	if err != nil {
		log.Println("failed to read client config", err)
		return "7657"
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if strings.Contains(line, "clientApp.0.args=") {
			trimmedline := strings.Replace(line, "clientApp.0.args=", "", 1)
			final := strings.Split(trimmedline, " ")[0]
			return final
		}
	}
	return "7657"
}

func consoleURL() string {
	port := consoleURLPort()
	return "http://" + *host + ":" + port
}

func onReady() {

	systray.SetTemplateIcon(icon.Icon, icon.Icon)
	systray.SetTitle("I2P Controller")
	systray.SetTooltip("Freestanding Invisisble Internet Router Control Appliance")

	mStartOrig := systray.AddMenuItem("Start I2P", "Start the I2P Service")
	mStopOrig := systray.AddMenuItem("Stop I2P", "Stop the I2P Service")
	mRestartOrig := systray.AddMenuItem("Restart I2P", "Restart the I2P Service")
	systray.AddSeparator()
	mUsabilitySwitch := systray.AddMenuItem(usabilityMode(), "Toggle browser configurations")
	mConsoleURL := systray.AddMenuItem("Console is available on: "+consoleURL(), "Show console")
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
				<-mUsabilitySwitch.ClickedCh
				usability = !usability
				mUsabilitySwitch.SetTitle(usabilityMode())
			}()

			go func() {
				<-smConsole.ClickedCh
				go browse(consoleURL() + "/console")
			}()

			go func() {
				<-mConsoleURL.ClickedCh
				go browse(consoleURL() + "/console")
			}()

			go func() {
				<-smTorrent.ClickedCh
				go browse(consoleURL() + "/i2psnark/")
			}()

			go func() {
				<-smEmail.ClickedCh
				go browse(consoleURL() + "/susimail/")
			}()

			go func() {
				<-smServices.ClickedCh
				go browse(consoleURL() + "/i2ptunnel/")
			}()

			go func() {
				<-smDNS.ClickedCh
				go browse(consoleURL() + "/susidns/")

			}()

			go func() {
				<-mBrowseOrig.ClickedCh
				log.Println("Launching an I2P Browser")
				go browse(consoleURL())
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
