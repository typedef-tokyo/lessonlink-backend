table "data_campuses" {
  schema = schema.lessonlink
  column "id" {
    null           = false
    type           = int
    auto_increment = true
  }
  column "campus" {
    null = false
    type = varchar(16)
  }
  column "campus_name" {
    null = false
    type = varchar(32)
  }
  column "order_index" {
    null = false
    type = int
  }
  column "created_at" {
    null    = false
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null      = false
    type      = datetime
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "campus" {
    unique  = true
    columns = [column.campus]
  }
}
table "data_lessons" {
  schema = schema.lessonlink
  column "id" {
    null           = false
    type           = int
    auto_increment = true
  }
  column "campus" {
    null = false
    type = varchar(16)
  }
  column "name" {
    null = false
    type = varchar(32)
  }
  column "duration" {
    null = false
    type = int
  }
  column "created_at" {
    null    = false
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null      = false
    type      = datetime
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "data_lessons_ibfk_1" {
    columns     = [column.campus]
    ref_columns = [table.data_campuses.column.campus]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "campus" {
    unique  = true
    columns = [column.campus, column.name]
  }
}
table "data_roles" {
  schema = schema.lessonlink
  column "id" {
    null           = false
    type           = int
    auto_increment = true
  }
  column "role_key" {
    null = false
    type = varchar(16)
  }
  column "role_name" {
    null = false
    type = varchar(16)
  }
  column "created_at" {
    null    = false
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null      = false
    type      = datetime
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "role_key" {
    unique  = true
    columns = [column.role_key]
  }
}
table "data_rooms" {
  schema = schema.lessonlink
  column "id" {
    null           = false
    type           = int
    auto_increment = true
  }
  column "campus" {
    null = false
    type = varchar(16)
  }
  column "room_index" {
    null = false
    type = int
  }
  column "name" {
    null = false
    type = varchar(32)
  }
  column "created_at" {
    null    = false
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null      = false
    type      = datetime
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "data_rooms_ibfk_1" {
    columns     = [column.campus]
    ref_columns = [table.data_campuses.column.campus]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "campus" {
    unique  = true
    columns = [column.campus, column.room_index]
  }
}
table "sys_sessions" {
  schema = schema.lessonlink
  column "session_id" {
    null = false
    type = varchar(128)
  }
  column "user_id" {
    null = false
    type = int
  }
  column "value" {
    null = false
    type = text
  }
  column "updated_at" {
    null      = false
    type      = datetime
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  column "created_at" {
    null    = false
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.session_id]
  }
  index "user_id" {
    columns = [column.user_id]
  }
}
table "tbl_schedule_invisible_rooms" {
  schema = schema.lessonlink
  column "id" {
    null           = false
    type           = int
    auto_increment = true
  }
  column "schedule_id" {
    null = false
    type = int
  }
  column "room_index" {
    null = false
    type = int
  }
  column "created_at" {
    null    = false
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
}
table "tbl_schedule_items" {
  schema = schema.lessonlink
  column "id" {
    null           = false
    type           = int
    auto_increment = true
  }
  column "schedule_id" {
    null = false
    type = int
  }
  column "history_index" {
    null = false
    type = int
  }
  column "lesson_id" {
    null = false
    type = int
  }
  column "identifier" {
    null = false
    type = varchar(36)
  }
  column "duration" {
    null = false
    type = int
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "tbl_schedule_items_ibfk_1" {
    columns     = [column.schedule_id]
    ref_columns = [table.tbl_schedules.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "tbl_schedule_items_ibfk_2" {
    columns     = [column.lesson_id]
    ref_columns = [table.data_lessons.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "lesson_id" {
    columns = [column.lesson_id]
  }
  index "schedule_id" {
    columns = [column.schedule_id]
  }
  index "schedule_id_2" {
    unique  = true
    columns = [column.schedule_id, column.history_index, column.identifier]
  }
}
table "tbl_schedule_room_items" {
  schema = schema.lessonlink
  column "id" {
    null           = false
    type           = int
    auto_increment = true
  }
  column "schedule_id" {
    null = false
    type = int
  }
  column "history_index" {
    null = false
    type = int
  }
  column "item_tag" {
    null = false
    type = varchar(32)
  }
  column "lesson_id" {
    null = false
    type = int
  }
  column "identifier" {
    null = false
    type = varchar(36)
  }
  column "duration" {
    null = false
    type = int
  }
  column "start_time_hour" {
    null = false
    type = int
  }
  column "start_time_minutes" {
    null = false
    type = int
  }
  column "end_time_hour" {
    null = false
    type = int
  }
  column "end_time_minutes" {
    null = false
    type = int
  }
  column "room_index" {
    null = false
    type = int
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "tbl_schedule_room_items_ibfk_1" {
    columns     = [column.schedule_id]
    ref_columns = [table.tbl_schedules.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "lesson_id" {
    columns = [column.lesson_id]
  }
  index "schedule_id" {
    columns = [column.schedule_id]
  }
  index "schedule_id_2" {
    unique  = true
    columns = [column.schedule_id, column.history_index, column.identifier]
  }
}
table "tbl_schedules" {
  schema = schema.lessonlink
  column "id" {
    null           = false
    type           = int
    auto_increment = true
  }
  column "campus" {
    null = false
    type = varchar(16)
  }
  column "title" {
    null = false
    type = varchar(64)
  }
  column "history_index" {
    null = false
    type = int
  }
  column "start_time" {
    null = false
    type = int
  }
  column "end_time" {
    null = false
    type = int
  }
  column "create_user" {
    null = false
    type = int
  }
  column "last_update_user" {
    null = false
    type = int
  }
  column "created_at" {
    null    = false
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null      = false
    type      = datetime
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "tbl_schedules_ibfk_1" {
    columns     = [column.campus]
    ref_columns = [table.data_campuses.column.campus]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "tbl_schedules_ibfk_2" {
    columns     = [column.create_user]
    ref_columns = [table.tbl_users.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "tbl_schedules_ibfk_3" {
    columns     = [column.last_update_user]
    ref_columns = [table.tbl_users.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "campus" {
    columns = [column.campus]
  }
  index "create_user" {
    columns = [column.create_user]
  }
  index "last_update_user" {
    columns = [column.last_update_user]
  }
}
table "tbl_users" {
  schema = schema.lessonlink
  column "id" {
    null           = false
    type           = int
    auto_increment = true
  }
  column "role_key" {
    null = false
    type = varchar(16)
  }
  column "user_name" {
    null = false
    type = varchar(64)
  }
  column "password" {
    null = false
    type = text
  }
  column "name" {
    null = false
    type = varchar(64)
  }
  column "update_user_id" {
    null = false
    type = int
  }
  column "delete_flag" {
    null = false
    type = int
  }
  column "created_at" {
    null    = false
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null      = false
    type      = datetime
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "tbl_users_ibfk_1" {
    columns     = [column.role_key]
    ref_columns = [table.data_roles.column.role_key]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "tbl_users_ibfk_2" {
    columns     = [column.update_user_id]
    ref_columns = [table.tbl_users.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "role_key" {
    columns = [column.role_key]
  }
  index "update_user_id" {
    columns = [column.update_user_id]
  }
  index "user_name" {
    unique  = true
    columns = [column.user_name]
  }
}
schema "lessonlink" {
  charset = "utf8mb4"
  collate = "utf8mb4_0900_ai_ci"
}
