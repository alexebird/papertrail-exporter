package papertrail

import (
	"encoding/json"
	"time"
	//"github.com/davecgh/go-spew/spew"
)

type Group struct {
	Name    string
	Systems []System
}

type System struct {
	Name           string
	LastEventAtRaw string `json:"last_event_at"`
	LastEventAt    time.Time
}

func ListGroups() []Group {
	body, err := getJson("/groups.json")

	if err != nil {
		panic(err)
	}

	var dat []Group

	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}

	for _, group := range dat {
		for i, system := range group.Systems {
			t, err := time.Parse(time.RFC3339, system.LastEventAtRaw)
			if err != nil {
				panic(err)
			}
			system.LastEventAt = t
			group.Systems[i] = system
		}
	}

	return dat
}
