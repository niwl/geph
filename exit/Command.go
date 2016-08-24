package exit

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/ProjectNiwl/natrium"
	"github.com/google/subcommands"
	"golang.org/x/net/context"
)

// Command is the exit subcommand.
type Command struct {
	idSeed string

	identity natrium.EdDSAPrivate
	edb      *entryDB
}

// Name returns the name "exit".
func (*Command) Name() string { return "exit" }

// Synopsis returns a description of the subcommand.
func (*Command) Synopsis() string { return "Run as an exit" }

// Usage returns a string describing usage.
func (*Command) Usage() string { return "" }

// SetFlags sets the flag on the binder subcommand.
func (cmd *Command) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.idSeed, "idSeed", "", "seed to use to generate private key")
}

// Execute executes the exit subcommand.
func (cmd *Command) Execute(_ context.Context,
	f *flag.FlagSet,
	args ...interface{}) subcommands.ExitStatus {
	// validate
	if cmd.idSeed == "" {
		panic("idSeed must be given")
	}
	// generate the real stuff from the flags
	cmd.identity = natrium.EdDSADeriveKey([]byte(cmd.idSeed))
	b64, _ := json.Marshal(cmd.identity.PublicKey())
	log.Println("Exit started; public key is", string(b64))
	cmd.edb = newEntryDB()
	// run the stuff
	go cmd.doProxy()
	http.HandleFunc("/update-node", cmd.handUpdateNode)
	http.HandleFunc("/get-nodes", cmd.handGetNodes)
	http.ListenAndServe(":8080", nil)
	return 0
}