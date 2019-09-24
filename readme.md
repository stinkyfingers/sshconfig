## Read
Read an ssh config into a `Config` object:

```
file, err := os.Open("~/.ssh/config")
if err != nil {
	return err
}
config, err := Read(file)
if err != nil {
	return err
}
fmt.Println(config)
```

## Write
Write a `Config` object to an ssh config file:

```
config := &Config{
	HostBlocks: []HostBlock{
		{
			Host:           "bastion",
			Hostname:       "1.1.1.1",
			User:           "bob-johnson",
			IdentitiesOnly: "yes",
			IdentityFile:   "~/.ssh/identity",
			ProxyCommand:   "ssh prebastion 'nc %h %p'",
		},
	},
}
file, err := os.Create("myNewSSHConfig")
if err != nil {
	return err
}
defer file.Close()

err = config.Write(file)
if err != nil {
	return err
}
```
