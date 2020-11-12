package constants

const (
	// EmailPasswordReset email template for password reset
	EmailPasswordReset = `<span style="color: #848484; font-size: 18px; line-height: 21px;">
		Hi <strong>{{.FullName}}!</srong>
	</span>
	<br/>
	<br/>
	<span style="color: #848484; font-size: 14px; line-height: 21px;">
		<a href="{{.PasswordResetLink}}" rel="noopener" style="text-decoration: underline; color: #5ad8ff;" target="_blank">Reset password</a>
	</span>`
)
