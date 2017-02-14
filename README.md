# Remote Tail
Simple go package to tail log remotely.

### Getting Started
```bash
go get -u github.com/alileza/rtail
```

```bash
usage: rtail [<servers>] [-F | -f | -r] [-q] [-b # | -c # | -n #] <file>

# single server
rtail root@192.168.100.160 -f /var/log/nginx/access.log

# multiple server
rtail root@192.168.100.160,root@192.168.100.161,root@192.168.100.162 -f /var/log/nginx/access.log
```

# Use Config File
```sh
rtail -config.file=example.json
```

**Config Example**
```json
[
    {
        "name" : "nginx",
        "server_addresses" : [
            "root@192.168.100.160",
            "root@192.168.100.161",
            "root@192.168.100.162"
        ],
        "options" : ["-f", "-n 1"],
        "file" : "/var/log/nginx/access.log"
    },
    {
        "name" : "apps",
        "server_addresses" : [
            "root@192.168.100.170",
            "root@192.168.100.171"
        ],
        "options" : ["-f", "-n 1"],
        "file" : "/var/log/apps/apps.log"
    }
]
```
