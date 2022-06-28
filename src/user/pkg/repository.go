package pkg

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

var users Users
var usersById map[string]int
var usersByUsername map[string]int
var usersByIdentityId map[string]int
var usersBySegment map[string][]int
var usersClaimedByIdentityId map[int]bool

var passwords Passwords

// Init
func InitUser() {
	loadedUsers, err := loadUsers("data/users.json.gz")
	if err != nil {
		log.Panic().Err(err).Msg("Unable to load users file")
	}
	users = loadedUsers
}

func loadUsers(filename string) (Users, error) {

	log.Printf("Attempting to load users file: %s", filename)

	var r Users
	usersById = make(map[string]int)
	usersByUsername = make(map[string]int)
	usersByIdentityId = make(map[string]int)
	usersBySegment = make(map[string][]int)
	usersClaimedByIdentityId = make(map[int]bool)

	passwords = make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		return r, err
	}

	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return r, err
	}

	defer gz.Close()

	dec := json.NewDecoder(gz)

	err = dec.Decode(&r)
	if err != nil {
		return r, err
	}

	// Load maps with user array index
	for i, u := range r {
		usersById[u.ID] = i
		usersByUsername[u.Username] = i
		usersBySegment[u.Segment] = append(usersBySegment[u.Segment], i)
	}

	log.Printf("Users successfully loaded into memory structures")

	return r, nil
}

// containsInt returns a bool indicating whether the given []int contained the given int
func containsInt(slice []int, value int) bool {
	for _, v := range slice {
		if value == v {
			return true
		}
	}
	return false
}

// RepoFindUsersIdBySegment Function
func RepoFindUsersIdBySegment(segment string) []int {
	return usersBySegment[segment]
}

// RepoFindRandomUsersBySegment Function
func RepoFindRandomUsersBySegment(Segment string, count int) Users {
	var unclaimedUsers Users
	var segmentFilteredUserIds = RepoFindUsersIdBySegment(Segment)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(segmentFilteredUserIds), func(i, j int) {
		segmentFilteredUserIds[i], segmentFilteredUserIds[j] = segmentFilteredUserIds[j], segmentFilteredUserIds[i]
	})
	for _, idx := range segmentFilteredUserIds {
		if len(unclaimedUsers) >= count {
			break
		}
		if containsInt(segmentFilteredUserIds, idx) && !(usersClaimedByIdentityId[idx]) {
			if users[idx].SelectableUser {
				log.Printf("User found matching filter criteria: %d", idx)
				unclaimedUsers = append(unclaimedUsers, users[idx])
			}
		}
	}
	return unclaimedUsers
}

// RepoClaimUser Function
// Function used to map which shopper user ids have been claimed by the user Id.
func RepoClaimUser(userId int) bool {
	log.Printf("An identity has claimed the user id:%s", userId)
	usersClaimedByIdentityId[userId] = true
	return true
}

func RepoFindRandomUser(count int) Users {
	rand.Seed(time.Now().UnixNano())
	var randomUserId int
	var randomUsers Users
	if len(users) > 0 {
		for len(randomUsers) < count {
			randomUserId = rand.Intn(len(users))
			log.Printf("Random number Selected:%d", randomUserId)
			if randomUserId != 0 {
				if !(usersClaimedByIdentityId[randomUserId]) {
					if users[randomUserId].SelectableUser {
						log.Printf("Random user id selected:%d", randomUserId)
						randomUsers = append(randomUsers, RepoFindUserByID(strconv.Itoa(randomUserId)))
						log.Printf("Random users :%v", randomUsers)
					}
				}
			}
		}
	}
	return randomUsers
}

// RepoFindUserByID Function
func RepoFindUserByID(id string) User {
	if idx, ok := usersById[id]; ok {
		return users[idx]
	} else {
		return User{}
	}
}

// RepoFindUserByUsername Function
func RepoFindUserByUsername(username string) User {
	if idx, ok := usersByUsername[username]; ok {
		return users[idx]
	} else {
		return User{}
	}
}

// RepoFindUserByIdentityID Function
func RepoFindUserByIdentityID(identityID string) User {
	if idx, ok := usersByIdentityId[identityID]; ok {
		return users[idx]
	} else {
		return User{}
	}
}

// RepoUpdateUser Function
func RepoUpdateUser(t User) User {
	if idx, ok := usersById[t.ID]; ok {
		u := &users[idx]
		u.FirstName = t.FirstName
		u.LastName = t.LastName
		u.Email = t.Email
		u.SignUpDate = t.SignUpDate
		u.LastSignInDate = t.LastSignInDate
		u.PhoneNumber = t.PhoneNumber

		if len(u.IdentityId) > 0 && u.IdentityId != t.IdentityId {
			delete(usersByIdentityId, u.IdentityId)
		}

		u.IdentityId = t.IdentityId

		if len(t.IdentityId) > 0 {
			usersByIdentityId[t.IdentityId] = idx
		}

		return RepoFindUserByID(t.ID)
	}

	// return empty User if not found
	return User{}
}

// RepoCreateUser Function
func RepoCreateUser(t User) (User, error) {
	if _, ok := usersByUsername[t.Username]; ok {
		return User{}, errors.New("이미 존재하는 아이디입니다.")
	}

	idx := len(users)

	if len(t.ID) > 0 {
		// ID provided by caller (provisionally created on storefront) so make
		// sure it's not already taken.
		if _, ok := usersById[t.ID]; ok {
			return User{}, errors.New("이미 존재하는 고유 아이디입니다.")
		}
	} else {
		t.ID = strconv.Itoa(idx)
	}

	users = append(users, t)
	usersById[t.ID] = idx
	usersByUsername[t.Username] = idx
	if len(t.IdentityId) > 0 {
		usersByIdentityId[t.IdentityId] = idx
	}

	return t, nil
}
