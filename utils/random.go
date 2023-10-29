package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "qwertyuiopasdfghjklzxcvbnm"
const numbers = "0123456789"

type UtilRandom struct {
	random *rand.Rand
}

func NewUtilRandom() *UtilRandom {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &UtilRandom{random: random}
}

// RandomInt generates a random integer between min and max
func (s *UtilRandom) RandomInt(min, max int64) int64 {
	return min + s.random.Int63n(max-min+1)
}

func (s *UtilRandom) RandomIntP(min, max int64) *int64 {
	n := s.RandomInt(min, max)
	return &n
}

// RandomString generates a random string of length n
func (s *UtilRandom) RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[s.random.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func (s *UtilRandom) RandomStringP(n int) *string {
	str := s.RandomString(n)
	return &str
}

func (s *UtilRandom) RandomNumberString(n int) string {
	var sb strings.Builder
	k := len(numbers)

	for i := 0; i < n; i++ {
		c := numbers[s.random.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func (s *UtilRandom) RandomNumberStringP(n int) *string {
	str := s.RandomNumberString(n)
	return &str
}

// RandomOwner generates a random name
func (s *UtilRandom) RandomName() string {
	nameList := []string{"London Kim", "Roman Dominguez", "Raegan Ortiz", "Landon Lawrence", "Lauren Coleman", "Micah Boone", "Mariam Santiago", "Beckham Booth", "Zariyah Reed", "Easton Daniels", "Ember Moore", "Levi Evans", "Eliana Quintana", "Kelvin Smith", "Olivia Powers", "Sean Lynch", "Malia Wilson", "Daniel Massey", "Clementine Ellison", "Kye Foster", "Brielle Doyle", "Kashton Vasquez", "Rose Nolan", "Maximo Grant", "Alaina Newton", "Santino Tran", "Kylie Morris", "Christian Corona", "Marianna Miles", "Jared Stokes", "Miranda Simon", "Zayne Hobbs", "Lacey Alvarado", "Andres Ford", "Alexandra Mann", "Nehemiah Quintana", "Kenia Aguilar", "Milo Singh", "Vivienne Lester", "Lee Whitney", "Madalynn Ware", "Tadeo Elliott", "Noelle Craig", "Odin Hebert", "Kyleigh Kelly", "Cooper Murillo", "Mikaela Miles", "Jared Friedman", "Aspyn Rogers", "Colton Delarosa"}
	nameLen := len(nameList)
	index := s.RandomInt(0, int64(nameLen-1))
	return nameList[index]
}

// RandomMoney generates a random amount of money
func (s *UtilRandom) RandomMoney() int64 {
	return s.RandomInt(0, 1000)
}

// RandomMoney generates a random phone number with prefix 09xx xxx xxx
func (s *UtilRandom) RandomPhone() string {
	var sb strings.Builder
	for i := 0; i < 8; i++ {
		c := s.RandomInt(0, 9)
		sb.WriteString(fmt.Sprint(c))
	}
	return fmt.Sprintf(`+849%s`, sb.String())
}

// RandomEmail generates a random email
func (s *UtilRandom) RandomEmail() string {
	return fmt.Sprintf("%s@email.com", s.RandomString(6))
}

// RandomBirthday generates a random email
func (s *UtilRandom) RandomBirthday() time.Time {
	randInt := s.RandomInt(20, 50)
	expectedYear := int(randInt)
	birTime := time.Now().AddDate(-expectedYear, 0, 0).In(time.UTC)

	return birTime
}

func (s *UtilRandom) RandomTime() time.Time {
	randInt := s.RandomInt(1, 100)
	expectedDate := int(randInt)
	birTime := time.Now().AddDate(0, 0, -expectedDate).In(time.UTC)
	return birTime
}

func (s *UtilRandom) RandomPastTime() time.Time {
	randInt := s.RandomInt(1, 100)
	expectedDate := int(randInt)
	birTime := time.Now().AddDate(0, 0, -expectedDate).In(time.UTC)
	return birTime
}

func (s *UtilRandom) RandomPastTimeP() *time.Time {
	birTime := s.RandomPastTime()
	return &birTime
}

func (s *UtilRandom) RandomFutureTime() time.Time {
	randInt := s.RandomInt(1, 100)
	expectedDate := int(randInt)
	birTime := time.Now().AddDate(0, 0, expectedDate).In(time.UTC)
	return birTime
}

func (s *UtilRandom) RandomFutureTimeP() *time.Time {
	birTime := s.RandomFutureTime()
	return &birTime
}
