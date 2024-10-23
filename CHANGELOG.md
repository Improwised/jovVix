<a name="v1.1.1"></a>
## v1.1.1 - 2024-10-18

### Bug Fixes
- responsive issue for invitation code display on waiting area ([#360](/Improwised/quizz-app/issues/360))
- continue redis channel after 'no player found' error. ([#356](/Improwised/quizz-app/issues/356))
- handle redis pubsub proper and remove global answer submission channel ([#335](/Improwised/quizz-app/issues/335))
- negative time for participants who join quiz after start. ([#330](/Improwised/quizz-app/issues/330))
- add ping for admin to prevent socket disconnection (Issue [#314](/Improwised/quizz-app/issues/314)) ([#328](/Improwised/quizz-app/issues/328))
- not store user data during password recovery ([#326](/Improwised/quizz-app/issues/326))
- check user count before quiz start ([#323](/Improwised/quizz-app/issues/323))
- change nginx config for debug socket close issue ([#321](/Improwised/quizz-app/issues/321))
- add ping for continues web socket connection ([#318](/Improwised/quizz-app/issues/318))
- duplicate user show in scoreboard with inconsistent score and accuracy. ([#315](/Improwised/quizz-app/issues/315))

### Code Refactoring
- change image of quiz background ([#378](/Improwised/quizz-app/issues/378))
- change images of winner UI ([#373](/Improwised/quizz-app/issues/373))
- increase width of bar and answer submission component not render in 5_sec_conter ([#368](/Improwised/quizz-app/issues/368))
- show questions order wise in reports ([#362](/Improwised/quizz-app/issues/362))
- create generic function ProcessAnalyticsData ([#358](/Improwised/quizz-app/issues/358))
- malfunction error ([#336](/Improwised/quizz-app/issues/336))

### Features
- add music in running quiz ([#380](/Improwised/quizz-app/issues/380))
- add some ui enhancement ([#375](/Improwised/quizz-app/issues/375))
- add edit questions of quiz ([#374](/Improwised/quizz-app/issues/374))
- add swagger documentation ([#371](/Improwised/quizz-app/issues/371))
- add avatar in all required area and change winner UI ([#370](/Improwised/quizz-app/issues/370))
- add verification of email ([#369](/Improwised/quizz-app/issues/369))
- add avatar and winner UI ([#366](/Improwised/quizz-app/issues/366))
- UI enhancement ([#365](/Improwised/quizz-app/issues/365))
- add eslint formating and remove unused code. ([#364](/Improwised/quizz-app/issues/364))
- add code block in question ([#361](/Improwised/quizz-app/issues/361))
- handle user inactivity while waiting for start quiz ([#340](/Improwised/quizz-app/issues/340))
- add questionAnalysis of user in report participants page ([#339](/Improwised/quizz-app/issues/339))
- add order in questions list of report analysis ([#334](/Improwised/quizz-app/issues/334))
- add pong message from server and add reconnect feature if user removed while waiting for start quiz ([#337](/Improwised/quizz-app/issues/337))
- display number of total question in question component ([#322](/Improwised/quizz-app/issues/322))
- add filteration and pagination in quiz analysis list ([#320](/Improwised/quizz-app/issues/320))


<a name="v1.1.0"></a>
## v1.1.0 - 2024-08-23


### Bug Fixes
- break functions
- added correct context
- terminate quiz on unexpected exit of admin and guest user join issue (issue [#282](/Improwised/quizz-app/issues/282)) ([#287](/Improwised/quizz-app/issues/287))
- take time period for each question from env rather than csv ([#180](/Improwised/quizz-app/issues/180))
- struct printing err, test cases
- changing docker envs
- Restrict Points Allocation to maximum 20 Points per question ([#163](/Improwised/quizz-app/issues/163))
- Add up and down migration to alter old data and make points 0 where points were 1 ([#243](/Improwised/quizz-app/issues/243))
- username and password not showing
- cleanup invitationcode store and userlist store when waiting space component unmounts ([#242](/Improwised/quizz-app/issues/242))
- changed runner to ubuntu-latest
- changed runner to ubuntu-latest
- award points for any correct option chosen in multi-answer questions (except survey questions) ([#213](/Improwised/quizz-app/issues/213))
- errors-caused-by-cookies are resolved ([#208](/Improwised/quizz-app/issues/208), [#204](/Improwised/quizz-app/issues/204)) ([#216](/Improwised/quizz-app/issues/216))
- modify admin analysis ([#207](/Improwised/quizz-app/issues/207))
- modify user anaytics ([#205](/Improwised/quizz-app/issues/205))
- admin-analysis ([#199](/Improwised/quizz-app/issues/199))
- added text color for buttons
- changed the static words to the camel case
- username should be maximum of 12 characters long and should not consider blankspace as character ([#174](/Improwised/quizz-app/issues/174))
- docker file is back to its previous state as wanted by the CI flow ([#181](/Improwised/quizz-app/issues/181))
- the docker compose functionality is now working properly with these changes as there was unwanted envs and introduced docker volume for redis and also changed ip of the containers ([#164](/Improwised/quizz-app/issues/164))
- username is now send by GetUserMeta whether it is kratos or guest user ([#278](/Improwised/quizz-app/issues/278))
- also count Not Attempted survey in accuracy ([#290](/Improwised/quizz-app/issues/290))
- Quesiton and Score Space being cut in the mobile view ([#197](/Improwised/quizz-app/issues/197))
- removed console.log
- modified changes in design
- changed final scoreboard route for admin
- fixed bug while submitting answer when admin refreshes the page
- cq and pre-commit
- authentication added for scoreboard api
- modified user scoreboard api
- remove extra variable and else condition
- change email in kratos also ([#291](/Improwised/quizz-app/issues/291))
- error
- session get after session completed
- remove href attribute
- change port in nginx.conf from 3000 to 4000 as per .env.example [#32](/Improwised/quizz-app/issues/32)
- added selfhosted runner to api ci
- added path condition in CI
- **accuracy-logic-to-correct-incorrect:** Modify accuracy formula to calculate accuracy based on the correct and incorrect answers count ([#274](/Improwised/quizz-app/issues/274))
- **accuracy-logic-userwise:** fix userwise accuracy logic to count accuracy based on points instead of correct/incorrect ([#240](/Improwised/quizz-app/issues/240))
- **participants-disconnect-issue:** Issue of users getting disconnected in Waiting area ([#272](/Improwised/quizz-app/issues/272))
- **score-logic-for-survey:** Add logic to award score for survey questions where score is provided and do not award for survey questions where score is not provided ([#236](/Improwised/quizz-app/issues/236))
- **type-column-migrations:** Add migrations to add new column 'type' and constraints and all related to that ([#241](/Improwised/quizz-app/issues/241))

### Code Refactoring
- use $fetch instead of usefetch in create quiz ([#303](/Improwised/quizz-app/issues/303))
- helper interface ([#302](/Improwised/quizz-app/issues/302))
- remove logic of storing user_played_quiz in cookie ([#298](/Improwised/quizz-app/issues/298))
- remove middleware from frontend and create quick user properly and remove unnecessary api calling for user data ([#297](/Improwised/quizz-app/issues/297))
- parse csv by header instead of row and column ([#288](/Improwised/quizz-app/issues/288))

### Features
- add-mutex and wrong answer highlight ([#309](/Improwised/quizz-app/issues/309))
- add kratos login flow in quiz app ([#271](/Improwised/quizz-app/issues/271))
- add kratos schema ([#275](/Improwised/quizz-app/issues/275))
- making the font of inviatation code larger ([#266](/Improwised/quizz-app/issues/266))
- show-username-when-user-submits-answer ([#233](/Improwised/quizz-app/issues/233))
- analytics board api for the admin and user ([#179](/Improwised/quizz-app/issues/179))
- test case for csv auth
- added test cases for quiz start functionality
- added routes for list-quiz and modified '/admin/quiz' route as 'admin/quiz/create-quiz'
- skip button added for admin to skip timer
- user score board api
- added layout for final scoreboard
- skip event
- admin will redirect when he/she tring to join their own quiz
- added  modified dockerfile and also added the ci for build and push for qizz-api [#28](/Improwised/quizz-app/issues/28)
- added s6-overlay to the quiz-app also modified dockerfile and also added the ci for build and push [#18](/Improwised/quizz-app/issues/18)

### Reverts
- debug logs for admin disconnection issue ([#312](/Improwised/quizz-app/issues/312))
- add debug log ([#313](/Improwised/quizz-app/issues/313))
- adding the sort by latest and paginations to reports page ([#311](/Improwised/quizz-app/issues/311))
- fix: some ui and add logger in quiz_operation to solve all-param-required
- Recheck staging, checking for valid approach

### Merged features
- Merge pull request [#155](/Improwised/quizz-app/issues/155) from Improwised/fix/button-text-color
- Merge pull request [#153](/Improwised/quizz-app/issues/153) from Improwised/fix/admin-scoreboard-route
- Merge pull request [#152](/Improwised/quizz-app/issues/152) from Improwised/feature/final-scoreboard-rank-calculation
- Merge pull request [#140](/Improwised/quizz-app/issues/140) from Improwised/fix/bug-answer-submit
- Merge pull request [#135](/Improwised/quizz-app/issues/135) from Improwised/test-cases/quiz-start
- Merge pull request [#136](/Improwised/quizz-app/issues/136) from Improwised/feature/display-user-to-admin-while-join
- Merge pull request [#131](/Improwised/quizz-app/issues/131) from Improwised/feature/uploaded-quiz-list
- Merge pull request [#128](/Improwised/quizz-app/issues/128) from Improwised/fix/next-question-button
- Merge pull request [#121](/Improwised/quizz-app/issues/121) from Improwised/fix/question-scoreboard
- Merge pull request [#117](/Improwised/quizz-app/issues/117) from Improwised/feature/finalscoreboard-time
- Merge pull request [#118](/Improwised/quizz-app/issues/118) from Improwised/feature/skip-timer-button
- Merge pull request [#108](/Improwised/quizz-app/issues/108) from Improwised/feature/finalscoreboard-admin
- Merge pull request [#105](/Improwised/quizz-app/issues/105) from Improwised/feature/finalscoreboard-user-score
- Merge pull request [#101](/Improwised/quizz-app/issues/101) from Improwised/final-scoreboard
- Merge pull request [#95](/Improwised/quizz-app/issues/95) from Improwised/start-demo
- Merge pull request [#93](/Improwised/quizz-app/issues/93) from Improwised/56-feat-create-a-page-where-admin-can-enter-question-csv
- Merge pull request [#91](/Improwised/quizz-app/issues/91) from Improwised/fix-null-found-in-duratin
- Merge pull request [#89](/Improwised/quizz-app/issues/89) from Improwised/45-user-and-admin-should-see-rank-board-after-submitting-the-answer
- Merge pull request [#88](/Improwised/quizz-app/issues/88) from Improwised/43-handle-skip-event
- Merge pull request [#87](/Improwised/quizz-app/issues/87) from Improwised/40-handle-timeout-of-question
- Merge pull request [#82](/Improwised/quizz-app/issues/82) from Improwised/39-handle-start-event
- Merge pull request [#81](/Improwised/quizz-app/issues/81) from Improwised/fix-all-param-required
- Merge pull request [#80](/Improwised/quizz-app/issues/80) from Improwised/revert-79-fix-all-param-required
- Merge pull request [#79](/Improwised/quizz-app/issues/79) from Improwised/fix-all-param-required
- Merge pull request [#78](/Improwised/quizz-app/issues/78) from Improwised/75-fix-change-nuxtsession-with-active-library
- Merge pull request [#67](/Improwised/quizz-app/issues/67) from Improwised/user_admin_playground_ui
- Merge pull request [#69](/Improwised/quizz-app/issues/69) from Improwised/revert-68-recheck-staging
- Merge pull request [#70](/Improwised/quizz-app/issues/70) from Improwised/recheck-staging
- Merge pull request [#68](/Improwised/quizz-app/issues/68) from Improwised/recheck-staging
- Merge pull request [#66](/Improwised/quizz-app/issues/66) from Improwised/fix/remove-export-env
- Merge pull request [#65](/Improwised/quizz-app/issues/65) from Improwised/add-systemEnvConfig
- Merge pull request [#36](/Improwised/quizz-app/issues/36) from Improwised/create-login-page
- Merge pull request [#52](/Improwised/quizz-app/issues/52) from Improwised/51-change-configs-after-rename-envexample-to-envdocker
- Merge pull request [#48](/Improwised/quizz-app/issues/48) from Improwised/swagger-and-route-change
- Merge pull request [#47](/Improwised/quizz-app/issues/47) from Improwised/46-reset-ports-in-envs
- Merge pull request [#34](/Improwised/quizz-app/issues/34) from Improwised/fix-migration-downs
- Merge pull request [#35](/Improwised/quizz-app/issues/35) from Improwised/change-env-as-per-new-urls
- Merge pull request [#33](/Improwised/quizz-app/issues/33) from Improwised/feat/nginx-port-changes
- Merge pull request [#31](/Improwised/quizz-app/issues/31) from Improwised/change-dev-server-config
- Merge pull request [#30](/Improwised/quizz-app/issues/30) from Improwised/quiz-migrations
- Merge pull request [#22](/Improwised/quizz-app/issues/22) from Improwised/create-auth-middleware
- Merge pull request [#26](/Improwised/quizz-app/issues/26) from Improwised/23-bug-fix-the-admin-creation-on-zsh-as-well-make-file-generate-for-migration
- Merge pull request [#29](/Improwised/quizz-app/issues/29) from Improwised/feat/app-api-dockerfile-ci
- Merge pull request [#27](/Improwised/quizz-app/issues/27) from Improwised/fix/change-runner
- Merge pull request [#25](/Improwised/quizz-app/issues/25) from Improwised/fix/change-context-ci
- Merge pull request [#24](/Improwised/quizz-app/issues/24) from Improwised/feat/add-s6-overlay-with-ci
- Merge pull request [#20](/Improwised/quizz-app/issues/20) from Improwised/create-admin-command
- Merge pull request [#10](/Improwised/quizz-app/issues/10) from Improwised/init-setup


[Unreleased]: https://git.pride.improwised.dev/Improwised/quizz-app/compare/v1.1.0...HEAD
[v1.1.0]: https://git.pride.improwised.dev/Improwised/quizz-app/compare/v1.0.0...v1.1.0
