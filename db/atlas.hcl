env "dev" {
  url = "mysql://root:root@mysql/lessonlink"
  migration {
    dir = "file://migrations"
  }
}
