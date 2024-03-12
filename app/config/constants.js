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

  // Components
  Question: "Question",

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
  AdminDisconnectedMessage: "admin is disconnected, please wait...",

  // Errors
  Unauthorized: "unauthorized to access resource",
  CodeNotFound: "invitation code not found",
  ReloadRequired: "there was some error, please reload the page!!!",
  InvitationCodeNotFound: "invitation code not found",
  QuizSessionValidationFailed: "quiz-session-validation-failed",
  SessionWasCompleted: "session was completed",
};
