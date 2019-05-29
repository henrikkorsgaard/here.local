package configuration

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/henrikkorsgaard/here.local/logging"
	"github.com/henrikkorsgaard/wifi"
)

var (
	monitorInterface *net.Interface
	wlanInterface    *wifi.Interface
	mainPhyInterface *wifi.PHY
//	scanInterface    *wifi.Interface
)

func configureNetworkDevices() {
	logging.Info("Configuring network.")

	mainPhyInterface, err := detectCompatibleNetworkDevice()
	logging.Fatal(err)

	logging.Info("Found compatible network device " + mainPhyInterface.Name)

	monitorInterface, err = detectMonitorInterface(mainPhyInterface)
	logging.Fatal(err)

	logging.Info("Setting up here-monitor interface")

	wlanInterface, err = detectWLANInterface(mainPhyInterface)
	logging.Fatal(err)

//	scanInterface, err = detectScanInterface(wlanInterface)
//	logging.Fatal(err)

	logging.Info("Setting up wlan interface")

	ssid := configViper.GetString("network.ssid")
	if ssid == "" {
		setupAccessPoint()
	} else if ok := isNetworkAvailable(ssid, wlanInterface); ok {
		setupWifiConnection()
	} else {
		setupAccessPoint()
	}
}

func detectCompatibleNetworkDevice() (phy *wifi.PHY, err error) {
	c, err := wifi.New()
	defer c.Close()
	if err != nil {
		return
	}

	phys, err := c.PHYs()
	if err != nil {
		return
	}

	for _, phyIface := range phys {
		var monitor bool
		var accessPoint bool

		//we need to check if the interface supports access point and monitor interface types
		// we do not check possible interface combinations as we assume:
		// that access point will always run isolated
		// and that we can run monitor and station simultaniusly
		for _, supportedIfaceType := range phyIface.SupportedIftypes {
			if supportedIfaceType == wifi.InterfaceTypeAP {
				accessPoint = true
			}

			if supportedIfaceType == wifi.InterfaceTypeMonitor {
				monitor = true
			}
		}

		if monitor && accessPoint {
			phy = phyIface
		}
	}

	if phy == nil {
		err = fmt.Errorf("Unable to detect compantible physical wifi device")
	}

	return
}

func detectScanInterface(wlan *wifi.Interface) (scanIface *wifi.Interface, err error) {
	c, err := wifi.New()
	defer c.Close()
	if err != nil {
		return
	}

	var scanNetIface *net.Interface

	ifaces, err := c.Interfaces()
	if err != nil {
		return
	}

	for _, iface := range ifaces {
		if iface.Name != wlan.Name && iface.Name != "here-monitor" {
			scanNetIface, err = net.InterfaceByName(iface.Name)
			if err != nil {
				return
			}

			scanIface = iface

			break
		}
	}

	if scanNetIface.Flags&net.FlagUp == 0 {
		_, _, err = runCommand("sudo ifconfig "+ scanNetIface.Name + " up") 
		if err != nil || scanNetIface.Flags&net.FlagUp == 0 {
			scanIface = nil
			return
		}

	}

	return
}

func detectMonitorInterface(phy *wifi.PHY) (monIface *net.Interface, err error) {

	monIface, _ = net.InterfaceByName("here-monitor")

	if monIface == nil {
		_, _, err = runCommand("sudo iw phy " + phy.Name + " interface add here-monitor type monitor")
		if err != nil {
			return
		}

		monIface, err = net.InterfaceByName("here-monitor")
		if err != nil || monIface == nil {
			return
		}
	}

	if monIface.Flags&net.FlagUp == 0 {

		_, _, err = runCommand("sudo ifconfig here-monitor up")
		if err != nil {
			return
		}

		if monIface.Flags&net.FlagUp == 0 {
			return
		}
	}

	return
}

func detectWLANInterface(phy *wifi.PHY) (wlanIface *wifi.Interface, err error) {
	c, err := wifi.New()
	defer c.Close()
	if err != nil {
		return
	}

	ifaces, err := c.Interfaces()
	if err != nil {
		return
	}

	var wlanNetIface *net.Interface

	for _, iface := range ifaces {
		if iface.PHY == phy.Index && iface.Name != "here-monitor" {
			wlanIface = iface
			wlanNetIface, err = net.InterfaceByName(iface.Name)
			if err != nil {
				return
			}

			break
		}
	}

	if wlanNetIface == nil {
		_, _, err = runCommand("sudo iw phy " + phy.Name + " interface add here-wlan type managed")
		if err != nil {
			return
		}

		wlanNetIface, err = net.InterfaceByName("here-wlan")
		if err != nil || wlanNetIface == nil {
			return
		}
	}

	if wlanNetIface.Flags&net.FlagUp == 0 {
		_, _, err = runCommand("sudo ifconfig " + wlanNetIface.Name + " up")
		if err != nil {
			return
		}

		if wlanNetIface.Flags&net.FlagUp == 0 {
			return
		}
	}

	return
}

