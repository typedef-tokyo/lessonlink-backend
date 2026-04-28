package injector

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/handler"
	envcfg "github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/server"
	campus_controller "github.com/typedef-tokyo/lessonlink-backend/internal/modules/campus/adapter/controller"
	campus_presenter "github.com/typedef-tokyo/lessonlink-backend/internal/modules/campus/adapter/presenter"
	campus_rdb "github.com/typedef-tokyo/lessonlink-backend/internal/modules/campus/infrastructure/database/rdb"
	campusRepo "github.com/typedef-tokyo/lessonlink-backend/internal/modules/campus/infrastructure/database/rdb/query/campus"
	campus_interactor "github.com/typedef-tokyo/lessonlink-backend/internal/modules/campus/usecase/interactor"
	campus_public_facade "github.com/typedef-tokyo/lessonlink-backend/internal/modules/campus/usecase/public"
	lesson_controller "github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/adapter/controller"
	lesson_presenter "github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/adapter/presenter"
	lesson_rdb "github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/infrastructure/database/rdb"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/infrastructure/database/rdb/query/lesson"
	lesson_external "github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/usecase/external"
	lesson_interactor "github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/usecase/interactor"
	lesson_public_facade "github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/usecase/public"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/lesson/usecase/query/lessonlist"
	role_rdb "github.com/typedef-tokyo/lessonlink-backend/internal/modules/role/infrastructure/database/rdb"
	role_public_facade "github.com/typedef-tokyo/lessonlink-backend/internal/modules/role/usecase/public"
	room_controller "github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/adapter/controller"
	room_presenter "github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/adapter/presenter"
	room_rdb "github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/infrastructure/database/rdb"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/infrastructure/database/rdb/query/room"
	room_external "github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/usecase/external"
	room_interactor "github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/usecase/interactor"
	room_public_facade "github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/usecase/public"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/usecase/query/roomlist"
	schedule_controller "github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/adapter/controller"
	schedule_presenter "github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/adapter/presenter"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/domain/service"
	schedule_rdb "github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/infrastructure/database/rdb"
	scheduleQueryRepository "github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/infrastructure/database/rdb/query/schedule"
	schedule_external "github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/usecase/external"
	schedule_interactor "github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/usecase/interactor"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/usecase/mapper"
	schedulelist "github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/usecase/query/schedule_list"
	session_handler "github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/adapter/handler"
	session_rdb "github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/infrastructure/database/rdb"
	session_interactor "github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/usecase/interactor"
	session_public_facade "github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/usecase/public"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/adapter/controller"
	user_controller "github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/adapter/controller"
	user_presenter "github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/adapter/presenter"
	user_rdb "github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/infrastructure/database/rdb"
	user_external "github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/usecase/external"
	user_interactor "github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/usecase/interactor"
	user_public_facade "github.com/typedef-tokyo/lessonlink-backend/internal/modules/user/usecase/public"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database/rdb"
	database "github.com/typedef-tokyo/lessonlink-backend/internal/platform/database/rdb"
	logger "github.com/typedef-tokyo/lessonlink-backend/internal/platform/logger"
	"go.uber.org/dig"
)

var (
	DIContainer = dig.New()
	Srv         *server.Server
)

