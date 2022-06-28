package pkg

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
)

var users Users
var usersById map[string]int
var usersByUsername map[string]int
var usersBySegment map[string][]int

var products Products
var productsBySegment map[string][]int
var productsByGender map[string][]int

var PurchaseRatio map[string]float64
var FluctuationRatio map[int]float64

// Init
func init() {
	loadedUsers, err := loadUsers("data/users.json")
	if err != nil {
		log.Panic().Err(err).Msg("Unable to load users file")
	}
	loadedProducts, err := loadProducts("data/products.json")
	if err != nil {
		log.Panic().Err(err).Msg("Unable to load products file")
	}

	users = loadedUsers
	products = loadedProducts

	PurchaseRatio = make(map[string]float64)
	PurchaseRatio["M10"] = 10
	PurchaseRatio["M20"] = 11
	PurchaseRatio["M30"] = 13
	PurchaseRatio["M40"] = 12
	PurchaseRatio["M50"] = 9
	PurchaseRatio["F10"] = 11
	PurchaseRatio["F20"] = 12
	PurchaseRatio["F30"] = 16
	PurchaseRatio["F40"] = 14
	PurchaseRatio["F50"] = 13

	FluctuationRatio = make(map[int]float64)
	FluctuationRatio[0] = 0.8
	FluctuationRatio[1] = 1
	FluctuationRatio[2] = 1.1
	FluctuationRatio[3] = 1.1
	FluctuationRatio[4] = 1.4
	FluctuationRatio[5] = 1.5
	FluctuationRatio[6] = 1.4
	FluctuationRatio[7] = 1.3
	FluctuationRatio[8] = 1.1
	FluctuationRatio[9] = 1
	FluctuationRatio[10] = 1
	FluctuationRatio[11] = 0.9
	FluctuationRatio[12] = 1
	FluctuationRatio[13] = 1.1
	FluctuationRatio[14] = 1.2
	FluctuationRatio[15] = 1.3
	FluctuationRatio[16] = 1.1
	FluctuationRatio[17] = 0.9
	FluctuationRatio[18] = 0.5
	FluctuationRatio[19] = 0.3
	FluctuationRatio[20] = 0.2
	FluctuationRatio[21] = 0.2
	FluctuationRatio[22] = 0.2
	FluctuationRatio[23] = 0.3
}

func AutoPurchaseOrder(segment string, cooldown time.Duration) {
	for {
		h := time.Now().UTC().Hour()
		users := RepoFindRandomUsersBySegment(segment, int(PurchaseRatio[segment]*FluctuationRatio[h]))
		for _, u := range users {
			var p Product
			rand.Seed(time.Now().UnixNano())
			n := rand.Intn(10)

			if n >= 0 && n < 4 {
				// users are likely to purchase items that has affinity to users's segment. (40%)
				p = getRandomSegmentProducts(segment, 1)[0]
			} else if n >= 4 && n < 7 {
				// users are likely to purchase items that has gender affinity (30%)
				p = getRandomGenderProduct(u.Gender)
			} else {
				// users are likely to purchase random items. (30%)
				p = getRandomProduct()
			}

			log.Info().Str("Personalize", "true").Str("currency", "KRW").Str("email", u.Email).Str("user_firstname", u.FirstName).Str("user_lastname", u.LastName).Str("user_username", u.Username).Str("user_segment", u.Segment).Int("user_age", u.Age).Int("order_id", rand.Intn(100000000)).Str("user_gender", u.Gender).Interface("products", p).Msg("user purchased order.")
		}
		time.Sleep(cooldown)
	}
}

func getRandomGenderProduct(gender string) Product {
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(productsByGender[gender]))

	return products[productsByGender[gender][idx]]
}

// RepoFindProductsBySegment Function
func RepoFindProductsBySegment(segment string) []int {
	return productsBySegment[segment]
}

func getRandomSegmentProducts(segment string, count int) Products {
	var ret Products
	var segmentFilteredProducts = RepoFindProductsBySegment(segment)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(segmentFilteredProducts), func(i, j int) {
		segmentFilteredProducts[i], segmentFilteredProducts[j] = segmentFilteredProducts[j], segmentFilteredProducts[i]
	})
	for _, idx := range segmentFilteredProducts {
		if len(ret) >= 1 {
			break
		}
		if containsInt(segmentFilteredProducts, idx) {
			ret = append(ret, products[idx])
		}
	}
	return ret
}

func getRandomProduct() Product {
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(products))
	return products[idx]
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
		if containsInt(segmentFilteredUserIds, idx) {
			unclaimedUsers = append(unclaimedUsers, users[idx])
		}
	}
	return unclaimedUsers
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

func loadUsers(filename string) (Users, error) {
	start := time.Now()

	log.Printf("Loading users from file: %s", filename)

	var r Users
	usersById = make(map[string]int)
	usersByUsername = make(map[string]int)
	usersBySegment = make(map[string][]int)

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(bytes, &r)
	if err != nil {
		return r, err
	}

	// Load maps with user array index
	for i, u := range r {
		usersById[u.ID] = i
		usersByUsername[u.Username] = i
		usersBySegment[u.Segment] = append(usersBySegment[u.Segment], i)
	}

	log.Printf("Users loaded in %s", time.Since(start))

	return r, nil
}

func loadProducts(filename string) (Products, error) {
	start := time.Now()

	log.Printf("Loading products from file: %s", filename)

	var r Products
	productsBySegment = make(map[string][]int)
	productsByGender = make(map[string][]int)

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(bytes, &r)
	if err != nil {
		return r, err
	}

	// Load maps with user array index
	for i, p := range r {
		if len(p.Affinity) > 0 {
			for _, a := range p.Affinity {
				productsBySegment[a] = append(productsBySegment[a], i)
			}
		}

		if p.GenderAffinity != "" {
			productsByGender[p.GenderAffinity] = append(productsByGender[p.GenderAffinity], i)
		}
	}

	log.Printf("Products loaded in %s", time.Since(start))

	return r, nil
}
