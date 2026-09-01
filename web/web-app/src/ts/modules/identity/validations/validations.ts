export function isValidEmail(email: string): boolean {
	if (!email.trim()) {
		return false;
	}

	const emailRegex =
		/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$/;

	return emailRegex.test(email);
}

export function isValidPassword(password: string): boolean {
	if (!password || password.length < 8) {
		return false;
	}

	let hasUpper = false;
	let hasLower = false;
	let hasNumber = false;

	for (const char of password) {
		if (/[A-Z]/.test(char)) {
			hasUpper = true;
		}

		if (/[a-z]/.test(char)) {
			hasLower = true;
		}

		if (/[0-9]/.test(char)) {
			hasNumber = true;
		}
	}

	return hasUpper && hasLower && hasNumber;
}

export function isFormValidSignUp(
	email: string,
	password: string,
	confirm_password: string,
): boolean {
	const emailValid = isValidEmail(email);

	const passwordValid = isValidPassword(password);

	const confirmPasswordValid =
		isValidPassword(confirm_password) &&
		confirm_password === password;

	return emailValid && passwordValid && confirmPasswordValid;
}

export function isFormValidSignIn(
	email: string,
	password: string,
): boolean {
	const emailValid = isValidEmail(email);

	const passwordValid = password != "";

	return emailValid && passwordValid;
}