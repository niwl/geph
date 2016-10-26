package client

import (
	"encoding/json"
	"log"
)

// smFindEntry is the FindEntry state.
// => QueryBinder if cache is stale
// => ConnEntry if cache is fresh
func (cmd *Command) smFindEntry() {
	log.Println("** => FindEntry **")
	defer log.Println("** <= FindEntry **")
	// attempt to read cache entries from disk
	if cmd.cdb != nil {
		var bts []byte
		err := cmd.cdb.QueryRow("SELECT v FROM main WHERE k='bst.entries'").Scan(&bts)
		if err == nil {
			json.Unmarshal(bts, &cmd.entryCache)
		}
	}

	if cmd.entryCache == nil {
		cmd.smState = cmd.smQueryBinder
	} else {
		// save cache entries for next time
		if cmd.cdb != nil {
			bts, err := json.Marshal(&cmd.entryCache)
			if err != nil {
				panic(err.Error())
			}
			cmd.cdb.Exec("INSERT OR REPLACE INTO main VALUES('bst.entries', $1)", bts)
		}
		cmd.smState = cmd.smConnEntry
	}
}