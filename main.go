package main

import (
    // import standard libraries
    "fmt"
    "log"
    "strings"
    "os"
    "strconv"

    // import third party libraries
    "github.com/PuerkitoBio/goquery"
    "github.com/gocarina/gocsv"
)

//Struct of masjid can be adjust from website kemenag
type Masjid struct {
    No string `csv:"masjid_no"`
    Wilayah string `csv:"masjid_wilayah"`
    Kecamatan string `csv:"masjid_kecamatan"`
    Nama string `csv:"masjid_nama"`
    Id string `csv:"masjid_id"`
    Tipologi string `csv:"masjid_tipologi"`
    Alamat string `csv:"masjid_alamat"`
    LuasTanah string `csv:"masjid_luas_tanah"`
    Status string `csv:"masjid_status"`
    LuasBangunan string `csv:"masjid_luas_bangunan"`
    TahunBerdiri string `csv:"masjid_tahun_berdiri"`
    Jamaah string `csv:"masjid_jamaah"`
    Imam string `csv:"masjid_imam"`
    Khatib string `csv:"masjid_khatib"`
    Muazin string `csv:"masjid_muazin"`
    Remaja string `csv:"masjid_remaja"`
    Hp string `csv:"masjid_hp"`
    Keterangan string `csv:"masjid_keterangan"`
}

//Scrape function return array of struct Masjid and string of result for tracing process
func scrape(url string, masjids []*Masjid) ([]*Masjid, string) {
    //scraping a new document
    doc, err := goquery.NewDocument(url)
    if err != nil {
        log.Fatal(err)
    }

    //Looping for rows
    doc.Find("#the-list tr").Each(func(indexRow int, row *goquery.Selection) {
        var masjid [20]string

        //Temporary rule must change dynamix
        if indexRow < 10 {
            //Looping for columns
            row.Find("td").Each(func(indexColumn int, column *goquery.Selection){
                value := column.Text()
                if value == "" { value = "-" }

                //clear whitespace
                value = strings.TrimSpace(value)
                value = strings.Replace(value, "\"", "", -1)
                masjid[indexColumn] = value
            })

            //append array masjid to array of struct Masjid
            masjids = append(masjids, &Masjid{No: masjid[0], 
                                          Wilayah: masjid[1], 
                                          Kecamatan: masjid[2], 
                                          Nama: masjid[3], 
                                          Id: masjid[4], 
                                          Tipologi: masjid[5], 
                                          Alamat: masjid[6], 
                                          LuasTanah: masjid[7], 
                                          Status: masjid[8], 
                                          LuasBangunan: masjid[9], 
                                          TahunBerdiri: masjid[10], 
                                          Jamaah: masjid[11], 
                                          Imam: masjid[12], 
                                          Khatib: masjid[13], 
                                          Muazin: masjid[14], 
                                          Remaja: masjid[15], 
                                          Hp: masjid[16], 
                                          Keterangan: masjid[17]}) 
        }
    })

    //Get string array of struct Masjid
    csvContent, err := gocsv.MarshalString(&masjids)  
    if err != nil {
        panic(err)
    }

    return masjids, csvContent
}

func main() {
    //Create masjid file csv
    masjidFile, err := os.OpenFile("masjids.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
    if err != nil {
        panic(err)
    }
    defer masjidFile.Close()

    if _, err := masjidFile.Seek(0, 0); err != nil {
        panic(err)
    }

    //Initialization variabel
    masjids := []*Masjid{}
    baseUrl := "http://simas.kemenag.go.id/index.php/profil/masjid/page/"
    var csvContent string

    //Looping per page kemenag for scraping
    for page := 0;page < 219840;page += 10 {
        url := baseUrl + strconv.Itoa(page) //append string url + page
        masjids, csvContent = scrape(url, masjids) //append masjid array of struct Masjid
        fmt.Println(csvContent) //print for tracing on console
    }

    //Save file masjids.csv
    err = gocsv.MarshalFile(&masjids, masjidFile)
}