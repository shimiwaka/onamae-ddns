package main

import (
	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"

	"crypto/tls"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Phase int

//go:embed config.json
var configRawData []byte

const (
	LOGIN Phase = iota
	MODIP
	LOGOUT
	END
)

type configData struct {
	Url       string `json:"url"`
	Port      string `json:"port"`
	User_id   string `json:"user_id"`
	Password  string `json:"password"`
	Hostname  string `json:"hostname"`
	Domname   string `json:"domname"`
	IPaddress string
}

type caller struct {
	Config configData
}

func sendCmd(c caller, w telnet.Writer, phase Phase) Phase {
	switch phase {
	case LOGIN:
		fmt.Println("LOGIN")
		oi.LongWriteString(w, "LOGIN\n")
		oi.LongWriteString(w, "USERID:"+c.Config.User_id+"\n")
		oi.LongWriteString(w, "PASSWORD:"+c.Config.Password+"\n")
		oi.LongWriteString(w, ".\n")
	case MODIP:
		fmt.Println("MODIP")
		oi.LongWriteString(w, "MODIP\n")
		oi.LongWriteString(w, "HOSTNAME:"+c.Config.Hostname+"\n")
		oi.LongWriteString(w, "DOMNAME:"+c.Config.Domname+"\n")
		oi.LongWriteString(w, "IPV4:"+c.Config.IPaddress+"\n")
		oi.LongWriteString(w, ".\n")
	case LOGOUT:
		fmt.Println("LOGOUT")
		oi.LongWriteString(w, "LOGOUT\n")
		oi.LongWriteString(w, ".\n")
	default:
		return END
	}
	phase++

	return phase
}

func (c caller) CallTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
	p := make([]byte, 1, 256)
	buffer := ""
	phase := LOGIN

	for {
		n, err := r.Read(p)
		if n > 0 {
			bytes := p[:n]
			fmt.Print(string(bytes))
			buffer += string(bytes)

			switch buffer {
			case "000 COMMAND SUCCESSFUL\n.\n":
				phase = sendCmd(c, w, phase)
				buffer = ""
			case "001 COMMAND ERROR\n.\n":
				fmt.Println("processing command failed")
				phase = END
			case "002 LOGIN ERROR\n.\n":
				fmt.Println("login failed")
				phase = END
			}

			if phase >= END {
				break
			}
		}

		if err != nil {
			fmt.Println("reading response failed")
			break
		}
	}
}

func main() {
	var c configData
	err := json.Unmarshal(configRawData, &c)
	if err != nil {
		fmt.Println("parsing json file failed")
		return
	}

	resp, err := http.Get("https://ifconfig.me/")
	if err != nil {
		fmt.Println("obtaining IP Address failed")
		return
	}
	defer resp.Body.Close()

	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("obtaining IP Address failed")
		return
	}

	c.IPaddress = string(buffer)
	fmt.Println(c.IPaddress)

	tlsConfig := &tls.Config{}
	err = telnet.DialToAndCallTLS(c.Url+":"+c.Port, caller{Config: c}, tlsConfig)

	if err != nil {
		fmt.Println("connection failed")
		return
	}
}
