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
  google_id UUID,
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
  created_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ,
  deleted BOOL DEFAULT FALSE,
  updated_by UUID NOT NULL REFERENCES "user" (id)
);
CREATE TABLE IF NOT EXISTS question (
  id UUID PRIMARY KEY NOT NULL,
  quiz_id UUID NOT NULL REFERENCES quiz (id),
  parent_id UUID REFERENCES question (id),
  type TEXT,
  "order" INT,
  content TEXT,
  note TEXT,
  media TEXT,
  time_limit INT,
  have_time_factor BOOL,
  time_factor INT,
  font_size INT,
  layout_idx INT,
  selected_up_to INT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS question_history (
  id UUID PRIMARY KEY NOT NULL,
  question_id UUID NOT NULL REFERENCES question (id),
  quiz_id UUID NOT NULL REFERENCES quiz_history (id),
  parent_id UUID REFERENCES question_history (id),
  type TEXT,
  "order" INT,
  content TEXT,
  note TEXT,
  media TEXT,
  time_limit INT,
  have_time_factor BOOL,
  time_factor INT,
  font_size INT,
  layout_idx INT,
  selected_up_to INT,
);
CREATE TABLE IF NOT EXISTS option_choice (
  id UUID PRIMARY KEY NOT NULL,
  question_id UUID NOT NULL REFERENCES question (id),
  "order" INT,
  content TEXT,
  mark INT,
  color TEXT,
  correct BOOL,
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
  correct BOOL,
);
CREATE TABLE IF NOT EXISTS option_text (
  id UUID PRIMARY KEY NOT NULL,
  question_id UUID NOT NULL REFERENCES question (id),
  "order" INT,
  content TEXT,
  mark INT,
  case_sensitive BOOL,
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
  case_sensitive BOOL,
);
CREATE TABLE IF NOT EXISTS option_matching (
  id UUID PRIMARY KEY NOT NULL,
  question_id UUID NOT NULL REFERENCES question (id),
  prompt_id UUID,
  option_id UUID,
  mark INT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS option_matching_history (
  id UUID PRIMARY KEY NOT NULL,
  option_matching_id UUID NOT NULL REFERENCES option_matching (id),
  question_id UUID NOT NULL REFERENCES question_history (id),
  prompt_id UUID,
  option_id UUID,
  mark INT,
);
CREATE TABLE IF NOT EXISTS option_matching_prompt (
  id UUID PRIMARY KEY NOT NULL,
  option_matching_id UUID NOT NULL REFERENCES option_matching (id),
  content TEXT,
  "order" INT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS option_matching_prompt_history (
  id UUID PRIMARY KEY NOT NULL,
  option_matching_prompt_id UUID NOT NULL REFERENCES option_matching_prompt (id),
  option_matching_id UUID NOT NULL REFERENCES option_matching_history (id),
  content TEXT,
  "order" INT,
);
CREATE TABLE IF NOT EXISTS option_matching_option (
  id UUID PRIMARY KEY NOT NULL,
  option_matching_id UUID NOT NULL REFERENCES option_matching (id),
  content TEXT,
  "order" INT,
  eliminated BOOL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS option_matching_option_history (
  id UUID PRIMARY KEY NOT NULL,
  option_matching_option_id UUID NOT NULL REFERENCES option_matching_option (id),
  option_matching_id UUID NOT NULL REFERENCES option_matching_history (id),
  content TEXT,
  "order" INT,
  eliminated BOOL,
);
CREATE TABLE IF NOT EXISTS option_pin (
  id UUID PRIMARY KEY NOT NULL,
  question_id UUID NOT NULL REFERENCES question (id),
  x_axis INT,
  y_axis INT,
  mark INT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS option_pin_history (
  id UUID PRIMARY KEY NOT NULL,
  option_pin_id UUID NOT NULL REFERENCES option_pin (id),
  question_id UUID NOT NULL REFERENCES question_history (id),
  x_axis INT,
  y_axis INT,
  mark INT,
);
CREATE TABLE IF NOT EXISTS live_quiz_session (
  id UUID PRIMARY KEY NOT NULL,
  host_id UUID NOT NULL REFERENCES "user" (id),
  quiz_id UUID NOT NULL REFERENCES quiz (id),
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
  marks INT
);
CREATE TABLE IF NOT EXISTS response_choice (
  id UUID PRIMARY KEY NOT NULL,
  participant_id UUID NOT NULL REFERENCES participant (id),
  option_choice_id UUID NOT NULL REFERENCES option_choice (id)
);
CREATE TABLE IF NOT EXISTS response_text (
  id UUID PRIMARY KEY NOT NULL,
  participant_id UUID NOT NULL REFERENCES participant (id),
  option_text_id UUID NOT NULL REFERENCES option_text (id),
  content TEXT
);
CREATE TABLE IF NOT EXISTS response_matching (
  id UUID PRIMARY KEY NOT NULL,
  participant_id UUID NOT NULL REFERENCES participant (id),
  option_matching_id UUID NOT NULL REFERENCES option_matching (id),
  option_matching_prompt_id UUID NOT NULL REFERENCES option_matching_prompt (id),
  option_matching_option_id UUID NOT NULL REFERENCES option_matching_option (id)
);
CREATE TABLE IF NOT EXISTS response_pin (
  id UUID PRIMARY KEY NOT NULL,
  participant_id UUID NOT NULL REFERENCES participant (id),
  option_pin_id UUID NOT NULL REFERENCES option_pin (id),
  x_axis INT NOT NULL,
  y_axis INT NOT NULL
);
CREATE TABLE IF NOT EXISTS admin (
  id UUID PRIMARY KEY NOT NULL,
  email TEXT UNIQUE,
  password TEXT
);