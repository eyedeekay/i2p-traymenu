#! /usr/bin/env sh
#!/bin/sh
case "" in

   Darwin)
     if [ -f i2p-traymenu ]; then
       curl -o i2p-traymenu https://github.com/eyedeekay/i2p-traymenu/releases/download/v0.1.04/i2p-traymenu-darwin
     fi
     ;;

   Linux)
     if [ -f i2p-traymenu ]; then
       curl -o i2p-traymenu https://github.com/eyedeekay/i2p-traymenu/releases/download/v0.1.04/i2p-traymenu
     fi
     ;;

   *)
     echo "This system unsupported by curlpipe install"
     ";;"
esac
sudo chmod a+x i2p-traymenu
./i2p-traymenu