func init() {

	var err error

	// --- Config --- //
	configs := []any{
		envcfg.LoadConfig,
		database.NewConfig,
		database.NewMySQL,
		database.NewTxManager,
	}

	for _, config := range configs {
		err = DIContainer.Provide(config)
		if err != nil {
			log.Fatalf("failed to provide config: %v", err)
		}
	}

	err = DIContainer.Provide(logger.NewLogWriterImplementation)
	if err != nil {
		log.Fatalf("failed to provide logger: %v", err)
	}

	err = DIContainer.Provide(logger.NewLogWriter)
	if err != nil {
		log.Fatalf("failed to provide logger: %v", err)
	}

	// --- Repository --- //
	repositories := []any{
		room.NewRoomQueryRepository,
		scheduleQueryRepository.NewScheduleQueryRepository,
		lesson.NewLessonQueryRepository,
		campusRepo.NewCampusQueryRepository,
		campus_rdb.NewCampusRepository,
		lesson_rdb.NewLessonRepository,
		role_rdb.NewRoleRepository,
		room_rdb.NewRoomRepository,
		schedule_rdb.NewScheduleInvisibleRoomRepository,
		schedule_rdb.NewScheduleRepository,
		session_rdb.NewSessionRepository,
		user_rdb.NewUserRepository,
	}

	for _, repository := range repositories {
		err = DIContainer.Provide(repository)
		if err != nil {
			log.Fatalf("failed to provide repository: %v", err)
		}
	}

	// --- Service --- //
	services := []any{
		service.NewScheduleEditPermissionService,
	}

	for _, service := range services {
		err = DIContainer.Provide(service)
		if err != nil {
			log.Fatalf("failed to provide service: %v", err)
		}
	}

	// --- Facade --- //
	err = DIContainer.Provide(
		session_public_facade.NewSessionSaveFacade,
		dig.As(new(user_external.ISessionSaveFacade)),
	)
	if err != nil {
		log.Fatalf("failed to provide session save facade: %v", err)
	}

	err = DIContainer.Provide(
		session_public_facade.NewSessionDeleteFacade,
		dig.As(new(user_external.ISessionDeleteFacade)),
	)
	if err != nil {
		log.Fatalf("failed to provide session delete facade: %v", err)
	}

	err = DIContainer.Provide(
		role_public_facade.NewRoleGetFacade,
		dig.As(new(user_external.IRoleGetFacade)),
	)
	if err != nil {
		log.Fatalf("failed to provide role get facade: %v", err)
	}

	err = DIContainer.Provide(
		room_public_facade.NewRoomGetFacade,
		dig.As(new(schedule_external.IRoomGetFacade)),
	)
	if err != nil {
		log.Fatalf("failed to provide room get facade: %v", err)
	}

	err = DIContainer.Provide(
		room_public_facade.NewRoomExistsFacade,
		dig.As(new(schedule_external.IRoomExistsFacade)),
	)
	if err != nil {
		log.Fatalf("failed to provide room exists facade: %v", err)
	}

	err = DIContainer.Provide(
		lesson_public_facade.NewLessonGetFacade,
		dig.As(new(schedule_external.ILessonGetFacade)),
	)
	if err != nil {
		log.Fatalf("failed to provide lesson get facade: %v", err)
	}

	err = DIContainer.Provide(
		campus_public_facade.NewCampusExistFacade,
		dig.As(
			new(lesson_external.ICampusExistFacade),
			new(room_external.ICampusExistFacade),
			new(schedule_external.ICampusExistFacade),
		),
	)
	if err != nil {
		log.Fatalf("failed to provide campus exist facade: %v", err)
	}

	err = DIContainer.Provide(
		user_public_facade.NewUserGetFacade,
		dig.As(new(schedule_external.IUserGetFacade)),
	)
	if err != nil {
		log.Fatalf("failed to provide user get facade: %v", err)
	}

	// // --- Usecase --- //
	usecases := []any{
		lessonlist.NewLessonListQueryInteractor,
		roomlist.NewRoomListQueryInteractor,
		schedulelist.NewScheduleListQueryInteractor,
		mapper.NewScheduleItemEditOutputMapper,
		campus_interactor.NewCampusListInteractor,
		lesson_interactor.NewLessonAddInteractor,
		lesson_interactor.NewLessonEditInteractor,
		room_interactor.NewRoomEditInteractor,
		schedule_interactor.NewInvisibleRoomSaveInteractor,
		schedule_interactor.NewScheduleCreateInteractor,
		schedule_interactor.NewScheduleDeleteInteractor,
		schedule_interactor.NewScheduleDuplicateInteractor,
		schedule_interactor.NewScheduleGetInteractor,
		schedule_interactor.NewScheduleItemDivideInteractor,
		schedule_interactor.NewScheduleItemJoinInteractor,
		schedule_interactor.NewScheduleItemMoveInteractor,
		schedule_interactor.NewScheduleItemReturnListInteractor,
		schedule_interactor.NewScheduleItemShiftInteractor,
		schedule_interactor.NewScheduleSaveTitleInteractor,
		schedule_interactor.NewScheduleSaveInteractor,
		schedule_interactor.NewScheduleTimeEditEditInteractor,
		user_interactor.NewUserAddInteractor,
		user_interactor.NewUserDeleteInteractor,
		user_interactor.NewUserGetInteractor,
		user_interactor.NewUserListInteractor,
		user_interactor.NewUserLoginInteractor,
		user_interactor.NewUserUpdateInteractor,
		session_interactor.NewSessionGetInteractor,
	}

	for _, usecase := range usecases {
		err = DIContainer.Provide(usecase)
		if err != nil {
			log.Fatalf("failed to provide usecase: %v", err)
		}
	}

	// --- Handler --- //
	handlers := []any{
		session_handler.NewSessionGetHandler,
	}

	for _, controller := range handlers {
		err = DIContainer.Provide(controller)
		if err != nil {
			log.Fatalf("failed to provide handler: %v", err)
		}
	}

	// --- Controller --- //
	controllers := []any{
		campus_controller.NewCampusListController,
		lesson_controller.NewLessonAddController,
		lesson_controller.NewLessonEditController,
		lesson_controller.NewLessonListController,
		controller.NewLoginUserGetController,
		room_controller.NewRoomEditController,
		room_controller.NewRoomListController,
		schedule_controller.NewInvisibleRoomController,
		schedule_controller.NewScheduleCreateController,
		schedule_controller.NewScheduleDeleteController,
		schedule_controller.NewScheduleDuplicateController,
		schedule_controller.NewScheduleGetController,
		schedule_controller.NewScheduleItemDivideController,
		schedule_controller.NewScheduleItemJoinController,
		schedule_controller.NewScheduleItemMoveController,
		schedule_controller.NewScheduleItemReturnListController,
		schedule_controller.NewScheduleItemShiftController,
		schedule_controller.NewScheduleListController,
		schedule_controller.NewScheduleSaveController,
		schedule_controller.NewScheduleSaveTitleController,
		schedule_controller.NewScheduleTimeEditController,
		user_controller.NewUserAddController,
		user_controller.NewUserDeleteController,
		user_controller.NewUserGetController,
		user_controller.NewUserListController,
		user_controller.NewUserLoginController,
		user_controller.NewUserLogoutController,
		user_controller.NewUserUpdateController,
	}

	for _, controller := range controllers {
		err = DIContainer.Provide(controller)
		if err != nil {
			log.Fatalf("failed to provide controller: %v", err)
		}
	}

	// // --- Presenter --- //
	presenters := []any{
		campus_presenter.NewCampusListPresenter,
		lesson_presenter.NewLessonListPresenter,
		lesson_presenter.NewLessonAddPresenter,
		lesson_presenter.NewLessonEditPresenter,
		room_presenter.NewRoomListPresenter,
		room_presenter.NewRoomEditPresenter,
		schedule_presenter.NewScheduleListPresenter,
		schedule_presenter.NewInvisibleRoom,
		schedule_presenter.NewScheduleCreatePresenter,
		schedule_presenter.NewScheduleGet,
		schedule_presenter.NewScheduleItemEditPresenter,
		schedule_presenter.NewScheduleSaveTitlePresenter,
		schedule_presenter.NewScheduleSavePresenter,
		user_presenter.NewUserAddPresenter,
		user_presenter.NewUserDeletePresenter,
		user_presenter.NewUserGetPresenter,
		user_presenter.NewUserListPresenter,
		user_presenter.NewUserLoginPresenter,
		user_presenter.NewUpdateUserPresenter,
	}

	for _, presenter := range presenters {
		err = DIContainer.Provide(presenter)
		if err != nil {
			log.Fatalf("failed to provide presenter: %v", err)
		}
	}

	// --- Server --- //
	err = DIContainer.Provide(server.NewServer)
	if err != nil {
		log.Fatalln(err)
	}

	// --- Router --- //
	err = DIContainer.Provide(handler.NewRouter)
	if err != nil {
		log.Fatalln(err)
	}
}

