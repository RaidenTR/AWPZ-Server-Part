package device

type deviceContoler struct {
	rooms    map[string]bool
	lectures map[string]string
}

var instance *deviceContoler

//GetInstance returns instance of deviceControler
func GetInstance() *deviceContoler {
	if instance == nil {
		instance = new(deviceContoler)
		instance.rooms = make(map[string]bool)
		instance.lectures = make(map[string]string)
	}
	return instance
}

func (d *deviceContoler) StartWatching(room, id string) {
	instance.rooms[room] = true
	instance.lectures[room] = id
}

func (d *deviceContoler) FinishWatching(room, id string) {
	instance.rooms[room] = false
	delete(instance.lectures, id)
}

func (d deviceContoler) IsStreaming(room string) bool {
	return instance.rooms[room]
}

func (d deviceContoler) GetLectureId(room string) string {
	return instance.lectures[room]
}
