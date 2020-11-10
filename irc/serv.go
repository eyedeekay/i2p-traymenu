package trayirc

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/mmcloughlin/professor"
	"github.com/prologic/eris/irc"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

var motd string = `
.___               .__         .__ ___.    .__           .___ __________ _________
|   |  ____ ___  __|__|  ______|__|\_ |__  |  |    ____  |   |\______   \\_   ___ \
|   | /    \\  \/ /|  | /  ___/|  | | __ \ |  |  _/ __ \ |   | |       _//    \  \/
|   ||   |  \\   / |  | \___ \ |  | | \_\ \|  |__\  ___/ |   | |    |   \\     \____
|___||___|  / \_/  |__|/____  >|__| |___  /|____/ \___  >|___| |____|_  / \______  /
          \/                \/          \/            \/              \/         \/

        ___                                                 _____ __         __ 
       / _ |  ___  ___   ___  __ __ __ _  ___  __ __ ___   / ___// /  ___ _ / /_
      / __ | / _ \/ _ \ / _ \/ // //  ' \/ _ \/ // /(_-<  / /__ / _ \/ _ '// __/
     /_/ |_|/_//_/\___//_//_/\_, //_/_/_/\___/\_,_//___/  \___//_//_/\_,_/ \__/
     ==========================================================================     

`

// GenerateEncodedPassword generated a bcrypt hashed passwords
// Taken from github.com/prologic/mkpasswd
func GenerateEncodedPassword(passwd []byte) (encoded string, err error) {
	if passwd == nil {
		err = fmt.Errorf("empty password")
		return
	}
	bcrypted, err := bcrypt.GenerateFromPassword(passwd, bcrypt.MinCost)
	if err != nil {
		return
	}
	encoded = base64.StdEncoding.EncodeToString(bcrypted)
	return
}

func OutputServerConfigFile(dir, configfile string) (string, error) {
	operatorpassword, err := password.Generate(14, 2, 2, false, false)
	if err != nil {
		return "", err
	}
	operator, err := GenerateEncodedPassword([]byte(operatorpassword))
	if err != nil {
		return "", err
	}
	accountpassword, err := password.Generate(14, 2, 2, false, false)
	if err != nil {
		return "", err
	}
	account, err := GenerateEncodedPassword([]byte(accountpassword))
	if err != nil {
		return "", err
	}

	var serverconfigfile string = `mutex: {}
  network:
    name: Local
  server:
    password: ""
    listen: []
    tlslisten: {}
    i2plisten:
      invisibleirc:
        i2pkeys: ` + filepath.Join(dir, "iirc.i2pkeys") + `
        samaddr: 127.0.0.1:7656
    #torlisten:
      #hiddenirc:
        #torkeys: ` + filepath.Join(dir, "tirc.torkeys") + `
        #controlport: 0
    log: ""
    motd: ` + filepath.Join(dir, "ircd.motd") + `
    name: myinvisibleirc.i2p
    description: Hidden IRC Services
  operator:
    admin:
      password: ` + operator + `
  account:
    admin:
      password: ` + account

	err = ioutil.WriteFile(filepath.Join(dir, configfile), []byte(serverconfigfile), 0644)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(filepath.Join(dir, "operator-admin-passwd.txt"), []byte(operatorpassword), 0644)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(filepath.Join(dir, "account-admin-passwd.txt"), []byte(accountpassword), 0644)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(filepath.Join(dir, "ircd.motd"), []byte(motd), 0644)
	if err != nil {
		return "", err
	}
	return serverconfigfile, nil
}

func IRCServerMain(version, debug bool, dir, configfile string) {

	if version {
		fmt.Printf(irc.FullVersion())
		os.Exit(0)
	}

	if debug {
		go professor.Launch(":6060")
	}

	config, err := irc.LoadConfig(filepath.Join(dir, configfile))
	if err != nil {
		log.Fatal("IRC Server Config file did not load successfully:", err.Error())
	}

	irc.NewServer(config).Run()
}
