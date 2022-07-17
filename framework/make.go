package framework

import (
	"html/template"
	"log"
	"os"
	"runtime"
)

const templ = `
# := the variable on the right side of := is evaluated once at assignment time
# = the variable on the right side 0f = is evaluated each time it is used

# VPS (Server) variables
Username={{.User}}
IP={{.IP}}
PSW={{.Password}}

# Go Variables
PWD="$PWD"
App=$(shell basename $(PWD))
GOARCH=GOARCH=amd64

# Multiline Variables
define SERVICE
[Unit]
Description=$(App) service
After=network.target

[Service]
Type=simple
Restart=always
User=$(Username)
Group=$(Username)
WorkingDirectory=/home/$(Username)/$(App)
ExecStart=/home/$(Username)/$(App)/$(App)

# make sure log directory exists and owned by syslog
PermissionsStartOnly=true
ExecStartPre=/bin/mkdir -p /var/log/sleepservice
ExecStartPre=/bin/chown syslog:adm /var/log/sleepservice
ExecStartPre=/bin/chmod 755 /var/log/sleepservice
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=sleepservice

[Install]
WantedBy=multi-user.target
endef
export SERVICE

# Commands

info:
	@echo "The following commands can be used:"
	@echo "-deploy"
	@echo "-delete"
	@echo "-stop"
	@echo "-start"
	@echo "-status"

hello:
	@echo "🐒 Project will be deployed on: $(Username)@$(IP)"

build:
	@echo "🐒 Building executable: $(App)..."
	@go build -o $(App) main.go

transfer:
	@echo "🐝 Transfer project to $(Username)@$(IP):/home/$(Username)/$(App)"
	@rsync -ar . $(Username)@$(IP):/home/$(Username)/$(App)

vps:
	@echo "🐒 Making systemd/$(App).service"
	@mkdir -p systemd
	@touch systemd/$(App).service
	@echo "$$SERVICE" > systemd/$(App).service
	@echo "🐝 Transfer systemd/$(App).service to $(Username)@$(IP):~"
	@rsync -P systemd/$(App).service $(Username)@$(IP):~
	@ rm -fr systemd/
	@echo "😰 Move $(App).service, set cap and enable $(App).service"
	@ssh -t $(Username)@$(IP) '\
		echo $(PSW) | \
		sudo -S mv ~/$(App).service /etc/systemd/system/$(App).service \
		&& sudo -S setcap CAP_NET_BIND_SERVICE=+eip /home/$(Username)/$(App)/$(App) \
		&& sudo -S systemctl enable $(App) \
		&& sudo -S systemctl restart $(App) \
	'
	@echo "😊 Done"

delete:
	@echo "😰 Removing $(App) and $(App).service"
	@ssh -t $(Username)@$(IP) '\
		echo $(PSW) | \
		sudo -S systemctl stop $(App) \
		&& sudo -S systemctl disable $(App) \
		&& sudo -S rm /etc/systemd/system/$(App).service \
		&& sudo -S systemctl daemon-reload \
		&& sudo -S systemctl reset-failed \
		&& sudo -S rm -r $(App) \
	'
	@echo "😊 Done"

stop:
	@echo "😰 Stopping $(App).service"
	@ssh -t $(Username)@$(IP) 'echo $(PSW) | sudo -S systemctl stop $(App)'
	@echo "😊 Done"

start:
	@echo "😰 Starting $(App).service"
	@ssh -t $(Username)@$(IP) 'echo $(PSW) | sudo -S systemctl start $(App)'
	@echo "😊 Done"

status:
	@echo "😰 Status of $(App).service"
	@ssh -t $(Username)@$(IP) 'echo $(PSW) | sudo -S systemctl status $(App)'
	@echo "😊 Done"

deploy: hello build transfer vps

.PHONY: info compile hello build transfer vps
`

func CreateMakeFile() error {
	t := template.Must(template.New("Makefile").Parse(templ))

	if _, err := os.Stat("Makefile"); !os.IsNotExist(err) {
		err := os.RemoveAll("Makefile")
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile("Makefile", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := t.Execute(file, Server.Prod); err != nil {
		return err
	}
	return err
}

func LogPrintErr(err error) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, filename, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)
	}
}

func LogFatalErr(err error) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, filename, line, _ := runtime.Caller(1)
		log.Fatalf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)
	}
}
