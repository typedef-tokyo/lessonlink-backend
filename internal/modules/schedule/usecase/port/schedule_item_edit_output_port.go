package port

type ScheduleItemEditOutput struct {
	ScheduleItem ScheduleItemEditOutputDTO
}

type (
	ScheduleItemEditOutputDTO struct {
		HistoryIndex   int
		LessonItemList []ScheduleLessonItem
		RoomLessonList []ScheduleRoomLesson
	}

	ScheduleLessonItem struct {
		LessonID   int
		Identifier string
		LessonName string
		Duration   int
	}

	ScheduleRoomLesson struct {
		ItemTag    string
		LessonID   int
		LessonName string
		Identifier string
		Duration   int
		StartTime  ScheduleItemEditRoomLessonTime
		EndTime    ScheduleItemEditRoomLessonTime
		RoomIndex  int
	}

	ScheduleItemEditRoomLessonTime struct {
		ScheduleItemTimeHour    int
		ScheduleItemTimeMinutes int
	}
)
