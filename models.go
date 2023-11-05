package main

type ImageObj struct {
	Path   string `json:"path"`
	Name   string `json:"name"`
	Base64 string `json:"base64"`
}

type HotkeyObj struct {
	Key    string `json:"key"`
	KeyHex uint16 `json:"keyHex"`
}

type Config struct {
	LastBrowseFolder string        `json:"lastBrowseFolder"`
	Config_Hotkey    Config_Hotkey `json:"hotkey"`
}

type Config_Hotkey struct {
	Hotkey_ShiftMod bool   `json:"shiftmod"`
	Hotkey_CtrlMod  bool   `json:"ctrlmod"`
	Hotkey_Key      string `json:"key"`
}
