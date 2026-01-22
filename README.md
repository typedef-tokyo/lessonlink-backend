## 1. プロジェクト概要

本プロジェクトは、Go を用いて実装した **DDD（Domain-Driven Design）+ Clean Architecture 構成のバックエンドアプリケーション**です。  
対象ドメインは「講座と教室のスケジュール作成」というシンプルな領域ですが、**技術的関心事の分離とドメインモデル中心の設計を実践すること**を目的としています。

本プロジェクトで重視しているポイントは以下です。

- **戦術的 DDD パターン**（Entity / ValueObject / Aggregate / Repository / DomainService）の適用  
- **Clean Architecture** による依存関係の統制と疎結合な構造  
- **ドメインロジックの入出力層からの完全分離**  
- **Application Service** によるユースケース単位の振る舞いの定義

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

## 3. プロジェクト構成（Go + Clean Architecture）

本プロジェクトは、Clean Architecture をベースにしつつ、実務で過度な抽象化を避けるために構成を調整している。  
依存関係は内向きに限定し、ドメイン中心の構造を維持しながらも、運用時に扱いやすいレイヤ分割としている。

```
.
├── cmd/                 # エントリポイント
├── internal/
│   ├── adapter/         # Interface Adapter 層（入出力境界：Handler / Presenter など）
│   ├── configs/         # 設定情報
│   ├── domain/          # Entity / ValueObject / DomainService
│   ├── infrastructure/  # 永続化 / 外部サービス / 実装詳細
│   ├── injector/        # DI（依存解決）
│   ├── pkg/             # 共通ユーティリティ（ドメイン非依存）
│   ├── third_party/     # サードパーティー製のツール
│   └── usecase/         # Application Service（ユースケース）
└── go.mod
```

### アーキテクチャ概要図（レイヤ構造）
```

         +-------------------+
         |      domain       |
         +---------+---------+
                   ↑
             +-----+------+
             |   usecase  |
             +-----+------+
                   ↑
           +-------+---------+
           |     adapter     |
           +-------+---------+
                   ↑
           +-------+---------+
           |  infrastructure |
           +-----------------+
```

## 各レイヤの責務

- **domain/**  
  ドメインモデルとドメインロジックを定義する中心層。  
  内部にのみ依存し、外部への依存を持たない。

- **usecase/**  
  ユースケースの調整を担当し、ドメインを呼び出すためのアプリケーションロジックを提供する。

- **adapter/**  
  Clean Architecture における *Interface Adapter 層* に相当する。  
  入出力の境界を扱い、HTTP ハンドラやルーティングなど外部のリクエストを usecase 層が扱える形式へ変換する。

- **infrastructure/**  
  永続化や外部サービスなど、実装詳細に関わる処理を担う。  
  Repository の実装などが配置される。

- **configs/**  
  設定情報の管理を行う。

- **injector/**  
  依存解決（DI）を行い、起動時に各レイヤを組み立てる。

- **pkg/**  
  ドメイン非依存の共通ユーティリティ。

- **third_party/**  
  サードパーティー製のツールやラッパーを格納する。

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

### コンテナ起動

```bash
docker-compose up -d
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

- 境界づけられたコンテキストを明確化するため、モジュラモノリス構成への移行を検討している。
- ドメインごとにモジュールを分割し、依存方向と責務境界をより厳密に管理できるようにする。
- ユースケースの独立性を高め、変更の影響範囲を最小化する構造へ改善する。
- 現在利用している SQLBoiler がメンテナンスモードに移行したため、将来的には ent への移行を検討している。より柔軟な型安全性とコードベースでのスキーマ管理により、保守性と開発効率の向上を目指す。
- テストコードを追加し振る舞いの保証とリファクタリング耐性を高める。テストの方法については要検討。テーブルドリブンによる各ケースのテストはコストに対して得られる効果のバランスが悪いと感じているため効果的なアプローチを模索したい。