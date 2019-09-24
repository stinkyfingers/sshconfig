package connect

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteConfig(t *testing.T) {
	expected := `
Host bastion
	Hostname 1.1.1.1
	IdentitiesOnly yes
	IdentityFile ~/.ssh/identity
	ProxyCommand ssh prebastion 'nc %h %p'
	User bob-johnson

Host resource
	Hostname 2.2.2.2
	IdentitiesOnly yes
	IdentityFile ~/.ssh/identity
	ProxyCommand ssh -F ~/.ssh/identity bastion 'nc %h %p'
	User bob-johnson

Include ~/.ssh/config
`
	c := &Config{
		HostBlocks: []HostBlock{
			{
				Host:           "bastion",
				Hostname:       "1.1.1.1",
				User:           "bob-johnson",
				IdentitiesOnly: "yes",
				IdentityFile:   "~/.ssh/identity",
				ProxyCommand:   "ssh prebastion 'nc %h %p'",
			}, {
				Host:           "resource",
				Hostname:       "2.2.2.2",
				User:           "bob-johnson",
				IdentitiesOnly: "yes",
				IdentityFile:   "~/.ssh/identity",
				ProxyCommand:   "ssh -F ~/.ssh/identity bastion 'nc %h %p'",
			},
		},
		Include: "~/.ssh/config",
	}

	temp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(temp.Name())

	err = c.Write(temp)
	if err != nil {
		t.Error(err)
	}
	b, err := ioutil.ReadFile(temp.Name())
	if err != nil {
		t.Error(err)
	}
	if string(b) != expected {
		t.Error("file not as expected; got/expected")
		t.Log(string(b))
		t.Log(expected)
	}
}

func TestWriteBlock(t *testing.T) {
	block := &HostBlock{
		Host:           "bastion",
		Hostname:       "1.1.1.1",
		User:           "bob-johnson",
		IdentitiesOnly: "yes",
		IdentityFile:   filepath.Join("~", ".ssh", "identity"),
		ProxyCommand:   "ssh prebastion 'nc %h %p'",
		Ciphers:        []string{"aes192-cbc", "aes256-cbc"},
	}
	expectedBlock := `
Host bastion
	Hostname 1.1.1.1
	Ciphers aes192-cbc,aes256-cbc
	IdentitiesOnly yes
	IdentityFile ~/.ssh/identity
	ProxyCommand ssh prebastion 'nc %h %p'
	User bob-johnson
`
	b := make([]byte, 1024)
	buf := bytes.NewBuffer(b)
	err := block.write(buf)
	if err != nil {
		t.Error(err)
	}

	if strings.Trim(buf.String(), "\x00") != expectedBlock {
		t.Error("block not as expected")
		t.Log(buf.String())
		t.Log(expectedBlock)
	}

}

func TestReadConfig(t *testing.T) {
	reader := bytes.NewBuffer([]byte(`
Host bastion
	Hostname 1.1.1.1
	User bob-johnson
	IdentitiesOnly yes
	IdentityFile ~/.ssh/identity
	ProxyCommand ssh prebastion 'nc %h %p'

Host resource
	Hostname 2.2.2.2
	User bob-johnson
	IdentitiesOnly yes
	IdentityFile ~/.ssh/identity
	ProxyCommand ssh -F ~/.ssh/identity bastion 'nc %h %p'

Include ~/.ssh/config
`))
	config, err := Read(reader)
	if err != nil {
		t.Error(err)
	}
	t.Log(config)
}
