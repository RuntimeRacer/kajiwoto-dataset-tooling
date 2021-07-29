package query

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

/*
	ToCSV converts a Dataset entry into a String array used for writing it to a CSV file.

	Mapping:
	- ID
	- UserMessage
	- Message
	- ASM
	- Attachment Key
	- Daytime Key
	- Last Seen Key
	- Deleted
	- History, JSON FIXME: When the whole Dataset is known, the entries can be mapped to IDs

*/
func (e *AITrained) ToCSV() []string {
	result := make([]string, 9)

	result[0] = string(e.ID)
	result[1] = string(e.UserMessage)
	result[2] = string(e.Message)
	result[3] = string(e.ASM)

	// Condition Split
	attachmentVars := strings.Split(string(e.Condition), "")
	result[4] = attachmentVars[2]
	result[5] = attachmentVars[0]
	result[6] = attachmentVars[1]

	result[7] = strconv.FormatBool(bool(e.Deleted))

	// History as JSON
	if b, err := json.Marshal(e.History); err != nil {
		fmt.Println(fmt.Sprintf("Error marshalling History for Entry %v", e.ID))
	} else {
		result[8] = string(b)
	}

	return result
}
