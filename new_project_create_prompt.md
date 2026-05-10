# Project Architecture Template

この Markdown は、新規プロジェクトで AI に読み込ませる実装規約・アーキテクチャテンプレートである。AI は本ドキュメントを最優先の実装ルールとして扱い、構成、依存方向、命名、エラー処理、テスト方針を可能な限りこのリポジトリと同じ思想で再現すること。

この Markdown とサンプルプロジェクトが同封されている場合、AI はサンプルプロジェクトの実装を必ず参照し、類似機能がある場合はその構成、命名、依存方向、エラー処理、ログ設計、テスト方針を模倣すること。参照すべきサンプルプロジェクトがどれか不明な場合、または複数存在する場合は、必ずユーザーに質問して確認すること。

本ドキュメントとサンプルプロジェクトの実装が矛盾する場合は、原則としてサンプルプロジェクトの実装を優先すること。ただし、どちらを優先すべきか判断できない場合は、推測せず必ずユーザーに質問すること。

AI は判断に迷う場合、既存ルールと本ドキュメントの実例を優先すること。実例がなく判断できない場合は推測で実装せず、必ず質問すること。

AI は実装前に、既存コードから類似実装を探索すること。既存実装に類似パターンが存在する場合、必ずその実装を最優先で模倣すること。

AI は既存実装に存在しない新規アーキテクチャパターン、ライブラリ、レイヤー分離、抽象化を勝手に導入してはならない。既存実装を最優先で模倣し、改善より一貫性を優先すること。

既存コードと異なる実装を行う場合は、既存実装では対応できない理由を説明した上で、必ず質問し、承認されるまで実装してはならない。新規ライブラリ追加、アーキテクチャ変更、レイヤ追加を行う場合も、既存実装では対応できない理由を説明した上で、必ず質問し、承認されるまで実装してはならない。

実装対象と関係のない既存コードを変更してはならない。import 整理、rename、formatter 差分、無関係な refactor を実装に混ぜてはならない。

改善提案は許可する。ただし、ユーザーが承認するまで勝手に実装してはならない。

## 1. Overview

このプロジェクトは、Go で実装された DDD + Clean Architecture ベースのモジュラモノリス構成バックエンドである。新規プロジェクトでも、`internal/modules/<module>` を境界づけられた業務モジュールとして扱い、各モジュール内を `adapter` / `usecase` / `domain` / `infrastructure` に分割すること。

設計思想は、DDD や Clean Architecture を厳密に再現することではなく、実務で運用しやすい範囲で戦術的 DDD と依存関係の統制を適用することである。README では、過剰抽象化を避けるために Port/Gateway を一部省略し、Echo の adapter 層利用や、infrastructure 層による Repository 実装を許容している。

この構成では、ドメインモデルと Value Object を中心に業務ルールを表現し、usecase/interactor が認可、入力検証、外部モジュール連携、Repository 呼び出し、トランザクションを調整する。HTTP 層は Echo の controller と presenter で構成し、DI は `go.uber.org/dig` を用いて `internal/injector/injector.go` に集約する。

## 2. Directory Structure

新規プロジェクトでは、以下の構成を基準にすること。

```text
.
├── main.go
├── main_test.go
├── go.mod
├── Dockerfile
├── docker-compose.yaml
├── build/
│   └── app/
│       └── Dockerfile.local
├── db/
│   ├── atlas.hcl
│   ├── lessonlink_schema.hcl
│   └── migrations/
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── entrypoint/
│   └── mysql/
│       └── conf.d/
│           └── my.cnf
├── internal/
│   ├── adapter/
│   │   ├── handler/
│   │   └── utility/
│   ├── configs/
│   ├── infrastructure/
│   │   └── server/
│   ├── injector/
│   ├── modules/
│   │   ├── campus/
│   │   ├── lesson/
│   │   ├── role/
│   │   ├── room/
│   │   ├── schedule/
│   │   ├── session/
│   │   └── user/
│   ├── pkg/
│   │   ├── constants/
│   │   ├── hash/
│   │   ├── log/
│   │   └── utility/
│   ├── platform/
│   │   ├── database/
│   │   ├── logger/
│   │   └── utility/
│   └── third_party/
│       └── sqlboiler/
└── testdata/
```

各ディレクトリの責務は以下の通り。

