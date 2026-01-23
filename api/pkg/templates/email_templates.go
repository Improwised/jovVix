package templates

import "fmt"

// GenerateQuizShareEmail generates the HTML content for the quiz share email
func GenerateQuizShareEmail(quizLink, sharedBy string, permissions string) string {
	permissionsText := fmt.Sprintf("Permissions granted: <strong>%s</strong>", fmt.Sprintf("%v", permissions))

	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #eaeaea; border-radius: 8px; }
				.button { display: inline-block; padding: 10px 20px; margin: 20px 0; background-color: #007bff; color: #ffffff; text-decoration: none; border-radius: 5px; }
				.footer { font-size: 12px; color: #888; margin-top: 20px; }
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Hello,</h2>
				<p>You have been invited by <strong>%s</strong> to participate in the quiz.</p>
				<p>%s</p>
				<p>To participate in the quiz, click the link below:</p>
				<p><a href="%s" class="button" style="color:#ffffff !important; text-decoration:none;">Take the Quiz</a></p>
				<p>If the button doesn’t work, copy and paste this link into your browser:</p>
				<p><a href="%s">%s</a></p>
				<p class="footer">This invitation was sent to you by %s via the Quiz Platform.</p>
			</div>
		</body>
		</html>`, sharedBy, permissionsText, quizLink, quizLink, quizLink, sharedBy)
}
