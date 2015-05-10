# aws-utilities-go
aws utilities written by go

## Instance Operations

start, stop, and show status, your ec2 instances.

### Command Options

`--nametag`, `-n` [*required]

"Name" tag to search instances.
you can use "*"

`--mode`, `-m` [*required]

start or stop or status

### Example

command 

```
$ aws-utilities -p my_main_profile op -n *test -m status
```

outputs

(Nametag PrivateIP PublicIP InstanceID InstanceState)

```
simeji-test 172.30.0.200 54.xxx.xxx.201 i-xxxaaaa1 running
simeji-test2 172.30.0.201  i-xxxbbbb2 stopped
```

## Global Options

`--profile`, `-p`

your profile. default: "default"