- `main.go`: アプリケーション起動エントリーポイント。ローカル環境の `.env` 読み込み、環境変数検証、DI 経由のサーバ起動を行う。
- `main_test.go`: API 全体を対象にした結合テスト。テスト DB 作成、Atlas migration、HTTP golden test、Schemathesis を実行する。
- `db/`: Atlas のスキーマ定義、マイグレーション、migration 設定を置く。
- `docs/`: Swag で生成された Swagger ファイルを置く。
- `internal/adapter/handler`: アプリ全体の Echo router、共通 middleware、Swagger 初期化を置く。
- `internal/adapter/utility`: Echo context など adapter 層に依存する共通処理を置く。
- `internal/configs`: `envconfig` による環境変数定義と読み込み処理を置く。
- `internal/infrastructure/server`: Echo サーバ、CORS、session store など HTTP サーバの共通初期化を置く。
- `internal/injector`: `dig` の DI 登録とサーバ起動処理を集約する。
- `internal/modules/<module>`: 業務モジュール単位の実装を置く。
- `internal/pkg`: ドメイン非依存の共通処理、定数、エラー、ハッシュなどを置く。
- `internal/platform`: DB、TxManager、logger などの共通基盤を置く。
- `internal/third_party`: SQLBoiler など外部ツールの設定や生成関連コードを置く。
- `testdata`: API golden test の request / response JSON を置く。

モジュール内は以下を基準にすること。

```text
internal/modules/<module>/
├── adapter/
│   ├── controller/
│   ├── handler/
│   └── presenter/
├── domain/
│   ├── model/
│   ├── repository/
│   ├── service/
│   └── vo/
├── infrastructure/
│   └── database/
│       └── rdb/
│           └── query/
├── usecase/
│   ├── external/
│   ├── interactor/
│   ├── mapper/
│   ├── port/
│   ├── public/
│   └── query/
```

存在しないサブディレクトリは無理に作らなくてよい。実例では `session` に `adapter/handler` があり、通常の HTTP controller は `adapter/controller` に置かれている。

## 3. Layer Architecture

各モジュールは以下の責務を持つ。

- `adapter/controller`: Echo の `echo.Context` を受け取り、path parameter、request body、session data を取得し、usecase input DTO に変換する。usecase の戻り値または presenter の結果を JSON レスポンスに変換する。
- `adapter/presenter`: usecase の結果または固定メッセージを HTTP response DTO に整形する。
- `adapter/handler`: controller ではない adapter 境界を置く。実例では session 取得 handler が存在する。
- `usecase/interactor`: Application Service としてユースケースを実行する。認可、VO 生成、外部 facade 呼び出し、Repository 呼び出し、トランザクション制御、ドメインモデル生成を担当する。
- `usecase/query`: 参照系ユースケースを置く。ドメインモデルを経由せず、query repository と DTO で一覧取得などを実装する実例がある。
- `usecase/external`: 他モジュールの public facade に依存するための interface を、利用側モジュール内に定義する。
- `usecase/public`: 他モジュールに公開する facade 実装を置く。
- `usecase/mapper`: usecase の出力変換ロジックを置く。
- `usecase/port`: 出力ポートなど usecase 境界の型を置く。実例は限定的であり、README では Port は全面採用していないと明記されている。
- `domain/model`: 集約、エンティティ、ドメインモデルを置く。フィールドは非公開にし、コンストラクタと getter で扱う。
- `domain/vo`: Value Object を置く。入力値の trim、長さ検証、範囲検証などの不変条件をコンストラクタで保証する。
- `domain/repository`: Repository interface を定義する。
- `domain/service`: 単一モデルに置きづらいドメインサービスを置く。
- `infrastructure/database/rdb`: SQLBoiler DTO を使った Repository 実装を置く。
- `infrastructure/database/rdb/query`: 参照系 query repository 実装を置く。

依存方向は原則として adapter -> usecase -> domain とし、infrastructure は domain repository interface または usecase query repository interface を実装する。DI により runtime で concrete implementation を注入する。

README 上の意図的な逸脱もルールとして扱うこと。

- Echo は adapter/controller 層で利用してよい。
- Repository 実装は gateway 層を挟まず、infrastructure 層が直接実装してよい。
- interactor は戻り値を持ってよい。
- トランザクションは集約単位ではなく usecase 層で扱ってよい。
- DB の外部キーは、実務上のデータ破壊防止として集約間にも設定してよい。ただし実装は外部キーだけに依存せず、usecase 側で存在確認などを行う。

