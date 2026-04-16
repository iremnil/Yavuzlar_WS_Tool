package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gocolly/colly/v2"
)

func main() {
	asciiBaslik()
	//eger kullancı fıltreleme yaparsa dıye eklendi
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n" + color.BlueString("== YAVUZLAR WEB SCRAPER MENU =="))
		fmt.Println("[-1] The Hacker News(LATEST 5 NEWS)")
		fmt.Println("[-2] Bleeping Computer (LATEST 5 NEWS)")
		fmt.Println("[-3] KrebsOnSecurity (LATEST 5 NEWS)")
		color.Red("[-4] EXİT")
		fmt.Print("\nOptions: -date (hide date), -description (hide summary)(Ex: -1 -date   Ex: -2 -description): ")
		fmt.Print("\nYour selection:")

		scanner.Scan()
		input :=scanner.Text()
		
		// Girdi parcalara ayrıldı
		parcalar :=strings.Fields(input)
		if len(parcalar) == 0 { continue }

		secim := parcalar[0]
		
		// Parametre kontrolü
		tarihGizle := strings.Contains(input, "-date")
		aciklamaGizle := strings.Contains(input, "-description")

		switch secim {
		case "-1":
			fmt.Println(color.MagentaString("\n[!] Scraping The Hacker News..."))
			hackerNewsScraping(tarihGizle, aciklamaGizle)
			color.Green("\n[+] Success: Data saved to 'scan_results.txt'")
		case "-2":
			fmt.Println(color.MagentaString("\n[!] Scraping Bleeping Computer..."))
			bleepingComputerScraping(tarihGizle, aciklamaGizle)
			color.Green("\n[+] Success: Data saved to 'scan_results.txt'")
		case "-3":
    		fmt.Println(color.MagentaString("\n[!] Scraping KrebsonOnSecurity..."))
    		krebsonSecurityScraping(tarihGizle,aciklamaGizle)
			color.Green("\n[+] Success: Data saved to 'scan_results.txt'")
		case "-4":
			color.Yellow("System exit.Stay safe")
			os.Exit(0)
		default:
			color.Red("Invalid selection! Please use -1, -2, or -3.")
		}
	}
}

