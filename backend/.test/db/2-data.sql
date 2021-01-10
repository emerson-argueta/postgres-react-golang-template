-- password is password
DELETE FROM IDENTITY.USER;
INSERT INTO IDENTITY.USER(ID,ROLE,EMAIL,HASHPASSWORD) 
values
(
    gen_random_uuid(),
    'administrator',
    'test0@test.com',
    '$2a$10$o4OS7VZwSxf8qsjg5/JmCeAgLHD/KtnY5/YX4U4tRZ77tq3TNN7FG'
),
(
    gen_random_uuid(),
    'user',
    'test1@test.com',
    '$2a$10$o4OS7VZwSxf8qsjg5/JmCeAgLHD/KtnY5/YX4U4tRZ77tq3TNN7FG'
),
(
    gen_random_uuid(),
    'user',
    'test2@test.com',
    '$2a$10$o4OS7VZwSxf8qsjg5/JmCeAgLHD/KtnY5/YX4U4tRZ77tq3TNN7FG'
) RETURNING ID;


DELETE FROM COMMUNITY_GOAL_TRACKER.GOAL;
INSERT INTO COMMUNITY_GOAL_TRACKER.GOAL(NAME, ACHIEVERS)
values
(
    'get-degree',
    (
    '{
        "'||(SELECT ID FROM IDENTITY.USER WHERE EMAIL='test0@test.com')||'":{
            "state":"inprogress","progress":50,"messages":{"2008-09-01 12:35:45":"enrolled at caltech","2008-09-05 12:35:45":"registered for classes"}
        }
    }'
    )::jsonb
),
(
    'build-house',
    (
    '{
        "'||(SELECT ID FROM IDENTITY.USER WHERE EMAIL='test0@test.com')||'":{
            "state":"inprogress","progress":50,"messages":{"2008-09-01 12:35:45":"got approval from city council","2008-10-01 12:35:45":"put up walls"}
        },
        "'||(SELECT ID FROM IDENTITY.USER WHERE EMAIL='test2@test.com')||'":{
            "state":"inprogress","progress":50,"messages":{"2008-09-01 12:35:45":"laid foundation","2008-10-01 12:35:45":"set up supports"}
        }
    }'
    )::jsonb
),
(
    'create-business',
    (
    '{
        "'||(SELECT ID FROM IDENTITY.USER WHERE EMAIL='test0@test.com')||'":{
            "state":"inprogress","progress":50,"messages":{"2008-09-01 12:35:45":"created website","2008-10-01 12:35:45":"saved money for startup costs"}
        },
        "'||(SELECT ID FROM IDENTITY.USER WHERE EMAIL='test1@test.com')||'":{
            "state":"inprogress","progress":50,"messages":{"2008-09-01 12:35:45":"created business plan","2008-10-01 12:35:45":"applied for loan"}
        }
    }'
    )::jsonb
);

DELETE FROM COMMUNITY_GOAL_TRACKER.ACHIEVER;
INSERT INTO COMMUNITY_GOAL_TRACKER.ACHIEVER(ID,USERID,ROLE,FIRSTNAME, LASTNAME, ADDRESS, PHONE, GOALS)
values
(   
    gen_random_uuid(),
    (SELECT ID FROM IDENTITY.USER WHERE EMAIL='test0@test.com'),
    'user',
    'test admin0 firstname',
    'test admin0 lastname',
    'test admin0 address',
    '123-1234-1234',
    (
        '{
            "'||(SELECT NAME FROM COMMUNITY_GOAL_TRACKER.GOAL WHERE NAME='get-degree')||'": true,
            "'||(SELECT NAME FROM COMMUNITY_GOAL_TRACKER.GOAL WHERE NAME='build-house')||'": true,
            "'||(SELECT NAME FROM COMMUNITY_GOAL_TRACKER.GOAL WHERE NAME='create-business')||'": true
        }'
    )::jsonb
),
(   
    gen_random_uuid(),
    (SELECT ID FROM IDENTITY.USER WHERE EMAIL='test1@test.com'),
    'user',
    'test admin1 firstname',
    'test admin1 lastname',
    'test admin1 address',
    '123-1234-1234',
    (
        '{
            "'||(SELECT NAME FROM COMMUNITY_GOAL_TRACKER.GOAL WHERE NAME='create-business')||'":true
        }'
    )::jsonb
),
(
    gen_random_uuid(),
    (SELECT ID FROM IDENTITY.USER WHERE EMAIL='test2@test.com'),
    'user',
    'test admin2 firstname',
    'test admin2 lastname',
    'test admin2 address',
    '123-1234-1234',
    (
        '{
            "'||(SELECT NAME FROM COMMUNITY_GOAL_TRACKER.GOAL WHERE NAME='build-house')||'":true     
        }'
    )::jsonb
);