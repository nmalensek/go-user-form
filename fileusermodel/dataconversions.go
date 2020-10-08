package fileusermodel

import (
	"encoding/json"
	"sort"

	"github.com/nmalensek/go-user-form/model"
)

//JSONToUsers takes a JSON string of users, puts them in a map
//keyed on ID, then sorts by ID. Returns a slice and a map of users.
func JSONToUsers(sourceBytes []byte) ([]model.User, map[int]model.User, error) {
	userMap, err := JSONToUserMap(sourceBytes)
	if err != nil {
		return nil, nil, err
	}

	userSlice := make([]model.User, len(userMap))
	i := 0
	for _, val := range userMap {
		userSlice[i] = val
		i++
	}

	sort.Slice(userSlice, func(j, k int) bool {
		return userSlice[j].ID < userSlice[k].ID
	})

	return userSlice, userMap, nil
}

//JSONToUserMap takes a JSON string of users, allocates a new map, and populates the map with the JSON string contents.
func JSONToUserMap(sourceBytes []byte) (map[int]model.User, error) {
	var userMap map[int]model.User
	err := json.Unmarshal(sourceBytes, &userMap)
	if err != nil {
		return nil, err
	}

	return userMap, nil
}
