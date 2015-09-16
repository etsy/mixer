package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"text/template"

	. "github.com/etsy/mixer/config"
	. "github.com/etsy/mixer/db"

	"github.com/etsy/mixer/debug"
	"github.com/etsy/mixer/mail"
)

var dbg debug.Debug

const debugon = false
const insert_matches = true
const email_matches = true
const override_email = false // used to test the email

var current_week int = 0
var avail_pool []Person
var num_participants int = 0
var odd_person_out Person
var pairs map[Person]Person
var mailer *mail.EmailUser = nil

func main() {
	dbg = debugon
	pairs = make(map[Person]Person)
	dbg.Printf("matching...\n")

	// get flags
	groupname := flag.String("group", "Managers", "a string representing the group name")
	flag.Parse()

	// get the current week
	current_week = GetLastWeek(*groupname) + 1
	fmt.Printf("week : %d\n", current_week)

	participants := GetRandomPeopleParticipating(*groupname)
	avail_pool = make([]Person, len(participants))
	copy(avail_pool, participants)
	num_participants = len(participants)
	fmt.Printf("# people: %d\n", num_participants)
	if num_participants < 2 {
		fmt.Printf("Sorry, we need at least two participants to mix it up\n")
		os.Exit(0)
	}

	for _, v := range participants {
		dbg.Printf("\ncurrent person: %#v\n", v)
		// check if they are already chosen
		if isPaired(v) {
			dbg.Printf("  already paired\n")
			continue
		} else if odd_person_out.Id == 0 && num_participants%2 == 1 {
			dbg.Printf("  odd person\n")
			odd_person_out = v
			removeFromPool(v)
			continue
		}

		matchPair(v)
	}

	// loop through last of avail pool if pair exists
	if len(avail_pool) > 1 {
		dbg.Printf("!!cleaning out avail pool%v\n", avail_pool)
		for _, v := range avail_pool {
			matchPair(v)
		}
	}

	if insert_matches {
		w := InsertWeek(current_week, *groupname)

		for k, v := range pairs {
			fmt.Printf("inserting pair %v : %v\n", k, v)
			InsertPair(w, k, v)
			if email_matches {
				emailPair(k, v, *groupname)
			}
		}
		if odd_person_out.Id != 0 {
			fmt.Printf("inserting odd person %v\n", odd_person_out)
			InsertOddPerson(w, odd_person_out)
			if email_matches {
				emailOddPersonOut(odd_person_out, *groupname)
			}
		}
	}
}

func emailPair(p1 Person, p2 Person, groupname string) {

	mdir := Config.GetRootDir()
	log.Println(mdir)
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/mail/match_email.tmpl", mdir))
	if err != nil {
		log.Fatal(err)
	}

	emailData := struct {
		DirectoryURL string
		Person1LDAP  string
		Person1Name  string
		Person2LDAP  string
		Person2Name  string
		MixerName    string
		ServerUrl    string
	}{
		Config.Staff.DirectoryUrl,
		p1.Username,
		p1.Name,
		p2.Username,
		p2.Name,
		groupname,
		Config.Server.Url,
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, emailData)
	if err != nil {
		log.Fatal(err)
	}

	subject := fmt.Sprintf("%s Mixer", groupname)
	// send to both people
	emailMsg(buffer.Bytes(), subject, []Person{p1, p2})
}

func emailMsg(msg []byte, subject string, p []Person) {

	var to []string
	for _, v := range p {
		email := fmt.Sprintf("%s@%s", v.Username, Config.Mail.Domain)
		fmt.Printf("sending mail to %s\n", email)
		if override_email {
			// override to field for testing
			to = append(to, fmt.Sprintf("%s+%s@%s", Config.Mail.AdminUsername, v.Username, Config.Mail.Domain))
		} else {
			to = append(to, email)
		}

		// check if there is an assistant
		if len(v.Assistant) > 0 {
			assistant_email := fmt.Sprintf("%s@%s", v.Assistant, Config.Mail.Domain)
			fmt.Printf("  sending assistant mail to %s\n", assistant_email)

			if override_email {
				// override recipient field for testing
				to = append(to, fmt.Sprintf("%s+%s@%s", Config.Mail.AdminUsername, v.Assistant, Config.Mail.Domain))
			} else {
				to = append(to, assistant_email)
			}

		}
	}

	// create the mail instance once
	if nil == mailer {
		mailer, _ = mail.NewMail()
	}
	mailer.Mail(msg, subject, to)
}

