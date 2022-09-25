I2P System Tray Control Panel, Cross-Platform
=============================================

**Status:** Maintained. Not recieving new features. Updated at about the same
time as `i2p.firefox`, because `i2p.plugins.firefox` is it's main source of
browser profiles.

This is a very simple system tray application for interacting with I2P. It
can start, stop, or restart an I2P router, or it can launch a browser. It
depends on the presence of an I2P router implementing the i2pcontrol jsonrpc
interface. You must have Firefox or one of it's variants installed.

To enable jsonrpc on the Java I2P router, go to this the webapps config page:
[http://localhost:7657/configwebapps](http://localhost:7657/configwebapps) and
enable the `jsonrpc` app as seen in the screenshot below:

![enable i2pcontrol](i2pcontrol.png)
