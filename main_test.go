package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/samber/lo"

	"github.com/joho/godotenv"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
)

const LOCAL_TEST_BASE_POINT = "http://localhost:8080/api"

var sharedClient *http.Client

func Test_main(t *testing.T) {

	t.Setenv("ENVIRONMENT", "local_test")
	t.Setenv("TEST_DB_NAME", TEST_DATABASE_NAME)
	t.Setenv("LOCAL_TEST", "true")

	// サーバーを起動 mainを実行する
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		main()
	}()

	// サーバーが起動するまで待機
	wg.Wait()
	time.Sleep(2 * time.Second)

	/////////////////////

	jar, _ := cookiejar.New(nil)
	sharedClient = &http.Client{
		Jar: jar,
	}

	// ログイン実行
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  LOCAL_TEST_BASE_POINT,
		Client:   sharedClient,
		Reporter: httpexpect.NewAssertReporter(nil), // TestMain内なので一旦 nil で初期化
	})

	e.POST("/user/login").
		WithJSON(map[string]string{"user_name": "admin@admin.com", "password": "admin"}).
		Expect().
		Status(http.StatusOK)

	u, _ := url.Parse(LOCAL_TEST_BASE_POINT)
	cookies := sharedClient.Jar.Cookies(u)

	var cookieHeader string
	for _, c := range cookies {
		cookieHeader += c.Name + "=" + c.Value + "; "
	}

	// 講座登録のテスト
	runGolden(t, "/lesson/shibuya", "POST", false, "lesson")

	// 教室の登録テスト
	runGolden(t, "/room/shibuya/edit", "POST", false, "room")

	// スケジュール作成テスト
	runGolden(t, "/schedule/create/shibuya", "POST", false, "schedule/create")
	// runGolden(t, "/schedule/create/shinagawa", "POST", false)

	// スケジュール一覧取得
	runGolden(t, "/schedule/list/shibuya", "GET", true, "schedule/list")
	// runGolden(t, "/schedule/list/shinagawa", "GET", false)

	// スケジュール取得
	runGolden(t, "/schedule/1", "GET", false, "schedule/get")

	// スケジュール編集 アイテム移動 リストからルームへ
	runGolden(t, "/schedule/1/item-move", "POST", false, "schedule/item-move")

	// スケジュール編集 アイテム移動 ルームからリストへ
	runGolden(t, "/schedule/1/item-return-list", "POST", false, "schedule/item-return-list")

	// スケジュール編集 アイテム分割
	runGolden(t, "/schedule/1/item-divide", "POST", false, "schedule/item-divide")

	// アイテム結合用データ更新
	_, err := db.Exec("update tbl_schedule_items set identifier = 'identifier_lesson_1_from' where schedule_id = 1 and history_index = 5 and lesson_id = 1 and identifier != 'identifier_lesson_1' ")
	if err != nil {
		panic(err)
	}

	// スケジュール編集 アイテム結合
	runGolden(t, "/schedule/1/item-join", "POST", false, "schedule/item-join")

	// アイテムシフトを行うためにアイテムを配置
	runGolden(t, "/schedule/1/item-move", "POST", false, "schedule/item-shift-ready")

	// スケジュール編集 アイテムシフト
	runGolden(t, "/schedule/1/item-shift", "POST", false, "schedule/item-shift")

	// スケジュール編集 タイトル変更
	runGolden(t, "/schedule/1/title", "PATCH", false, "schedule/title")

	// スケジュール複製
	runGolden(t, "/schedule/1/duplicate", "POST", false, "schedule/duplicate")

	// スケジュール編集 ルーム非表示設定
	runGolden(t, "/schedule/1/room/invisible", "PUT", false, "schedule/room/invisible")

	// スケジュール編集 スケジュール時間変更
	runGolden(t, "/schedule/3/time", "PATCH", false, "schedule/time")

	// スケジュール削除
	runGolden(t, "/schedule/1", "DELETE", false, "schedule/delete")

	// schemathesisテスト
	runSchemathesis(t, cookieHeader)
	runSchemathesisOne(t, cookieHeader, "DELETE /schedule/{schedule_id}")
	runSchemathesisOne(t, cookieHeader, "DELETE /user/{userid}")
	runSchemathesisOne(t, cookieHeader, "POST /user/login")
	runSchemathesisOne(t, cookieHeader, "POST /user/logout")
}

func runGolden(t *testing.T, apiPath string, method string, ordered bool, testDataDir string) {

	dir := "testdata/" + testDataDir
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir failed: %v", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), "case_req.json") {
			continue
		}

		reqBody, err := os.ReadFile(dir + "/" + file.Name())
		if err != nil {
			t.Fatalf("read req failed: %v", err)
		}

		var req map[string]any
		if err := json.Unmarshal(reqBody, &req); err != nil {
			t.Fatalf("invalid req json: %v", err)
		}

		comment, _ := req["comment"].(string)

		t.Run(
			fmt.Sprintf("%s|%s|%s|%s", apiPath, method, fmt.Sprintf("%s/%s", dir, file.Name()), comment),
			func(t *testing.T) {

				e := httpexpect.WithConfig(httpexpect.Config{
					BaseURL:  LOCAL_TEST_BASE_POINT,
					Client:   sharedClient,
					Reporter: httpexpect.NewAssertReporter(t),
				})

				resPath := dir + "/" + strings.Replace(file.Name(), "_req.json", "_res.json", 1)
				resBody, err := os.ReadFile(resPath)
				if err != nil {
					t.Fatalf("read res failed: %v", err)
				}

				var exp map[string]any
				if err := json.Unmarshal(resBody, &exp); err != nil {
					t.Fatalf("invalid res json: %v", err)
				}

				status := int(exp["http_status"].(float64))
				ignore := toStringSlice(exp["_ignore"])

				delete(exp, "http_status")
				delete(exp, "_ignore")

				resp := e.Request(method, apiPath).
					WithHeader("Content-Type", "application/json").
					WithBytes(reqBody).
					Expect().
					Status(status)

				var act any

				if hasNoBody(status) {
					act = map[string]any{}
				} else {
					act = resp.JSON().Raw()
				}

				applyIgnore(act, ignore)
				applyIgnore(exp, ignore)

				if !ordered {
					normalize(act)
					normalize(exp)
				}

				if !reflect.DeepEqual(act, exp) {

					ab, _ := json.MarshalIndent(act, "", "  ")
					eb, _ := json.MarshalIndent(exp, "", "  ")

					t.Fatalf(
						"response mismatch\nactual:\n%s\nexpect:\n%s",
						ab, eb,
					)
				}
			})
	}
}

