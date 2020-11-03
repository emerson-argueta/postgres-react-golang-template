DELETE FROM IDENTITY.DOMAIN;
INSERT INTO IDENTITY.DOMAIN(name) 
values
('identity'),
('community-goal-tracker') RETURNING ID;

-- password mocanela12
DELETE FROM IDENTITY.USER;
INSERT INTO IDENTITY.USER(email,password,domains) 
values
(
    'jargueta1964@gmail.com','$2a$10$fPnqG7PMDRfQMZTqKkr6Iem5Z7/QKSY67Vi53trqX.D5t35AJBWSy',
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
    'test1@test.com','$2a$10$fPnqG7PMDRfQMZTqKkr6Iem5Z7/QKSY67Vi53trqX.D5t35AJBWSy',
    (
        '{
            "'||(SELECT ID FROM IDENTITY.DOMAIN WHERE NAME='community-goal-tracker')||'":{
                "role":"user"
            }
        }'
    )::jsonb
),
(
    'test2@test.com','$2a$10$fPnqG7PMDRfQMZTqKkr6Iem5Z7/QKSY67Vi53trqX.D5t35AJBWSy',
    (
        '{
            "'||(SELECT ID FROM IDENTITY.DOMAIN WHERE NAME='community-goal-tracker')||'":{
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
        "'||(SELECT UUID FROM IDENTITY.USER WHERE EMAIL='jargueta1964@gmail.com')||'":{
            "state":"inprogress","progress":50,"messages":{"2008-09-01 12:35:45":"enrolled at caltech","2008-09-05 12:35:45":"registered for classes"}
        }
    }'
    )::jsonb
),
(
    'build-house',
    (
    '{
        "'||(SELECT UUID FROM IDENTITY.USER WHERE EMAIL='jargueta1964@gmail.com')||'":{
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
        "'||(SELECT UUID FROM IDENTITY.USER WHERE EMAIL='jargueta1964@gmail.com')||'":{
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
    (SELECT UUID FROM IDENTITY.USER WHERE EMAIL='jargueta1964@gmail.com'),
    'user',
    'Jorge',
    'Argueta',
    '18851 Benicia St, Hesperia CA,92345',
    '909-644-5114',
    (
        '[
            '||(SELECT ID FROM COMMUNITY_GOAL_TRACKER.GOAL WHERE NAME='get-degree')||',
            '||(SELECT ID FROM COMMUNITY_GOAL_TRACKER.GOAL WHERE NAME='build-house')||',
            '||(SELECT ID FROM COMMUNITY_GOAL_TRACKER.GOAL WHERE NAME='create-business')||'
        ]'
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
        '[
            '||(SELECT ID FROM COMMUNITY_GOAL_TRACKER.GOAL WHERE NAME='create-business')||'
        ]'
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
        '[
            '||(SELECT ID FROM COMMUNITY_GOAL_TRACKER.GOAL WHERE NAME='build-house')||'     
        ]'
    )::jsonb
);