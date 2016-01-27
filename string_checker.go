package validity

import (
	"log"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type StringValidityChecker struct {
	Key   string
	Rules []string
	Item  string
}

func (v StringValidityChecker) GetKey() string {
	return v.Key
}

func (v StringValidityChecker) GetItem() interface{} {
	return v.Item
}

func (v StringValidityChecker) GetRules() []string {
	return v.Rules
}

func (v StringValidityChecker) GetErrors() []string {
	return GetCheckerErrors(v.Rules[1:], &v)
}

func (v StringValidityChecker) toInt(s string) int {
	out, _ := strconv.ParseInt(s, 10, 64)

	return int(out)
}

func (v StringValidityChecker) checkRegexp(r string) bool {
	expression, _ := regexp.Compile(r)

	return expression.MatchString(v.Item)
}

func (v StringValidityChecker) parseIP() net.IP {
	return net.ParseIP(v.Item)
}

//----------------------------------------------------------------------------------------------------------------------
// For explanation involving validation rules, checkout the first huge comment in validity.go.
//----------------------------------------------------------------------------------------------------------------------

func (v StringValidityChecker) ValidateCnp() bool {

	rawCNP := v.Item

	var (
		bigSum    int
		ctrlDigit int
		digits    = []int{}
		year      = 0
		control   = []int{2, 7, 9, 1, 4, 6, 3, 5, 8, 2, 7, 9}
	)

	// iterate

	for i := 0; i < 12; i++ {
		current, errCurrent := strconv.Atoi(string(rawCNP[i]))
		if errCurrent != nil {
			log.Println("The character at position " + strconv.Itoa(i) + "[" + string(rawCNP[i]) + "] is not a digit")
			return false
		}
		bigSum += control[i] * current
		digits[i] = current
	}

	// Sex -  allowed only 1 -> 9

	if digits[0] == 0 {
		log.Println("Sex can not be 0")
		return false
	}

	// year
	year = digits[1]*10 + digits[2]

	switch digits[0] {
	case 1, 2:
		year += 1900
		break
	case 3, 4:
		year += 1800
		break
	case 5, 6:
		year += 2000
		break
		// TODO to check
	case 7, 8, 9:
		year += 2000
		now := time.Now()
		if year > now.Year()-14 {
			year -= 100
		}
		break
	}

	if year < 1800 || year > 2099 {
		log.Println("Wrong year: " + strconv.Itoa(year))
		return false
	}

	// Month - allowed only 1 -> 12
	month := digits[3]*10 + digits[4]
	if month < 1 || month > 12 {
		log.Println("Wrong month: " + strconv.Itoa(month))
		return false
	}

	day := digits[5]*10 + digits[6]

	// check date
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	if int(t.Year()) != year || int(t.Month()) != month || t.Day() != day {
		log.Println("The date does not exist: " + strconv.Itoa(year) + "/" + strconv.Itoa(month) + "/" + strconv.Itoa(day))
	}

	// County - allowed only 1 -> 52

	county := digits[7]*10 + digits[8]
	if county < 1 || county > 52 {
		log.Println("Wrong county: " + strconv.Itoa(county))
		return false
	}

	// Number - allowed only 001 --> 999

	number := digits[9]*100 + digits[10]*10 + digits[11]
	if number < 1 || number > 999 {
		log.Println("Wrong number: " + strconv.Itoa(number))
		return false
	}

	// Check control

	ctrlDigit = bigSum % 11

	if ctrlDigit == 10 {
		ctrlDigit = 1
	}

	return strconv.Itoa(ctrlDigit) == string(rawCNP[12])

}

func (v StringValidityChecker) ValidateAccepted() bool {
	return v.Item == "yes" || v.Item == "on" || v.Item == "1"
}

func (v StringValidityChecker) ValidateAlpha() bool {
	return v.checkRegexp("^[A-Za-z]*$")
}

func (v StringValidityChecker) ValidateAlphaDash() bool {
	return v.checkRegexp("^[A-Za-z0-9\\-_]*$")
}

func (v StringValidityChecker) ValidateAlphaNum() bool {
	return v.checkRegexp("^[A-Za-z0-9]*$")
}

func (v StringValidityChecker) ValidateBetween(min string, max string) bool {
	length := len([]rune((v.Item)))

	return length > v.toInt(min) && length < v.toInt(max)
}

func (v StringValidityChecker) ValidateBetweenInclusive(min string, max string) bool {
	length := len([]rune(v.Item))
	return length >= v.toInt(min) && length <= v.toInt(max)
}

func (v StringValidityChecker) ValidateDate() bool {
	_, err := time.Parse("Jan 2, 2006 at 3:04pm (MST)", v.Item)

	return err == nil
}

func (v StringValidityChecker) ValidateEmail() bool {
	return v.checkRegexp("^.+\\@.+\\..+$")
}

func (v StringValidityChecker) ValidateIpv4() bool {
	parsed := v.parseIP()

	return parsed != nil && parsed.To4() != nil
}

func (v StringValidityChecker) ValidateIpv6() bool {
	parsed := v.parseIP()

	return parsed != nil && parsed.To16() != nil
}

func (v StringValidityChecker) ValidateIp() bool {
	return v.parseIP() != nil
}

func (v StringValidityChecker) ValidateLen(length string) bool {
	return len([]rune(v.Item)) != v.toInt(length)
}

func (v StringValidityChecker) ValidateFullName() bool {
	return v.checkRegexp(`^[A-Za-z0-9\s\.]*$`)
}

func (v StringValidityChecker) ValidateMax(length string) bool {
	return len([]rune(v.Item)) <= v.toInt(length)
}

func (v StringValidityChecker) ValidateMin(length string) bool {
	return len([]rune(v.Item)) >= v.toInt(length)
}

func (v StringValidityChecker) ValidateRegexp(r string) bool {
	return v.checkRegexp(r)
}

func (v StringValidityChecker) ValidateUrl() bool {
	_, err := url.ParseRequestURI(v.Item)

	return err == nil
}
