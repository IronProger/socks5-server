package main

import (
	"bufio"
	"flag"
	"github.com/armon/go-socks5"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {

	var (
		proxyPort      uint
		proxyUser      string
		proxyPassword  string
		proxyUsersFile string
	)

	flag.UintVar(&proxyPort, "port", 1080, "custom port")
	flag.StringVar(&proxyUser, "login", "", "login name, needs if you are too lazy to creat a file which contains user data")
	flag.StringVar(&proxyPassword, "password", "", "password, please read about login name")
	flag.StringVar(&proxyUsersFile, "users", "", "a file which include a list of users, in which has set one user per line, every line include USERNAME:PASSWORD separated by colons")

	flag.Parse()

	//Initialize socks5 config
	socsk5conf := &socks5.Config{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	//creds := socks5.StaticCredentials{
	//	proxyUser: proxyPassword,
	//	"dsa": "dsa",
	//}
	//var creds socks5.CredentialStore

	userData := make(map[string]string)

	if (proxyUser != "") != (proxyPassword != "") {
		log.Fatal("Specify password and login name just together a time")
	}
	if proxyUser != "" && proxyPassword != "" {
		userData[proxyUser] = proxyPassword
	}

	// reads users from file
	if proxyUsersFile != "" {
		file, err := os.Open(proxyUsersFile)
		if err != nil {
			log.Print(err)
			log.Print("I cannot read the file, I ignore it.")
		}

		r := bufio.NewReader(file)
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				break
			}
			data := strings.Split(string(line), ":")
			if len(data) != 2 {
				log.Print("line \"" + string(line) + "\" is invalid")
				break
			}
			correct, err := regexp.MatchString("[A-z0-9_]+", data[0])
			if !correct || err != nil {
				log.Print("error at parsing line: " + string(line))
				break
			}
			correct, err = regexp.MatchString("[A-z0-9_]+", data[1])
			if !correct || err != nil {
				log.Print("error at parsing line: " + string(line))
				break
			}

			// all right
			userData[data[0]] = data[1]
		}
	}

	cator := socks5.UserPassAuthenticator{Credentials: socks5.StaticCredentials(userData)}
	socsk5conf.AuthMethods = []socks5.Authenticator{cator}

	server, err := socks5.New(socsk5conf)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Start listening proxy service on port %d\n", proxyPort)
	if err := server.ListenAndServe("tcp", ":"+strconv.Itoa(int(proxyPort))); err != nil {
		log.Fatal(err)
	}
}
