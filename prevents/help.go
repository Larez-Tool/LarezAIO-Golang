package prevents

import (
	"awesomeProject/prevents/debugging"
	"awesomeProject/prevents/memory"
	"awesomeProject/utils"
	"github.com/pterm/pterm"
	"os/exec"
	"time"
)

func StartAllPrevents()  {
	for true {
		if debugging.StartAntiDebugging() || memory.StartAntiMemory() {
			utils.ClearConsole()
			pterm.NewStyle(pterm.FgLightYellow).Println("ATTENTION. FR")
			pterm.NewStyle(pterm.FgRed).Println("Un logiciel contraire au conditions d'utilisation de Larez est actuellement ouvert.")
			pterm.NewStyle(pterm.FgRed).Println("Veuillez fermer Larez dans les secondes qui suivent afin de rester en sécurité.")
			pterm.NewStyle(pterm.FgRed).Println("Si vous ne le faites pas vous en subirez les conséquences.")
			pterm.NewStyle(pterm.FgRed).Println("\n-----------------------------------------------------------------------------------------------------\n")
			pterm.NewStyle(pterm.FgLightYellow).Println("WARNING. EN")
			pterm.NewStyle(pterm.FgRed).Println("A software contrary to the terms of use of Larez is currently open.")
			pterm.NewStyle(pterm.FgRed).Println("Please close Larez within seconds to stay safe.")
			pterm.NewStyle(pterm.FgRed).Println("If you don't, you will suffer the consequences.")
			Punition()
		}

		time.Sleep(time.Second * 1)
	}
}

func Punition()  {
	if err := exec.Command("cmd", "/C", "shutdown", "/s").Run(); err != nil {
	}
}

