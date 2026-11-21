export function accuracy(correct, questions) {
  return questions ? Math.round((correct / questions) * 100) : 0;
}

// Accuracy over every answer on the merged leaderboard, not an average of averages.
export function overallAccuracy(leaderboard = []) {
  const correct = leaderboard.reduce((sum, p) => sum + p.correct, 0);
  const questions = leaderboard.reduce((sum, p) => sum + p.questions, 0);
  return accuracy(correct, questions);
}

// Sessions are numbered by start time, the order the API returns them in.
export function sessionNumbers(sessions = []) {
  return new Map(sessions.map((s, i) => [s.id, i + 1]));
}

// ScoreBoardTable rows; the detail line says where each player's best play came from.
export function combinedRows(report) {
  const numbers = sessionNumbers(report?.sessions);
  return (report?.leaderboard || []).map((p) => ({
    ...p,
    detail: [
      `Session ${numbers.get(p.session_id) ?? "?"}`,
      `${p.correct}/${p.questions} correct (${accuracy(p.correct, p.questions)}%)`,
      p.plays > 1 ? `best of ${p.plays} plays` : "",
    ]
      .filter(Boolean)
      .join(" · "),
  }));
}