func RunInjectedServer() error {

	var httpServer *http.Server
	var dbConn *sql.DB

	defer func() {
		if dbConn != nil {
			if err := dbConn.Close(); err != nil {
				log.Println("failed to close DB:", err)
			} else {
				log.Println("DB connection closed")
			}
		}
	}()

	errChan := make(chan error, 1)

	go func() {
		err := DIContainer.Invoke(func(e *echo.Echo, env envcfg.EnvConfig, mysql rdb.IMySQL) error {

			e.HTTPErrorHandler = func(err error, c echo.Context) {
				if he, ok := err.(*echo.HTTPError); ok {
					if he.Code == http.StatusNotFound && c.Path() != "" {

						// このルートはGETのみ許可
						c.Response().Header().Set(echo.HeaderAllow, http.MethodGet)
						_ = c.JSON(http.StatusMethodNotAllowed, map[string]string{
							"message": "Method Not Allowed",
						})
						return
					}
				}
				e.DefaultHTTPErrorHandler(err, c)
			}

			e.Routes()
			dbConn = mysql.GetConn()

			httpServer = &http.Server{
				Addr:    env.ServerBindAddress,
				Handler: e,
			}

			err := httpServer.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				return fmt.Errorf("server error: %w", err)
			}
			return nil
		})
		errChan <- err
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return err
	case <-quit:
		log.Println("graceful shutdown start")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if httpServer != nil {
			if err := httpServer.Shutdown(ctx); err != nil {
				return fmt.Errorf("http server graceful shutdown error: %w", err)
			}
		}

		log.Println("Server graceful shutdown")
		return nil
	}
}