## 4. File and Naming Rules

ファイル名は snake_case を基本にし、役割を suffix で表すこと。

- controller: `<feature>_controller.go`
- presenter: `<feature>_presenter.go`
- interactor/usecase: `<feature>_usecase.go`
- repository interface: `<resource>_repository.go`
- RDB repository implementation: `<resource>.go`
- query interactor: `<feature>_query_Interactor.go` の実例があるが、Go の一般的な命名としては小文字の `<feature>_query_interactor.go` を優先してよい。不明な場合は既存モジュールの近い命名に合わせる。
- Value Object: `<value_name>.go`
- root model: `root_<aggregate>_model.go`
- child model: `<name>_model.go`
- domain service: `<ServiceName>.go` の実例があるが、Go ファイル名としては snake_case を優先する。

パッケージ名はディレクトリ末尾に合わせること。

- `adapter/controller` は `package controller`
- `adapter/presenter` は `package presenter`
- `usecase/interactor` は `package interactor`
- `usecase/external` は `package external`
- `usecase/public` は `package public`
- `domain/vo` は `package vo`
- `domain/repository` は `package repository`
- `infrastructure/database/rdb` は `package rdb`

型名は以下の規則に従うこと。

- Controller interface: `I<Feature>Controller`
- Controller struct: `<Feature>Controller`
- Controller constructor: `New<Feature>Controller(...) I<Feature>Controller`
- Presenter interface: `I<Feature>Presenter`
- Presenter struct: `<Feature>Presenter`
- Presenter constructor: `New<Feature>Presenter(...) I<Feature>Presenter`
- Usecase input port interface: `I<Feature>InputPort`
- Interactor struct: `<Feature>Interactor`
- Interactor constructor: `New<Feature>Interactor(...) I<Feature>InputPort`
- Request DTO: `<Feature>RequestData`
- Input DTO: `<Feature>InputDTO`
- Query output DTO: `<Feature>QueryOutput`
- Query item DTO: `Query<Resource>DTO`
- Facade interface: `I<Resource><Action>Facade`
- Facade struct: `<Resource><Action>Facade`
- Repository interface: `<Resource>Repository`
- Root aggregate: `Root<Resource>Model`
- Slice helper: `Root<Resource>ModelSlice`
- VO constructor: `New<ValueName>(...) (<ValueName>, error)`
- VO value getter: `Value() <primitive>`

JSON field は snake_case を使うこと。例: `lesson_name`, `schedule_id`。

## 5. Implementation Rules

- Controller は Echo の `echo.Context` を受け取り、`Execute(c echo.Context) error` を公開する。
- Controller は request body を `c.Bind(&requestData)` で bind し、bind 失敗時は `400` と `{"msg": "リクエスト形式が不正です"}` の形式で返す。
- Controller は path parameter の必須値が空の場合、`400` と `{"msg": "...が不正です"}` の形式で返す。
- 認証済み API では `session_util.GetSessionData(c)` から user ID と role key を取得する。
- usecase は primitive input を直接扱い続けず、冒頭で VO に変換する。
- VO 生成に失敗した場合は `log.WrapErrorWithStackTraceBadRequest(err)` で返す。
- 権限チェックは usecase 内で実施する。実例では `role.IsOwner()` を満たさない場合に forbidden error を返す。
- 他モジュールの存在確認は、対象モジュールの `public` facade を利用側の `usecase/external` interface 経由で呼び出す。
- 集約を作るときは `NewRoot<Resource>Model(...)` を使い、フィールドを直接公開しない。
- 複数 VO の組み立てでは `errors.Join` と `utility.SetVOConstructor` を使って検証エラーをまとめる実例がある。
- Repository は `context.Context` を第一引数に取る。
- 更新系 Repository は `*sql.Tx` を受け取り、usecase の `txManager.Do(ctx, func(tx *sql.Tx) error { ... })` 内で呼び出す。
- 参照系 query repository は必要に応じて domain model を経由せず、query DTO を直接返してよい。
- DB access は SQLBoiler の generated DTO を使う。
- SQLBoiler generated files は `internal/platform/database/rdb/dto` に置く。
- HTTP response は `map[string]string{"msg": ...}` または presenter response struct を返す。
- Swagger annotation は controller の `Execute` 直前に記述する。
- 起動時の依存登録は `internal/injector/injector.go` に集約し、config、logger、repository、service、facade、usecase、handler、controller、presenter、server、router の順に登録する実例に合わせる。
- Graceful shutdown は `SIGINT` / `SIGTERM` を受け、30 秒 timeout の `http.Server.Shutdown` を行う。
- DB connection は `RunInjectedServer` の defer で close する。

