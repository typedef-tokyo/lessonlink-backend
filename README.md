## 1. プロジェクト概要

本プロジェクトは、Go を用いて実装した **DDD（Domain-Driven Design）+ Clean Architecture ベースのモジュラモノリス構成バックエンドアプリケーション**です。  
対象ドメインは「講座と教室のスケジュール作成」というシンプルな領域ですが、**技術的関心事の分離とドメインモデル中心の設計を実践すること**を目的としています。

本プロジェクトで重視しているポイントは以下です。

- **戦術的 DDD パターン**（Entity / ValueObject / Aggregate / Repository / DomainService）の適用  
- ~~**Clean Architecture** による依存関係の統制と疎結合な構造~~  
- ~~**ドメインロジックの入出力層からの完全分離**~~  
- ~~**Application Service** によるユースケース単位の振る舞いの定義~~
- **モジュール単位（`internal/modules/*`）での責務分離**  
- **モジュール内レイヤ分割（adapter / usecase / domain / infrastructure）**  
- **Application Service（usecase）によるユースケース単位の振る舞い定義**

このプロジェクトは、DDD と Clean Architecture の“理論を鵜呑みにしない”視点から、**実務で成立するバランスを取った Go バックエンド**として設計されています。
また、本プロジェクトには **DDD や Clean Architecture の基本に忠実ではない箇所も意図的に存在**しており、実務での過剰抽象化や運用負荷を避けるために、必要十分な設計へ調整しています。

## 2. 戦術的 DDD 要素

本プロジェクトでは、DDD の戦術的パターンを実務向けに再構成した形で採用している。  
各要素は過度な抽象化を避け、**現場で運用可能な責務分割**を重視して構成している。

- **Entity**  
  ・識別子を持ち、ライフサイクルと状態を管理するドメインオブジェクト。

- **Value Object**  
  ・値としての等価性、不変性、完全性を重視するオブジェクト。

- **Aggregate**  
  ・整合性を一貫して保証する単位として定義されたモデル。

- **Domain Service**  
  ・オブジェクト単独では表現できないドメイン固有の振る舞いを扱う。

- **Repository**  
  ・Aggregate の永続化を抽象化するインターフェース。

- **Application Service**  
  ・ユースケースを実行し、入出力や永続化を調整するアプリケーション層。

これらの要素は、基本的には理論に忠実に再現する方針で実装しているが、  
実務的な観点や判断により、**無自覚に理論から外れている箇所が存在する可能性がある**ことも前提としている。

## 3. プロジェクト構成（モジュラモノリス / Go + Clean Architecture）

本プロジェクトは、Clean Architecture をベースにしつつ、実務で過度な抽象化を避けるために構成を調整している。  
依存関係は内向きに限定し、ドメイン中心の構造を維持しながらも、運用時に扱いやすいレイヤ分割としている。  
その上で、`internal/modules` 配下に境界づけられたコンテキスト相当のモジュールを配置するモジュラモノリス構成を採用している。

```
.
├── main.go              # API起動エントリポイント
├── internal/
│   ├── modules/         # 業務モジュール（境界づけられたコンテキスト）
│   │   ├── campus/      # キャンパス管理モジュール
│   │   ├── lesson/      # 講座管理モジュール
│   │   ├── role/        # 権限管理モジュール
│   │   ├── room/        # 教室管理モジュール
│   │   ├── schedule/    # スケジュール管理モジュール
│   │   ├── session/     # セッション管理モジュール
│   │   └── user/        # ユーザー管理モジュール
│   ├── adapter/         # アプリ全体の共通ルータ / 共通ミドルウェア
│   ├── configs/         # 環境変数設定（.env / env.go）
│   ├── infrastructure/  # 共通サーバ実装（HTTPサーバ初期化など）
│   ├── injector/        # DI（依存解決）
│   ├── pkg/             # 共通ユーティリティ（ドメイン非依存）
│   ├── platform/        # DB / Logger など共通プラットフォーム実装
│   └── third_party/     # サードパーティー製ツール連携コード
└── go.mod               # Go依存管理
```

