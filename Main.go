package main

import (
	"fmt"
	_ "golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type siteRoad struct{
	status bool
	url string
}

func main() {

	start := time.Now()

	readyFile := readingFile()
	continueProcessingFile(readyFile)
	fmt.Println("Program end")
	time.Sleep(time.Second * 2)
	elapsed := time.Since(start)
	fmt.Printf("page took %s", elapsed)
}

func readingFile() []string{

	file, err := os.Open("Site.txt")
	_check(err)
	defer file.Close()

	// получить размер файла
	stat, err := file.Stat()
	_check(err)
	// чтение файла
	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	_check(err)

	str := string(bs)

	readFile := strings.Split(str, "\r\n")
	return readFile
}

func continueProcessingFile(readyFile []string){

	for i :=0; i<len(readyFile);i++ {

		checkHTTP, err := regexp.MatchString(`^https?.*`, readyFile[i])
		checkURL, err := regexp.MatchString(`[^в]контакты?|\bcontacts?|joindre|\bkontakte?|contactos?|contacta|contacter|kontakty`, strings.ToLower(readyFile[i]))
		_check(err)
		if checkURL{
			fmt.Println(siteRoad{status: true, url: readyFile[i]})
			continue
		}

		if checkHTTP{
			getHTML := getHTML(readyFile[i])
			if getHTML ==""{
				continue
			}


			//a := regexp.MustCompile(`<a\s[^>]*href=\"([^\"]*)\"(.*?)?\>(контакты?|joindre|kontakte?(.*?)?|contactos?|contacta?|contacter|kontakty|contacts?(.)?(us)?(.{1,2})?)<\/a>`)
			//url := a.FindString(strings.ToLower(getHTML))
			//b := regexp.MustCompile(`(.*)((<a\s[^>]*href=\"([^\"]*)\"(.*?)?\>(контакты?|joindre|kontakte?(.*?)?|contactos?|contacta?|contacter|kontakty|contacts?(.)?(us)?(.{1,2})?)(.*)(\<\/a\>)))`)
			//urlF := b.FindStringSubmatch(url)

			regexpString := regexp.MustCompile(`(.*)((<a\s[^>]*href=\"([^\"]*)\"(.*?)?\>(контакт(ы)?|joindre|kontakt(e)?|contacto(s)?|contact(a)?|contacter|kontakty|contacts?(\s)?(us)?)((\s|\n){1,2})?(\<\/a\>)))`)
			url := regexpString.FindStringSubmatch(strings.ToLower(getHTML))

			//url != "" &&

			if url !=nil {
				//falseWork := falseWork(url[4])
				//if falseWork==false{
					link := condition(url[4], readyFile[i])
					fmt.Println(siteRoad{status: true, url: link})
				//}else{
				//	fmt.Println(siteRoad{status: false, url: readyFile[i]})
				//}
			}else{
				fmt.Println(siteRoad{status: false, url: readyFile[i]})
			}

		}
	}

}

func condition(link string, url1 string) string{

	reg, err := regexp.MatchString(`https?|/?www\.`, strings.ToLower(link))
	_check(err)

	if reg ==true{
		return link
	} else {
		re := regexp.MustCompile(".*://|/.*")
		cleanLink := re.ReplaceAllString(url1, "")
		link = "http://"+cleanLink + "/"+link
		return  link
	}

}

func falseWork(URL string) bool {
	checkURLFalse, err := regexp.MatchString(`contactboissiere|контактная|контактный|линз|пар|педал|linsen|len|grill|börse|reiniger|câble|
			thermom|blut|change|pay|board|zahlung|cuota|pago|pai|wikipedia|vikidia|googl|(<\/a>) `, strings.ToLower(URL))
	_check(err)
	return checkURLFalse
}

func getHTML(URL string) string{

	var html_page string

	receivedURL, err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer receivedURL.Body.Close()

	if receivedURL.StatusCode != 200 {
		fmt.Println("Get ", URL, "status code error: ", receivedURL.StatusCode, receivedURL.Status)
		return ""
	}else {
		buffer := make([]byte, 1014)
		for true {

			length, err := receivedURL.Body.Read(buffer)
			html_page += string(buffer[:length])
			if length == 0 || err != nil {
				break
			}
		}

		codding, err := regexp.MatchString(`meta(.*?)utf-8`, strings.ToLower(html_page))
		_check(err)
		if codding == false{
			return decodingHTML(html_page)
		}
		return html_page
	}

}

func decodingHTML(html_page string)  string{
	dec := charmap.Windows1251.NewDecoder()
	newBody := make([]byte, len(html_page)*2)
	n, _, err := dec.Transform(newBody, []byte(html_page), false)
	if err != nil {
		panic(err)
	}
	newBody = newBody[:n]
	htmlDecoding := string(newBody)
	return htmlDecoding
}

func _check(err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
}

//<a\s+([^>]*)href="(.*?)"(.*?)>Контакты?|kontakt<\/a>