## 6. Dependency Rules

許可される依存方向は以下。

- `main.go` -> `internal/configs`, `internal/injector`
- `internal/injector` -> すべての concrete constructor
- `internal/adapter/handler` -> server、configs、logger、各 controller、session handler
- `adapter/controller` -> usecase/interactor interface、presenter、logger、adapter utility、Echo
- `adapter/presenter` -> 原則として標準型または usecase output DTO
- `usecase/interactor` -> domain model、domain repository interface、domain VO、usecase external interface、platform database TxManager、pkg/log
- `usecase/public` -> 自モジュールの domain repository、domain VO、pkg/log
- `usecase/query` -> query repository interface、pkg/log
- `domain/model` -> domain/vo、および必要な汎用ライブラリ
- `domain/vo` -> pkg/log、標準ライブラリ
- `domain/repository` -> domain model、domain VO、`database/sql` の `*sql.Tx`
- `infrastructure/database/rdb` -> domain repository interface、domain model、domain VO、platform database/rdb、SQLBoiler dto、pkg/log
- `infrastructure/database/rdb/query` -> usecase/query repository interface、query DTO、platform database/rdb、SQLBoiler dto、pkg/log
- `platform/database/rdb` -> configs、`database/sql`、MySQL driver
- `platform/logger` -> configs、pkg/log、Echo

禁止または避ける依存は以下。

- domain 層から adapter、infrastructure、configs、Echo、SQLBoiler dto に依存してはいけない。
- VO や domain model から DB、HTTP、環境変数、logger 実装に依存してはいけない。
- usecase/interactor から別モジュールの concrete implementation を直接 import してはいけない。利用側の `usecase/external` interface と提供側の `usecase/public` facade を使う。
- controller から infrastructure repository を直接呼び出してはいけない。
- router 以外に URL routing を分散させてはいけない。
- DI 登録を各モジュールに分散させてはいけない。実例では `internal/injector/injector.go` に集約されている。
- SQLBoiler generated DTO を domain model として扱ってはいけない。infrastructure で domain model に変換する。

## 7. Configuration Rules

環境変数は `internal/configs/env.go` の `EnvConfig` に追加し、`required:"true"` を付与すること。読み込みは `envconfig.Process("", &config)` で行い、失敗時は `log.Fatalln(err)` で起動を止める。

既存の環境変数は以下。

- `ENVIRONMENT`
- `LOG_ERROR_REQUEST_DUMP`
- `SERVER_BIND_ADDRESS`
- `DB_ADDRESS`
- `DB_USER`
- `DB_PASSWORD`
- `LOG_LEVEL`
- `DB_NAME`
- `SESSION_NAME`
- `SESSION_SECRET_KEY`
- `TEST_DB_NAME`。`local_test` 時に `DB_NAME` へ反映される。`EnvConfig` には定義されていない。
- `LOCAL_TEST`。cookie secure/samesite やテストで参照される。`EnvConfig` には定義されていない。

`.env` は `internal/configs/.env` に置き、サンプルは `internal/configs/.env.example` に置く。`main.go` は `ENVIRONMENT` が `local` または `local_test` の場合だけ `godotenv.Overload("./internal/configs/.env")` を実行する。

DB 接続は `internal/platform/database/rdb/mysql.go` で初期化する。DSN は `parseTime=true&loc=Asia%2FTokyo` を付与する。connection pool は `MaxIdleConnection=10`, `MaxOpenConnection=20`, `ConnectionMaxLifeTime=10s` を既定値にする。

HTTP server は `internal/infrastructure/server/server.go` で Echo を生成し、session middleware と CORS を設定する。local 環境では `http://localhost:3031` を CORS origin に追加する。Swagger UI は `ENVIRONMENT=local` の場合だけ `/swagger/*` に公開する。

