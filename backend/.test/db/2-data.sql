DELETE FROM IDENTITY.DOMAIN;
INSERT INTO IDENTITY.DOMAIN(name) 
values
('identity'),
('community-goal-tracker') RETURNING ID;

-- password is password
DELETE FROM IDENTITY.USER;
INSERT INTO IDENTITY.USER(role,email,password,domains) 
values
(
    'administrator',
    'test0@test.com','$2a$10$o4OS7VZwSxf8qsjg5/JmCeAgLHD/KtnY5/YX4U4tRZ77tq3TNN7FG',
    (
        '{
            "'||(SELECT ID FROM IDENTITY.DOMAIN WHERE NAME='community-goal-tracker')||'":{
                "role":"user"
            },
            "'||(SELECT ID FROM IDENTITY.DOMAIN WHERE NAME='identity')||'":{
                "role":"administrator"
            }
        }'
    )::jsonb
),
(
    'user',
    'test1@test.com','$2a$10$o4OS7VZwSxf8qsjg5/JmCeAgLHD/KtnY5/YX4U4tRZ77tq3TNN7FG',
    (
        '{
            "'||(SELECT ID FROM IDENTITY.DOMAIN WHERE NAME='community-goal-tracker')||'":{
                "role":"user"
            },
            "'||(SELECT ID FROM IDENTITY.DOMAIN WHERE NAME='identity')||'":{
                "role":"user"
            }
        }'
    )::jsonb
),
(
    'user',
    'test2@test.com','$2a$10$o4OS7VZwSxf8qsjg5/JmCeAgLHD/KtnY5/YX4U4tRZ77tq3TNN7FG',
    (
        '{
            "'||(SELECT ID FROM IDENTITY.DOMAIN WHERE NAME='community-goal-tracker')||'":{
                "role":"user"
            },
            "'||(SELECT ID FROM IDENTITY.DOMAIN WHERE NAME='identity')||'":{
                "role":"user"
            }
        }'
    )::jsonb
) RETURNING UUID;


DELETE FROM COMMUNITY_GOAL_TRACKER.GOAL;
INSERT INTO COMMUNITY_GOAL_TRACKER.GOAL(NAME, ACHIEVERS)
values
(
    'get-degree',
    (
    '{
        "'||(SELECT UUID FROM IDENTITY.USER WHERE EMAIL='test0@test.com')||'":{
            "state":"inprogress","progress":50,"messages":{"2008-09-01 12:35:45":"enrolled at caltech","2008-09-05 12:35:45":"registered for classes"}
        }
    }'
    )::jsonb
),
(
    'build-house',
    (
    '{
        "'||(SELECT UUID FROM IDENTITY.USER WHERE EMAIL='test0@test.com')||'":{
            "state":"inprogress","progress":50,"messages":{"2008-09-01 12:35:45":"got approval from city council","2008-10-01 12:35:45":"put up walls"}
        },
        "'||(SELECT UUID FROM IDENTITY.USER WHERE EMAIL='test2@test.com')||'":{
            "state":"inprogress","progress":50,"messages":{"2008-09-01 12:35:45":"laid foundation","2008-10-01 12:35:45":"set up supports"}
        }
    }'
    )::jsonb
),
(
    'create-business',
    (
    '{
        "'||(SELECT UUID FROM IDENTITY.USER WHERE EMAIL='test0@test.com')||'":{
            "state":"inprogress","progress":50,"messages":{"2008-09-01 12:35:45":"created website","2008-10-01 12:35:45":"saved money for startup costs"}
        },
        "'||(SELECT UUID FROM IDENTITY.USER WHERE EMAIL='test1@test.com')||'":{
            "state":"inprogress","progress":50,"messages":{"2008-09-01 12:35:45":"created business plan","2008-10-01 12:35:45":"applied for loan"}
        }
    }'
    )::jsonb
);

DELETE FROM COMMUNITY_GOAL_TRACKER.ACHIEVER;
INSERT INTO COMMUNITY_GOAL_TRACKER.ACHIEVER(UUID,ROLE,FIRSTNAME, LASTNAME, ADDRESS, PHONE, GOALS)
values
(   
    (SELECT UUID FROM IDENTITY.USER WHERE EMAIL='test0@test.com'),
    'user',
    'test admin0 firstname',
    'test admin0 lastname',
    'test admin0 address',
    'test admin0 phone',
    (
        '{
            "'||(SELECT ID FROM COMMUNITY_GOAL_TRACKER.GOAL WHERE NAME='get-degree')||'": true,
            "'||(SELECT ID FROM COMMUNITY_GOAL_TRACKER.GOAL WHERE NAME='build-house')||'": true,
            "'||(SELECT ID FROM COMMUNITY_GOAL_TRACKER.GOAL WHERE NAME='create-business')||'": true
        }'
    )::jsonb
),
(   
    (SELECT UUID FROM IDENTITY.USER WHERE EMAIL='test1@test.com'),
    'user',
    'test admin1 firstname',
    'test admin1 lastname',
    'test admin1 address',
    'test admin1 phone',
    (
        '{
            "'||(SELECT ID FROM COMMUNITY_GOAL_TRACKER.GOAL WHERE NAME='create-business')||'":true
        }'
    )::jsonb
),
(
    (SELECT UUID FROM IDENTITY.USER WHERE EMAIL='test2@test.com'),
    'user',
    'test admin2 firstname',
    'test admin2 lastname',
    'test admin2 address',
    'test admin2 phone',
    (
        '{
            "'||(SELECT ID FROM COMMUNITY_GOAL_TRACKER.GOAL WHERE NAME='build-house')||'":true     
        }'
    )::jsonb
);