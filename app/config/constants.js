export default {
  // Middleware
  EmailOrPasswordMissing: "email and password are required",
  IncorrectCredentials: "Error or password incorrect, try again",

  // Defaults
  InvitationCode: "",
  UserIdentifier: "user",
  CountTill: 3000,
  ReadyMessage: "Ready to go...",

  // Status
  Fail: "fail",
  Success: "success",
  Error: "error",

  // Events
  StartQuiz: "start_quiz",
  Authentication: "authentication",
  SentInvitaionCode: "send_invitation_code",
  GetQuestion: "send_question",
  Counter: "5_sec_counter",
  TerminateQuiz: "terminate_quiz",
  InvitationCodeValidation: "invitation_code_validation",
  RedirectToAdmin: "redirect_to_admin",

  // Errors
  Unauthorized: "unauthorized to access resource",
  CodeNotFound: "invitation code not found",
  ReloadRequired: "there was some error, please reload the page!!!",
  InvitationCodeNotFound: "invitation code not found",
  QuizSessionValidationFailed: "quiz-session-validation-failed",
};