func runSchemathesis(t *testing.T, cookieHeader string) {

	cmd := exec.Command(
		"/opt/venv/bin/schemathesis",
		"run",
		"docs/swagger.yaml",
		"--url=http://localhost:8080/api",
		"--header", "Cookie: "+cookieHeader,
		"--exclude-name", "DELETE /schedule/{schedule_id}",
		"--exclude-name", "POST /user/login",
		"--exclude-name", "POST /user/logout",
		"--exclude-name", "DELETE /user/{userid}",
		"--checks=status_code_conformance,not_a_server_error",
		"--exclude-checks=unsupported_method",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("schemathesis failed: %v", err)
	}

	/*
		/opt/venv/bin/schemathesis run docs/swagger.yaml --url=http://localhost:8080/api --header "Cookie: $COOKIE" --include-name "GET /schedule/list/{campus}" --checks=status_code_conformance,not_a_server_error
	*/
}

func runSchemathesisOne(t *testing.T, cookieHeader, operation string) {
	cmd := exec.Command(
		"/opt/venv/bin/schemathesis",
		"run",
		"docs/swagger.yaml",
		"--url=http://localhost:8080/api",
		"--header", "Cookie: "+cookieHeader,
		"--include-name", operation,
		"--checks=status_code_conformance,not_a_server_error",
		"--exclude-checks=unsupported_method",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("schemathesis failed (%s): %v", operation, err)
	}
}

func hasNoBody(status int) bool {
	if status >= 100 && status < 200 {
		return true
	}
	var noBodyStatuses = []int{
		204,
		205,
		304,
	}
	return lo.Contains(noBodyStatuses, status)
}

func applyIgnore(v any, ignores []string) {
	for _, ig := range ignores {
		deleteByPath(v, strings.Split(ig, "."))
	}
}

func deleteByPath(v any, path []string) {
	if len(path) == 0 {
		return
	}

	switch cur := v.(type) {
	case map[string]any:
		key := path[0]
		if len(path) == 1 {
			delete(cur, key)
			return
		}
		if next, ok := cur[key]; ok {
			deleteByPath(next, path[1:])
		}

	case []any:
		if path[0] == "[]" {
			for _, e := range cur {
				deleteByPath(e, path[1:])
			}
		}
	}
}

func normalize(v any) {
	switch x := v.(type) {
	case map[string]any:
		for _, vv := range x {
			normalize(vv)
		}
	case []any:
		for _, e := range x {
			normalize(e)
		}
		sort.Slice(x, func(i, j int) bool {
			bi, _ := json.Marshal(x[i])
			bj, _ := json.Marshal(x[j])
			return string(bi) < string(bj)
		})
	}
}

func toStringSlice(v any) []string {
	a, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(a))
	for _, e := range a {
		if s, ok := e.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

const TEST_DATABASE_NAME = "lessonlink_test"

var db *sql.DB

func TestMain(m *testing.M) {

	environment := os.Getenv("ENVIRONMENT")
	if environment == "local" {
		err := godotenv.Overload("./internal/configs/.env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	// 環境変数のチェック
	configs.LoadConfig()

	dbAddress := os.Getenv("DB_ADDRESS")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")

	// DBに接続し新しくスキーマを作成する
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?parseTime=true&loc=Asia%%2FTokyo", dbUser, dbPass, dbAddress)

	// データベースに接続
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// 接続を確認
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// テスト用スキーマを作成
	db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", TEST_DATABASE_NAME))
	createDatabaseSQL := fmt.Sprintf("CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;", TEST_DATABASE_NAME)
	_, err = db.Exec(createDatabaseSQL)
	if err != nil {
		log.Fatalf("Failed to execute CREATE DATABASE: %v", err)
	}

	result := 1

	// 最後に削除する
	defer func() {
		dropSQL := fmt.Sprintf("drop database %s", TEST_DATABASE_NAME)
		_, err = db.Exec(dropSQL)
		if err != nil {
			fmt.Printf("Failed to execute drop database: %v", err)
		}
		os.Exit(result)
	}()

	_, err = db.Exec(fmt.Sprintf("use %s", TEST_DATABASE_NAME))
	if err != nil {
		log.Printf("Failed to execute DDL: %v", err)
		return
	}

	// マイグレーションを実行する
	cmd := exec.Command(
		"atlas",
		"migrate",
		"apply",
		"--url", fmt.Sprintf("mysql://%s:%s@%s:3306/%s", dbUser, dbPass, dbAddress, TEST_DATABASE_NAME),
		"--dir", "file://db/migrations",
		"--allow-dirty",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// テストを実行
	result = m.Run()
}