DB migration は Atlas を使い、`db/atlas.hcl` の `dev` env に migration dir `file://db/migrations` を設定する。Docker Compose では MySQL 8.0、phpMyAdmin、app-dev を用意する。

## 8. Error Handling Rules

エラーは `internal/pkg/log` の独自 `Error` 型で扱うこと。

```go
type Error struct {
    StatusCode int
    Message    string
    StackTrace string
}
```

新規エラーは `log.Errorf(...)` または標準 error を作成し、以下の wrapper で HTTP status と stack trace を付与する。

- `WrapErrorWithStackTrace(err)`
- `WrapErrorWithStackTraceInternalServerError(err)`
- `WrapErrorWithStackTraceBadRequest(err)`
- `WrapErrorWithStackTraceNotFound(err)`
- `WrapErrorWithStackTraceUnauthorized(err)`
- `WrapErrorWithStackTraceConflict(err)`
- `WrapErrorWithStackTraceForbidden(err)`

status code と log severity は `LogLevelMap` / `LogSeverityMap` に定義する。既定の対応は以下。

- `500`: `ERROR`
- `400`: `WARNING`
- `404`: `WARNING`
- `401`: `WARNING`
- `409`: `WARNING`
- `403`: `WARNING`

Controller では usecase から返った error を `logger.WriteErrLog(c, err)` に渡し、返却された `status` と `msg` を `{"msg": msg}` 形式で JSON response に変換する。

local 環境では error message、stack trace、transaction ID を標準出力に表示する。local 以外では `slog.NewJSONHandler(os.Stdout)` で JSON log を出力し、`transactionID`, `severity`, `stacktrace` を含める。

`LOG_ERROR_REQUEST_DUMP=true` の場合、router で Echo `BodyDump` middleware を有効にし、`status >= 400` かつ `401` 以外の response に対して、transaction ID、user ID、URL、method、query params、status、request body を JSON log として出力する。

不明な点: エラーメッセージの国際化、多言語対応、クライアント向けエラーコード体系はリポジトリ内に実例がない。

## 9. Testing Rules

このリポジトリでは、モジュール単位の `*_test.go` は確認できない。主なテストは `main_test.go` に集約された API 結合テストである。

新規プロジェクトでも、少なくとも以下のテスト方針を再現すること。

- `TestMain` でテスト用 DB を作成する。
- テスト DB 名は `lessonlink_test` のような専用名を使い、テスト終了時に drop する。
- Atlas migration を `atlas migrate apply --dir file://db/migrations --allow-dirty` で適用する。
- `Test_main` で `ENVIRONMENT=local_test`, `LOCAL_TEST=true`, `TEST_DB_NAME=<test db>` を設定する。
- goroutine で `main()` を起動し、HTTP API に対して結合テストを行う。
- `httpexpect` を使って API request を送信する。
- 最初に `/api/user/login` を実行し、cookie jar を共有して認証必須 API をテストする。
- request / expected response は `testdata/<feature>/<case>_case_req.json` と `<case>_case_res.json` に分ける。
- expected response JSON には `http_status` を含める。
- expected response JSON の `_ignore` 配列に dot path を指定し、比較から除外できるようにする。
- response 配列の順序が本質でない場合は normalize して比較する。
- Swagger schema に対して Schemathesis を実行し、`status_code_conformance` と `not_a_server_error` を検証する。

実例上、unit test、mock 生成、Repository mock、domain model 単体テストの規約は不明である。新規追加する場合は、このリポジトリの既存方針を壊さない範囲で、domain と usecase を優先して小さく追加すること。

## 10. Code Generation Instructions

新規機能を追加するとき、AI は以下の手順に従うこと。

