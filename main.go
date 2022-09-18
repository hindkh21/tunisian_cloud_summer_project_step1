package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	//"log"
	//"strconv"
	"encoding/csv"
	"encoding/json"

	"github.com/gocolly/colly"
)

type Entreprise_Result struct {
	Type       string
	Company    string
	Coordinate int
	Link       string
}

type Domaine_Et_URL_Entreprise struct {
	Domaine_Full_URL string `json:"domaine_url"`
	Domaine_Nom      string `json:"domaine_nom"`
}
type Pagination struct {
	Page_URL     string `json:"page_url"`
	Pages_Nombre string `json:"pages_nombre"`
}
type Entreprise struct {
	Source  string `json:"source"`
	Nom     string `json:"nom"`
	Adresse string `json:"adresse"`
	Tel     string `json:"tel"`
	Fax     string `json:"fax"`
	Domaine string `json:"domaine"`
	Site    string `json:"site"`
	Email   string `json:"email"`
}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
func convert_array_from_string_to_int(t []string, t2 []int) []int {

	for _, i := range t {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		t2 = append(t2, j)
	}
	return t2

}
func findMax(a []int) (max int) {
	max = a[0]
	for _, value := range a {

		if value > max {
			max = value
		}
	}
	return max
}
func convertJSONToCSV(test, destination string) error {
	// 2. Read the JSON file into the struct array
	testFile, err := os.Open(test)
	if err != nil {
		return err
	}

	// remember to close the file at the end of the function
	defer testFile.Close()

	var Entreprises []Entreprise
	if err := json.NewDecoder(testFile).Decode(&Entreprises); err != nil {
		return err
	}

	// 3. Create a new file to store CSV data
	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// 4. Write the header of the CSV file and the successive rows by iterating through the JSON struct array
	writer := csv.NewWriter(outputFile)
	writer.Comma = ';'
	defer writer.Flush()

	header := []string{"Source", "Nom", "Adresse", "Tel", "Fax", "Domaine", "Site", "Email"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, r := range Entreprises {
		var csvRow []string
		csvRow = append(csvRow, r.Source, r.Nom, r.Adresse, r.Tel, r.Fax, r.Domaine, r.Site, r.Email)
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}
func main() {

	var premiere_moitie_url = "https://www.malipages.com/annuaire/banque-et-finance/page/"
	var Scraping_URL_Source = "https://www.malipages.com/annuaire/banque-et-finance/"

	fmt.Println()
	fmt.Println("===============================================================")
	fmt.Printf("Starting Scraping URL: %s", Scraping_URL_Source)
	fmt.Println()
	fmt.Println("===============================================================")

	var DData1 Pagination

	liste_pagination_string := make([]string, 0)
	liste_pagination_int := make([]int, 0)
	liste_pagination := make([]Pagination, 0)

	collector_pagination := colly.NewCollector()
	collector_pagination.OnHTML("a.page-numbers", func(element *colly.HTMLElement) {

		if isInt(strings.TrimSpace(element.Text)) == true {

			Page_Number := element.Text
			liste_pagination_string = append(liste_pagination_string, Page_Number)

		}
	})

	collector_pagination.Visit(Scraping_URL_Source)

	liste_pagination_int = convert_array_from_string_to_int(liste_pagination_string, liste_pagination_int)
	nbre_pages := 10
	fmt.Printf("le nombre de pages est %d", nbre_pages)
	fmt.Println("")
	for i := 1; i < nbre_pages+1; i++ {
		nbr := strconv.Itoa(i)
		DData1.Page_URL = premiere_moitie_url + nbr + "/"
		DData1.Pages_Nombre = nbr
		fmt.Printf("Page_Number: %s Page_Full_URL:%s ", DData1.Pages_Nombre, DData1.Page_URL)
		fmt.Println("")
		liste_pagination = append(liste_pagination, DData1)
	}
	fmt.Printf("Page Number : %d ", nbre_pages)
	fmt.Println("=====================================================")
	fmt.Printf("Number of page to scrap = %s ", liste_pagination[0].Pages_Nombre)
	fmt.Printf("URL Page To Scrap = %s", liste_pagination[0].Page_URL)
	fmt.Println("=====================================================")

	Liste_entreprises := make([]Entreprise, 0)
	c := colly.NewCollector()

	c.OnError(func(r *colly.Response, err error) {
		log.Fatal("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	i := 0

	for k := 0; k < nbre_pages; k++ {

		var Entreprise_Data Entreprise
		fmt.Printf("Page_Number: '%s' Page_Full_URL: %s ", liste_pagination[k].Pages_Nombre, liste_pagination[k].Page_URL)
		fmt.Println("")
		c.OnHTML("div.article_inner", func(element *colly.HTMLElement) {
			fmt.Println("------------------------------------------------")
			Entreprise_Nom := element.ChildText("h1.entry-title")
			Entreprise_Domaine := element.ChildText("a")
			Entreprise_Tel := element.ChildText(".meta-numero_telephone")
			Entreprise_Fax := element.ChildText(".meta-numero_fax")
			Entreprise_Addr1 := element.ChildText(".meta-adresse_rue")
			Entreprise_Addr2 := element.ChildText(".meta-adresse_ville")
			Entreprise_Addr := Entreprise_Addr1 + " " + Entreprise_Addr2
			Entreprise_Data.Adresse = Entreprise_Addr
			Entreprise_Data.Domaine = Entreprise_Domaine
			Entreprise_Data.Email = ""
			Entreprise_Data.Fax = Entreprise_Fax
			Entreprise_Data.Nom = Entreprise_Nom
			Entreprise_Data.Site = ""
			Entreprise_Data.Source = Scraping_URL_Source
			Entreprise_Data.Tel = Entreprise_Tel
			Liste_entreprises = append(Liste_entreprises, Entreprise_Data)
			i++
			fmt.Println()
			fmt.Println("____________________________________________________________")
			fmt.Printf("Entreprise numero : %d  Domaine: %s", i, Entreprise_Data.Domaine)
			fmt.Println()
			fmt.Printf("Entreprise numero : %d  Nom: %s", i, Entreprise_Data.Nom)
			fmt.Println()
			fmt.Printf("Entreprise numero : %d  Adresse: %s", i, Entreprise_Data.Adresse)
			fmt.Println()
			fmt.Printf("Entreprise numero : %d  Tel: %s", i, Entreprise_Data.Tel)
			fmt.Println()
			fmt.Printf("Entreprise numero : %d  Fax: %s", i, Entreprise_Data.Fax)
			fmt.Println()
			fmt.Printf("Entreprise numero : %d  Web Site: %s", i, Entreprise_Data.Site)
			fmt.Println()
		})
		c.Visit(liste_pagination[k].Page_URL)

	} //end k loop

	file, _ := json.MarshalIndent(Liste_entreprises, "", " ")
	_ = ioutil.WriteFile("test.json", file, 0644)

	if err := convertJSONToCSV("test.json", "data.csv"); err != nil {
		log.Fatal(err)
	}
}
