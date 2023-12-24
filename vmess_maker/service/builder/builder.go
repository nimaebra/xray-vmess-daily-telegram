package builder

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"xray-telegram/entity"

	"github.com/google/uuid"
)

// SetServerIP  Get preferred outbound ip of this machine
func (b *Builder) SetServerIP() *Builder {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println("error during the SetServerIP ", err)
		return nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	b.ServerIP = localAddr.IP.String()
	return b
}

// SetSettingsFile returns the settings file
func (b *Builder) SetSettingsFile() *Builder {

	// Open our jsonFile
	jsonFile, err := os.Open("./setting.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println("error open setting file", err)
		return nil
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("error ReadAll json file", err)
		return nil
	}

	setting := entity.Setting{}
	// we unmarshal our byteArray which contains our
	err = json.Unmarshal(byteValue, &setting)
	if err != nil {
		fmt.Println("error unmarshal setting", err)
		return nil
	}

	b.Setting = setting
	return b

}

// SetConfigurations sets the xray configuration
func (b *Builder) SetConfigurations() *Builder {

	if b.privateKey == "" || b.publicKey == "" {
		fmt.Println("private key or public key is empty")
		return nil
	}

	b.newVmess.Inbounds = make([]entity.Inbound, 1)

	var inbound entity.Inbound
	inbound.Listen = "null"
	inbound.Port = 443
	inbound.Protocol = "vmess"
	inbound.Settings.Clients = make([]entity.Client, 1)
	inbound.Settings.Clients[0].Email = b.Setting.ChannelName
	inbound.Settings.Clients[0].ID = uuid.New().String()

	inbound.Sniffing.Enabled = true
	inbound.Sniffing.DestOverride = []string{"http", "tls", "quic", "fakedns"}

	inbound.StreamSettings.Network = "tcp"
	inbound.StreamSettings.Security = "none"

	inbound.StreamSettings.Sockopt.AcceptProxyProtocol = false
	inbound.StreamSettings.Sockopt.Mark = 0
	inbound.StreamSettings.Sockopt.Tproxy = "off"
	inbound.StreamSettings.Sockopt.TCPFastOpen = true

	inbound.StreamSettings.TCPSettings.AcceptProxyProtocol = false
	inbound.StreamSettings.TCPSettings.Header.Request.Headers.Host = make([]string, 1)
	inbound.StreamSettings.TCPSettings.Header.Request.Headers.Host[0] = "mashhad1.irancell.ir,shiraz1.irancell.ir,tabriz1.irancell.ir,speedtest1.irancell.ir,ahvaz1.irancell.ir,esfahan1.irancell.ir"
	inbound.StreamSettings.TCPSettings.Header.Request.Method = "GET"
	inbound.StreamSettings.TCPSettings.Header.Request.Path = []string{"/speedtest"}
	inbound.StreamSettings.TCPSettings.Header.Response.Headers.Connection = make([]string, 1)
	inbound.StreamSettings.TCPSettings.Header.Response.Headers.Connection[0] = "keep-alive"

	inbound.StreamSettings.TCPSettings.Header.Response.Headers.ContentLength = make([]string, 1)
	inbound.StreamSettings.TCPSettings.Header.Response.Headers.ContentLength[0] = "109"
	inbound.StreamSettings.TCPSettings.Header.Response.Headers.ContentType = make([]string, 1)
	inbound.StreamSettings.TCPSettings.Header.Response.Headers.ContentType[0] = "text/html"

	inbound.StreamSettings.TCPSettings.Header.Response.Reason = "OK"
	inbound.StreamSettings.TCPSettings.Header.Response.Status = "200"
	inbound.StreamSettings.TCPSettings.Header.Response.Version = "HTTP/1.1"

	port := strconv.Itoa(inbound.Port)

	inbound.Tag = "inbound-" + port

	code := "{\"add\":\"" + b.ServerIP + "\",\"aid\":\"0\",\"host\":\"" + inbound.StreamSettings.TCPSettings.Header.Request.Headers.Host[0] + "\",\"id\":\"" + inbound.Settings.Clients[0].ID + "\",\"net\":\"tcp\",\"path\":\"/speedtest\",\"port\":\"" + port + "\",\"ps\":\"" + b.Setting.ChannelName + "\",\"scy\":\"auto\",\"sni\":\"\",\"tls\":\"\",\"type\":\"http\",\"v\":\"2\"}"
	base64code := base64.StdEncoding.EncodeToString([]byte(code))
	StringConfig := "vmess://" + base64code

	b.StringConfigZero = StringConfig
	b.newVmess.Inbounds[0] = inbound

	return b

}

// SetPublicKeyAndPrivateKey read public and private key from key pair
func (b *Builder) SetPublicKeyAndPrivateKey() *Builder {

	dat, err := os.ReadFile("./key_pair.txt")
	if err != nil {
		fmt.Println("error during the key pair")
		return nil
	}
	allData := string(dat)

	allData = strings.TrimSpace(allData)
	allData = strings.ReplaceAll(allData, " ", "")

	privateKeyFirst := RemoveRightPart(allData, "Publickey:")
	privateKey := RemoveLeftPart(privateKeyFirst, "Privatekey:")

	pubAns := strings.SplitAfter(allData, "Publickey:")
	publicKey := pubAns[1]

	b.privateKey = privateKey
	b.publicKey = publicKey

	return b

}

// SetBlock block Iranian and Chinese and porn websites
func (b *Builder) SetBlock() *Builder {

	b.newVmess.Log.Loglevel = "warning"
	b.newVmess.Routing.DomainStrategy = "IPOnDemand"
	b.newVmess.Routing.Rules = make([]entity.Rule, 2)
	b.newVmess.Routing.Rules[0] = entity.Rule{
		Type:        "field",
		IP:          []string{"geoip:cn", "geoip:ir"},
		OutboundTag: "block",
	}

	b.newVmess.Routing.Rules[1] = entity.Rule{
		Type:        "field",
		Domain:      []string{"geosite:category-porn"},
		OutboundTag: "block",
	}

	b.newVmess.Outbounds = make([]entity.Outbound, 2)
	b.newVmess.Outbounds[0] = entity.Outbound{
		Tag:      "direct",
		Protocol: "freedom",
		Settings: struct{}{},
	}
	b.newVmess.Outbounds[1] = entity.Outbound{
		Tag:      "block",
		Protocol: "blackhole",
		Settings: struct{}{},
	}

	return b

}

func (b *Builder) Save() *Builder {

	if b.StringConfigZero == "" {
		fmt.Println(" string config zero is empty")
		return nil
	}

	//save new Reality in file
	err := WriteFile("./config.json", b.newVmess)
	if err != nil {
		log.Fatal("error during the WriteFile ", err)
		return nil
	}

	return b

}