1. 追加対象を既存モジュールに入れるか、新規 `internal/modules/<module>` を作るかを決める。
2. 更新系機能の場合、`domain/vo` に必要な Value Object を作り、コンストラクタで trim、長さ、範囲、必須チェックを行う。
3. 集約が必要な場合、`domain/model/<aggregate>/root_<aggregate>_model.go` を作り、非公開フィールド、constructor、getter、振る舞いメソッドを定義する。
4. 永続化が必要な場合、`domain/repository/<resource>_repository.go` に interface を定義する。
5. RDB 実装を `infrastructure/database/rdb/<resource>.go` に作り、SQLBoiler dto と domain model の相互変換を `toModel` / `toDTO` に閉じ込める。
6. 参照系一覧 API の場合は `usecase/query/<feature>` に query interactor、query repository interface、query DTO を置き、`infrastructure/database/rdb/query/<resource>` に実装を置く。
7. 他モジュールの情報が必要な場合、利用側に `usecase/external/I...Facade` を定義し、提供側に `usecase/public/...Facade` を実装する。
8. usecase/interactor を作り、interface `I<Feature>InputPort`、input DTO、constructor、`Execute(ctx, ...)` を定義する。
9. usecase 内で primitive input を VO に変換し、権限チェック、存在確認、重複チェック、トランザクション、Repository 呼び出しを行う。
10. HTTP controller を `adapter/controller/<feature>_controller.go` に作り、Echo request を input DTO に変換する。
11. response presenter を `adapter/presenter/<feature>_presenter.go` に作る。
12. Swagger annotation を controller に追加する。
13. `internal/injector/injector.go` に repository、facade、usecase、controller、presenter を登録する。
14. `internal/adapter/handler/router.go` に route を追加する。認証必須 API は `auth := api.Group("")` 配下に追加する。
15. DB 変更がある場合は `db/lessonlink_schema.hcl` と `db/migrations` を更新し、SQLBoiler DTO を再生成する。
16. API テスト用に `testdata/<feature>` 配下へ `*_case_req.json` と `*_case_res.json` を追加し、`main_test.go` の golden test 呼び出しを追加する。
17. Swagger を再生成し、Schemathesis の対象として矛盾がないことを確認する。

新規 controller の雛形は以下に合わせること。

```go
type (
    IFeatureController interface {
        Execute(c echo.Context) error
    }

    FeatureController struct {
        inputPort interactor.IFeatureInputPort
        presenter presenter.IFeaturePresenter
        logger    logger.ILogWriter
    }
)
```

新規 interactor の雛形は以下に合わせること。

```go
type (
    IFeatureInputPort interface {
        Execute(ctx context.Context, input FeatureInputDTO) error
    }
)

type FeatureInteractor struct {
    txManager database.TxManager
    repositoryX repository.XRepository
}

func NewFeatureInteractor(...) IFeatureInputPort {
    return &FeatureInteractor{...}
}
```

## 11. Do / Don't

Do:

- モジュールを `internal/modules/<module>` 単位で分ける。
- モジュール内に `adapter` / `usecase` / `domain` / `infrastructure` を置く。
- ドメインルールは VO、domain model、domain service に寄せる。
- usecase で認可、VO 生成、外部 facade、Repository、TxManager を調整する。
- 他モジュール連携は `public` facade と利用側 `external` interface で行う。
- DB 更新は usecase の `txManager.Do` 内で行う。
- SQLBoiler dto と domain model の変換を infrastructure に閉じ込める。
- エラーは `internal/pkg/log` の wrapper で status code と stack trace を付ける。
- Controller では `logger.WriteErrLog` の戻り値を HTTP response に変換する。
- DI 登録は `internal/injector/injector.go` に追加する。
- request / response の JSON field は snake_case にする。
- Swagger annotation を controller に書く。
- local のみ Swagger UI を公開する。
- 結合テストは `testdata` の golden file と `httpexpect` を使う。

Don't:

- domain 層から Echo、SQLBoiler、DB、環境変数、logger 実装に依存させない。
- controller から repository や SQLBoiler dto を直接呼ばない。
- usecase から他モジュールの infrastructure や concrete facade を直接 import しない。
- SQLBoiler dto を domain model として扱わない。
- VO を作らず primitive のままドメインルールを散らさない。
- DI 登録を各ファイルに分散させない。
- route 定義を controller に分散させない。
- エラーをそのまま `c.JSON` に返さない。usecase error は logger を通す。
- DB の外部キーだけに整合性保証を任せない。usecase 側でも存在確認する。
- DDD/Clean Architecture のためだけに gateway や Port を過剰に増やさない。
- README にある意図的な実務上の逸脱を無視して、理論上の純粋性だけで構成を変更しない。

不明:

- unit test と mock の詳細規約は不明。
- 外部 API client の実装規約は、リポジトリ内に具体例が見当たらないため不明。
- 非同期処理、ジョブ、メッセージキューの規約は不明。
- リリース環境の本番設定値、CI/CD、デプロイ手順は不明。
