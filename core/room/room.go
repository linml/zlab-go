package room

type Room struct {
	ID        string
	MaxPlayer int32
	MinPlayer int32
	Players   RoomSeatSlice //[座位号]plyid
	IsSystem  bool          //是否是系统创建
	IsStream  bool          //是否是流动场 /.无固定座位
}

type RoomSeat struct {
	PlayerID string //用户id
	Seat     uint8  //座位号
}
type RoomSeatSlice []RoomSeat

func (s RoomSeatSlice) Less(i, j int) bool {
	return s[i].Seat < s[j].Seat
}

func (s RoomSeatSlice) Len() int {
	return len(s)
}
func (s RoomSeatSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *RoomSeatSlice) AddSort(plyid string, seat uint8) {
	for i, v := range *s {
		if v.Seat > seat {
			var tmp = (*s)[i:]
			*s = append((*s)[:i], RoomSeat{
				PlayerID: plyid,
				Seat:     seat,
			})
			*s = append((*s), tmp...)
			break
		}
	}
}
func (s *RoomSeatSlice) Add(plyid string, seat uint8) {
	*s = append(*s, RoomSeat{
		PlayerID: plyid,
		Seat:     seat,
	})
}

func (s *RoomSeatSlice) Remove(seat int) {

}
