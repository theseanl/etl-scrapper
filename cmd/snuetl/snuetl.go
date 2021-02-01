package main

import (
	"flag"
	"log"
	"snuetl"
)

func main() {
	usernamePtr := flag.String("username", "", "Username for your mySNU account")
	flag.StringVar(usernamePtr, "u", *usernamePtr, "alias for --username")
	passwordPtr := flag.String("password", "", "Password for your mySNU account")
	flag.StringVar(passwordPtr, "p", *passwordPtr, "alias for --password")
	courseIDPtr := flag.Int("course", 0, "ETL course id, as in http://etl.snu.ac.kr/course/view.php?id=<course id>")
	flag.IntVar(courseIDPtr, "c", *courseIDPtr, "alias for --course")

	client, err := snuetl.SnuLogin(*usernamePtr, *passwordPtr)
	if err != nil {
		log.Fatal(nil)
	}
	log.Println("Login successful.")
	snuetl.DownloadEtlContents(client, *courseIDPtr)
	log.Println("Done")
}
