<a name="Final MVP"></a>
## Final MVP

<a name="v1.1.0"></a>
## [v1.1.0] - 2024-08-23

<a name="v1.0.0"></a>
## v1.0.0 - 2024-08-23
### Build
- Modify analytics_board_admin and analytics_board_user to add type and typeString ([#234](y/issues/234))

### Feat
- add-mutex and wrong answer highlight ([#309](y/issues/309))
- add kratos login flow in quiz app ([#271](y/issues/271))
- add kratos schema ([#275](y/issues/275))
- making the font of inviatation code larger ([#266](y/issues/266))
- show-username-when-user-submits-answer ([#233](y/issues/233))
- analytics board api for the admin and user ([#179](y/issues/179))
- test case for csv auth
- added test cases for quiz start functionality
- added routes for list-quiz and modified '/admin/quiz' route as 'admin/quiz/create-quiz'
- skip button added for admin to skip timer
- user score board api
- added layout for final scoreboard
- skip event
- admin will redirect when he/she tring to join their own quiz
- added  modified dockerfile and also added the ci for build and push for qizz-api [#28](y/issues/28)
- added s6-overlay to the quiz-app also modified dockerfile and also added the ci for build and push [#18](y/issues/18)

### Feature
- Add accouont recovery and fix UI (Issue [#300](y/issues/300)) ([#308](y/issues/308))
- hide input for username field for existing user(Issue [#277](y/issues/277)) ([#280](y/issues/280))
- modify-header ([#202](y/issues/202))
- design-ui-landing ([#200](y/issues/200))
- design-admin-analytics-questionwise ([#198](y/issues/198))
- display user to admin when user join the quiz
- api for admin to show final scoreboard
- **add-type-survey:** Add type column in database for questions ([#231](y/issues/231))
- **class-accuracy-bar:** Show Class accuracy bar on the admin analysis page ([#228](y/issues/228))
- **design-join-page:** Design Join Page and Modify Header and Footer ([#232](y/issues/232))

### Feature
- add reports api and ui pages for admin ([#296](y/issues/296))
- Display user's name on their game area ([#167](y/issues/167))

### Fix
- break functions
- added correct context
- terminate quiz on unexpected exit of admin and guest user join issue (issue [#282](y/issues/282)) ([#287](y/issues/287))
- take time period for each question from env rather than csv ([#180](y/issues/180))
- struct printing err, test cases
- changing docker envs
- Restrict Points Allocation to maximum 20 Points per question ([#163](y/issues/163))
- Add up and down migration to alter old data and make points 0 where points were 1 ([#243](y/issues/243))
- username and password not showing
- cleanup invitationcode store and userlist store when waiting space component unmounts ([#242](y/issues/242))
- changed runner to ubuntu-latest
- changed runner to ubuntu-latest
- award points for any correct option chosen in multi-answer questions (except survey questions) ([#213](y/issues/213))
- errors-caused-by-cookies are resolved ([#208](y/issues/208), [#204](y/issues/204)) ([#216](y/issues/216))
- modify admin analysis ([#207](y/issues/207))
- modify user anaytics ([#205](y/issues/205))
- admin-analysis ([#199](y/issues/199))
- added text color for buttons
- changed the static words to the camel case
- username should be maximum of 12 characters long and should not consider blankspace as character ([#174](y/issues/174))
- docker file is back to its previous state as wanted by the CI flow ([#181](y/issues/181))
- the docker compose functionality is now working properly with these changes as there was unwanted envs and introduced docker volume for redis and also changed ip of the containers ([#164](y/issues/164))
- username is now send by GetUserMeta whether it is kratos or guest user ([#278](y/issues/278))
- also count Not Attempted survey in accuracy ([#290](y/issues/290))
- Quesiton and Score Space being cut in the mobile view ([#197](y/issues/197))
- removed console.log
- modified changes in design
- changed final scoreboard route for admin
- fixed bug while submitting answer when admin refreshes the page
- cq and pre-commit
- authentication added for scoreboard api
- modified user scoreboard api
- remove extra variable and else condition
- change email in kratos also ([#291](y/issues/291))
- error
- session get after session completed
- remove href attribute
- change port in nginx.conf from 3000 to 4000 as per .env.example [#32](y/issues/32)
- added selfhosted runner to api ci
- added path condition in CI
- **accuracy-logic-to-correct-incorrect:** Modify accuracy formula to calculate accuracy based on the correct and incorrect answers count ([#274](y/issues/274))
- **accuracy-logic-userwise:** fix userwise accuracy logic to count accuracy based on points instead of correct/incorrect ([#240](y/issues/240))
- **participants-disconnect-issue:** Issue of users getting disconnected in Waiting area ([#272](y/issues/272))
- **score-logic-for-survey:** Add logic to award score for survey questions where score is provided and do not award for survey questions where score is not provided ([#236](y/issues/236))
- **type-column-migrations:** Add migrations to add new column 'type' and constraints and all related to that ([#241](y/issues/241))

### Refactor
- use $fetch instead of usefetch in create quiz ([#303](y/issues/303))
- helper interface ([#302](y/issues/302))
- remove logic of storing user_played_quiz in cookie ([#298](y/issues/298))
- remove middleware from frontend and create quick user properly and remove unnecessary api calling for user data ([#297](y/issues/297))
- parse csv by header instead of row and column ([#288](y/issues/288))

### Reverts
- debug logs for admin disconnection issue ([#312](y/issues/312))
- add debug log ([#313](y/issues/313))
- adding the sort by latest and paginations to reports page ([#311](y/issues/311))
- fix: some ui and add logger in quiz_operation to solve all-param-required
- Recheck staging, checking for valid approach

### Pull Requests
- Merge pull request [#155](y/issues/155) from Improwised/fix/button-text-color
- Merge pull request [#153](y/issues/153) from Improwised/fix/admin-scoreboard-route
- Merge pull request [#152](y/issues/152) from Improwised/feature/final-scoreboard-rank-calculation
- Merge pull request [#140](y/issues/140) from Improwised/fix/bug-answer-submit
- Merge pull request [#135](y/issues/135) from Improwised/test-cases/quiz-start
- Merge pull request [#136](y/issues/136) from Improwised/feature/display-user-to-admin-while-join
- Merge pull request [#131](y/issues/131) from Improwised/feature/uploaded-quiz-list
- Merge pull request [#128](y/issues/128) from Improwised/fix/next-question-button
- Merge pull request [#121](y/issues/121) from Improwised/fix/question-scoreboard
- Merge pull request [#117](y/issues/117) from Improwised/feature/finalscoreboard-time
- Merge pull request [#118](y/issues/118) from Improwised/feature/skip-timer-button
- Merge pull request [#108](y/issues/108) from Improwised/feature/finalscoreboard-admin
- Merge pull request [#105](y/issues/105) from Improwised/feature/finalscoreboard-user-score
- Merge pull request [#101](y/issues/101) from Improwised/final-scoreboard
- Merge pull request [#95](y/issues/95) from Improwised/start-demo
- Merge pull request [#93](y/issues/93) from Improwised/56-feat-create-a-page-where-admin-can-enter-question-csv
- Merge pull request [#91](y/issues/91) from Improwised/fix-null-found-in-duratin
- Merge pull request [#89](y/issues/89) from Improwised/45-user-and-admin-should-see-rank-board-after-submitting-the-answer
- Merge pull request [#88](y/issues/88) from Improwised/43-handle-skip-event
- Merge pull request [#87](y/issues/87) from Improwised/40-handle-timeout-of-question
- Merge pull request [#82](y/issues/82) from Improwised/39-handle-start-event
- Merge pull request [#81](y/issues/81) from Improwised/fix-all-param-required
- Merge pull request [#80](y/issues/80) from Improwised/revert-79-fix-all-param-required
- Merge pull request [#79](y/issues/79) from Improwised/fix-all-param-required
- Merge pull request [#78](y/issues/78) from Improwised/75-fix-change-nuxtsession-with-active-library
- Merge pull request [#67](y/issues/67) from Improwised/user_admin_playground_ui
- Merge pull request [#69](y/issues/69) from Improwised/revert-68-recheck-staging
- Merge pull request [#70](y/issues/70) from Improwised/recheck-staging
- Merge pull request [#68](y/issues/68) from Improwised/recheck-staging
- Merge pull request [#66](y/issues/66) from Improwised/fix/remove-export-env
- Merge pull request [#65](y/issues/65) from Improwised/add-systemEnvConfig
- Merge pull request [#36](y/issues/36) from Improwised/create-login-page
- Merge pull request [#52](y/issues/52) from Improwised/51-change-configs-after-rename-envexample-to-envdocker
- Merge pull request [#48](y/issues/48) from Improwised/swagger-and-route-change
- Merge pull request [#47](y/issues/47) from Improwised/46-reset-ports-in-envs
- Merge pull request [#34](y/issues/34) from Improwised/fix-migration-downs
- Merge pull request [#35](y/issues/35) from Improwised/change-env-as-per-new-urls
- Merge pull request [#33](y/issues/33) from Improwised/feat/nginx-port-changes
- Merge pull request [#31](y/issues/31) from Improwised/change-dev-server-config
- Merge pull request [#30](y/issues/30) from Improwised/quiz-migrations
- Merge pull request [#22](y/issues/22) from Improwised/create-auth-middleware
- Merge pull request [#26](y/issues/26) from Improwised/23-bug-fix-the-admin-creation-on-zsh-as-well-make-file-generate-for-migration
- Merge pull request [#29](y/issues/29) from Improwised/feat/app-api-dockerfile-ci
- Merge pull request [#27](y/issues/27) from Improwised/fix/change-runner
- Merge pull request [#25](y/issues/25) from Improwised/fix/change-context-ci
- Merge pull request [#24](y/issues/24) from Improwised/feat/add-s6-overlay-with-ci
- Merge pull request [#20](y/issues/20) from Improwised/create-admin-command
- Merge pull request [#10](y/issues/10) from Improwised/init-setup


[Unreleased]: y/compare/v1.1.0...HEAD
[v1.1.0]: y/compare/v1.0.0...v1.1.0