func setupAccessPoint() {

	logging.Info("Setting up Access Point")
	fmt.Println("Setting up Access Point")

	str := "interface=" + wlanInterface.Name + "\n"
	str += "domain-needed\n"
	str += "bogus-priv\n"
	str += "dhcp-range=10.0.10.2,10.0.10.25,255.255.255.0,1h\n"
	str += "address=/#/10.0.10.1\n"
	str += "no-resolv\n"

	err := ioutil.WriteFile("/etc/dnsmasq.conf", []byte(str), 0766)
	logging.Fatal(err)

	str = "interface=" + wlanInterface.Name + "\n"
	str += "ssid=" + configViper.GetString("node.location") + "\n"
	str += "driver=nl80211\n"
	str += "hw_mode=g\n"
	str += "channel=6\n"
	str += "auth_algs=1\n"
	str += "wmm_enabled=0\n"

	err = ioutil.WriteFile("/etc/hostapd/hostapd.conf", []byte(str), 0766)
	logging.Fatal(err)

	str = "auto eth0\nallow-hotplug eth0\niface eth0 inet dhcp\n\n"
	str += "auto " + wlanInterface.Name + "\niface " + wlanInterface.Name + " inet static\n"
	str += "\taddress 10.0.10.1\n\tnetmask 255.255.255.0\n\tnetwork 10.0.10.0\n"

	err = ioutil.WriteFile("/etc/network/interfaces", []byte(str), 0766)
	logging.Fatal(err)

	str = "DAEMON_CONF=\"/etc/hostapd/hostapd.conf\"\n"
	err = ioutil.WriteFile("/etc/default/hostapd", []byte(str), 0766)
	logging.Fatal(err)

	err = restartNetworkService()
	logging.Fatal(err)
	err = startAccessPointServices()
	logging.Fatal(err)

	/*
		//systemctl unmask name.service
		//https://askubuntu.com/a/1017315

	*/
	envViper.Set("ip", "10.0.10.1")
	envViper.Set("station", wlanInterface.HardwareAddr.String())
	envViper.Set("mode", "CONFIG")
	fmt.Println("this far")
	logging.Info("Access Point configured with ssid: " + configViper.GetString("node.location") + ", ip: 10.0.10.1, in CONFIG mode!")
}

func setupWifiConnection() {
	logging.Info("Setting up WLAN conncection")
	password := configViper.GetString("network.password")
	ssid := configViper.GetString("network.ssid")

	if password == "" {
		wpa := "network={\n\tssid=\"" + ssid + "\"\n\tkey_mgmt=NONE\n}"
		err := ioutil.WriteFile("/etc/wpa_supplicant/wpa_supplicant.conf", []byte(wpa), 0766)
		logging.Fatal(err)
	} else if len(password) > 7 {
		_, _, err := runCommand("sudo wpa_passphrase " + ssid + " " + password + " > /etc/wpa_supplicant/wpa_supplicant.conf")
		logging.Fatal(err)
	} else {
		logging.Info("Unable to connect to network named " + ssid + ". Password to short for WPA (HERE.LOCAL do not support WEB as is)")
		setupAccessPoint()
		return
	}

	str := "allow-hotplug eth0\nauto eth0\niface eth0 inet dhcp\n\n"
	str += "allow-hotplug " + wlanInterface.Name + "\nauto " + wlanInterface.Name + "\niface " + wlanInterface.Name + " inet dhcp\n"
	str += "\twpa_conf /etc/wpa_supplicant/wpa_supplicant.conf\n\n"
	err := ioutil.WriteFile("/etc/network/interfaces", []byte(str), 0766)
	logging.Fatal(err)
	err = restartNetworkService()
	logging.Fatal(err)

	station, err := detectLinkAddress(wlanInterface)

	if err != nil || station == nil {
		logging.Info("Unable to associate WLAN with link " + ssid + "! Aborting WLAN configuration.")
		setupAccessPoint()
		return
	}

	ip, err := detectIP(wlanInterface)
	logging.Fatal(err)

	if ip == "" {
		logging.Info("Unable to accuire IP! Aborting WLAN configuration.")
		setupAccessPoint()
		return
	}

	envViper.Set("ip", ip)
	envViper.Set("station", station.String())
	masterDetected, err := detectMasterMode()
	logging.Fatal(err)
	if masterDetected {
		envViper.Set("mode", "SLAVE")
	} else {
		envViper.Set("mode", "MASTER")
	}

	logging.Info("WLAN configured and connected to " + ssid + " with ip " + ip + " in " + envViper.GetString("mode") + " mode.")
}

