# Go Package Manager Toolkit

Manage Go applications in `$GOBIN`.

## 1. Install

```bash
go install github.com/fsgo/gopt@latest
```

## 2. Usage

### 2.1 List Go applications
```
# gopt list -help
Usage of list:
  -t int
    	list timeout, seconds (default 10)
  -dev string
    	filter devel. 'yes': only devel; 'no': no devel; default is '': no filter
  -e	filter only expired
  -json
    	print JSON result
  -l	get latest version info (default true)
```

```
# gopt list
  1 /Users/work/go/bin/bin-auto-switcher
           Path : github.com/fsgo/bin-auto-switcher
             Go : go1.19.3
   Install Time : 2022-11-03 09:52:08
        Version : (devel)
 Latest Version : v0.1.3
    Latest Time : 2022-10-30 09:53:19
    
  2 /Users/work/go/bin/dlv
           Path : github.com/go-delve/delve/cmd/dlv
             Go : go1.19.2
   Install Time : 2022-10-21 13:10:04
        Version : v1.9.1
 Latest Version : v1.9.1
    Latest Time : 2022-08-23 14:35:35   
```

### 2.2 Update Go applications
```
#gopt update -help
Usage of list:
  -t int
    	update timeout, seconds (default 60)
```

Update All Go applications:
```bash
# gopt update
```

Update with given name:
```bash
# gopt update bin-auto-switcher
# gopt update dlv
```


### 2.3 Install Go applications
```
#gopt install -help
Usage of install:
  -t duration
    	install timeout (default 2m0s)
```

Install with given name:
```bash
# gopt install dlv
```