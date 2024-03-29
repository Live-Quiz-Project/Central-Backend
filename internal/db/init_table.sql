CREATE TABLE IF NOT EXISTS "user" (
  id UUID PRIMARY KEY NOT NULL,
  name TEXT,
  email TEXT UNIQUE,
  password TEXT,
  image TEXT,
  display_name TEXT,
  display_emoji TEXT,
  display_color TEXT,
  account_status TEXT,
  google_id TEXT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS quiz (
  id UUID PRIMARY KEY NOT NULL,
  creator_id UUID NOT NULL REFERENCES "user" (id),
  title TEXT DEFAULT 'Untitled',
  description TEXT,
  cover_image TEXT,
  visibility TEXT,
  time_limit INT,
  have_time_factor BOOLEAN,
  time_factor INT,
  font_size INT,
  mark INT,
  select_min INT,
  select_max INT,
  case_sensitive BOOLEAN,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS quiz_history (
  id UUID PRIMARY KEY NOT NULL,
  quiz_id UUID NOT NULL REFERENCES quiz (id),
  creator_id UUID NOT NULL REFERENCES "user" (id),
  title TEXT DEFAULT 'Untitled',
  description TEXT,
  cover_image TEXT,
  visibility TEXT,
  time_limit INT,
  have_time_factor BOOLEAN,
  time_factor INT,
  font_size INT,
  mark INT,
  select_min INT,
  select_max INT,
  case_sensitive BOOLEAN,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS question_pool (
  id UUID PRIMARY KEY NOT NULL,
  quiz_id UUID NOT NULL REFERENCES quiz (id),
  "order" INT,
  pool_order INT,
  content TEXT,
  note TEXT,
  media TEXT,
  media_type TEXT,
  time_limit INT,
  have_time_factor BOOLEAN,
  time_factor INT,
  font_size INT,
  mark INT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS question_pool_history (
  id UUID PRIMARY KEY NOT NULL,
  question_pool_id UUID NOT NULL REFERENCES question_pool (id),
  quiz_id UUID NOT NULL REFERENCES quiz_history (id),
  "order" INT,
  pool_order INT,
  content TEXT,
  note TEXT,
  media TEXT,
  media_type TEXT,
  time_limit INT,
  have_time_factor BOOLEAN,
  time_factor INT,
  font_size INT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS question (
  id UUID PRIMARY KEY NOT NULL,
  quiz_id UUID NOT NULL REFERENCES quiz (id),
  question_pool_id UUID REFERENCES question_pool (id),
  type TEXT,
  "order" INT,
  pool_order INT,
  pool_required BOOLEAN,
  content TEXT,
  note TEXT,
  media TEXT,
  media_type TEXT,
  use_template BOOLEAN,
  time_limit INT,
  have_time_factor BOOLEAN,
  time_factor INT,
  font_size INT,
  layout_idx INT,
  select_min INT,
  select_max INT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS question_history (
  id UUID PRIMARY KEY NOT NULL,
  question_id UUID NOT NULL REFERENCES question (id),
  quiz_id UUID NOT NULL REFERENCES quiz_history (id),
  question_pool_id UUID REFERENCES question_pool_history (id),
  type TEXT,
  "order" INT,
  pool_order INT,
  pool_required BOOLEAN,
  content TEXT,
  note TEXT,
  media TEXT,
  media_type TEXT,
  use_template BOOLEAN,
  time_limit INT,
  have_time_factor BOOLEAN,
  time_factor INT,
  font_size INT,
  layout_idx INT,
  select_min INT,
  select_max INT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS option_choice (
  id UUID PRIMARY KEY NOT NULL,
  question_id UUID NOT NULL REFERENCES question (id),
  "order" INT,
  content TEXT,
  mark INT,
  color TEXT,
  correct BOOLEAN,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS option_choice_history (
  id UUID PRIMARY KEY NOT NULL,
  option_choice_id UUID NOT NULL REFERENCES option_choice (id),
  question_id UUID NOT NULL REFERENCES question_history (id),
  "order" INT,
  content TEXT,
  mark INT,
  color TEXT,
  correct BOOLEAN,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS option_text (
  id UUID PRIMARY KEY NOT NULL,
  question_id UUID NOT NULL REFERENCES question (id),
  "order" INT,
  content TEXT,
  mark INT,
  case_sensitive BOOLEAN,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS option_text_history (
  id UUID PRIMARY KEY NOT NULL,
  option_text_id UUID NOT NULL REFERENCES option_text (id),
  question_id UUID NOT NULL REFERENCES question_history (id),
  "order" INT,
  content TEXT,
  mark INT,
  case_sensitive BOOLEAN,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS option_matching (
  id UUID PRIMARY KEY NOT NULL,
  question_id UUID NOT NULL REFERENCES question (id),
  "type" TEXT,
  "order" INT,
  content TEXT,
  color TEXT,
  eliminate BOOLEAN,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS option_matching_history (
  id UUID PRIMARY KEY NOT NULL,
  option_matching_id UUID NOT NULL REFERENCES option_matching (id),
  question_id UUID NOT NULL REFERENCES question_history (id),
  "type" TEXT,
  "order" INT,
  content TEXT,
  color TEXT,
  eliminate BOOLEAN,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS answer_matching (
  id UUID PRIMARY KEY NOT NULL,
  question_id UUID NOT NULL REFERENCES question (id),
  prompt_id UUID,
  option_id UUID,
  mark INT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS answer_matching_history (
  id UUID PRIMARY KEY NOT NULL,
  answer_matching_id UUID NOT NULL REFERENCES answer_matching (id),
  question_id UUID NOT NULL REFERENCES question_history (id),
  prompt_id UUID,
  option_id UUID,
  mark INT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS live_quiz_session (
  id UUID PRIMARY KEY NOT NULL,
  host_id UUID NOT NULL REFERENCES "user" (id),
  quiz_id UUID NOT NULL REFERENCES quiz_history (id),
  status TEXT NOT NULL,
  exempted_question_ids TEXT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS participant (
  id UUID PRIMARY KEY NOT NULL,
  user_id UUID REFERENCES "user" (id),
  live_quiz_session_id UUID NOT NULL REFERENCES live_quiz_session (id),
  status TEXT NOT NULL,
  name TEXT,
  marks INT,
  emoji TEXT,
  color TEXT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS answer_response (
  id UUID PRIMARY KEY NOT NULL,
  live_quiz_session_id UUID NOT NULL REFERENCES live_quiz_session (id),
  participant_id UUID NOT NULL REFERENCES participant (id),
  "type" TEXT,
  question_id UUID NOT NULL REFERENCES question_history (id),
  answer TEXT,
  use_time INT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS admin (
  id UUID PRIMARY KEY NOT NULL,
  email TEXT UNIQUE,
  password TEXT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);