package main

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"time"

	"github.com/coreos/go-systemd/v22/journal"
	"github.com/coreos/go-systemd/v22/login1"
)

var (
	conn   login1.Conn
	sess   login1.Session
	keyID  string
	unlock = false
)

func main() {
	if !journal.Enabled() {
		return
	}
	journal.Print(journal.PriAlert, "Starting pamlockd")
	setUp()
	go timeOutLoop()

	fmt.Printf("%+v\n%+v\n", sess, keyID)

}

func timeOutLoop() {
	//TODO: Implement an actual check for unlock status
	for {
		time.Sleep(30 * time.Second)
		if !unlock {
			journal.Print(journal.PriInfo, "Timeout reached waiting for unlock. Locking session!")
			conn.LockSession(sess.ID)
		}
		time.Sleep(5 * time.Minute)
	}
}

func setUp() {
	journal.Print(journal.PriInfo, "Opening new connection to logind.")
	conn, err := login1.New()
	if err != nil {
		journal.Print(journal.PriCrit, "Unable to open login1 connection: %+v\n", err)
		return
	}
	journal.Print(journal.PriInfo, "Getting current user.")
	usr, err := user.Current()
	if err != nil {
		journal.Print(journal.PriCrit, "Unable to determine current user: %+v\n", err)
		return
	}

	journal.Print(journal.PriInfo, "Getting current session from Active Sessions.")
	dop, err := conn.GetActiveSession()
	sessions, err := conn.ListSessions()
	for _, v := range sessions {
		if v.Path == dop {
			sess = v
		}
	}
	content, err := ioutil.ReadFile(fmt.Sprintf("%+v/.pgplockd", usr.HomeDir))
	if err != nil {
		journal.Print(journal.PriCrit, "Unable to read ~/.pgplockd config: %+v\n", err)
		return
	}
	keyID = string(content)
}
