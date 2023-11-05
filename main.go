package main

import (
	"encoding/json"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/hotkey"

	fynex "fyne.io/x/fyne/layout"
	fynexTheme "fyne.io/x/fyne/theme"
	"github.com/gen2brain/beeep"
	"github.com/sqweek/dialog"
)

var hk = hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
var a = app.New()
var w = a.NewWindow("Image Browser")

func writeConfig(conf Config) {
	file, _ := json.MarshalIndent(conf, "", " ")
	_ = os.WriteFile("config.json", file, 0644)
}

func readConfig() Config {
	file, err := os.ReadFile("./config.json")
	if err != nil {
		file, _ := json.MarshalIndent(Config{LastBrowseFolder: "",
			Config_Hotkey: Config_Hotkey{
				Hotkey_ShiftMod: true,
				Hotkey_CtrlMod:  true,
				Hotkey_Key:      "P",
			},
		}, "", " ")
		_ = os.WriteFile("config.json", file, 0644)
		return readConfig()
	} else {
		data := Config{}
		json.Unmarshal(file, &data)
		return data
	}
}

func parseImageToGrid(g *fyne.Container, path string) {
	images := GetImages(path)

	// Reset List
	g.Objects = []fyne.CanvasObject{}

	for _, img := range images {
		// Fix GO foreach variables bug
		imgFix := img

		im := canvas.NewImageFromFile(imgFix.Path)
		im.FillMode = canvas.ImageFillOriginal

		button := widget.NewButton("", func() {
			err := beeep.Notify("Quick Image Browser", "Image copied to clipboard!", imgFix.Path)
			if err != nil {
				panic(err)
			}
			CopyImageToClipboard(imgFix.Path)
		})

		container := container.New(
			layout.NewStackLayout(),
			button,
			im,
		)

		g.Objects = append(g.Objects, container)
	}

	g.Refresh()

	curConfig := readConfig()
	curConfig.LastBrowseFolder = path
	writeConfig(curConfig)
}

// hk = *hotkey.New(modifiers, hotkey.Key(keyAsHex))

func updateHotkey(currentConfig Config) {
	// hk.Unregister()

	modifiers := []hotkey.Modifier{}
	if currentConfig.Config_Hotkey.Hotkey_CtrlMod {
		modifiers = append(modifiers, 0x25)
	}
	if currentConfig.Config_Hotkey.Hotkey_ShiftMod {
		modifiers = append(modifiers, 0x32)
	}

	keyAsHex := getUintKey(currentConfig.Config_Hotkey.Hotkey_Key)

	go func() {
		hk = hotkey.New(modifiers, hotkey.Key(keyAsHex))
		fmt.Println(hk.String())

		err := hk.Register()
		if err != nil {
			fmt.Println("ERROR")
			return
		}

		for range hk.Keydown() {
			fmt.Println("It's doesn't work :P")
		}
	}()

	writeConfig(currentConfig)
}

func getOptionTab(currentConfig Config) *fyne.Container {
	ctrlModCheckbox := widget.NewCheck("Ctrl", func(b bool) {
		currentConfig.Config_Hotkey.Hotkey_CtrlMod = b
		updateHotkey(currentConfig)
	})
	ctrlModCheckbox.SetChecked(currentConfig.Config_Hotkey.Hotkey_CtrlMod)

	shiftModCheckbox := widget.NewCheck("Shift", func(b bool) {
		currentConfig.Config_Hotkey.Hotkey_ShiftMod = b
		updateHotkey(currentConfig)
	})
	shiftModCheckbox.SetChecked(currentConfig.Config_Hotkey.Hotkey_ShiftMod)

	comboBox := widget.NewSelect(getKeyListMap(), func(s string) {
		currentConfig.Config_Hotkey.Hotkey_Key = s
		updateHotkey(currentConfig)
	})
	comboBox.SetSelected(currentConfig.Config_Hotkey.Hotkey_Key)

	settingtab := container.NewVBox(
		widget.NewLabel("Hotkey (Just for show, currently broken :P)"), fynex.NewResponsiveLayout(
			fynex.Responsive(ctrlModCheckbox, 0.2),  // all sizes to 100%
			fynex.Responsive(shiftModCheckbox, 0.2), // all sizes to 100%  // all sizes to 100%
			fynex.Responsive(comboBox, 0.4),         // all sizes to 100% // small to 50%, medium to 75%, all others to 100%
		))

	return settingtab
}

func main() {
	// currentDirectory := ""
	currentConfig := readConfig()

	a.Settings().SetTheme(fynexTheme.AdwaitaTheme())

	w.SetIcon(theme.BrokenImageIcon())
	w.Resize(fyne.NewSize(500, 500))

	// ░█▄█░█▀█░▀█▀░█▀█░░░▀█▀░█▀█░█▀▄
	// ░█░█░█▀█░░█░░█░█░░░░█░░█▀█░█▀▄
	// ░▀░▀░▀░▀░▀▀▀░▀░▀░░░░▀░░▀░▀░▀▀░

	g := container.NewGridWrap(fyne.NewSize(100, 100))
	s := container.NewVScroll(g)
	s.SetMinSize(fyne.NewSize(100, 500))

	input := widget.NewEntry()
	input.SetText(currentConfig.LastBrowseFolder)
	input.SetPlaceHolder("Enter path...")

	goButton := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		parseImageToGrid(g, input.Text)
	})

	browseButton := widget.NewButtonWithIcon("", theme.FolderIcon(), func() {
		directory, err := dialog.Directory().Title("Load images").Browse()
		if err != nil {
			// panic(err)
		} else {
			fmt.Println(directory)
			input.SetText(directory)
			parseImageToGrid(g, directory)
		}
	})

	inputContainer := fynex.NewResponsiveLayout(
		fynex.Responsive(input, 0.8, 0.9),
		fynex.Responsive(goButton, 0.1, 0.05),
		fynex.Responsive(browseButton, 0.1, 0.05),
	)

	// ░█▀▀░█▀▀░▀█▀░▀█▀░▀█▀░█▀█░█▀▀░░░▀█▀░█▀█░█▀▄
	// ░▀▀█░█▀▀░░█░░░█░░░█░░█░█░█░█░░░░█░░█▀█░█▀▄
	// ░▀▀▀░▀▀▀░░▀░░░▀░░▀▀▀░▀░▀░▀▀▀░░░░▀░░▀░▀░▀▀░

	settingtab := getOptionTab(currentConfig)

	mainTab := container.NewVBox(
		inputContainer,
		s,
	)
	tabs := container.NewAppTabs(&container.TabItem{Text: "Home", Content: mainTab}, &container.TabItem{Text: "Settings", Content: settingtab})

	w.SetContent(tabs)

	go func() {
		if currentConfig.LastBrowseFolder != "" {
			parseImageToGrid(g, currentConfig.LastBrowseFolder)
		}
		updateHotkey(currentConfig)
	}()

	if desk, ok := a.(desktop.App); ok {
		desk.SetSystemTrayIcon(theme.FileIcon())

		m := fyne.NewMenu("Quick Image Browser",
			fyne.NewMenuItem("Show", func() {
				w.Show()
			}))

		desk.SetSystemTrayMenu(m)

		w.SetCloseIntercept(func() {
			w.Hide()
			// hk.Unregister()
		})
	}

	w.ShowAndRun()
}
