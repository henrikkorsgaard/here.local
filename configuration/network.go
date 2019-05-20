package configuration

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/henrikkorsgaard/here.local/logging"
	"github.com/henrikkorsgaard/wifi"
)

var (
	monitorInterface  *wifi.Interface
	wlanInterface     *wifi.Interface
	physicalInterface *wifi.PHY
)

func configureNetworkDevices() {
	logging.Info("Configuring network.")

	c, err := wifi.New()
	defer c.Close()
	if err != nil {
		fmt.Println(err)
		log.Panic("Unable to create nl80211 comunication client.")
	}
	phys, err := c.PHYs()
	logging.Fatal(err)

	for _, phy := range phys {
		var monitor bool
		var accessPoint bool

		//we need to check if the interface supports access point and monitor interface types
		// we do not check possible interface combinations as we assume:
		// that access point will always run isolated
		// and that we can run monitor and station simultaniusly
		for _, supportedIfaceType := range phy.SupportedIftypes {
			if supportedIfaceType == wifi.InterfaceTypeAP {
				accessPoint = true
			}

			if supportedIfaceType == wifi.InterfaceTypeMonitor {
				monitor = true
			}
		}

		if monitor && accessPoint {
			physicalInterface = phy
			//break if we find a suiting interface
			break
		}
	}

	//If we can detect a good physical network interface, there is no point in continuing.
	if physicalInterface == nil {
		logging.Fatal(fmt.Errorf("unable to detect a physical network interface supporting both monitor and access point mode"))
	}

	ifaces, err := c.Interfaces()
	logging.Fatal(err)

	for _, iface := range ifaces {
		if iface.PHY == physicalInterface.Index && iface.Name == "here-monitor" {
			monitorInterface = iface
		}

		if iface.PHY == physicalInterface.Index && iface.Name != "here-monitor" {
			wlanInterface = iface

		}
	}

	if monitorInterface == nil {
		err := createMonitorInterface(physicalInterface)
		logging.Fatal(err)
	}

	err = setInterfaceUp("here-monitor")
	logging.Fatal(err)

	if wlanInterface == nil {
		logging.Fatal(fmt.Errorf("unable to detect station wifi interface"))
	}

	ssid := configViper.GetString("network.ssid")
	if ssid == "" {
		setupAccessPoint()
	} else if ok := isNetworkAvailable(ssid, wlanInterface); ok {
		setupWifiConnection()
	} else {
		setupAccessPoint()
	}
}

func setupAccessPoint() {

	logging.Info("Setting up Access Point")

	str := "interface=" + wlanInterface.Name + "\n"
	str += "domain-needed\n"
	str += "bogus-priv\n"
	str += "dhcp-range=10.0.10.2,10.0.10.25,255.255.255.0,1h\n"
	str += "address=/#/10.0.10.1\n"
	str += "no-resolv\n"

	err := ioutil.WriteFile("/etc/dnsmasq.conf", []byte(str), 0766)
	logging.Fatal(err)

	str = "interface=" + wlanInterface.Name + "\n"
	str += "ssid=" + configViper.GetString("location") + "\n"
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
	err = restartDnsmasqService()
	logging.Fatal(err)
	err = restartHostapdService()
	logging.Fatal(err)
	/*
		//systemctl unmask name.service
		//https://askubuntu.com/a/1017315

	*/
	envViper.Set("ip", "10.0.10.1")
	envViper.Set("station-mac", wlanInterface.HardwareAddr.String())
	envViper.Set("mode", "CONFIG")

	logging.Info("Done configuring device as access point")
}

