// Recover direct play links without mistaking name or authentication errors for
// a completed quiz. The join form preserves the code and collects a name.
export function joinRecoveryQuery(code, error) {
  const status = error?.status || error?.statusCode || error?.response?.status;
  const detail = error?.data?.data;
  let message = "Please enter your name to join this quiz.";
  if (status === 409) {
    message =
      "That name is already taken. Choose another name and submit to see suggestions.";
  } else if (status === 400 && detail === "invitation code not found") {
    message = "invitation code not found";
  } else if (status === 400) {
    message = "Please enter a valid name for this quiz.";
  } else if (status === 403) {
    message =
      typeof detail === "string" ? detail : "You cannot join this quiz.";
  } else if (status && status !== 401) {
    message = "We couldn't reconnect to the quiz. Please try again.";
  } else if (!status && error) {
    message = "We couldn't reconnect to the quiz. Please try again.";
  }
  return { code: String(code), error: message };
}