### モジュール内ディレクトリ例（`schedule`）
```
internal/modules/schedule/
├── adapter/             # 入出力境界（コントローラー / プレゼンター）
│   ├── controller/      # HTTPリクエスト受け口
│   └── presenter/       # ユースケース出力整形
├── domain/              # ドメインルールの中心
│   ├── model/           # 集約 / エンティティ
│   │   ├── invisible/   # 非表示教室モデル
│   │   └── schedule/    # スケジュールモデル
│   ├── repository/      # リポジトリインターフェース
│   ├── service/         # ドメインサービス
│   └── vo/              # バリューオブジェクト
├── infrastructure/      # 実装詳細（永続化）
│   └── database/        # DBアクセス実装
│       └── rdb/         # RDB実装
│           └── query/   # クエリ実装
│               └── schedule/    # スケジュール向けクエリ
└── usecase/             # アプリケーション層
    ├── external/        # 他モジュール依存インターフェース
    ├── interactor/      # ユースケース実行ロジックの本体
    ├── mapper/          # DTO変換
    ├── port/            # 入出力DTO定義
    └── query/           # 参照系ユースケース
        └── schedule_list/   # スケジュール一覧クエリ
```

### アーキテクチャ概要図（レイヤ構造）
```

      +-----------------------------------------+
      |      internal/modules/<module>          |
      |                                         |
      |  +-----------+                          |
      |  |  adapter  |                          |
      |  +-----+-----+                          |
      |        |                                |
      |        v                                |
      |  +-----+-----+      +---------------+   |
      |  |  usecase  +----->+  external     |   |
      |  +-----+-----+      +---------------+   |
      |        |                                |
      |        v                                |
      |  +-----+-----+                          |
      |  |  domain   |                          |
      |  +-----+-----+                          |
      |        |                                |
      |        v                                |
      |  +-----+-----------+                    |
      |  | infrastructure  |                    |
      |  +-----------------+                    |
      +-----------------------------------------+
```

## 各レイヤの責務

- **`internal/modules/<module>/domain/`**  
  ドメインモデルとドメインロジックを定義する中心層。  
  内部にのみ依存し、外部への依存を持たない。

- **`internal/modules/<module>/usecase/`**  
  ユースケースの調整を担当し、ドメインを呼び出すためのアプリケーションロジックを提供する。  
  `interactor` / `query` / `public` / `external` などの役割で分割している。

- **`internal/modules/<module>/adapter/`**  
  Clean Architecture における *Interface Adapter 層* に相当する。  
  入出力の境界を扱い、HTTP ハンドラやルーティングなど外部のリクエストを usecase 層が扱える形式へ変換する。

- **`internal/modules/<module>/infrastructure/`**  
  永続化や外部サービスなど、実装詳細に関わる処理を担う。  
  Repository の実装などが配置される。

- **`internal/adapter/`**  
  アプリ全体の共通ルーティングや共通ミドルウェアを扱う。

- **`internal/injector/`**  
  依存解決（DI）を行い、起動時に各レイヤを組み立てる。

- **`internal/platform/`**  
  DB接続、トランザクション管理、ロガーなどの共通基盤を扱う。

- **`internal/configs/`, `internal/pkg/`, `internal/third_party/`**  
  設定、共通ユーティリティ、サードパーティ連携コードを管理する。

DDD、Clean Architecture の原則は尊重しつつ、  
実務運用で扱いやすい構造として **“必要な部分のみを採用する”** 方針で設計している。

## DDD において意図的に外している点

- 本来、DDD では外部キー制約が必要となる関係は同一集約として扱うべきとされ、集約境界を跨いだ強い整合性を持たせることは推奨されていない。しかし本プロジェクトでは実務上のデータ破壊防止のため、集約間にも外部キー制約を設定している。ただし実装としては外部キーの存在に依存しておらず、外部キーが存在しない場合でも正しく動作するように整合性チェックをユースケース側で行う方針を取っている。

- トランザクションは集約単位で張るべきとされるが、実務上は複数集約を跨ぐ処理が一般的であるため、本プロジェクトではユースケース層でトランザクションを扱う方式を採用している。

## Clean Architecture において意図的に外している点

- 本来は Port（interface）を用いてユースケースと外部を明確に切り離すが、本プロジェクトでは採用していない。初期に導入を試みたものの制御が複雑化し、メリットを感じなかったため削除した。そのためinteractorが戻り値を持っている。

