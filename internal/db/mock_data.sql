-- Import uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

SELECT * FROM "quiz"


SELECT * FROM "quiz_history"

SELECT * FROM "quiz" WHERE id = 'd90fe3d4-56c4-4cb3-9204-07b2c201d644';

-- Insert mock data in live_quiz_session
INSERT INTO "live_quiz_session" (id, host_id, quiz_id, status, exempted_question_ids, created_at, updated_at)
VALUES
  (uuid_generate_v4(), (SELECT id FROM "user" WHERE name = 'TestProfile 1'), (SELECT id FROM "quiz_history" ORDER BY random() LIMIT 1), 'active', NULL, NOW(), NOW());

SELECT * FROM "live_quiz_session";

-- Insert mock data user 2,3 in participant;
INSERT INTO "participant" (id, user_id, live_quiz_session_id, status, name, marks, created_at, updated_at)
VALUES
	(uuid_generate_v4(), (SELECT id FROM "user" WHERE name = 'TestProfile 2'), (SELECT id FROM "live_quiz_session" WHERE status = 'active' ), 'complete', 'Participant 2', 20, NOW(),NOW()),
	(uuid_generate_v4(), (SELECT id FROM "user" WHERE name = 'TestProfile 3'), (SELECT id FROM "live_quiz_session" WHERE status = 'active' ), 'complete', 'Participant 3', 20, NOW(),NOW());

SELECT * FROM "participant";

-- Insert mock response of user 2 in answer_response;
INSERT INTO answer_response (id, live_quiz_session_id, participant_id, "type", question_id, answer, use_time, created_at, updated_at)
VALUES
	(uuid_generate_v4(), (SELECT id FROM "live_quiz_session" WHERE status = 'active'), (SELECT id FROM "participant" WHERE name = 'Participant 2'), 'CHOICE', (SELECT id FROM "question_history" WHERE quiz_id = 'f1ff2d37-2cda-434c-9103-c10d84e30666' AND "order" = 1 ), 'Nah fam', 3, NOW(), NOW()),
	(uuid_generate_v4(), (SELECT id FROM "live_quiz_session" WHERE status = 'active'), (SELECT id FROM "participant" WHERE name = 'Participant 2'), 'SHORT_TEXT', (SELECT id FROM "question_history" WHERE quiz_id = 'f1ff2d37-2cda-434c-9103-c10d84e30666' AND "order" = 3 ), '2', 5, NOW(), NOW()),
	(uuid_generate_v4(), (SELECT id FROM "live_quiz_session" WHERE status = 'active'), (SELECT id FROM "participant" WHERE name = 'Participant 2'), 'PARAGRAPH', (SELECT id FROM "question_history" WHERE quiz_id = 'f1ff2d37-2cda-434c-9103-c10d84e30666' AND "order" = 5 ),'Main character died at the end', 14, NOW(), NOW()),
	(uuid_generate_v4(), (SELECT id FROM "live_quiz_session" WHERE status = 'active'), (SELECT id FROM "participant" WHERE name = 'Participant 2'), 'TRUE_FALSE', (SELECT id FROM "question_history" WHERE quiz_id = 'f1ff2d37-2cda-434c-9103-c10d84e30666' AND "order" = 6 ),'TRUE', 4, NOW(), NOW()),
	(uuid_generate_v4(), (SELECT id FROM "live_quiz_session" WHERE status = 'active'), (SELECT id FROM "participant" WHERE name = 'Participant 2'), 'MATCHING', (SELECT id FROM "question_history" WHERE quiz_id = 'f1ff2d37-2cda-434c-9103-c10d84e30666' AND "order" = 7 ),'Fish:Swim,Bird:Fly,Dog:Run', 18, NOW(), NOW());

-- Insert mock response of user 3 in answer_response;
INSERT INTO answer_response (id, live_quiz_session_id, participant_id, "type", question_id, answer, use_time, created_at, updated_at)
VALUES
	(uuid_generate_v4(), (SELECT id FROM "live_quiz_session" WHERE status = 'active'), (SELECT id FROM "participant" WHERE name = 'Participant 3'), 'CHOICE', (SELECT id FROM "question_history" WHERE quiz_id = 'f1ff2d37-2cda-434c-9103-c10d84e30666' AND "order" = 1 ), 'Yeah sure', 5, NOW(), NOW()),
	(uuid_generate_v4(), (SELECT id FROM "live_quiz_session" WHERE status = 'active'), (SELECT id FROM "participant" WHERE name = 'Participant 3'), 'SHORT_TEXT', (SELECT id FROM "question_history" WHERE quiz_id = 'f1ff2d37-2cda-434c-9103-c10d84e30666' AND "order" = 3 ), '4', 4 , NOW(), NOW()),
	(uuid_generate_v4(), (SELECT id FROM "live_quiz_session" WHERE status = 'active'), (SELECT id FROM "participant" WHERE name = 'Participant 3'), 'PARAGRAPH', (SELECT id FROM "question_history" WHERE quiz_id = 'f1ff2d37-2cda-434c-9103-c10d84e30666' AND "order" = 5 ),'No one dead', 7, NOW(), NOW()),
	(uuid_generate_v4(), (SELECT id FROM "live_quiz_session" WHERE status = 'active'), (SELECT id FROM "participant" WHERE name = 'Participant 3'), 'TRUE_FALSE', (SELECT id FROM "question_history" WHERE quiz_id = 'f1ff2d37-2cda-434c-9103-c10d84e30666' AND "order" = 6 ), 'FALSE', 2 ,NOW(), NOW()),
	(uuid_generate_v4(), (SELECT id FROM "live_quiz_session" WHERE status = 'active'), (SELECT id FROM "participant" WHERE name = 'Participant 3'), 'MATCHING', (SELECT id FROM "question_history" WHERE quiz_id = 'f1ff2d37-2cda-434c-9103-c10d84e30666' AND "order" = 7 ), 'Fish:Fly,Bird:Swim,Dog:Run', 22 ,NOW(), NOW());

SELECT * FROM answer_response

SELECT * FROM live_quiz_session WHERE id = '95ed504a-12d8-4c56-be6d-23b848682d0a'