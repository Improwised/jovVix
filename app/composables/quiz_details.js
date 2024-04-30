export function useGetTime(targetTime) {
  const quizCreationTime = new Date(targetTime);

  const currentTimeInMs = Date.now();

  const timeSinceCreation = currentTimeInMs - quizCreationTime.getTime();

  const secondsSinceCreation = Math.round(timeSinceCreation / 1000);
  const minutesSinceCreation = Math.round(secondsSinceCreation / 60);
  const hoursSinceCreation = Math.round(minutesSinceCreation / 60);
  const daysSinceCreation = Math.round(hoursSinceCreation / 24);
  const monthsSinceCreation = Math.round(daysSinceCreation / 30);
  const yearsSinceCreation = Math.round(monthsSinceCreation / 12);

  let message;
  if (secondsSinceCreation < 60) {
    return "This quiz was created just now.";
  } else if (minutesSinceCreation < 60) {
    message = `This quiz was created ${minutesSinceCreation} minutes ago`;
  } else if (hoursSinceCreation < 24) {
    message = `This quiz was created ${hoursSinceCreation} hours ago`;
  } else if (daysSinceCreation < 30) {
    message = `This quiz was created ${daysSinceCreation} days ago`;
  } else if (monthsSinceCreation < 12) {
    message = `This quiz was created ${monthsSinceCreation} months ago`;
  } else {
    message = `This quiz was created ${yearsSinceCreation} years ago`;
  }

  return message;
}