func emailOddPersonOut(p Person, groupname string) {
	mdir := Config.GetRootDir()
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/mail/match_odd_email.tmpl", mdir))
	if err != nil {
		log.Fatal(err)
	}

	emailData := struct {
		Name      string
		MixerName string
	}{
		p.Name,
		groupname,
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, emailData)
	if err != nil {
		log.Fatal(err)
	}

	subject := fmt.Sprintf("%s Mixer", groupname)
	emailMsg(buffer.Bytes(), subject, []Person{p})
}

func matchPair(p Person) {
	random, err := pickRandom(p)
	if err != nil {
		dbg.Printf("  failed to match, calling swapPairedMatch\n")
		swapPairedMatch(p)
	} else {
		dbg.Printf("  randomly paired person: %#v\n\n", random)
		pairs[p] = *random
		// remove this person from available pool (random removed in pickRandom)
		removeFromPool(p)
	}
}

func removeFromPool(p Person) {
	for k, v := range avail_pool {
		if v == p {
			avail_pool = append(avail_pool[:k], avail_pool[k+1:]...)
			break
		}
	}
}

func swapPairedMatch(p Person) {
	for k, v := range pairs {
		if isValid(p, k, false) {
			// swap out the value
			dbg.Printf("  key match on %#v, taking value: %#v\n\n", k, v)
			pairs[k] = p
			// add the person we pulled to avail and remove current
			avail_pool = append(avail_pool, v)
			removeFromPool(p)
			break
		} else if isValid(p, v, false) {
			// swap out the key
			dbg.Printf("  value match on %#v, taking key: %#v\n\n", v, k)
			delete(pairs, k)
			pairs[p] = v
			// add the person we pulled to avail and remove current
			avail_pool = append(avail_pool, k)
			removeFromPool(p)
			break
		}
	}
}

// isPaired checks the map of pairs to see if they are already paired
func isPaired(p Person) bool {
	for k, v := range pairs {
		if k == p || v == p {
			return true
		}
	}
	return false
}

// pickRandom will pick a random person but make sure they adhere to criteria
func pickRandom(cur_person Person) (*Person, error) {
	list := rand.Perm(len(avail_pool))
	for i := range list {
		var random Person = avail_pool[i]
		if isValid(cur_person, random, true) {
			avail_pool = append(avail_pool[:i], avail_pool[i+1:]...)
			return &random, nil
		}
	}

	return nil, fmt.Errorf("couldn't find an available participant")
}

// isValid means they can't get themself
// or have been paired with this person in the recent past
func isValid(cur_person Person, other_person Person, check_pairing bool) bool {
	if other_person.Id == cur_person.Id {
		dbg.Printf("  same person %v\n", other_person)
		return false
	}

	if pairedTooRecently(cur_person, other_person) {
		dbg.Printf("  too recent %v\n", other_person)
		return false
	}

	return true
}

// pairedTooRecently is false if they have never been paired together,
// or true if difference in weeks since last meet is greater than the
// total number of people in the pool
func pairedTooRecently(cur_person Person, other_person Person) bool {
	pair := GetLastPairing(cur_person.Id, other_person.Id)
	if pair.Id == 0 {
		return false
	} else {
		if (current_week - pair.Week) <= (num_participants - 1) {
			dbg.Printf("  vv paired recently in the past on week %d\n", pair.Week)
			return true
		}
	}

	return false
}
