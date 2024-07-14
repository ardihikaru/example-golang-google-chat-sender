package application

const (
	Msg1 string = "ayok guys, ditunggu standup daily-nya"
	Msg2 string = "Mari teman-teman, ditungguin standup daily-nya ya!"
	Msg3 string = "daily daily .."
	Msg4 string = "guys, daily standup sekarang"
	Msg5 string = "yok, daily standup guys"
	Msg6 string = "kindly reminder buat yang belum daily"
)

// GetScheduleMsg returns a raw message (string)
func GetScheduleMsg() map[int]string {
	return map[int]string{
		1: Msg1,
		2: Msg2,
		3: Msg3,
		4: Msg4,
		5: Msg5,
		6: Msg6,
	}
}

// GetScheduleList returns a raw message map
func GetScheduleList() []string {
	return []string{
		Msg1, Msg2, Msg3, Msg4, Msg5,
	}
}
