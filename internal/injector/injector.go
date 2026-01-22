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
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/controller"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/handler"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/presenter"
	envcfg "github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/service"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb"
	campusRepo "github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb/query/campus"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb/query/lesson"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb/query/room"
	scheduleQueryRepository "github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb/query/schedule"
	logger "github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/logger"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/server"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/mapper"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/query/lessonlist"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/query/roomlist"
	schedulelist "github.com/typedef-tokyo/lessonlink-backend/internal/usecase/query/schedule_list"
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
		rdb.NewConfig,
		rdb.NewMySQL,
		rdb.NewTxManager,
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
		rdb.NewCampusRepository,
		rdb.NewLessonRepository,
		rdb.NewRoleRepository,
		rdb.NewRoomRepository,
		rdb.NewScheduleInvisibleRoomRepository,
		rdb.NewScheduleRepository,
		rdb.NewSessionRepository,
		rdb.NewUserRepository,
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

	// // --- Usecase --- //
	usecases := []any{
		lessonlist.NewLessonListQueryInteractor,
		roomlist.NewRoomListQueryInteractor,
		schedulelist.NewScheduleListQueryInteractor,
		mapper.NewScheduleItemEditOutputMapper,
		usecase.NewCampusListInteractor,
		usecase.NewInvisibleRoomSaveInteractor,
		usecase.NewLessonAddInteractor,
		usecase.NewLessonEditInteractor,
		usecase.NewRoomEditInteractor,
		usecase.NewScheduleCreateInteractor,
		usecase.NewScheduleDeleteInteractor,
		usecase.NewScheduleDuplicateInteractor,
		usecase.NewScheduleGetInteractor,
		usecase.NewScheduleItemDivideInteractor,
		usecase.NewScheduleItemJoinInteractor,
		usecase.NewScheduleItemMoveInteractor,
		usecase.NewScheduleItemReturnListInteractor,
		usecase.NewScheduleItemShiftInteractor,
		usecase.NewScheduleSaveTitleInteractor,
		usecase.NewScheduleSaveInteractor,
		usecase.NewScheduleTimeEditEditInteractor,
		usecase.NewUserAddInteractor,
		usecase.NewUserDeleteInteractor,
		usecase.NewUserGetInteractor,
		usecase.NewUserListInteractor,
		usecase.NewUserLoginInteractor,
		usecase.NewUserUpdateInteractor,
	}

	for _, usecase := range usecases {
		err = DIContainer.Provide(usecase)
		if err != nil {
			log.Fatalf("failed to provide usecase: %v", err)
		}
	}

	// --- Controller --- //
	controllers := []any{
		controller.NewCampusListController,
		controller.NewInvisibleRoomController,
		controller.NewLessonAddController,
		controller.NewLessonEditController,
		controller.NewLessonListController,
		controller.NewLoginUserGetController,
		controller.NewRoomEditController,
		controller.NewRoomListController,
		controller.NewScheduleCreateController,
		controller.NewScheduleDeleteController,
		controller.NewScheduleDuplicateController,
		controller.NewScheduleGetController,
		controller.NewScheduleItemDivideController,
		controller.NewScheduleItemJoinController,
		controller.NewScheduleItemMoveController,
		controller.NewScheduleItemReturnListController,
		controller.NewScheduleItemShiftController,
		controller.NewScheduleListController,
		controller.NewScheduleSaveController,
		controller.NewScheduleSaveTitleController,
		controller.NewScheduleTimeEditController,
		controller.NewUserAddController,
		controller.NewUserDeleteController,
		controller.NewUserGetController,
		controller.NewUserListController,
		controller.NewUserLoginController,
		controller.NewUserLogoutController,
		controller.NewUserUpdateController,
	}

	for _, controller := range controllers {
		err = DIContainer.Provide(controller)
		if err != nil {
			log.Fatalf("failed to provide controller: %v", err)
		}
	}

	// // --- Presenter --- //
	presenters := []any{
		presenter.NewLessonListPresenter,
		presenter.NewRoomListPresenter,
		presenter.NewScheduleListPresenter,
		presenter.NewCampusListPresenter,
		presenter.NewInvisibleRoom,
		presenter.NewLessonAddPresenter,
		presenter.NewLessonEditPresenter,
		presenter.NewRoomEditPresenter,
		presenter.NewScheduleCreatePresenter,
		presenter.NewScheduleGet,
		presenter.NewScheduleItemEditPresenter,
		presenter.NewScheduleSaveTitlePresenter,
		presenter.NewScheduleSavePresenter,
		presenter.NewUserAddPresenter,
		presenter.NewUserDeletePresenter,
		presenter.NewUserGetPresenter,
		presenter.NewUserListPresenter,
		presenter.NewUserLoginPresenter,
		presenter.NewUpdateUserPresenter,
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
