package conversation

import (
	"notionboy/db/ent"
)

func checkRateLimit(acc *ent.Account, qt *ent.Quota) bool {
	return qt.TokenUsed >= qt.Token
}
