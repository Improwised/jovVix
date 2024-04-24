export default function useMillisToMinutesAndSeconds(millis, precision) {
  var minutes = Math.floor(millis / 60000);
  var seconds = (millis % 60000) / 1000;

  if (precision === undefined) {
    // Default precision is 0
    seconds = seconds.toFixed(0);
  } else {
    // Round to the specified precision
    seconds = seconds.toFixed(precision);
  }

  // If seconds rounded up to 60, carry over to minutes
  if (seconds === "60") {
    seconds = "00";
    minutes++;
  }

  if (minutes > 0) {
    return minutes + "m " + (seconds < 10 ? "0" : "") + seconds + "s";
  }
  return (seconds < 10 ? "0" : "") + seconds + "s";
}
