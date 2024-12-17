export default {
  // Middleware
  EmailOrPasswordMissing: "email and password are required",
  IncorrectCredentials: "Error or password incorrect, try again",

  // Defaults
  InvitationCode: "",
  UserIdentifier: "user",
  CurrentQuizIdentifier: "user_played_quiz",
  CountTill: 3000,
  ReadyMessage: "Ready to go...",
  AnswerSubmitted: "Answer submitted successfully",

  // Status
  Fail: "fail",
  Success: "success",
  Error: "error",

  // Quiz State
  Running: "running",
  Pause: "pause",

  // Components
  Question: "Question",
  Score: "Score",
  Waiting: "Waiting",

  // Events
  StartQuiz: "start_quiz",
  Authentication: "authentication",
  SentInvitaionCode: "send_invitation_code",
  GetQuestion: "send_question",
  Counter: "5_sec_counter",
  TerminateQuiz: "terminate_quiz",
  InvitationCodeValidation: "invitation_code_validation",
  RedirectToAdmin: "redirect_to_admin",
  NextQuestionAsk: "next_question",
  AdminDisconnected: "admin_is_disconnected",
  AskSkip: "ask_skip",
  AskForceSkip: "ask_force_skip",
  ShowScore: "show_score",
  AdminDisconnectedMessage: "admin is disconnected, please wait...",
  SkipTimer: "skip_timer",
  EventAnswerSubmittedByUser: "answer submitted by user",
  EventPing: "ping",
  PauseQuiz: "pause_quiz",
  PauseQuizMessage: "quiz is paused by admin, please wait",
  ResumeQuiz: "resume_quiz",
  ResumeQuizMessage: "quiz is resumed, you could play now",

  // Actions
  ActionAnserSubmittedByUser: "answer submitted by user",

  // Errors
  Unauthorized: "unauthorized to access resource",
  CodeNotFound: "invitation code not found",
  ReloadRequired: "there was some error, please reload the page!!!",
  InvitationCodeNotFound: "invitation code not found",
  QuizSessionValidationFailed: "quiz-session-validation-failed",
  SessionWasCompleted: "session was completed",
  NoAnswerFound: "please select an answer",

  // Messages
  CsvUploadSuccess: "file uploaded successfully",
};