- Echo が本来の境界を越えて adapter 層（Controller 相当）まで侵食している。本来はフレームワークを infrastructure 層に留めるべきだが、Echo を infrastructure 層に閉じ込めた構成を実現できなかったため、adapter 層での利用を許容している。

- Repository の実装は infrastructure 層が直接行っている。本来は gateway（永続化抽象）を挟むべきだが、抽象レイヤを増やすメリットが薄く、過剰な階層化であると判断したため省略した。

## データベースのトランザクション分離レベルについて

本プロジェクトでは、MySQL のトランザクション分離レベルをデフォルトの REPEATABLE READ から READ COMMITTED に変更している。

READ COMMITTED を採用する理由:
- 長時間トランザクションによるギャップロックを回避したい
- 実務上、REPEATABLE READ が不要な一貫性保証を過剰に要求するケースが多い
- 更新頻度が高いテーブルでデッドロックが発生しやすいため、リスクを低減するため

これらの判断から、実務的な運用を優先する形で READ COMMITTED を採用している。

## 4. 実行方法

本プロジェクトは Docker および Docker Compose を利用して、アプリケーションの実行環境（Go, MySQL など）を構築する。

### 環境変数ファイル作成
```bash
cp internal/configs/.env.example internal/configs/.env
```

### コンテナ起動
```bash
docker-compose up -d
```

### マイグレーション実行
アプリケーション起動前に、データベースマイグレーションを実行する。

```bash
docker-compose exec app-dev atlas migrate apply --env dev -c file://db/atlas.hcl
```

### アプリケーション実行  
コンテナは実行環境のみを提供するため、アプリケーション本体は手動で起動する。

```bash
docker-compose exec app-dev go run main.go
```

### Swagger UI
API ドキュメントは Swag により自動生成され、local 環境ではブラウザから確認できる。

http://localhost:3002/swagger/index.html

### 停止

```bash
docker-compose down
```

## 5. 技術スタック

### 言語・基盤
- Go 1.24

### Web Framework
- Echo (github.com/labstack/echo/v4)

### アーキテクチャ / DI
- Uber Dig (go.uber.org/dig)

### ORM / データアクセス
- SQLBoiler (github.com/aarondl/sqlboiler/v4)

### ユーティリティ / サポート
- Godotenv (github.com/joho/godotenv)
- Envconfig (github.com/kelseyhightower/envconfig)
- LO (github.com/samber/lo)

### API ドキュメント (Swagger)
- Swag (github.com/swaggo/swag)

## 6. Swagger ドキュメント生成

本プロジェクトでは Swag を使用して Swagger ドキュメントを生成している。
フロントエンド側でも API 仕様を参照できるようにするため、必ず --requiredByDefault を付与して生成している。

swag init --requiredByDefault

生成されたドキュメントは /swagger 以下に配置される。
起動後は以下の URL から確認できる。

http://localhost:3002/swagger/index.html


## 7. データベースマイグレーション

本プロジェクトでは Atlas を使用してデータベースマイグレーションを管理している。
マイグレーションファイルは db/migrations に配置されており、ここにある SQL が順番に実行される。

マイグレーション実行例（ローカル環境）:

atlas migrate apply \
  --dir "file://db/migrations" \
  --url "mysql://root:root@localhost:23306/lessonlink"


## 今後の展望

- ~~境界づけられたコンテキストを明確化するため、モジュラモノリス構成への移行を検討している。~~
- ~~ドメインごとにモジュールを分割し、依存方向と責務境界をより厳密に管理できるようにする。~~
- 現在利用している SQLBoiler がメンテナンスモードに移行したため、将来的には ent への移行を検討している。より柔軟な型安全性とコードベースでのスキーマ管理により、保守性と開発効率の向上を目指す。
- ~~テストコードを追加し振る舞いの保証とリファクタリング耐性を高める。テストの方法については要検討。テーブルドリブンによる各ケースのテストはコストに対して得られる効果のバランスが悪いと感じているため効果的なアプローチを模索したい。~~
- テスト方針はゴールデンファイルによる比較検証と、OpenAPI（Swagger）スキーマベースの自動テストを採用する。加えて、モデル層にはテーブルドリブンテストを追加して仕様の明文化を進める。