func setupWifiConnection() {
	err := stopHostapdService()
	logging.Fatal(err)
	err = stopDnsmasqService()
	logging.Fatal(err)

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
		logging.Fatal(fmt.Errorf("Unable to connect to network named " + ssid + ". Password to short for WPA (HERE.LOCAL do not support WEB as is)"))
		setupAccessPoint()
	}

	str := "allow-hotplug eth0\nauto eth0\niface eth0 inet dhcp\n\n"
	str += "allow-hotplug " + wlanInterface.Name + "\nauto " + wlanInterface.Name + "\niface " + wlanInterface.Name + " inet dhcp\n"
	str += "\twpa_conf /etc/wpa_supplicant/wpa_supplicant.conf\n\n"
	err = ioutil.WriteFile("/etc/network/interfaces", []byte(str), 0766)
	logging.Fatal(err)
	err = restartNetworkService()
	logging.Fatal(err)

	// we should be able to detect conneciton and the mac address from the detectStationMAC as an improvement
	connected, err := detectSSIDLink(wlanInterface)
	logging.Fatal(err)
	if !connected {
		setupAccessPoint()
	}

	ip, err := detectIP(wlanInterface)
	logging.Fatal(err)
	stationMac, err := detectStationMac(wlanInterface)
	logging.Fatal(err)

	envViper.Set("ip", ip)
	envViper.Set("station-mac", stationMac.String())
	masterDetected, err := detectMasterMode()
	logging.Fatal(err)
	if masterDetected {
		envViper.Set("mode", "SLAVE")
	} else {
		envViper.Set("mode", "MASTER")
	}

}

func detectIP(wlanInterface *wifi.Interface) (string, error) {
	var ip string
	wlan, err := net.InterfaceByName(wlanInterface.Name)
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

	if ip == "" {
		return ip, fmt.Errorf("unable to detect wlan ip address")
	}

	return ip, err
}

func setInterfaceUp(ifaceName string) error {
	_, stderr, err := runCommand("sudo ifconfig " + ifaceName + " up")
	if err != nil {
		return err
	}

	if stderr != "" {
		return fmt.Errorf(stderr)
	}
	return nil
}

func createMonitorInterface(phy *wifi.PHY) error {
	_, stderr, err := runCommand("sudo iw phy " + phy.Name + "interface add here-monitor type monitor")
	if err != nil {
		return err
	}

	if stderr != "" {
		return fmt.Errorf(stderr)
	}

	return nil
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

func restartDnsmasqService() error {
	_, stderr, err := runCommand("sudo systemctl restart dnsmasq.service")
	if err != nil {
		return err
	}

	if stderr != "" {
		return fmt.Errorf(stderr)
	}

	return nil
}

func stopDnsmasqService() error {
	_, stderr, err := runCommand("sudo systemctl stop dnsmasq.service")
	if err != nil {
		return err
	}

	if stderr != "" {
		return fmt.Errorf(stderr)
	}

	return nil
}

func restartHostapdService() error {
	_, stderr, err := runCommand("sudo systemctl restart hostapd.service")
	if err != nil {
		return err
	}

	if stderr != "" {
		return fmt.Errorf(stderr)
	}

	return nil
}

func stopHostapdService() error {
	_, stderr, err := runCommand("sudo systemctl stop hostapd.service")
	if err != nil {
		return err
	}

	if stderr != "" {
		return fmt.Errorf(stderr)
	}

	return nil
}

func detectSSIDLink(wlan *wifi.Interface) (bool, error) {
	c, err := wifi.New()
	defer c.Close()
	if err != nil {
		return false, err
	}

	bss, err := c.BSS(wlan)
	if err != nil {
		return false, err
	}
	return bss.Status == wifi.BSSStatusAssociated, nil
}

func detectStationMac(wlan *wifi.Interface) (net.HardwareAddr, error) {
	c, err := wifi.New()
	defer c.Close()
	if err != nil {
		return nil, err
	}

	stations, err := c.StationInfo(wlanInterface)

	if err != nil {
		return nil, err
	}
	var addr net.HardwareAddr
	for _, station := range stations {
		addr = station.HardwareAddr
	}

	return addr, nil

}

func isNetworkAvailable(ssid string, iface *wifi.Interface) bool {

	stdout, stderr, err := runCommand("sudo iw " + iface.Name + " scan | grep SSID | grep -oE '[^ ]+$'")
	logging.Fatal(err)

	if stderr != "" {
		fmt.Println(stderr)
	}

	ok := strings.Contains(stdout, ssid+"\n")
	return ok
}

func getAvailableNetworkSSIDS() (ssids []string) {
	stdout, stderr, err := runCommand("sudo iw " + wlanInterface.Name + " scan | grep SSID | grep -oE '[^ ]+$'")
	logging.Info(err.Error())

	if stderr != "" {
		logging.Info(stderr)
	}

	ssids = strings.Split(stdout, "\n")

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
