package connect

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// Config represents an ssh config with relevant fields
type Config struct {
	HostBlocks []HostBlock `ssh_config:"HostBlocks"`
	Include    string      `ssh_config:"Include"`
}

// HostBlock represents a host
type HostBlock struct {
	Host                             string   `ssh_config:"Host"`
	Hostname                         string   `ssh_config:"Hostname"`
	Ciphers                          []string `ssh_config:"Ciphers"`
	AddressFamily                    string   `ssh_config:"AddressFamily"`
	BatchMode                        string   `ssh_config:"BatchMode"`
	BindAddress                      string   `ssh_config:"BindAddressBindAddress"`
	ChallengeResponseAuthentication  string   `ssh_config:"ChallengeResponseAuthenticationChallengeResponseAuthentication"`
	CheckHostIP                      string   `ssh_config:"CheckHostIP"`
	Cipher                           string   `ssh_config:"Cipher"`
	ClearAllForwardings              string   `ssh_config:"ClearAllForwardings"`
	Compression                      string   `ssh_config:"Compression"`
	CompressionLevel                 string   `ssh_config:"CompressionLevel"`
	ConnectionAttempts               string   `ssh_config:"ConnectionAttempts"`
	ConnectTimeout                   string   `ssh_config:"ConnectTimeout"`
	ControlMaster                    string   `ssh_config:"ControlMaster"`
	ControlPath                      string   `ssh_config:"ControlPath"`
	DynamicForward                   string   `ssh_config:"DynamicForward"`
	EscapeChar                       string   `ssh_config:"EscapeChar"`
	ExitOnForwardFailure             string   `ssh_config:"ExitOnForwardFailure"`
	ForwardAgent                     string   `ssh_config:"ForwardAgent"`
	ForwardX11                       string   `ssh_config:"ForwardX11"`
	ForwardX11Trusted                string   `ssh_config:"ForwardX11Trusted"`
	GatewayPorts                     string   `ssh_config:"GatewayPorts"`
	GlobalKnownHostsFile             string   `ssh_config:"GlobalKnownHostsFile"`
	GSSAPIAuthentication             string   `ssh_config:"GSSAPIAuthentication"`
	GSSAPIKeyExchange                string   `ssh_config:"GSSAPIKeyExchange"`
	GSSAPIClientIdentity             string   `ssh_config:"GSSAPIClientIdentity"`
	GSSAPIDelegateCredentials        string   `ssh_config:"GSSAPIDelegateCredentials"`
	GSSAPIRenewalForcesRekey         string   `ssh_config:"GSSAPIRenewalForcesRekey"`
	GSSAPITrustDns                   string   `ssh_config:"GSSAPITrustDns"`
	HashKnownHosts                   string   `ssh_config:"HashKnownHosts"`
	HostbasedAuthentication          string   `ssh_config:"HostbasedAuthentication"`
	HostKeyAlgorithms                string   `ssh_config:"HostKeyAlgorithms"`
	HostKeyAlias                     string   `ssh_config:"HostKeyAlias"`
	IdentitiesOnly                   string   `ssh_config:"IdentitiesOnly"`
	IdentityFile                     string   `ssh_config:"IdentityFile"`
	KbdInteractiveAuthentication     string   `ssh_config:"KbdInteractiveAuthentication"`
	KbdInteractiveDevices            string   `ssh_config:"KbdInteractiveDevices"`
	LocalCommand                     string   `ssh_config:"LocalCommand"`
	LocalForward                     string   `ssh_config:"LocalForward"`
	LogLevel                         string   `ssh_config:"LogLevel"`
	MACs                             []string `ssh_config:"MACs"`
	Match                            string   `ssh_config:"Match"`
	NoHostAuthenticationForLocalhost string   `ssh_config:"HoNoHostAuthenticationForLocalhostst"`
	PreferredAuthentications         string   `ssh_config:"PreferredAuthentications"`
	Protocol                         string   `ssh_config:"Protocol"`
	ProxyCommand                     string   `ssh_config:"ProxyCommand"`
	PubkeyAuthentication             string   `ssh_config:"PubkeyAuthentication"`
	RemoteForward                    string   `ssh_config:"RemoteForward"`
	RhostsRSAAuthentication          string   `ssh_config:"RhostsRSAAuthentication"`
	RSAAuthentication                string   `ssh_config:"RSAAuthentication"`
	SendEnv                          string   `ssh_config:"SendEnv"`
	ServerAliveCountMax              string   `ssh_config:"ServerAliveCountMax"`
	ServerAliveInterval              string   `ssh_config:"ServerAliveInterval"`
	SmartcardDevice                  string   `ssh_config:"SmartcardDevice"`
	StrictHostKeyChecking            string   `ssh_config:"StrictHostKeyCheckingStrictHostKeyChecking"`
	TCPKeepAlive                     string   `ssh_config:"TCPKeepAlive"`
	Tunnel                           string   `ssh_config:"Tunnel"`
	TunnelDevice                     string   `ssh_config:"TunnelDevice"`
	User                             string   `ssh_config:"User"`
	UsePrivilegedPort                string   `ssh_config:"UsePrivilegedPort"`
	UserKnownHostsFile               string   `ssh_config:"UserKnownHostsFile"`
	VerifyHostKeyDNS                 string   `ssh_config:"VerifyHostKeyDNS"`
	VisualHostKey                    string   `ssh_config:"VisualHostKey"`
}

