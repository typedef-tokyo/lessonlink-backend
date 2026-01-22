package handler

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/controller"
	session_util "github.com/typedef-tokyo/lessonlink-backend/internal/adapter/utility"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	logWriter "github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/logger"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/server"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/repository"
)

func NewRouter(
	sever *server.Server,
	env configs.EnvConfig,
	logWriter *logWriter.LogWriter,
	sessionRepository repository.SessionRepository,
	campusListController controller.ICampusListController,
	lessonListController controller.ILessonListController,
	lessonAddController controller.ILessonAddController,
	lessonEditController controller.ILessonEditController,
	roomListController controller.IRoomListController,
	roomEditController controller.IRoomEditController,
	loginUserGetController controller.ILoginUserGetController,
	scheduleCreateController controller.IScheduleCreateController,
	scheduleDeleteController controller.IScheduleDeleteController,
	scheduleDuplicateController controller.IScheduleDuplicateController,
	scheduleGetController controller.IScheduleGetController,
	scheduleItemDivideController controller.IScheduleItemDivideController,
	scheduleItemJoinController controller.IScheduleItemJoinController,
	scheduleItemMoveController controller.IScheduleItemMoveController,
	scheduleItemReturnListController controller.IScheduleItemReturnListController,
	scheduleItemShiftController controller.IScheduleItemShiftController,
	scheduleListController controller.IScheduleListController,
	scheduleSaveController controller.IScheduleSaveController,
	scheduleSaveTitleController controller.IScheduleSaveTitleController,
	scheduleTimeEditController controller.IScheduleTimeEditController,
	invisibleRoomController controller.IInvisibleRoomController,
	userListController controller.IUserListController,
	userAddController controller.IUserAddController,
	userDeleteController controller.IUserDeleteController,
	userGetController controller.IUserGetController,
	userLoginController controller.IUserLoginController,
	userLogoutController controller.IUserLogoutController,
	userUpdateController controller.IUserUpdateController,
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
				_userID, _, err := session_util.GetSessionData(c)
				if err == nil {
					userID = _userID.Value()
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
	auth.Use(AuthRequiredMiddleware(env, sessionRepository, logWriter))

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
