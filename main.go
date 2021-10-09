package main

import (
	"awesomeProject/modules"
	"awesomeProject/modules/amazon"
	"awesomeProject/modules/coinbase"
	"awesomeProject/modules/debouncer"
	"awesomeProject/modules/disney"
	"awesomeProject/modules/managerone"
	"awesomeProject/modules/netflix"
	"awesomeProject/modules/paypal"
	"awesomeProject/modules/spotify"
	"awesomeProject/prevents"
	"awesomeProject/utils"
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/pterm/pterm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var userLvl string = "2"

func main() {
	go prevents.StartAllPrevents()

	if _, err := os.Stat("./lists"); os.IsNotExist(err) {
		_ = os.Mkdir("./lists", os.ModePerm)
	}

	if _, err := os.Stat("./results"); os.IsNotExist(err) {
		_ = os.Mkdir("./results/bad", os.ModePerm)
		_ = os.Mkdir("./results/good", os.ModePerm)
		_= os.Mkdir("./results/errors", os.ModePerm)
	}

	if _, err := os.Stat("./results/bad"); os.IsNotExist(err) {
		_= os.Mkdir("./results/bad", os.ModePerm)
	}

	if _, err := os.Stat("./results/good"); os.IsNotExist(err) {
		_ = os.Mkdir("./results/good", os.ModePerm)
	}

	if _, err := os.Stat("./results/errors"); os.IsNotExist(err) {
		_ = os.Mkdir("./results/errors", os.ModePerm)
	}

	_, _ = utils.SetConsoleTitle("Larez AIO - By @mlk667 on telegram")

	d := color.New(color.FgHiYellow, color.Italic)

	d.Println("                 yyyyyyyyh")
	d.Println("                 y.      `s")
	d.Println("           myym  +-    .-  o")
	d.Println("          y.  `+mmy   ``    d                    :::            :::     :::::::::  :::::::::: :::::::::")
	d.Println("        m/      `+    hmy::o                    :+:          :+: :+:   :+:    :+: :+:             :+:")
	d.Println("       m``-:/`    m   :h                       +:+         +:+   +:+  +:+    +:+ +:+            +:+")
	d.Println("          m+`     h.    .+h                   +#+        +#++:++#++: +#++:++#:  +#++:++#      +#+")
	d.Println("           yssssssds       /                 +#+        +#+     +#+ +#+    +#+ +#+          +#+")
	d.Println("                    s`      d               #+#        #+#     #+# #+#    #+# #+#         #+#")
	d.Println("            dhyyo++/:.      d              ########## ###     ### ###    ### ########## #########")
	d.Println("       mo/-.``            .s")
	d.Println("       myyyyyyyyyyyyyyyyyhm                Made by @mlk667 on telegram - #s/o witeCapz & Onirez")

	utils.ComplexTextAnim("\n_______________________________________________________________________________________________________________\n", 5, []*color.Color{
		color.New(color.FgHiYellow),
		color.New(color.FgYellow),
		color.New(color.FgWhite),
		color.New(color.FgHiBlack),
	}, utils.ModeLoop)

	utils.TextAnim("\nWelcome to Larez,", 50, color.New(color.FgHiYellow, color.Bold))
	color.New(color.FgWhite, color.Bold).Printf("                                                                         Don't have licence yet ?\n")
	utils.TextAnim("Current version: v2.0", 50, color.New(color.FgHiWhite, color.Bold))
	fmt.Print("                                                                        ")
	color.New(color.FgBlack, color.BgHiYellow).Printf("https://larez.eu\n\n")
	utils.TextAnim("Annoucement:\n", 50, color.New(color.FgHiYellow, color.Bold))
	color.New(color.FgHiWhite, color.Bold).Printf(getLarezAnnouncements() + "\n\n")
	time.Sleep(time.Second * 1)

	reader := bufio.NewReader(os.Stdin)
	color.New(color.FgHiYellow, color.Bold).Printf("Press any key to continue: \n>")
	_, _ = reader.ReadString('\n')

	utils.ClearConsole()
	listTypeInput()
}

func proxyTypeInput()  {
	utils.ClearConsole()
	pterm.NewStyle(pterm.FgLightYellow).Print("With this mode, you can use proxyless version\n")
	pterm.NewStyle(pterm.FgLightYellow).Print("Write 1 to active proxyless mode (not recommended)\n")

	pterm.NewStyle(pterm.FgLightYellow).Print("> ")

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = utils.CleanText(text)

	if text == "1" {
		modules.IsProxylessMode = true
	} else {
		utils.ClearConsole()
		d := color.New(color.FgHiYellow, color.Bold)
		d.Printf("Please enter your proxy type: \n")

		_ = pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{
			{Level: 2, Text: "HTTP(S)", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "1 >", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
			{Level: 2, Text: "SOCKS4", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "2 >", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
			{Level: 2, Text: "SOCKS5", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "3 >", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
		}).Render()

		pterm.NewStyle(pterm.FgLightYellow).Print("> ")

		reader2 := bufio.NewReader(os.Stdin)
		text, _ = reader2.ReadString('\n')
		text = utils.CleanText(text)

		fmt.Print(text)

		switch text {
			case "1":
				modules.ProxyType = "HTTP"
			case "2":
				modules.ProxyType = "SOCKS4"
			case "3":
				modules.ProxyType = "SOCKS5"
			default:
				d := color.New(color.FgRed, color.Bold)
				d.Printf("\n ⚠ Invalid proxy type !\n")

				time.Sleep(time.Second * 2)
				proxyTypeInput()

				break
		}

	}
	
	

}

func speedTypeInput()  {
	utils.ClearConsole()
	pterm.NewStyle(pterm.FgLightYellow).Print("Salut, je suis un texte en français^^ Je suis différent des autres !\n")
	pterm.NewStyle(pterm.FgLightYellow).Print("Bref, avant de lancer ce mode tu dois choisir la vitesse (en fonction de ton pc)\n")
	_ = pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{
		{Level: 2, Text: "Brouteur (lent)", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "1 ->" , BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
		{Level: 2, Text: "BMW M8 F92 cylindrée de 4,4 litres 625 ch (460 kW) (normal)", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "2 ->", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
		{Level: 2, Text: "Usain Bolt (rapide)", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "3 ->", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
	}).Render()

	pterm.NewStyle(pterm.FgLightYellow).Print("> ")

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = utils.CleanText(text)

		switch text {
		case "1":
			modules.Speed = 50
		case "2":
			modules.Speed = 10
		case "3":
			modules.Speed = 1
		default:
			d := color.New(color.FgRed, color.Bold)
			d.Printf("\n ⚠ Invalid speed type !\n")

			time.Sleep(time.Second * 2)
			speedTypeInput()

			break
		}

	}


var ListMail int64 = 0
var ListPN int64 = 1

func listTypeInput()  {
	utils.ClearConsole()

	pterm.NewStyle(pterm.FgLightYellow).Println("Select your list type:")
	_ = pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{
		{Level: 1, Text: "WITH MAILS", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "1 ->" , BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
		{Level: 1, Text: "WITH PHONE-NUMBERS", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "2 ->", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
	}).Render()
	pterm.NewStyle(pterm.FgLightYellow).Print("> ")

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = utils.CleanText(text)

	var file string

	switch text {
		case "1":
			file = pickListFile(ListMail)
			proxyFile := pickProxyListFile(ListMail)
			modeInput(ListMail, file, proxyFile)
			break
		case "2":
			file = pickListFile(ListPN)
			proxyFile := pickProxyListFile(ListPN)
			modeInput(ListPN, file, proxyFile)
			break

		default:
			d := color.New(color.FgRed, color.Bold)
			d.Printf("\n ⚠ Invalid mode !\n")

			time.Sleep(time.Second * 2)
			listTypeInput()
			break
	}

	time.Sleep(1000 * time.Second)
}

func pickListFile(listType int64) string {
	utils.ClearConsole()

	files, _ := ioutil.ReadDir("./lists")

	var bulletListItems []pterm.BulletListItem
	switch listType {
		case ListMail:
			pterm.NewStyle(pterm.FgLightYellow).Println("Select your mail list:")
			bulletListItems = append(
				bulletListItems,
				pterm.BulletListItem{
					Level: 1,
					Text: "WITH MAILS",
					TextStyle: pterm.NewStyle(pterm.FgLightWhite),
					Bullet: "  •",
					BulletStyle: pterm.NewStyle(pterm.FgLightYellow),
				},
			)
			break

		case ListPN:
			pterm.NewStyle(pterm.FgLightYellow).Println("Select your phone list:")
			bulletListItems = append(
				bulletListItems,
				pterm.BulletListItem{
					Level: 1,
					Text: "WITH PHONE-NUMBERS",
					TextStyle: pterm.NewStyle(pterm.FgLightWhite),
					Bullet: "  •",
					BulletStyle: pterm.NewStyle(pterm.FgLightYellow),
				},
			)
			break
	}

	for i, file := range files {
		bulletListItems = append(
			bulletListItems,
			pterm.BulletListItem{
				Level: 2,
				Text: file.Name(),
				TextStyle: pterm.NewStyle(pterm.FgLightWhite),
				Bullet: strconv.FormatInt(int64(i), 10) + " •",
				BulletStyle: pterm.NewStyle(pterm.FgLightYellow),
			},
		)
	}
	_ = pterm.DefaultBulletList.WithItems(bulletListItems).Render()

	pterm.NewStyle(pterm.FgLightYellow).Print("> ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = utils.CleanText(text)
	index, _ := strconv.Atoi(text)

	return files[index].Name()
}

func pickProxyListFile(listType int64) string {
	utils.ClearConsole()
	files, _ := ioutil.ReadDir("./proxylist")

	var bulletListItems []pterm.BulletListItem
	switch listType {
	case ListMail:
		pterm.NewStyle(pterm.FgLightYellow).Println("Select your proxy list:")

		bulletListItems = append(
			bulletListItems,
			pterm.BulletListItem{
				Level: 1,
				Text: "WITH MAILS",
				TextStyle: pterm.NewStyle(pterm.FgLightWhite),
				Bullet: "  •",
				BulletStyle: pterm.NewStyle(pterm.FgLightYellow),
			},
		)
		break

	case ListPN:
		pterm.NewStyle(pterm.FgLightYellow).Println("Select your proxy list:")

		bulletListItems = append(
			bulletListItems,
			pterm.BulletListItem{
				Level: 1,
				Text: "WITH PHONE-NUMBERS",
				TextStyle: pterm.NewStyle(pterm.FgLightWhite),
				Bullet: "  •",
				BulletStyle: pterm.NewStyle(pterm.FgLightYellow),
			},
		)
		break
	}

	for i, file := range files {
		bulletListItems = append(
			bulletListItems,
			pterm.BulletListItem{
				Level: 2,
				Text: file.Name(),
				TextStyle: pterm.NewStyle(pterm.FgLightWhite),
				Bullet: strconv.FormatInt(int64(i), 10) + " •",
				BulletStyle: pterm.NewStyle(pterm.FgLightYellow),
			},
		)
	}
	_ = pterm.DefaultBulletList.WithItems(bulletListItems).Render()

	pterm.NewStyle(pterm.FgLightYellow).Print("> ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = utils.CleanText(text)
	index, _ := strconv.Atoi(text)

	return files[index].Name()
}



func modeInput(listType int64, file string, proxyFile string)  {
	utils.ClearConsole()
	pterm.NewStyle(pterm.FgLightYellow).Println("Select your mode:")

	switch listType {
		case ListMail:
			_ = pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{
				{Level: 1, Text: "WITH MAILS", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "  •", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
				{Level: 1, Text: "WITH FILE " + file, TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "  •", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
				{Level: 2, Text: "DEBOUNCER", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "1 >", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
				{Level: 2, Text: "NETFLIX CHECKER", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "2 >", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
				{Level: 2, Text: "DISNEY+ CHECKER", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "3 >", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
				{Level: 2, Text: "PAYPAL CHECKER", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "4 >", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
				{Level: 2, Text: "AMAZON CHECKER", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "5 >", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
				{Level: 2, Text: "COINBASE CHECKER", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "6 >", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
				{Level: 2, Text: "SPOTIFY CHECKER", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "7 >", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
				{Level: 2, Text: "MANAGER ONE CHECKER", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "8 >", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
			}).Render()

			pterm.NewStyle(pterm.FgLightYellow).Print("> ")

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = utils.CleanText(text)

			switch text {
				case "1":
					if userLvl == "1" || userLvl == "0" || userLvl == "2"{
						utils.ClearConsole()
						emails := utils.ReadLineToArray("./lists/" + file)
						proxies := utils.ReadLineToArray("./proxylist/" + proxyFile)

						proxyTypeInput()
						speedTypeInput()

						debouncer.StartEmailChecker(emails, proxies, "")
						utils.ClearConsole()
					} else {
						d := color.New(color.FgRed, color.Bold)
						d.Printf("\n ⚠ You doesn't have access to this mode !\n")
					}
					break

				case "2":
					if userLvl == "0" || userLvl == "2" {
						utils.ClearConsole()
						emails := utils.ReadLineToArray("./lists/" + file)
						proxies := utils.ReadLineToArray("./proxylist/" + proxyFile)

						proxyTypeInput()
						speedTypeInput()

						netflix.StartEmailChecker(emails, proxies, "")
						utils.ClearConsole()
					} else {
						d := color.New(color.FgRed, color.Bold)
						d.Printf("\n ⚠ You doesn't have access to this mode !\n")
					}
					break
				case "3":
					if userLvl == "0" || userLvl == "2" {
						utils.ClearConsole()
						emails := utils.ReadLineToArray("./lists/" + file)
						proxies := utils.ReadLineToArray("./proxylist/" + proxyFile)

						proxyTypeInput()
						speedTypeInput()

						disney.StartEmailChecker(emails, proxies, "")
						utils.ClearConsole()
					} else {
						d := color.New(color.FgRed, color.Bold)
						d.Printf("\n ⚠ You doesn't have access to this mode !\n")
					}
					break

				case "4":
					if userLvl == "0" || userLvl == "2" {
						utils.ClearConsole()
						emails := utils.ReadLineToArray("./lists/" + file)
						proxies := utils.ReadLineToArray("./proxylist/" + proxyFile)

						proxyTypeInput()
						speedTypeInput()

						paypal.StartEmailChecker(emails, proxies, "")
						utils.ClearConsole()
					} else {
						d := color.New(color.FgRed, color.Bold)
						d.Printf("\n ⚠ You doesn't have access to this mode !\n")
					}

					break

				case "5":
					if userLvl == "0" || userLvl == "2" {
						utils.ClearConsole()
						emails := utils.ReadLineToArray("./lists/" + file)
						proxies := utils.ReadLineToArray("./proxylist/" + proxyFile)

						proxyTypeInput()
						speedTypeInput()

						amazon.StartEmailChecker(emails, proxies, "")
						utils.ClearConsole()
					} else {
						d := color.New(color.FgRed, color.Bold)
						d.Printf("\n ⚠ You doesn't have access to this mode !\n")
					}
					break

				case "6":
					if userLvl == "0" || userLvl == "2" {
						utils.ClearConsole()
						emails := utils.ReadLineToArray("./lists/" + file)
						proxies := utils.ReadLineToArray("./proxylist/" + proxyFile)

						proxyTypeInput()
						speedTypeInput()

						coinbase.StartEmailChecker(emails, proxies, "")
						utils.ClearConsole()
					} else {
						d := color.New(color.FgRed, color.Bold)
						d.Printf("\n ⚠ You doesn't have access to this mode !\n")
					}
					break

				case "7":
					if userLvl == "0" || userLvl == "2" {
						utils.ClearConsole()
						emails := utils.ReadLineToArray("./lists/" + file)
						proxies := utils.ReadLineToArray("./proxylist/" + proxyFile)

						proxyTypeInput()
						speedTypeInput()

						spotify.StartEmailChecker(emails, proxies, "")
						utils.ClearConsole()
					} else {
						d := color.New(color.FgRed, color.Bold)
						d.Printf("\n ⚠ You doesn't have access to this mode !\n")
					}
					break

			case "8":
				if userLvl == "0" || userLvl == "2" {
					utils.ClearConsole()
					emails := utils.ReadLineToArray("./lists/" + file)
					proxies := utils.ReadLineToArray("./proxylist/" + proxyFile)

					proxyTypeInput()
					speedTypeInput()

					managerone.StartEmailChecker(emails, proxies, "")
					utils.ClearConsole()
				} else {
					d := color.New(color.FgRed, color.Bold)
					d.Printf("\n ⚠ You doesn't have access to this mode !\n")
				}
				break

				default:
					d := color.New(color.FgRed, color.Bold)
					d.Printf("\n ⚠ Invalid mode !\n")

					time.Sleep(time.Second * 2)
					modeInput(ListMail, file, proxyFile)
					break
				}

				break

		case ListPN:
			_ = pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{
				{Level: 1, Text: "WITH PHONE-NUMBERS", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "  •", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
				{Level: 1, Text: "WITH FILE " + file, TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "  •", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
				{Level: 2, Text: "AMAZON CHECKER", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "1 >", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
				{Level: 2, Text: "CARRIER CHECKER", TextStyle: pterm.NewStyle(pterm.FgLightWhite), Bullet: "2 >", BulletStyle: pterm.NewStyle(pterm.FgLightYellow)},
			}).Render()

			pterm.NewStyle(pterm.FgLightYellow).Print("> ")

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = utils.CleanText(text)

			switch text {
				case "1":
					fmt.Println("SOON")
					time.Sleep(time.Second * 3)
					modeInput(listType, file, proxyFile)
					break

				case "2":
					fmt.Println("SOON")
					time.Sleep(time.Second * 3)
					modeInput(listType, file, proxyFile)
					break

				default:
					d := color.New(color.FgRed, color.Bold)
					d.Printf("\n ⚠ Invalid mode !\n")

					time.Sleep(time.Second * 2)
					modeInput(ListPN, file, proxyFile)
					break
				}
				break
	}
}

func getLarezAnnouncements() string {
	resp, err := http.Get("http://larez.eu/help/annoucement.txt")
	if err != nil {
		log.Fatalln("Unable to retrieve larez information")

	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln("Unable to retrieve larez information")
	}

	sb := string(body)

	return sb

}