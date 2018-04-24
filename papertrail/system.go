package papertrail

import (
	"encoding/json"
	"regexp"
	"time"
	//"github.com/davecgh/go-spew/spew"
)

type Group struct {
	Name    string
	Systems []System
}

type System struct {
	Id             int
	Name           string
	GroupName      string
	LastEventAtRaw string `json:"last_event_at"`
	LastEventAt    time.Time
}

func ListSystems() ([]System, error) {
	body, err := getJson("/groups.json")

	if err != nil {
		return nil, err
	}

	var groups []Group
	var systems []System = make([]System, 0)

	if err := json.Unmarshal(body, &groups); err != nil {
		return nil, err
	}

	for _, group := range groups {
		for _, system := range group.Systems {
			t, err := time.Parse(time.RFC3339, system.LastEventAtRaw)
			if err != nil {
				return nil, err
			}

			system.LastEventAt = t
			system.GroupName = group.Name
			systems = append(systems, system)
		}
	}

	return systems, nil
}

func FilterSystems(groupRegexp *regexp.Regexp, sysRegexp *regexp.Regexp) ([]System, error) {
	systems, err := ListSystems()
	if err != nil {
		return nil, err
	}

	filteredSystems := make([]System, 0)

	for _, sys := range systems {
		if groupRegexp.MatchString(sys.GroupName) && sysRegexp.MatchString(sys.Name) {
			filteredSystems = append(filteredSystems, sys)
		}
	}

	return filteredSystems, nil
}
