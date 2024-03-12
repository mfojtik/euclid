package timeago

func myLocale(_ float64, index int) (ago string, in string) {
	var res = [][]string{
		{"ယခုအတွင်း", "ယခု"},
		{"%d စက္ကန့် အကြာက", "%d စက္ကန့်အတွင်း"},
		{"1 မိနစ် အကြာက", "1 မိနစ်အတွင်း"},
		{"%d မိနစ် အကြာက", "%d မိနစ်အတွင်း"},
		{"1 နာရီ အကြာက", "1 နာရီအတွင်း"},
		{"%d နာရီ အကြာက", "%d နာရီအတွင်း"},
		{"1 ရက် အကြာက", "1 ရက်အတွင်း"},
		{"%d ရက် အကြာက", "%d ရက်အတွင်း"},
		{"1 ပတ် အကြာက", "1 ပတ်အတွင်း"},
		{"%d ပတ် အကြာက", "%d ပတ်အတွင်း"},
		{"1 လ အကြာက", "1 လအတွင်း"},
		{"%d လ အကြာက", "%d လအတွင်း"},
		{"1 နှစ် အကြာက", "1 နှစ်အတွင်း"},
		{"%d နှစ် အကြာက", "%d နှစ်အတွင်း"},
	}[index]
	return res[0], res[1]
}