func hackerNewsScraping(tGizle bool, aGizle bool) {
	dosya, _ := os.OpenFile("scan_results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer dosya.Close()

	// haber sayısı kontrolu icin sayac
	sayac := 0

	c := colly.NewCollector(colly.AllowedDomains("thehackernews.com", "www.thehackernews.com")) 
	c.UserAgent ="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"
	//onlem amaclı user agent
	c.OnHTML(".body-main .story-section, .body-main .article, a.story-link", func(e *colly.HTMLElement) {
		if sayac >= 5 {
			return
		}

		baslik := e.ChildText(".home-title")
		if baslik == "" { baslik = e.ChildText("h2") }
		if baslik == "" && e.DOM.HasClass("story-link") { baslik= e.Text }

		// Başlık icin kontrol cop bılgı olabılıyor
		if baslik != "" && len(baslik) > 10 {
			sayac++ 

			color.Cyan("\n>>> [%d] [HackerNews] News: %s", sayac,baslik)
			fmt.Fprintf(dosya,"Source: The Hacker News\nTitle: %s\n", baslik)

			if !tGizle {
    			tarihRaw :=e.ChildText(".item-label")
    			if tarihRaw != "" {
        			// 1. İkon karakteri vardı temızlemek lazım
        			tarihRaw = strings.ReplaceAll(tarihRaw, "", "")
        			tarihRaw = strings.TrimSpace(tarihRaw)

        			// Tarıh sorunu ıcın : Virgülden sonrasına bakmmak
        			// Tarih formatı: "April 16, 2026" (Yani her zaman bir virgül ve 4 haneli yıl var)
        			if strings.Contains(tarihRaw, ",") {
            			parts :=strings.Split(tarihRaw, ",") // Virgülden böl
            			if len(parts) > 1 {
                			// Virgülden sonraki ilk 5 karakteri al (Bosluk+ 4 hane yıl:" 2026")
                			// Boylece bozlmaz
                			yilVeSonrasi := strings.TrimSpace(parts[1])
                			if len(yilVeSonrasi) >= 4 {
                    			yil := yilVeSonrasi[:4]
                    			tarih :=parts[0] + ", " + yil
                    			color.Yellow("Date: %s", tarih)
                    			fmt.Fprintf(dosya,"Date: %s\n", tarih)
                			}
            			}
        			} else {
            			// Eğer virgül yoksa bu muhtemelen sadece bir kategori
            			color.Yellow("Date: No inf.") 
            			fmt.Fprintf(dosya, "Date:No inf.\n")
        			}
    			}
			}

			if !aGizle {
				desc := e.ChildText(".home-desc")
				if desc != "" {
					fmt.Printf("Description: %s\n", desc)
					fmt.Fprintf(dosya, "Description: %s\n", desc)
				}
			}
			fmt.Fprintln(dosya, "-------------------------")
			dosya.Sync()
		}
	})
	c.Visit("https://thehackernews.com/")
}

func bleepingComputerScraping(tGizle bool, aGizle bool) {
    dosya, _ := os.OpenFile("scan_results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    defer dosya.Close()
    sayac := 0

    c := colly.NewCollector(colly.AllowedDomains("bleepingcomputer.com", "www.bleepingcomputer.com"))
    c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"

    c.OnHTML(".bc_latest_news_text", func(e *colly.HTMLElement) {
        if sayac >=5 {
            return
        }

        baslik := e.ChildText("h4 a")
        // 10 karakter sınırı
        if baslik != "" && len(baslik) > 10 {
            sayac++ 
            color.Cyan("\n>>> [%d] [Bleeping] News: %s", sayac, baslik)
            fmt.Fprintf(dosya, "Source: Bleeping Computer\nTitle: %s\n", baslik)

            if !tGizle {
                tarih := ""
                e.ForEach("li", func(_ int, el *colly.HTMLElement) {
                    if strings.Contains(el.Text, "202") || strings.Contains(el.Text, ",") {
                        tarih = el.Text
                    }
                })
                if tarih != "" {
                    color.Yellow("Date: %s", tarih)
                    fmt.Fprintf(dosya, "Date: %s\n", tarih)
                }
            }

            if !aGizle {
                desc := e.ChildText("p")
                if desc != "" {
                    fmt.Printf("Description: %s\n", desc)
                    fmt.Fprintf(dosya, "Description: %s\n", desc)
                }
            }
            fmt.Fprintln(dosya, "-------------------------")
            dosya.Sync()
        }
    })
    c.Visit("https://www.bleepingcomputer.com/")
}
func krebsonSecurityScraping(tGizle bool, aGizle bool) {
	dosya, _ := os.OpenFile("scan_results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer dosya.Close()
	sayac := 0
	c := colly.NewCollector(colly.AllowedDomains("krebsonsecurity.com"))
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"
	//div
	c.OnHTML(".post", func(e *colly.HTMLElement) {
		if sayac >= 5 {
			return
		}

		//h2 içindeki link
		baslik := e.ChildText("h2 a")
		if baslik == "" {
			baslik = e.ChildText("h2")
		}
		baslik = strings.TrimSpace(baslik)

		if baslik != "" {
			sayac++
			color.Cyan("\n>>> [%d] [Krebs] News: %s", sayac, baslik)
			fmt.Fprintf(dosya, "Source: KrebsOnSecurity\nTitle: %s\n", baslik)

			if !tGizle {
				tarih := e.ChildText(".date.updated, .adt span, .post-date")
				if tarih == "" {
					tarih = e.ChildText("small")
				}
				// "Posted on..." kısmını temizleyelim
				tarih = strings.TrimPrefix(strings.TrimSpace(tarih), "Posted on ")
				
				color.Yellow("Date: %s", tarih)
				fmt.Fprintf(dosya, "Date: %s\n", tarih)
			}

			if !aGizle {
				// entry içindeki ilk p
				desc := e.ChildText(".entry p")
				if desc == "" {
					desc = e.ChildText("p")
				}
				
				desc = strings.TrimSpace(desc)
				if len(desc) > 250 {
					desc = desc[:250] + "..."
				}
				fmt.Printf("Description: %s\n", desc)
				fmt.Fprintf(dosya, "Description: %s\n", desc)
			}
			fmt.Fprintln(dosya, "-------------------------")
			dosya.Sync()
		}
	})

	c.Visit("https://krebsonsecurity.com/")
}
func asciiBaslik() {
	c1 := color.New(color.FgMagenta, color.Bold)       //figlet4go kutupahesınde sorun cıktıgı ıcın manuel
	c2 := color.New(color.FgCyan, color.Bold)
	c1.Println(" __     __                      _             ")
	c1.Println(" \\ \\   / /                     | |            ")
	c2.Println("  \\ \\_/ /_ ___   ___   _ ____  | | __ _ _ __ ")
	c2.Println("   \\   / _` \\ \\ / / | | |_  /  | |/ _` | '__|")
	c1.Println("    | | (_| |\\ V /| |_| |/ /   | | (_| | |    ")
	c1.Println("    |_|\\__,_| \\_/  \\__,_/___|  |_|\\__,_|_|    ")
	fmt.Println("")
	c2.Println("           IREM WS TOOL - Demo")
	c1.Println("----------------------------------------------")
}