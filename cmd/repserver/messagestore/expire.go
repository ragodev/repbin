package messagestore

import (
	log "github.com/repbin/repbin/deferconsole"
	"github.com/repbin/repbin/utils"
)

// ExpireFromIndex reads the expire index and expires messages as they are recorded
func (store Store) ExpireFromIndex() {
	// ExpireRun
	delMessages, err := store.db.SelectMessageExpire(CurrentTime())
	if err != nil {
		log.Errorf("ExpireFromIndex, SelectMessageExpire: %s\n", err)
		return
	}
	for _, msg := range delMessages {
		err := store.db.DeleteBlob(&msg.MessageID)
		if err != nil {
			log.Errorf("ExpireFromIndex, DeleteBlob: %s %s\n", err, utils.B58encode(msg.MessageID[:]))
			continue
		}
		err = store.db.DeleteMessageByID(&msg.MessageID)
		if err != nil {
			log.Errorf("ExpireFromIndex, DeleteMessageByID: %s %s\n", err, utils.B58encode(msg.MessageID[:]))
		}
	}
	_, _, err = store.db.ExpireSigners(MaxAgeSigners)
	if err != nil {
		log.Errorf("ExpireFromIndex, ExpireSigners: %s\n", err)
	}
	err = store.db.ExpireMessageCounter(MaxAgeRecipients)
	if err != nil {
		log.Errorf("ExpireFromIndex, ExpireMessageCounter: %s\n", err)
	}
	err = store.db.ForgetMessages(CurrentTime() - MaxAgeRecipients)
	if err != nil {
		log.Errorf("ExpireFromIndex, ForgetMessages: %s\n", err)
	}
}
