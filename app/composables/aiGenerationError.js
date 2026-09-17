const errorStatus = (error) =>
  error?.statusCode ||
  error?.status ||
  error?.response?.status ||
  error?.data?.statusCode;

const errorText = (error) =>
  [error?.data?.data, error?.data?.message, error?.message]
    .filter(Boolean)
    .join(" ")
    .toLowerCase();

export const readAIGenerationError = (error) => {
  const status = errorStatus(error);
  const text = errorText(error);

  if (status === 429) {
    return {
      message: "The AI provider is busy. Wait a moment, then try again.",
      needsSettings: false,
    };
  }

  if (status === 504 || /timed out|timeout|deadline exceeded/.test(text)) {
    return {
      message:
        "The AI took too long to respond. Try again, or generate fewer questions.",
      needsSettings: false,
    };
  }

  if (
    /returned an empty response|could not read the questions|did not return any usable questions/.test(
      text
    )
  ) {
    return {
      message: "The AI did not return usable questions. Try again.",
      needsSettings: false,
    };
  }

  if ([400, 401, 403].includes(status)) {
    return {
      message:
        "Your AI settings need attention. Check the API key and selected model, then try again.",
      needsSettings: true,
    };
  }

  if (
    status === 502 ||
    status === 503 ||
    /could not reach the ai service/.test(text)
  ) {
    return {
      message:
        "The AI provider is temporarily unavailable. Try again in a moment.",
      needsSettings: false,
    };
  }

  return {
    message: "Could not generate questions. Try again.",
    needsSettings: false,
  };
};