func detectIP(wlanIface *wifi.Interface) (ip string, err error) {
	wlan, err := net.InterfaceByName(wlanIface.Name)
	if err != nil {
		return ip, err
	}
	addrs, err := wlan.Addrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}

	return ip, err
}

func restartNetworkService() error {
	_, stderr, err := runCommand("sudo systemctl restart networking.service")
	if err != nil {
		return err
	}

	if stderr != "" {
		return fmt.Errorf(stderr)
	}

	return nil
}

func startAccessPointServices() error {

	_, stderr, err := runCommand("sudo systemctl start dnsmasq.service")
	if err != nil {
		return err
	}

	if stderr != "" {
		return fmt.Errorf(stderr)
	}

	_, stderr, err = runCommand("sudo systemctl start hostapd.service")
	if err != nil {
		fmt.Println(err)
		return err
	}

	if stderr != "" {
		return fmt.Errorf(stderr)
	}

	return nil

}
func stopAccessPointServices() error {
	_, stderr, err := runCommand("sudo systemctl stop dnsmasq.service")
	if err != nil {
		return err
	}

	if stderr != "" {
		return fmt.Errorf(stderr)
	}

	_, stderr, err = runCommand("sudo systemctl stop hostapd.service")
	if err != nil {
		return err
	}

	if stderr != "" {
		return fmt.Errorf(stderr)
	}

	return nil
}

func detectLinkAddress(wlan *wifi.Interface) (addr *net.HardwareAddr, err error) {
	c, err := wifi.New()
	defer c.Close()
	if err != nil {
		return
	}

	bss, err := c.BSS(wlan)
	if err != nil {
		return
	}
	if bss.Status == wifi.BSSStatusAssociated {
		addr = &bss.BSSID
	}
	return
}

func isNetworkAvailable(ssid string, iface *wifi.Interface) bool {
	stdout, _, _ := runCommand("sudo iw " + iface.Name + " scan | grep SSID | grep -oE '[^ ]+$'")
	ok := strings.Contains(stdout, ssid+"\n")
	return ok
}

func getSSIDList() (ssids []string) {
	mode := envViper.GetString("mode")

	var stdout, stderr string
	var err error

	fmt.Println("we got this far, right?")

	if mode == "CONFIG" {
		//if scanInterface != nil {
		//	stdout, stderr, err = runCommand("sudo iw " + scanInterface.Name + " scan | grep SSID | grep -oE '[^ ]+$'")
		//} else {
			stopAccessPointServices()
			stdout, stderr, err = runCommand("sudo iw " + wlanInterface.Name + " scan | grep SSID | grep -oE '[^ ]+$'")
			startAccessPointServices()
		//}

	} else {
		stdout, stderr, err = runCommand("sudo iw " + wlanInterface.Name + " scan | grep SSID | grep -oE '[^ ]+$'")
	}


	logging.Fatal(err)
	if stderr != "" {
		logging.Info(stderr)
	}

	raw := strings.Split(stdout, "\n")
	for _, ssidRaw := range raw {
		if ssidRaw != "" {
			exist := false

			for _, ssid := range ssids {
				if ssidRaw == ssid {
					exist = true
					break
				}
			}

			if !exist {
				ssids = append(ssids, ssidRaw)
			}
		}
	}

	return
}

func runCommand(command string) (stdout string, stderr string, err error) {
	cmd := exec.Command("/bin/sh", "-c", command)
	var stdoutBuffer, stderrBuffer bytes.Buffer
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stderrBuffer
	err = cmd.Run()
	if err != nil {
		return "", "", err
	}

	return string(stdoutBuffer.Bytes()), string(stderrBuffer.Bytes()), err
}

func detectMasterMode() (bool, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return false, err
	}

	foundService := false

	entries := make(chan *zeroconf.ServiceEntry)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			if entry.Instance == "here.local.context.server" {
				foundService = true
				break
			}
		}
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err = resolver.Browse(ctx, "_http._tcp", "local.", entries)
	if err != nil {
		return false, err
	}

	<-ctx.Done()
	return foundService, nil
}

/*
	go utils.ExecuteSystemCommand("avahi-publish -a -R here.local 10.0.10.1") //This need to run in the background
	_, err = zeroconf.Register("go-proxi-context-server", "_http._tcp", "local.", 1337, []string{"txtv=0", "lo=1", "la=2"}, nil)
	go utils.ExecuteSystemCommand("avahi-publish -a -R " + location + ".local " + ip.String()) //Run in the background as it blocks
	_, err = zeroconf.Register("go-proxi-context-server-node", "_http._tcp", "local.", 80, []string{"txtv=0", "lo=1", "la=2"}, nil)
	isAP = true
*/
