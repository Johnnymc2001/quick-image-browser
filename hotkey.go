package main

import "github.com/samber/lo"

func getHotkeyDict() []HotkeyObj {
	return []HotkeyObj{
		{"Space", 0x0020},
		{"1", 0x0030},
		{"2", 0x0031},
		{"3", 0x0032},
		{"4", 0x0033},
		{"5", 0x0034},
		{"6", 0x0035},
		{"7", 0x0036},
		{"8", 0x0037},
		{"9", 0x0038},
		{"0", 0x0039},
		{"A", 0x0061},
		{"B", 0x0062},
		{"C", 0x0063},
		{"D", 0x0064},
		{"E", 0x0065},
		{"F", 0x0066},
		{"G", 0x0067},
		{"H", 0x0068},
		{"I", 0x0069},
		{"J", 0x006a},
		{"K", 0x006b},
		{"L", 0x006c},
		{"M", 0x006d},
		{"N", 0x006e},
		{"O", 0x006f},
		{"P", 0x0070},
		{"Q", 0x0071},
		{"R", 0x0072},
		{"S", 0x0073},
		{"T", 0x0074},
		{"U", 0x0075},
		{"V", 0x0076},
		{"W", 0x0077},
		{"X", 0x0078},
		{"Y", 0x0079},
		{"Z", 0x007a},
		{"Return", 0xff0d},
		{"Escape", 0xff1b},
		{"Delete", 0xffff},
		{"Tab", 0xff1b},
		{"Left", 0xff51},
		{"Right", 0xff53},
		{"Up", 0xff52},
		{"Down", 0xff54},
		{"F1", 0xffbe},
		{"F2", 0xffbf},
		{"F3", 0xffc0},
		{"F4", 0xffc1},
		{"F5", 0xffc2},
		{"F6", 0xffc3},
		{"F7", 0xffc4},
		{"F8", 0xffc5},
		{"F9", 0xffc6},
		{"F10", 0xffc7},
		{"F11", 0xffc8},
		{"F12", 0xffc9},
		{"F13", 0xffca},
		{"F14", 0xffcb},
		{"F15", 0xffcc},
		{"F16", 0xffcd},
		{"F17", 0xffce},
		{"F18", 0xffcf},
		{"F19", 0xffd0},
		{"F20", 0xffd1},
	}
}

func getKeyListMap() []string {
	return lo.Map(getHotkeyDict(), func(hkObj HotkeyObj, _ int) string {
		return hkObj.Key
	})
}

func getUintKey(key string) uint16 {
	str, _ := lo.Find(getHotkeyDict(), func(item HotkeyObj) bool {
		return item.Key == key
	})

	return str.KeyHex
}
