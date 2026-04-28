package handler

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/server"
	campus_controller "github.com/typedef-tokyo/lessonlink-backend/internal/modules/campus/adapter/controller"
	lesson_controller "github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/adapter/controller"
	room_controller "github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/adapter/controller"
	schedule_controller "github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/adapter/controller"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/adapter/handler"
	user_controller "github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/adapter/controller"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	logWriter "github.com/typedef-tokyo/lessonlink-backend/internal/platform/logger"
)

func NewRouter(
	sever *server.Server,
	env configs.EnvConfig,
	logWriter *logWriter.LogWriter,
	sessionHandler handler.ISessionGetHandler,
	campusListController campus_controller.ICampusListController,
	lessonListController lesson_controller.ILessonListController,
	lessonAddController lesson_controller.ILessonAddController,
	lessonEditController lesson_controller.ILessonEditController,
	roomListController room_controller.IRoomListController,
	roomEditController room_controller.IRoomEditController,
	scheduleCreateController schedule_controller.IScheduleCreateController,
	scheduleDeleteController schedule_controller.IScheduleDeleteController,
	scheduleDuplicateController schedule_controller.IScheduleDuplicateController,
	scheduleGetController schedule_controller.IScheduleGetController,
	scheduleItemDivideController schedule_controller.IScheduleItemDivideController,
	scheduleItemJoinController schedule_controller.IScheduleItemJoinController,
	scheduleItemMoveController schedule_controller.IScheduleItemMoveController,
	scheduleItemReturnListController schedule_controller.IScheduleItemReturnListController,
	scheduleItemShiftController schedule_controller.IScheduleItemShiftController,
	scheduleListController schedule_controller.IScheduleListController,
	scheduleSaveController schedule_controller.IScheduleSaveController,
	scheduleSaveTitleController schedule_controller.IScheduleSaveTitleController,
	scheduleTimeEditController schedule_controller.IScheduleTimeEditController,
	invisibleRoomController schedule_controller.IInvisibleRoomController,
	loginUserGetController user_controller.ILoginUserGetController,
	userListController user_controller.IUserListController,
	userAddController user_controller.IUserAddController,
	userDeleteController user_controller.IUserDeleteController,
	userGetController user_controller.IUserGetController,
	userLoginController user_controller.IUserLoginController,
	userLogoutController user_controller.IUserLogoutController,
	userUpdateController user_controller.IUserUpdateController,
) *echo.Echo {

	api := sever.Engine.Group("/api")
	user := api.Group("/user")
	user.POST("/login", userLoginController.Execute)

	sever.Engine.Use(middleware.RequestID())
	if env.LogErrorRequestDump {

		sever.Engine.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			status := c.Response().Status
			if status >= 400 && status != http.StatusUnauthorized {

				url := c.Request().URL.String()
				method := c.Request().Method
				params := c.QueryParams()

				transactionID := c.Response().Header().Get(echo.HeaderXRequestID)

				userID := -1
				sessionUserID, _, err := session_util.GetSessionData(c)
				if err == nil {
					userID = sessionUserID
				}

				level := log.LogLevelMap[status]
				severity := log.LogSeverityMap[status]

				logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
					AddSource: false,
				}))
				logger.Log(
					c.Request().Context(),
					level,
					severity,
					slog.String("transactionID", transactionID),
					slog.String("userid", strconv.Itoa(userID)),
					slog.String("url", url),
					slog.String("method", method),
					slog.Any("params", params),
					slog.Int("status", status),
					slog.String("reqBody", string(reqBody)),
				)
			}
		}))
	}

	// -------認証必須-------------------
	auth := api.Group("")
	auth.Use(AuthRequiredMiddleware(env, sessionHandler, logWriter))

	campus := auth.Group("/campus")
	campus.GET("/list", campusListController.Execute)

	lesson := auth.Group("/lesson")
	lesson.GET("/:campus/list", lessonListController.Execute)
	lesson.POST("/:campus", lessonAddController.Execute)
	lesson.PATCH("/:lessonid", lessonEditController.Execute)

	room := auth.Group("/room")
	room.GET("/:campus/list", roomListController.Execute)
	room.POST("/:campus/edit", roomEditController.Execute)

	schedule := auth.Group("/schedule")
	schedule.GET("/list/:campus", scheduleListController.Execute)
	schedule.POST("/create/:campus", scheduleCreateController.Execute)
	schedule.GET("/:schedule_id", scheduleGetController.Execute)
	schedule.POST("/:schedule_id", scheduleSaveController.Execute)
	schedule.POST("/:schedule_id/item-move", scheduleItemMoveController.Execute)
	schedule.POST("/:schedule_id/item-return-list", scheduleItemReturnListController.Execute)
	schedule.POST("/:schedule_id/item-divide", scheduleItemDivideController.Execute)
	schedule.POST("/:schedule_id/item-join", scheduleItemJoinController.Execute)
	schedule.POST("/:schedule_id/item-shift", scheduleItemShiftController.Execute)
	schedule.PATCH("/:schedule_id/title", scheduleSaveTitleController.Execute)
	schedule.DELETE("/:schedule_id", scheduleDeleteController.Execute)
	schedule.POST("/:schedule_id/duplicate", scheduleDuplicateController.Execute)
	schedule.PUT("/:schedule_id/room/invisible", invisibleRoomController.Execute)
	schedule.PATCH("/:schedule_id/time", scheduleTimeEditController.Execute)

	authUser := auth.Group("/user")
	authUser.GET("/list", userListController.Execute)
	authUser.GET("/self", loginUserGetController.Execute)
	authUser.GET("/:userid", userGetController.Execute)
	authUser.POST("", userAddController.Execute)
	authUser.PUT("", userUpdateController.Execute)
	authUser.DELETE("/:userid", userDeleteController.Execute)
	authUser.POST("/logout", userLogoutController.Execute)

	initSwagger(env, sever)

	return sever.Engine
}