var (
	// ErrMalformedField is the error for a HostBlock field without a proper <key> <value>
	ErrMalformedField = errors.New("malformed field")
)

// Write writes a Config to filename
func (c *Config) Write(writer io.Writer) error {
	var err error
	for _, hostBlock := range c.HostBlocks {
		err = hostBlock.write(writer)
		if err != nil {
			return err
		}
	}
	if c.Include != "" {
		_, err = writer.Write([]byte(fmt.Sprintf("\nInclude %s\n", c.Include)))
	}
	return err
}

func (h *HostBlock) write(w io.Writer) error {
	var err error
	v := reflect.ValueOf(h).Elem()
	f := reflect.TypeOf(h).Elem()
	for i := 0; i < v.NumField(); i++ {
		value := f.Field(i).Tag.Get("ssh_config")
		if value == "" || value == "-" || v.Field(i).String() == "" {
			continue
		}

		switch value {
		case "Host":
			_, err = w.Write([]byte(fmt.Sprintf("\n%s %s\n", value, v.Field(i).String())))
		default:
			switch f.Field(i).Type.String() {
			case "string":
				_, err = w.Write([]byte(fmt.Sprintf("\t%s %s\n", value, v.Field(i).String())))
			case "[]string":
				var arr []string
				arr = v.Field(i).Interface().([]string)
				if len(arr) == 0 {
					break
				}
				_, err = w.Write([]byte(fmt.Sprintf("\t%s %s\n", value, strings.Join(arr, ","))))
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Read creates a Config from a reader
func Read(reader io.Reader) (*Config, error) {
	config := &Config{}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if len(strings.TrimSpace(scanner.Text())) == 0 {
			continue
		}
		arr := strings.SplitN(scanner.Text(), " ", 2)
		if len(arr) != 2 {
			return nil, ErrMalformedField
		}

		switch strings.TrimSpace(arr[0]) {
		case "Host":
			config.HostBlocks = append(config.HostBlocks, HostBlock{
				Host: strings.TrimSpace(arr[1]),
			})

		case "Include":
			config.Include = strings.TrimSpace(arr[1])

		default:
			value := reflect.ValueOf(&config.HostBlocks[len(config.HostBlocks)-1]).Elem()
			for i := 0; i < reflect.ValueOf(config.HostBlocks[len(config.HostBlocks)-1]).NumField(); i++ {
				if strings.TrimSpace(arr[0]) != reflect.ValueOf(config.HostBlocks[len(config.HostBlocks)-1]).Type().Field(i).Tag.Get("ssh_config") {
					continue
				}
				if !value.Field(i).CanSet() {
					return nil, fmt.Errorf("unsettable struct field %v", value.Field(i))
				}
				value.Field(i).SetString(strings.TrimSpace(arr[1]))
			}
		}
	}
	return config, nil
}
